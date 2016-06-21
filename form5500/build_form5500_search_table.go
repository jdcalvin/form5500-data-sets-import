package main

import (
	"fmt"
	"log"
	"strings"
  "database/sql"
)

const form5500Search string = "form_5500_search"

func buildTable(connection string, section string, years []string) {
	fmt.Println("Building form_5500_search table...")
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	executableStatements := make([]Statement, 0)

	for _, statement := range createSearchTable() {
		executableStatements = append(executableStatements, statement)
	}

	unionTables := make([]string, 0)
	for _, year := range years {
		unionTables = append(unionTables, selectLongFormTable(year, section))
		unionTables = append(unionTables, selectShortFormTable(year, section))
	}

	selectStatement := strings.Join(unionTables, "\n      UNION ALL\n")

	cols := ""
	for _, row := range tableMappings() {
		cols += row.alias + ","
	}

	cols += "table_origin"

	insertStatement := fmt.Sprintf("INSERT INTO form_5500_search (%[1]s) SELECT %[1]s FROM (\n%[2]s\n) as f_s;", cols, selectStatement)
	executableStatements = append(executableStatements, Statement{sql: insertStatement, description: "Inserting records into form_5500_search"})

	// - Set total assets on form_5500_search from schedule H, or I
	// - Set providers on form_5500_search from schedule C if applicable (long form only)
	//   based on service codes http://freeerisa.benefitspro.com/static/popups/legends.aspx#5500c09
	for _, year := range years {
		for _, statement := range updateFromSchedules(section, year) {
			executableStatements = append(executableStatements, statement)
		}
	}

	// - Create materialized view form5500_search_view
	for _, statement := range createMaterializedView() {
		executableStatements = append(executableStatements, Statement{sql: statement, description: "Creating materialized view"})
	}
	// - Create index for each column in form5500_search_view
	for _, row := range tableMappings() {
		executableStatements = append(executableStatements, buildIndexStatement(row.alias))
	}

	for _, statement := range executableStatements {
		fmt.Println(fmt.Sprintf("  - %s", statement.description))
		_, err = db.Exec(statement.sql)
		if err != nil {
      fmt.Println(statement)
			log.Fatal(err)
		}
	}
}

//private

func buildIndexStatement(field string) Statement {
	return Statement{
					sql: fmt.Sprintf("CREATE INDEX idx_%[1]s ON form5500_search_view (%[1]s);", field), 
					description: fmt.Sprintf("Creating index idx_%[1]s", field),
				}
}

func createSearchTable() []Statement {
	statements := make([]Statement, 0)
	statements = append(statements, Statement{sql: fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", form5500Search), 			description: "drop form5500_search table"})
	statements = append(statements, Statement{sql: fmt.Sprintf("CREATE TABLE %s (%s);", form5500Search, tableColumns()), 	description: "create form5500_search table"})
	return statements
}

func tableColumns() string {
	var cols string
	for _, row := range tableMappings() {
		cols += fmt.Sprintf("%s %s, ", row.alias, row.dataType)
	}
	var providerCols = []string{
		"rk_name", "rk_ein", "tpa_name", "tpa_ein", "advisor_name", "advisor_ein",
	}
	for _, col := range providerCols {
		cols += col + " text,"
	}
	cols += "table_origin text"
	return cols
}

func selectLongFormTable(year string, section string) string {
	statement := "   SELECT "
	for _, row := range tableMappings() {
		statement += fmt.Sprintf("%s as %s, ", row.longForm, row.alias)
	}
	statement += fmt.Sprintf("'%[1]s_%[2]s' as table_origin from f_5500_%[1]s_%[2]s as f_%[1]s", year, section)
	return statement
}

func selectShortFormTable(year string, section string) string {
	statement := "   SELECT "
	for _, row := range tableMappings() {
		statement += fmt.Sprintf("%s as %s, ", row.shortForm, row.alias)
	}
	statement += fmt.Sprintf("'sf_%[1]s_%[2]s' as table_origin from f_5500_sf_%[1]s_%[2]s as f_%[1]s_sf", year, section)
	return statement
}

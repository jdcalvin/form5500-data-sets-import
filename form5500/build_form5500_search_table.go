package form5500

import (
	"fmt"
	"log"
	"strings"
  "database/sql"
	utils "github.com/jdcalvin/form5500-data-sets-import/form5500/internal/utils"
)

const form5500Search string = "form_5500_search"

func buildTable(connection string, section string, years []string) {
	fmt.Println("Building form_5500_search table...")
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, statement := range buildStatements(section, years) {
		fmt.Println(fmt.Sprintf("  - %s", statement.description))
		_, err = db.Exec(statement.sql)
		if err != nil {
      fmt.Println(statement)
			log.Fatal(err)
		}
	}
}

//private

func buildStatements(section string, years []string) []Statement {
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
	for _, row := range utils.TableMappings() {
		cols += row.Alias + ","
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
	for _, row := range utils.TableMappings() {
		executableStatements = append(executableStatements, buildIndexStatement(row))
	}
	return executableStatements
	
}

func buildIndexStatement(mapping utils.Mapping) Statement {
	return Statement{
					sql: fmt.Sprintf("CREATE INDEX ON form5500_search_view (%[1]s);", mapping.IndexName()), 
					description: fmt.Sprintf("Creating index %[1]s", mapping.IndexName()),
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
	for _, row := range utils.TableMappings() {
		cols += fmt.Sprintf("%s %s, ", row.Alias, row.DataType)
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
	for _, row := range utils.TableMappings() {
		statement += fmt.Sprintf("%s as %s, ", row.LongForm, row.Alias)
	}
	statement += fmt.Sprintf("'%[1]s_%[2]s' as table_origin from f_5500_%[1]s_%[2]s as f_%[1]s", year, section)
	return statement
}

func selectShortFormTable(year string, section string) string {
	statement := "   SELECT "
	for _, row := range utils.TableMappings() {
		statement += fmt.Sprintf("%s as %s, ", row.ShortForm, row.Alias)
	}
	statement += fmt.Sprintf("'sf_%[1]s_%[2]s' as table_origin from f_5500_sf_%[1]s_%[2]s as f_%[1]s_sf", year, section)
	return statement
}

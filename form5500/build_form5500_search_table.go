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

	executableStatements := make([]string, 0)

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
	executableStatements = append(executableStatements, insertStatement)

	for _, year := range years {
		for _, statement := range updateFromSchedules(section, year) {
			executableStatements = append(executableStatements, statement)
		}
	}

	for _, statement := range createMaterializedView() {
		executableStatements = append(executableStatements, statement)
	}

	for _, statement := range executableStatements {
		_, err = db.Exec(statement)
		if err != nil {
      fmt.Println(statement)
			log.Fatal(err)
		}
	}
}

//private

func createSearchTable() []string {
	sqlLines := make([]string, 0)
	sqlLines = append(sqlLines, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", form5500Search))
	sqlLines = append(sqlLines, fmt.Sprintf("CREATE TABLE %s (%s);", form5500Search, tableColumns()))
	return sqlLines
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

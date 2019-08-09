package main

import (
	"fmt"
	"strings"

	utils "github.com/fiduciary-benchmarks/form5500/internal/utils"
)

const form5500Search string = "form_5500_search"

func rebuildSearchTable(section string, years []string) {
	fmt.Println("Building form_5500_search table...")

	for _, statement := range getRebuildStatements(section, years) {
		statement.Exec()
	}
}

//private

func getRebuildStatements(section string, years []string) []utils.SQLRunner {
	var executableStatements []utils.SQLRunner

	for _, statement := range getDropAndCreateSearchTableStatements() {
		executableStatements = append(executableStatements, statement)
	}

	var unionTables []string

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
	executableStatements = append(executableStatements, utils.SQLRunner{Statement: insertStatement, Description: "Inserting records into form_5500_search"})

	// - Set total assets on form_5500_search from schedule H, or I
	// - Set providers on form_5500_search from schedule C if applicable (long form only)
	//   based on service codes http://freeerisa.benefitspro.com/static/popups/legends.aspx#5500c09
	for _, year := range years {
		for _, statement := range getUpdateFromSchedulesStatements(section, year) {
			executableStatements = append(executableStatements, statement)
		}
	}

	// - Create materialized view form5500_search_view

	executableStatements = append(executableStatements, getCreateMaterializedViewStatement())

	// - Create index for each column in form5500_search_view
	for _, row := range utils.TableMappings() {
		executableStatements = append(executableStatements, getCreateIndexStatement(row))
	}
	return executableStatements
}

func getCreateIndexStatement(mapping utils.Mapping) utils.SQLRunner {
	return utils.SQLRunner{
		Statement:   fmt.Sprintf("CREATE INDEX %[1]s ON form5500_search_view (%[2]s);", mapping.IndexName(), mapping.Alias),
		Description: fmt.Sprintf("Creating index %[1]s", mapping.IndexName()),
	}
}

func getDropAndCreateSearchTableStatements() []utils.SQLRunner {
	var statements []utils.SQLRunner
	statements = append(statements, utils.SQLRunner{Statement: fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", form5500Search), Description: "drop form5500_search table"})
	statements = append(statements, utils.SQLRunner{Statement: fmt.Sprintf("CREATE TABLE %s (%s);", form5500Search, getSearchTableColumns()), Description: "create form5500_search table"})
	return statements
}

func getSearchTableColumns() string {
	var cols string
	for _, row := range utils.TableMappings() {
		cols += fmt.Sprintf("%s %s, ", row.Alias, row.DataType)
	}

	cols += "rk_id int, "
	var providerCols = []string{
		"rk_name", "rk_ein", "tpa_name", "tpa_ein", "advisor_name", "advisor_ein",
	}
	for _, col := range providerCols {
		cols += col + " text,"
	}

	var investmentTypes = []string{
		"inv_collective_trusts",
		"inv_separate_accounts",
		"inv_mutual_funds",
		"inv_general_accounts",
		"inv_company_stock",
	}
	for _, col := range investmentTypes {
		cols += col + " boolean,"
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

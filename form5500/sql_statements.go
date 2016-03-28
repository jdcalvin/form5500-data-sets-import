package main

import (
	"fmt"
)

func updateFromSchedules(section string, year string) []string {
	var executableStatements = []string{
		updateFromScheduleH(section, year),
		updateFromScheduleI(section, year),
		updateProviderFromScheduleCItem2(section, year, "rk", "'15','23', '60'"),
		updateProviderFromScheduleCItem3(section, year, "rk", "'15','23', '60'"),
		updateProviderFromScheduleCItem2(section, year, "advisor", "'26','27'"),
		updateProviderFromScheduleCItem3(section, year, "advisor", "'26','27'"),
	}
	return executableStatements
}

func createMaterializedView() []string {
	var executableStatements = []string{
		"DROP MATERIALIZED VIEW IF EXISTS form5500_search_view;",
		`CREATE MATERIALIZED VIEW form5500_search_view AS
  		SELECT *, to_tsvector(sponsor_name) || to_tsvector(sponsor_ein)
   			as sponsor_search
  		FROM (SELECT DISTINCT ON (sponsor_name, sponsor_ein) * FROM form_5500_search ORDER BY sponsor_name, sponsor_ein, date_received DESC) as foo;`,
		"CREATE INDEX idx_fts_sponsor ON form5500_search_view USING gin(sponsor_search);",
	}
	return executableStatements
}

//private

func updateFromScheduleH(section string, year string) string {
	joinField := "\"ACK_ID\""

	updateField := "\"TOT_ASSETS_EOY_AMT\""

	scheduleTable := fmt.Sprintf("f_sch_h_%s_%s", year, section)

	selectStatement := fmt.Sprintf("SELECT ack_id, %[5]s FROM form_5500_search  JOIN %[3]s ON %[3]s.%[4]s = form_5500_search.ack_id", year, section, scheduleTable, joinField, updateField)

	updateStatement := fmt.Sprintf("UPDATE form_5500_search as f SET total_assets=foo_1.%[1]s FROM (%[2]s) as foo_1 WHERE foo_1.ack_id=f.ack_id", updateField, selectStatement)
	return updateStatement
}

func updateFromScheduleI(section string, year string) string {
	joinField := "\"ACK_ID\""

	updateField := "\"SMALL_TOT_ASSETS_EOY_AMT\""

	scheduleTable := fmt.Sprintf("f_sch_i_%s_%s", year, section)

	selectStatement := fmt.Sprintf("SELECT ack_id, %[5]s FROM form_5500_search  JOIN %[3]s ON %[3]s.%[4]s = form_5500_search.ack_id", year, section, scheduleTable, joinField, updateField)

	updateStatement := fmt.Sprintf("UPDATE form_5500_search as f SET total_assets=foo_1.%[1]s FROM (%[2]s) as foo_1 WHERE foo_1.ack_id=f.ack_id", updateField, selectStatement)
	return updateStatement
}

func updateProviderFromScheduleCItem2(section string, year string, provider string, validCodes string) string {
	joinField := "\"ACK_ID\""
	joinField2 := "\"ROW_ORDER\""

	name := "\"PROVIDER_OTHER_NAME\""
	ein := "\"PROVIDER_OTHER_EIN\""

	scheduleTable := fmt.Sprintf("f_sch_c_part1_item2_%s_%s", year, section)

	whereClause := fmt.Sprintf("\"SERVICE_CODE\" IN (%s)", validCodes)

	codeTable := fmt.Sprintf("f_sch_c_part1_item2_codes_%s_%s", year, section)

	selectStatement := `SELECT ack_id, %[1]s, %[2]s FROM form_5500_search 
    JOIN %[3]s ON %[3]s.%[4]s = form_5500_search.ack_id
    -- codeTable.ack_id=scheduleTable.ack_id AND codeTable.row_order=scheduleTable.row_order
    JOIN %[5]s ON %[5]s.%[4]s = %[3]s.%[4]s AND %[5]s.%[6]s = %[3]s.%[6]s
    WHERE %[7]s 
  `
	selectStatement = fmt.Sprintf(selectStatement, name, ein, scheduleTable, joinField, codeTable, joinField2, whereClause)

	return fmt.Sprintf("UPDATE form_5500_search as f SET %[4]s_name=foo_1.%[1]s, %[4]s_ein=foo_1.%[2]s FROM (%[3]s) as foo_1 WHERE foo_1.ack_id=f.ack_id", name, ein, selectStatement, provider)
}

func updateProviderFromScheduleCItem3(section string, year string, provider string, validCodes string) string {
	joinField := "\"ACK_ID\""
	joinField2 := "\"ROW_ORDER\""

	name := "\"PROVIDER_INDIRECT_NAME\""
	ein := ""

	scheduleTable := fmt.Sprintf("f_sch_c_part1_item3_%s_%s", year, section)

	whereClause := fmt.Sprintf("\"SERVICE_CODE\" IN (%s)", validCodes)

	codeTable := fmt.Sprintf("f_sch_c_part1_item3_codes_%s_%s", year, section)

	selectStatement := `SELECT ack_id, %[1]s FROM form_5500_search 
    JOIN %[3]s ON %[3]s.%[4]s = form_5500_search.ack_id
    -- codeTable.ack_id=scheduleTable.ack_id AND codeTable.row_order=scheduleTable.row_order
    JOIN %[5]s ON %[5]s.%[4]s = %[3]s.%[4]s AND %[5]s.%[6]s = %[3]s.%[6]s
    WHERE %[7]s 
  `
	selectStatement = fmt.Sprintf(selectStatement, name, ein, scheduleTable, joinField, codeTable, joinField2, whereClause)

	return fmt.Sprintf("UPDATE form_5500_search as f SET %[3]s_name=foo_1.%[1]s FROM (%[2]s) as foo_1 WHERE foo_1.ack_id=f.ack_id", name, selectStatement, provider)
}

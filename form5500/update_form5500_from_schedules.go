package main

import (
	"database/sql"
	"fmt"
	"log"
)

func updateFromSchedules(connection string, section string, year string) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var executableStatements = []string{
		updateFromScheduleH(section, year),
		updateFromScheduleI(section, year),
		updateProviderFromScheduleC(section, year, "rk", "'14','15','21','23'"),
		updateProviderFromScheduleC(section, year, "advisor", "'26','27'"),
	}

	for _, statement := range executableStatements {

		_, err = db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//private

func updateFromScheduleH(section string, year string) string {
	fmt.Println("Updating total_assets from schedule H for", year)
	joinField := "\"ACK_ID\""

	updateField := "\"TOT_ASSETS_EOY_AMT\""

	scheduleTable := fmt.Sprintf("f_sch_h_%s_%s", year, section)

	selectStatement := fmt.Sprintf("SELECT ack_id, %[5]s FROM form_5500_search  JOIN %[3]s ON %[3]s.%[4]s = form_5500_search.ack_id", year, section, scheduleTable, joinField, updateField)

	updateStatement := fmt.Sprintf("UPDATE form_5500_search as f SET total_assets=foo_1.%[1]s FROM (%[2]s) as foo_1 WHERE foo_1.ack_id=f.ack_id", updateField, selectStatement)
	return updateStatement
}

func updateFromScheduleI(section string, year string) string {
	fmt.Println("Updating total_assets from schedule I for", year)
	joinField := "\"ACK_ID\""

	updateField := "\"SMALL_TOT_ASSETS_EOY_AMT\""

	scheduleTable := fmt.Sprintf("f_sch_i_%s_%s", year, section)

	selectStatement := fmt.Sprintf("SELECT ack_id, %[5]s FROM form_5500_search  JOIN %[3]s ON %[3]s.%[4]s = form_5500_search.ack_id", year, section, scheduleTable, joinField, updateField)

	updateStatement := fmt.Sprintf("UPDATE form_5500_search as f SET total_assets=foo_1.%[1]s FROM (%[2]s) as foo_1 WHERE foo_1.ack_id=f.ack_id", updateField, selectStatement)
	return updateStatement
}

func updateProviderFromScheduleC(section string, year string, provider string, validCodes string) string {
	fmt.Printf("Updating provider %s from schedule c table for year %s\n", provider, year)
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

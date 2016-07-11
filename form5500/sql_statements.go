package main

import (
	"fmt"
)

func updateFromSchedules(section string, year string) []Statement {
	var executableStatements = []Statement{
		{
      sql: updateFromScheduleH(section, year),
      description: fmt.Sprintf("Add info from schedule H %s", year),
    },
		{
      sql: updateFromScheduleI(section, year), 
      description: fmt.Sprintf("Add info from schedule I %s", year),
    },
		{
      sql: updateProviderFromScheduleCItem2(section, year, "rk", "'15','23', '60'"),
      description: fmt.Sprintf("Determining Recordkeeper from schedule C item 2 %s", year),
    },
		{
      sql: updateProviderFromScheduleCItem3(section, year, "rk", "'15','23', '60'"),
      description: fmt.Sprintf("Determining Recordkeeper from schedule C item 3 %s", year),
    },
		{
      sql: updateProviderFromScheduleCItem2(section, year, "advisor", "'26','27'"),
      description: fmt.Sprintf("Determining Advisor from schedule C item 2 %s", year),
    },
		{
      sql: updateProviderFromScheduleCItem3(section, year, "advisor", "'26','27'"),
      description: fmt.Sprintf("Determining Advisor from schedule C item 2 %s", year),
    },
	}
	return executableStatements
}

func createMaterializedView() []string {
	var executableStatements = []string{
		"DROP MATERIALIZED VIEW IF EXISTS form5500_search_view;",
		`CREATE MATERIALIZED VIEW form5500_search_view AS
  		SELECT *, to_tsvector(sponsor_name) || to_tsvector(sponsor_ein)
   			as sponsor_search
  		FROM (SELECT DISTINCT ON (plan_num, sponsor_ein) * FROM form_5500_search ORDER BY plan_num, sponsor_ein, date_received DESC) as foo;`,
		"CREATE INDEX idx_fts_sponsor ON form5500_search_view USING gin(sponsor_search);",
	}
	return executableStatements
}

func createZipCodesTable() []string {
  var executableStatements = []string {
    "DROP TABLE IF EXISTS zip_codes",
    `CREATE TABLE zip_codes (
      zip integer,
      city varchar(256),
      state varchar(2),
      latitude double precision,
      longitude double precision,
      timezone integer,
      dst integer
    )`,
  }
  return executableStatements
}

func createZipCodeSearchFunction() []string {
  // This may need to be installed from some other source on the server
  var executableStatements= []string {
    "CREATE EXTENSION IF NOT EXISTS cube",
    "CREATE EXTENSION IF NOT EXISTS earthdistance",
    `DROP FUNCTION IF EXISTS udf_distance_in_miles_from_zip(_miles integer, _zip integer);
    CREATE OR REPLACE FUNCTION udf_distance_in_miles_from_zip(_miles integer, _zip integer)
      RETURNS table (
        zip integer,
        city varchar(256),
        state varchar(2),
        latitude double precision,
        longitude double precision,
        distance double precision
      ) AS
      $$
        DECLARE
          _lon double precision;
          _lat double precision;
        BEGIN
          _lon := (SELECT zip_codes.longitude FROM zip_codes WHERE zip_codes.zip = _zip);
          _lat := (SELECT zip_codes.latitude FROM zip_codes WHERE zip_codes.zip = _zip);
          RETURN QUERY
          SELECT zip_codes.zip, zip_codes.state, zip_codes.city, zip_codes.latitude, zip_codes.longitude,
            point(_lon, _lat) <@> point(zip_codes.longitude, zip_codes.latitude)::point as distance FROM zip_codes
            WHERE (point(_lon, _lat) <@> point(zip_codes.longitude, zip_codes.latitude) < _miles);
        END
      $$
    LANGUAGE 'plpgsql' STABLE;`,
  }
  return executableStatements
}

//private

func updateFromScheduleH(section string, year string) string {
	joinField := "\"ACK_ID\""


	scheduleTable := fmt.Sprintf("f_sch_h_%s_%s", year, section)

	selectStatement := fmt.Sprintf("SELECT * FROM form_5500_search  JOIN %[3]s ON %[3]s.%[4]s = form_5500_search.ack_id", year, section, scheduleTable, joinField)

  // cast numeric value from investment types to a boolean true or NULL
  updateStatement := `
    UPDATE form_5500_search as f 
      SET total_assets = "TOT_ASSETS_EOY_AMT",
          inv_collective_trusts = NULLIF(substring(abs("INT_COMMON_TR_EOY_AMT")::varchar, 1,1),'')::int::boolean,
          inv_separate_accounts = NULLIF(substring(abs("INT_POOL_SEP_ACCT_EOY_AMT")::varchar, 1,1),'')::int::boolean,
          inv_mutual_funds = NULLIF(substring(abs("INT_REG_INVST_CO_EOY_AMT")::varchar, 1,1),'')::int::boolean,
          inv_general_accounts = NULLIF(substring(abs("INS_CO_GEN_ACCT_EOY_AMT")::varchar, 1,1),'')::int::boolean,
          inv_company_stock = NULLIF(substring(abs("EMPLR_SEC_EOY_AMT")::varchar, 1,1),'')::int::boolean
    FROM (%[1]s) as foo_1 WHERE foo_1.ack_id=f.ack_id
  `
	return fmt.Sprintf(updateStatement, selectStatement)
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

package main

import (
  "fmt"
  "strings"
)

func ExampleUpdateFromScheduleH() {
  s := updateFromScheduleH("latest", "2013")
  s = strings.Replace(s,"\n"," ", -1)
  fmt.Println(s)
  // Output:
  // UPDATE form_5500_search as f        SET total_assets = "TOT_ASSETS_EOY_AMT",           inv_collective_trusts = NULLIF(substring(abs("INT_COMMON_TR_EOY_AMT")::varchar, 1,1),'')::int::boolean,           inv_separate_accounts = NULLIF(substring(abs("INT_POOL_SEP_ACCT_EOY_AMT")::varchar, 1,1),'')::int::boolean,           inv_mutual_funds = NULLIF(substring(abs("INT_REG_INVST_CO_EOY_AMT")::varchar, 1,1),'')::int::boolean,           inv_general_accounts = NULLIF(substring(abs("INS_CO_GEN_ACCT_EOY_AMT")::varchar, 1,1),'')::int::boolean,           inv_company_stock = NULLIF(substring(abs("EMPLR_SEC_EOY_AMT")::varchar, 1,1),'')::int::boolean     FROM (SELECT * FROM form_5500_search  JOIN f_sch_h_2013_latest ON f_sch_h_2013_latest."ACK_ID" = form_5500_search.ack_id) as foo_1 WHERE foo_1.ack_id=f.ack_id
}

// Should return contents of file sql/schedule_c_providers/create_table.sql (concat for test)
func ExampleCreateScheduleCProvider() {
  s := createScheduleCProvider()
  fmt.Println(s[0:20])
  // Output: 
  // DROP TABLE schedule_
}
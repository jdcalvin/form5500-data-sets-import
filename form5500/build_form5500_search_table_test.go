package main

import (
  "fmt"
)

func ExampleTableColumns() {
  fmt.Println(tableColumns())
  // Output:
  // ack_id text, date_received text, effective_date timestamp, plan_name text, plan_num text, sponsor_name text, sponsor_address_1 text, sponsor_address_2 text, sponsor_city text, sponsor_state text, sponsor_zip text, sponsor_ein text, sponsor_phone text, sponsor_naics text, pension_benefit_codes text, admin_name text, admin_address_1 text, admin_address_2 text, admin_city text, admin_state text, admin_zip text, admin_ein text, admin_phone text, active_participants int, total_assets numeric(19,2), rk_name text,rk_ein text,tpa_name text,tpa_ein text,advisor_name text,advisor_ein text,inv_collective_trusts boolean,inv_separate_accounts boolean,inv_mutual_funds boolean,inv_general_accounts boolean,inv_company_stock boolean,table_origin text
}

func ExampleSelectLongFormTable() {
  fmt.Println(selectLongFormTable("2014", "latest"))
  // Output:
  // SELECT "ACK_ID" as ack_id, "DATE_RECEIVED" as date_received, "PLAN_EFF_DATE" as effective_date, "PLAN_NAME" as plan_name, "SPONS_DFE_PN" as plan_num, "SPONSOR_DFE_NAME" as sponsor_name, "SPONS_DFE_MAIL_US_ADDRESS1" as sponsor_address_1, "SPONS_DFE_MAIL_US_ADDRESS2" as sponsor_address_2, "SPONS_DFE_MAIL_US_CITY" as sponsor_city, "SPONS_DFE_MAIL_US_STATE" as sponsor_state, "SPONS_DFE_MAIL_US_ZIP" as sponsor_zip, "SPONS_DFE_EIN" as sponsor_ein, "SPONS_DFE_PHONE_NUM" as sponsor_phone, "BUSINESS_CODE" as sponsor_naics, "TYPE_PENSION_BNFT_CODE" as pension_benefit_codes, "ADMIN_NAME" as admin_name, "ADMIN_US_ADDRESS1" as admin_address_1, "ADMIN_US_ADDRESS2" as admin_address_2, "ADMIN_US_CITY" as admin_city, "ADMIN_US_STATE" as admin_state, "ADMIN_US_ZIP" as admin_zip, "ADMIN_EIN" as admin_ein, "ADMIN_PHONE_NUM" as admin_phone, "PARTCP_ACCOUNT_BAL_CNT" as active_participants, 0 as total_assets, '2014_latest' as table_origin from f_5500_2014_latest as f_2014
}

func ExampleSelectShortFormTable() {
  fmt.Println(selectShortFormTable("2014", "latest"))
  // Output:
  // SELECT "ACK_ID" as ack_id, "DATE_RECEIVED" as date_received, "SF_PLAN_EFF_DATE" as effective_date, "SF_PLAN_NAME" as plan_name, "SF_PLAN_NUM" as plan_num, "SF_SPONSOR_NAME" as sponsor_name, "SF_SPONS_US_ADDRESS1" as sponsor_address_1, "SF_SPONS_US_ADDRESS2" as sponsor_address_2, "SF_SPONS_US_CITY" as sponsor_city, "SF_SPONS_US_STATE" as sponsor_state, "SF_SPONS_US_ZIP" as sponsor_zip, "SF_SPONS_EIN" as sponsor_ein, "SF_SPONS_PHONE_NUM" as sponsor_phone, "SF_BUSINESS_CODE" as sponsor_naics, "SF_TYPE_PENSION_BNFT_CODE" as pension_benefit_codes, "SF_ADMIN_NAME" as admin_name, "SF_ADMIN_US_ADDRESS1" as admin_address_1, "SF_ADMIN_US_ADDRESS2" as admin_address_2, "SF_ADMIN_US_CITY" as admin_city, "SF_ADMIN_US_STATE" as admin_state, "SF_ADMIN_US_ZIP" as admin_zip, "SF_ADMIN_EIN" as admin_ein, "SF_ADMIN_PHONE_NUM" as admin_phone, "SF_PARTCP_ACCOUNT_BAL_CNT" as active_participants, "SF_NET_ASSETS_EOY_AMT" as total_assets, 'sf_2014_latest' as table_origin from f_5500_sf_2014_latest as f_2014_sf
}

func ExampleCreateSearchTable() {
  s := createSearchTable()
  fmt.Println(s[0].Description)
  fmt.Println(s[0].Sql)
  fmt.Println(s[1].Description)
  fmt.Println(s[1].Sql)
  // Output: 
  // drop form5500_search table
  // DROP TABLE IF EXISTS form_5500_search CASCADE;
  // create form5500_search table
  // CREATE TABLE form_5500_search (ack_id text, date_received text, effective_date timestamp, plan_name text, plan_num text, sponsor_name text, sponsor_address_1 text, sponsor_address_2 text, sponsor_city text, sponsor_state text, sponsor_zip text, sponsor_ein text, sponsor_phone text, sponsor_naics text, pension_benefit_codes text, admin_name text, admin_address_1 text, admin_address_2 text, admin_city text, admin_state text, admin_zip text, admin_ein text, admin_phone text, active_participants int, total_assets numeric(19,2), rk_name text,rk_ein text,tpa_name text,tpa_ein text,advisor_name text,advisor_ein text,inv_collective_trusts boolean,inv_separate_accounts boolean,inv_mutual_funds boolean,inv_general_accounts boolean,inv_company_stock boolean,table_origin text);
}

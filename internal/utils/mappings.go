package utils

import (
	"fmt"
)
// Mapping is used to map column names in short/long form tables to an alias
// stored on form5500_search_view
type Mapping struct {
	LongForm  string
	ShortForm string
	Alias     string
	DataType  string
}

// TableMappings returns array of Mapping instances
// Mapping{} should not be called in usage, declare all instances within this method
func TableMappings() []Mapping {
	
	mappings := []Mapping{
		{
			LongForm:  "\"ACK_ID\"",
			ShortForm: "\"ACK_ID\"",
			Alias:     "ack_id",
			DataType:  "text",
		},
		{
			LongForm:  "\"DATE_RECEIVED\"",
			ShortForm: "\"DATE_RECEIVED\"",
			Alias:     "date_received",
			DataType:  "text",
		},
		{
			LongForm:  "\"PLAN_EFF_DATE\"",
			ShortForm: "\"SF_PLAN_EFF_DATE\"",
			Alias:     "effective_date",
			DataType:  "timestamp",
		},
		{
			LongForm:  "\"PLAN_NAME\"",
			ShortForm: "\"SF_PLAN_NAME\"",
			Alias:     "plan_name",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_PN\"",
			ShortForm: "\"SF_PLAN_NUM\"",
			Alias:     "plan_num",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONSOR_DFE_NAME\"",
			ShortForm: "\"SF_SPONSOR_NAME\"",
			Alias:     "sponsor_name",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_SIGNED_NAME\"",
			ShortForm: "\"SF_SPONS_SIGNED_NAME\"",
			Alias:     "sponsor_signed_name",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_MAIL_US_ADDRESS1\"",
			ShortForm: "\"SF_SPONS_US_ADDRESS1\"",
			Alias:     "sponsor_address_1",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_MAIL_US_ADDRESS2\"",
			ShortForm: "\"SF_SPONS_US_ADDRESS2\"",
			Alias:     "sponsor_address_2",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_MAIL_US_CITY\"",
			ShortForm: "\"SF_SPONS_US_CITY\"",
			Alias:     "sponsor_city",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_MAIL_US_STATE\"",
			ShortForm: "\"SF_SPONS_US_STATE\"",
			Alias:     "sponsor_state",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_MAIL_US_ZIP\"",
			ShortForm: "\"SF_SPONS_US_ZIP\"",
			Alias:     "sponsor_zip",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_EIN\"",
			ShortForm: "\"SF_SPONS_EIN\"",
			Alias:     "sponsor_ein",
			DataType:  "text",
		},
		{
			LongForm:  "\"SPONS_DFE_PHONE_NUM\"",
			ShortForm: "\"SF_SPONS_PHONE_NUM\"",
			Alias:     "sponsor_phone",
			DataType:  "text",
		},
		{
			LongForm:  "\"BUSINESS_CODE\"",
			ShortForm: "\"SF_BUSINESS_CODE\"",
			Alias:     "sponsor_naics",
			DataType:  "text",
		},
		{
			LongForm:  "\"TYPE_PENSION_BNFT_CODE\"",
			ShortForm: "\"SF_TYPE_PENSION_BNFT_CODE\"",
			Alias:     "pension_benefit_codes",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_NAME\"",
			ShortForm: "\"SF_ADMIN_NAME\"",
			Alias:     "admin_name",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_SIGNED_NAME\"",
			ShortForm: "\"SF_ADMIN_SIGNED_NAME\"",
			Alias:     "admin_signed_name",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_US_ADDRESS1\"",
			ShortForm: "\"SF_ADMIN_US_ADDRESS1\"",
			Alias:     "admin_address_1",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_US_ADDRESS2\"",
			ShortForm: "\"SF_ADMIN_US_ADDRESS2\"",
			Alias:     "admin_address_2",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_US_CITY\"",
			ShortForm: "\"SF_ADMIN_US_CITY\"",
			Alias:     "admin_city",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_US_STATE\"",
			ShortForm: "\"SF_ADMIN_US_STATE\"",
			Alias:     "admin_state",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_US_ZIP\"",
			ShortForm: "\"SF_ADMIN_US_ZIP\"",
			Alias:     "admin_zip",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_EIN\"",
			ShortForm: "\"SF_ADMIN_EIN\"",
			Alias:     "admin_ein",
			DataType:  "text",
		},
		{
			LongForm:  "\"ADMIN_PHONE_NUM\"",
			ShortForm: "\"SF_ADMIN_PHONE_NUM\"",
			Alias:     "admin_phone",
			DataType:  "text",
		},
		{
			LongForm:  "\"PARTCP_ACCOUNT_BAL_CNT\"",
			ShortForm: "\"SF_PARTCP_ACCOUNT_BAL_CNT\"",
			Alias:     "active_participants",
			DataType:  "int",
		},
		{
			LongForm:  "0",
			ShortForm: "\"SF_NET_ASSETS_EOY_AMT\"",
			Alias:     "total_assets",
			DataType:  "numeric(19,2)",
		},
	}

	return mappings
}

// IndexName transforms alias as "idx_column_name"
func (m Mapping) IndexName() string {
	return fmt.Sprintf("idx_%s", m.Alias)
} 

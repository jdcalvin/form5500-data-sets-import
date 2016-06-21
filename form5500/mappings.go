package main

type Mapping struct {
	longForm  string
	shortForm string
	alias     string
	dataType  string
}

func tableMappings() []Mapping {
	mappings := []Mapping{
		{
			longForm:  "\"ACK_ID\"",
			shortForm: "\"ACK_ID\"",
			alias:     "ack_id",
			dataType:  "text",
		},
		{
			longForm:  "\"DATE_RECEIVED\"",
			shortForm: "\"DATE_RECEIVED\"",
			alias:     "date_received",
			dataType:  "text",
		},
		{
			longForm:  "\"PLAN_EFF_DATE\"",
			shortForm: "\"SF_PLAN_EFF_DATE\"",
			alias:     "effective_date",
			dataType:  "timestamp",
		},
		{
			longForm:  "\"PLAN_NAME\"",
			shortForm: "\"SF_PLAN_NAME\"",
			alias:     "plan_name",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_PN\"",
			shortForm: "\"SF_PLAN_NUM\"",
			alias:     "plan_num",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONSOR_DFE_NAME\"",
			shortForm: "\"SF_SPONSOR_NAME\"",
			alias:     "sponsor_name",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_MAIL_US_ADDRESS1\"",
			shortForm: "\"SF_SPONS_US_ADDRESS1\"",
			alias:     "sponsor_address_1",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_MAIL_US_ADDRESS2\"",
			shortForm: "\"SF_SPONS_US_ADDRESS2\"",
			alias:     "sponsor_address_2",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_MAIL_US_CITY\"",
			shortForm: "\"SF_SPONS_US_CITY\"",
			alias:     "sponsor_city",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_MAIL_US_STATE\"",
			shortForm: "\"SF_SPONS_US_STATE\"",
			alias:     "sponsor_state",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_MAIL_US_ZIP\"",
			shortForm: "\"SF_SPONS_US_ZIP\"",
			alias:     "sponsor_zip",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_EIN\"",
			shortForm: "\"SF_SPONS_EIN\"",
			alias:     "sponsor_ein",
			dataType:  "text",
		},
		{
			longForm:  "\"SPONS_DFE_PHONE_NUM\"",
			shortForm: "\"SF_SPONS_PHONE_NUM\"",
			alias:     "sponsor_phone",
			dataType:  "text",
		},
		{
			longForm:  "\"BUSINESS_CODE\"",
			shortForm: "\"SF_BUSINESS_CODE\"",
			alias:     "sponsor_naics",
			dataType:  "text",
		},
		{
			longForm:  "\"TYPE_PENSION_BNFT_CODE\"",
			shortForm: "\"SF_TYPE_PENSION_BNFT_CODE\"",
			alias:     "pension_benefit_codes",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_NAME\"",
			shortForm: "\"SF_ADMIN_NAME\"",
			alias:     "admin_name",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_US_ADDRESS1\"",
			shortForm: "\"SF_ADMIN_US_ADDRESS1\"",
			alias:     "admin_address_1",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_US_ADDRESS2\"",
			shortForm: "\"SF_ADMIN_US_ADDRESS2\"",
			alias:     "admin_address_2",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_US_CITY\"",
			shortForm: "\"SF_ADMIN_US_CITY\"",
			alias:     "admin_city",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_US_STATE\"",
			shortForm: "\"SF_ADMIN_US_STATE\"",
			alias:     "admin_state",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_US_ZIP\"",
			shortForm: "\"SF_ADMIN_US_ZIP\"",
			alias:     "admin_zip",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_EIN\"",
			shortForm: "\"SF_ADMIN_EIN\"",
			alias:     "admin_ein",
			dataType:  "text",
		},
		{
			longForm:  "\"ADMIN_PHONE_NUM\"",
			shortForm: "\"SF_ADMIN_PHONE_NUM\"",
			alias:     "admin_phone",
			dataType:  "text",
		},
		{
			longForm:  "\"PARTCP_ACCOUNT_BAL_CNT\"",
			shortForm: "\"SF_PARTCP_ACCOUNT_BAL_CNT\"",
			alias:     "active_participants",
			dataType:  "int",
		},
		{
			longForm:  "0",
			shortForm: "\"SF_NET_ASSETS_EOY_AMT\"",
			alias:     "total_assets",
			dataType:  "numeric(19,2)",
		},
	}

	return mappings
}

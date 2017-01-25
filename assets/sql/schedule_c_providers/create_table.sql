DROP TABLE schedule_c_providers;
CREATE TABLE schedule_c_providers
  (
    ack_id varchar(256),
    name varchar(256),
    ein varchar(256),
    code varchar(256),
    relation varchar(256),
    matched_ein integer
  );
CREATE INDEX ix_sch_c_provider_view_name ON schedule_c_providers (name);
CREATE INDEX ix_sch_c_provider_view_code ON schedule_c_providers (code);
CREATE INDEX ix_sch_c_provider_view_ein ON schedule_c_providers (ein);

CREATE MATERIALIZED VIEW schedule_c_providers AS
(
  (
    SELECT f_sch_c_part1_item3_2013_latest."ACK_ID" as ack_id, "PROVIDER_INDIRECT_NAME" as name, NULL as ein, "SERVICE_CODE" as code, NULL as relation
      FROM f_sch_c_part1_item3_2013_latest
      JOIN f_sch_c_part1_item3_codes_2013_latest ON f_sch_c_part1_item3_2013_latest."ACK_ID"=f_sch_c_part1_item3_codes_2013_latest."ACK_ID"
                                                AND f_sch_c_part1_item3_2013_latest."ROW_ORDER"=f_sch_c_part1_item3_codes_2013_latest."ROW_ORDER"
  )
  UNION ALL
  ( 
    SELECT f_sch_c_part1_item3_2014_latest."ACK_ID" as ack_id, "PROVIDER_INDIRECT_NAME" as name, NULL as ein,  "SERVICE_CODE" as code, NULL as relation
      FROM f_sch_c_part1_item3_2014_latest
      JOIN f_sch_c_part1_item3_codes_2014_latest ON f_sch_c_part1_item3_2014_latest."ACK_ID"=f_sch_c_part1_item3_codes_2014_latest."ACK_ID"
                                                AND f_sch_c_part1_item3_2014_latest."ROW_ORDER"=f_sch_c_part1_item3_codes_2014_latest."ROW_ORDER"
  )
  UNION ALL
  (
    SELECT f_sch_c_part1_item3_2015_latest."ACK_ID" as ack_id, "PROVIDER_INDIRECT_NAME" as name, NULL as ein,  "SERVICE_CODE" as code, NULL as relation
      FROM f_sch_c_part1_item3_2015_latest
      JOIN f_sch_c_part1_item3_codes_2015_latest ON f_sch_c_part1_item3_2015_latest."ACK_ID"=f_sch_c_part1_item3_codes_2015_latest."ACK_ID"
                                                AND f_sch_c_part1_item3_2015_latest."ROW_ORDER"=f_sch_c_part1_item3_codes_2015_latest."ROW_ORDER"
  )
  UNION ALL

  (
    SELECT f_sch_c_part1_item2_2013_latest."ACK_ID" as ack_id, "PROVIDER_OTHER_NAME" as name, "PROVIDER_OTHER_EIN" as ein, "SERVICE_CODE" as code, "PROVIDER_OTHER_RELATION" as relation
      FROM f_sch_c_part1_item2_2013_latest
      JOIN f_sch_c_part1_item2_codes_2013_latest ON f_sch_c_part1_item2_2013_latest."ACK_ID"=f_sch_c_part1_item2_codes_2013_latest."ACK_ID"
                                                AND f_sch_c_part1_item2_2013_latest."ROW_ORDER"=f_sch_c_part1_item2_codes_2013_latest."ROW_ORDER"
  )
  UNION ALL
  (
    SELECT f_sch_c_part1_item2_2014_latest."ACK_ID" as ack_id, "PROVIDER_OTHER_NAME" as name, "PROVIDER_OTHER_EIN" as ein, "SERVICE_CODE" as code, "PROVIDER_OTHER_RELATION" as relation
      FROM f_sch_c_part1_item2_2014_latest
      JOIN f_sch_c_part1_item2_codes_2014_latest ON f_sch_c_part1_item2_2014_latest."ACK_ID"=f_sch_c_part1_item2_codes_2014_latest."ACK_ID"
                                                AND f_sch_c_part1_item2_2014_latest."ROW_ORDER"=f_sch_c_part1_item2_codes_2014_latest."ROW_ORDER"
  )
  UNION ALL
  (
    SELECT f_sch_c_part1_item2_2014_latest."ACK_ID" as ack_id, "PROVIDER_OTHER_NAME" as name, "PROVIDER_OTHER_EIN" as ein, "SERVICE_CODE" as code, "PROVIDER_OTHER_RELATION" as relation
      FROM f_sch_c_part1_item2_2014_latest
      JOIN f_sch_c_part1_item2_codes_2014_latest ON f_sch_c_part1_item2_2014_latest."ACK_ID"=f_sch_c_part1_item2_codes_2014_latest."ACK_ID"
                                                AND f_sch_c_part1_item2_2014_latest."ROW_ORDER"=f_sch_c_part1_item2_codes_2014_latest."ROW_ORDER"

  )
);

CREATE INDEX ix_sch_c_provider_view_name ON schedule_c_providers (name);
CREATE INDEX ix_sch_c_provider_view_code ON schedule_c_providers (code);
CREATE INDEX ix_sch_c_provider_view_ein ON schedule_c_providers (ein);
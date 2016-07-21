(
  SELECT f_sch_c_part1_item2_%[1]s_latest."ACK_ID" as ack_id, "PROVIDER_OTHER_NAME" as name, "PROVIDER_OTHER_EIN" as ein, "SERVICE_CODE" as code, "PROVIDER_OTHER_RELATION" as relation
    FROM f_sch_c_part1_item2_%[1]s_latest
    JOIN f_sch_c_part1_item2_codes_%[1]s_latest ON f_sch_c_part1_item2_%[1]s_latest."ACK_ID"=f_sch_c_part1_item2_codes_%[1]s_latest."ACK_ID"
                                              AND f_sch_c_part1_item2_%[1]s_latest."ROW_ORDER"=f_sch_c_part1_item2_codes_%[1]s_latest."ROW_ORDER"
)
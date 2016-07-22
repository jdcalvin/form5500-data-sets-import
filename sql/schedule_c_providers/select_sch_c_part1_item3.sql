(
  SELECT f_sch_c_part1_item3_%[1]s_latest."ACK_ID" as ack_id, "PROVIDER_INDIRECT_NAME" as name, NULL as ein, "SERVICE_CODE" as code, NULL as relation
    FROM f_sch_c_part1_item3_%[1]s_latest
    JOIN f_sch_c_part1_item3_codes_%[1]s_latest ON f_sch_c_part1_item3_%[1]s_latest."ACK_ID"=f_sch_c_part1_item3_codes_%[1]s_latest."ACK_ID"
                                              AND f_sch_c_part1_item3_%[1]s_latest."ROW_ORDER"=f_sch_c_part1_item3_codes_%[1]s_latest."ROW_ORDER"
)
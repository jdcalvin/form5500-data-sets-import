DROP MATERIALIZED VIEW IF EXISTS form5500_search_view;
CREATE MATERIALIZED VIEW form5500_search_view AS
  (SELECT *, to_tsvector(sponsor_name) || to_tsvector(sponsor_ein)
    as sponsor_search
  FROM (SELECT DISTINCT ON (plan_num, sponsor_ein) * FROM form_5500_search ORDER BY plan_num, sponsor_ein, date_received DESC) as foo;)

CREATE INDEX idx_fts_sponsor ON form5500_search_view USING gin(sponsor_search);
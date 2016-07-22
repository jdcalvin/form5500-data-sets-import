DROP TABLE IF EXISTS zip_codes,
CREATE TABLE zip_codes (
  zip integer,
  city varchar(256),
  state varchar(2),
  latitude double precision,
  longitude double precision,
  timezone integer,
  dst integer
)
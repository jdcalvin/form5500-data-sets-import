CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;
DROP FUNCTION IF EXISTS udf_distance_in_miles_from_zip(_miles integer, _zip integer);
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
LANGUAGE 'plpgsql' STABLE;
# form5500
CLI tool to download and import the DOL Form5500 data sets into a postgres database.
https://www.dol.gov/agencies/ebsa/about-ebsa/our-activities/public-disclosure/foia/form-5500-datasets

### Install CLI
`go get github.com/jdcalvin/form5500`

`go install github.com/jdcalvin/form5500`

### Usage
##### Download csvs and store in specified database (host, password, port are optional)
  
  `form5500 -import -db $DB_NAME -host $HOST -user $USER -password $DB_PASSWORD  -port 5432 -section latest -years 2013,2014,2015`
  
##### Parse imported data from specified years and aggregate into form5500_search table 
  `form5500 -build -section latest -years 2013,2014,2015`
  
##### Options
```
Usage of form5500:
  -build
    	Builds form5500_search table from all long form and short form form5500 tables
  -db string
    	specify existing database to store data
  -host string
    	connection host (default "localhost")
  -import
    	Download csvs into database
  -password string
    	password (optional)
  -port string
    	connection port (default "5432")
  -section string
    	Specify form5500 section ('all' or 'latest')
  -ssl
    	connection sslmode enabled
  -user string
    	username (optional)
  -years string
    	Specify form5500 year
  ```

### Extensions
#### Search by zip code
Quickly retrieve all form5500 records by distance from sponsor or admin zip code

`form5500 -extension zip_codes -db dbname`

example query will return all form5500 records where the sponsor is 10 miles within 97202
```
SELECT * FROM form5500_search_view WHERE substring(sponsor_zip,1,5)::int IN (select zip from udf_distance_in_miles_from_zip(10,97202))
```

### Debugging
build `go build ./form5500`

test `go test ./form5500/...`

### Rebuilding sql files
Install go-bindata https://github.com/jteeuwen/go-bindata

`go get -u github.com/jteeuwen/go-bindata/...`

Run command to create bindata.go file in internal/utils:

`go-bindata -o internal/utils/bindata.go  -pkg utils assets/...`

### Testing Data Update

Run these queries to see what sort of data is actually available to connecting apps:

```
root@dc15e85f96ce:/opt/client_data_import/scripts/staging# PGPASSWORD=$FORM5500_RDS_PASSWORD psql -h $FORM5500_RDS_ENDPOINT -U $FORM5500_RDS_USER -d $FORM5500_RDS_NAME -a -w -c "select MIN(ack_id), MAX(ack_id) from form_5500_search;"
select MIN(ack_id), MAX(ack_id) from form_5500_search;
              min               |              max               
--------------------------------+--------------------------------
 20140102080855P030152717379001 | 20180625230109P040010367981001
(1 row)


root@dc15e85f96ce:/opt/client_data_import/scripts/staging# PGPASSWORD=$FORM5500_RDS_PASSWORD psql -h $FORM5500_RDS_ENDPOINT -U $FORM5500_RDS_USER -d $FORM5500_RDS_NAME -a -w -c "select count(*) from form_5500_search;"
select count(*) from form_5500_search;
  count  
---------
 3467376
(1 row)


root@dc15e85f96ce:/opt/client_data_import/scripts/staging# PGPASSWORD=$FORM5500_RDS_PASSWORD psql -h $FORM5500_RDS_ENDPOINT -U $FORM5500_RDS_USER -d $FORM5500_RDS_NAME -a -w -c "select count(*) from form5500_search_view;"
select count(*) from form5500_search_view;
  count  
---------
 1003885
(1 row)
```

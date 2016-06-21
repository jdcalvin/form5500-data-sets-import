# form5500-data-sets-import
CLI tool to download and import the DOL Form5500 data sets into a postgres database.
http://www.dol.gov/ebsa/foia/foia-5500.html

### Install CLI
`go get github.com/jdcalvin/form5500-data-sets-import/form5500`

`go install github.com/jdcalvin/form5500-data-sets-import/form5500`

### Usage
##### Download csvs and store in specified database (host, password, port are optional)
  
  `form5500 -import -db $DB_NAME -host $HOST -user $USER -password $DB_PASSWORD  -port 5432 -section latest -years 2013,2014,2015`
  
##### Parse imported data from specified years and aggregate into form5500_search table 
  `form5500 -build`
  
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

#### Reference
zip_code.csv provided by http://federalgovernmentzipcodes.us/
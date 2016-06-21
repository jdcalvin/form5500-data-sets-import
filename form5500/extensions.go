package main

import (
	"fmt"
	"log"
  "database/sql"
  "os"
  "path/filepath"
)

func callExtension(connection string, extension string) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	if (extension == "zip_codes") {
		fmt.Println("Adding zip codes extension")
		zipCodeSearchable(db, connection)
	} else {
		log.Fatal("Invalid extension")
	}
	
	defer db.Close()
}

//private

func zipCodeSearchable(db *sql.DB, connection string) error {

	fmt.Println("  - Creating zip codes table")
	for _, statement := range createZipCodesTable() {
		_, err := db.Exec(statement)
		if err != nil {
      fmt.Println(statement)
			log.Fatal(err)
		}
	}

	importZip := importZipCodes()
	fmt.Println(fmt.Sprintf("  - %s",importZip.description))
	_, err := db.Exec(importZip.sql)
	if err != nil {
		fmt.Println("test")
		fmt.Println(importZip)
		log.Fatal(err)
	}

	fmt.Println("  - Adding udf_distance_in_miles_from_zip(_miles integer, _zip integer)")
	for _, statement := range createZipCodeSearchFunction() {
		_, err := db.Exec(statement)
		if err != nil {
      fmt.Println(statement)
			log.Fatal(err)
		}
	}
	return nil
}

func importZipCodes() Statement {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		fmt.Println(dir)
	}
	return Statement{
		sql: fmt.Sprintf("COPY zip_codes FROM '%s/zipcode.csv' DELIMITER ',' CSV HEADER", dir),
		description: "Copying zip codes into table zip_codes",
	}
	// s := fmt.Sprintf(`\copy %s FROM '%s' DELIMITER ',' CSV HEADER`, tableName, "zipcode.csv")
	// fmt.Println("psql \"" + connection + "\" -c \"" + s + "\"")
	// cmd := exec.Command("psql", connection, "-c", s)
	// output, err := cmd.Output()
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(output))
	// return nil
}



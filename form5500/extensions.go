package form5500

import (
	"fmt"
	"log"
  "database/sql"
  "os"
  "net/http"
  "strings"
  "path/filepath"
  "io"
	"os/exec"
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

	importZipCodes(connection)

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

func importZipCodes(connection string) error{	
	csvFileName := downloadZipCodeCsv()
	s := fmt.Sprintf(`\copy zip_codes FROM '%s' DELIMITER ',' CSV HEADER`, csvFileName)
	fmt.Println("-  Executing: psql \"" + connection + "\" -c \"" + s + "\"")
	cmd := exec.Command("psql", connection, "-c", s)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func downloadZipCodeCsv() string {
	url := "https://raw.githubusercontent.com/jdcalvin/form5500-data-sets-import/master/form5500/zipcode.csv"
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("  - Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		fmt.Println(n)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
    fmt.Println(dir)
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/%s", dir, fileName)
}



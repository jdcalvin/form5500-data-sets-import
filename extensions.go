package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	utils "github.com/fiduciary-benchmarks/form5500/internal/utils"
)

func callExtension(extension string) {
	if extension == "zip_codes" {
		fmt.Println("Adding zip codes extension")
		createZipCodeTable := utils.SQLRunner{
			Statement:   utils.ReadAsset("sql/zip_codes/create_table.sql"),
			Description: "Create zip_codes table",
		}

		importZipCode := utils.SQLRunner{
			Statement:   fmt.Sprintf(`\copy zip_codes FROM '%s' DELIMITER ',' CSV HEADER`, downloadZipCodeCsv()),
			Description: "Importing zip codes into zip_codes table",
		}

		createZipCodeFunction := utils.SQLRunner{
			Statement:   utils.ReadAsset("sql/zip_codes/create_search_function.sql"),
			Description: "Create zip code search function",
		}

		createZipCodeTable.Exec()

		importZipCode.ExecCLI()
		createZipCodeFunction.Exec()

	} else if extension == "rk_mapping" {
		fmt.Println("Adding rk mapping extension")
		sqlErr := readExecSqlFile("assets/sql/rk_mapping/create_table.sql")
		if sqlErr != nil {
			fmt.Println(sqlErr)
		}
		importErr := importRkMappings("assets/sql/rk_mapping/rk_mapping.csv")
		if importErr != nil {
			fmt.Println(importErr)
		}
	} else {
		log.Fatal("Invalid extension")
	}
}

//private

func downloadZipCodeCsv() string {
	url := "https://raw.githubusercontent.com/jdcalvin/form5500/master/zipcode.csv"
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

func readExecSqlFile(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	sql := ""
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		sql += line
	}
	importSql := utils.SQLRunner{Statement: sql, Description: "Creating rk company id mapping table"}
	sqlErr := importSql.Exec()
	if sqlErr != nil {
		return sqlErr
	}
	return nil
}

func importRkMappings(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()
	bufReader := bufio.NewReader(file)
	csvReader := csv.NewReader(bufReader)
	var sql, errMsg string
	lines := 0
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading file: %v", err)
		}
		if len(line) != 2 {
			return fmt.Errorf("Unexpected line: %v", line)
		}
		if line[0] == "" && line[1] == "" { //quit when we find an empty line
			errMsg = "Found a line with no data, stopping reading now, repair your file if data was truncated."
			break
		}
		name := strings.ReplaceAll(line[0], "'", "''")
		id, _ := strconv.Atoi(line[1])
		if id != 0 { //conversion returns 0 if it's not a number, probably the header line, and in any case not valid
			sql += fmt.Sprintf("INSERT INTO sched_c_provider_to_fbi_rk_company_id_mappings (sched_c_provider_name, fbi_company_id) VALUES ('%v',%d);", name, id)
			lines++
		}
	}
	importSql := utils.SQLRunner{Statement: sql, Description: "Importing rk company id mapping"}
	importSql.Exec()
	if errMsg != "" {
		return fmt.Errorf("%v  %d lines of data were processed.", errMsg, lines)
	}
	return nil
}

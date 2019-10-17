package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	utils "github.com/fiduciary-benchmarks/form5500/internal/utils"
	_ "github.com/lib/pq"
)

const baseURL string = "http://askebsa.dol.gov/FOIA%20Files/"

// ImportResult represents information about an individual Form5500 import
type ImportResult struct {
	Status  string `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Source  string `json:"source,omitempty"`
	Section string `json:"section,omitempty"`
	Year    string `json:"year,omitempty"`
}

func runImport(section string, years []string) ([]ImportResult, error) {
	results := []ImportResult{}

	for _, year := range years {
		result := createAndPopulateTables(year, section)
		results = append(results, result...)
	}

	return results, nil
}

func createAndPopulateTables(year string, section string) []ImportResult {
	results := []ImportResult{}

	for _, name := range tableNames() {
		result := ImportResult{}
		tableName := fmt.Sprintf(name, year, section)
		runner, err := createTable(tableName, year, section)

		if err != nil {
			result = buildErrorResult(year, section, getURL(year, section, name), err)
			results = append(results, result)
			continue
		}

		err = runner.Exec()
		if err != nil {
			result = buildErrorResult(year, section, getURL(year, section, name), err)
			results = append(results, result)
			continue
		}

		csvFilename, sourceURL, err := downloadCSV(name, year, section)

		if err != nil {
			result = buildErrorResult(year, section, sourceURL, err)
			results = append(results, result)
			continue
		}

		defer os.Remove(csvFilename)
		fmt.Println("Created CSV file: " + csvFilename)

		err = importCSV(tableName, csvFilename)

		if err != nil {
			result = buildErrorResult(year, section, sourceURL, err)
			results = append(results, result)
			continue
		}

		result.Message = "Import successful"
		result.Status = "succeeded"
		result.Success = true
		result.Source = sourceURL
		result.Section = section
		result.Year = year
		results = append(results, result)
	}
	return results
}

// private

func tableNames() []string {
	tables := []string{
		"f_5500_%s_%s",
		"f_5500_sf_%s_%s",
		"f_sch_a_%s_%s",
		"f_sch_a_part1_%s_%s",
		"f_sch_c_%s_%s",
		"f_sch_c_part1_item1_%s_%s",
		"f_sch_c_part1_item2_%s_%s",
		"f_sch_c_part1_item2_codes_%s_%s",
		"f_sch_c_part1_item3_%s_%s",
		"f_sch_c_part1_item3_codes_%s_%s",
		"f_sch_c_part2_%s_%s",
		"f_sch_c_part2_codes_%s_%s",
		"f_sch_c_part3_%s_%s",
		// "f_sch_d_%s_%s",
		// "f_sch_d_part1_%s_%s",
		// "f_sch_d_part2_%s_%s",
		// "f_sch_g_%s_%s",
		// "f_sch_g_part1_%s_%s",
		// "f_sch_g_part2_%s_%s",
		// "f_sch_g_part3_%s_%s",
		// "f_sch_g_%s_%s",
		// "f_sch_g_%s_%s",
		// "f_sch_g_%s_%s",
		"f_sch_h_%s_%s",
		"f_sch_h_part1_%s_%s",
		"f_sch_i_%s_%s",
		// "f_sch_r_%s_%s",
		// "f_sch_r_part1_%s_%s",
		// "f_sch_mb_%s_%s",
		// "f_sch_mb_part1_%s_%s",
		// "f_sch_mb_part2_%s_%s",
		// "f_sch_sb_%s_%s",
		// "f_sch_sb_part1_%s_%s",
	}
	return tables
}

func importCSV(tableName string, csvFilename string) error {
	truncateTable := utils.SQLRunner{
		Statement:   fmt.Sprintf("TRUNCATE %s", tableName),
		Description: fmt.Sprintf("Truncating %s", tableName),
	}

	err := truncateTable.Exec()

	if err != nil {
		return err
	}

	copyCsv := utils.SQLRunner{
		Statement:   fmt.Sprintf(`\copy %s FROM '%s' DELIMITER ',' CSV HEADER`, tableName, csvFilename),
		Description: fmt.Sprintf("Copying %s into %s", csvFilename, tableName),
	}

	err = copyCsv.ExecCLI()

	if err != nil {
		return err
	}

	return nil
}

func downloadCSV(name string, year string, section string) (string, string, error) {
	name = fmt.Sprintf(name, year, section)
	url := getURL(year, section, name)

	fmt.Println("Dowloading ", url)

	zipFilename, err := downloadFile(name, url)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(zipFilename)

	reader, err := zip.OpenReader(zipFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	csvFilename := strings.ToLower(name) + ".csv"

	for _, f := range reader.File {
		if strings.ToLower(f.Name) == csvFilename {
			csvFile, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer csvFile.Close()

			tempFile, tempFilename, err := createTempFile(csvFilename)
			if err != nil {
				return "", "", err
			}
			defer tempFile.Close()

			_, err = io.Copy(tempFile, csvFile)
			if err != nil {
				log.Fatal(err)
			}

			return tempFilename, url, nil
		}
	}

	return "", "", errors.New("CSV not found in ZIP file at " + url)
}

func createTable(tableName string, year string, section string) (utils.SQLRunner, error) {
	url := baseURL + fmt.Sprintf("%s/%s/%s_layout.txt", year, section, tableName)
	fmt.Println("Downloading ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not resolve url: ", url)
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Unable to retrieve Form5500 Data File. Status: %s, Status Code: %s", resp.Status, strconv.Itoa(resp.StatusCode))
		err := fmt.Errorf("Unable to retrieve Form5500 Data File. Status: %s, Status Code: %s", resp.Status, strconv.Itoa(resp.StatusCode))
		return utils.SQLRunner{}, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	// eat first two header lines
	scanner.Scan()
	scanner.Scan()

	var sqlLines []string
	sqlLines = append(sqlLines, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", tableName))
	sqlLines = append(sqlLines, fmt.Sprintf("CREATE TABLE %s (", tableName))

	// first column is always ACK_ID
	scanner.Scan()
	sqlLines = append(sqlLines, `    "ACK_ID" varchar(30)`)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		sqlLines[len(sqlLines)-1] = sqlLines[len(sqlLines)-1] + ","

		switch {
		case strings.HasSuffix(parts[1], "_DATE") || strings.HasSuffix(parts[1], "_PRD"):
			sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" timestamp`, parts[1]))

		case len(parts) == 4:
			if parts[3] == "1" && strings.HasSuffix(parts[1], "_IND") {
				sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" int`, parts[1]))
			} else {
				if parts[3] == "0" {
					sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" char`, parts[1]))
				} else {
					sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" varchar(%s)`, parts[1], parts[3]))
				}
			}

		case parts[2] == "NUMERIC":
			if strings.HasSuffix(parts[1], "_ORDER") {
				sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" int`, parts[1]))
			} else if strings.HasSuffix(parts[1], "_AMT") {
				sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" numeric(19,2)`, parts[1]))
			} else if strings.HasSuffix(parts[1], "_CNT") {
				sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" int`, parts[1]))
			} else {
				sqlLines = append(sqlLines, fmt.Sprintf(`    "%s" numeric(19,6)`, parts[1]))
			}
		}
	}
	sqlLines = append(sqlLines, ");")

	if err := scanner.Err(); err != nil {
		return utils.SQLRunner{}, err
	}

	sql := ""
	for _, line := range sqlLines {
		sql = sql + line + "\n"
	}

	return utils.SQLRunner{
		Statement:   sql,
		Description: fmt.Sprintf("Creating table: %s", tableName),
	}, nil
}

func downloadFile(prefix string, url string) (string, error) {
	tempFile, tempFilename, err := createTempFile(prefix)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tempFilename, nil
}

func createTempFile(prefix string) (*os.File, string, error) {
	tempFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, "", err
	}
	return tempFile, tempFile.Name(), nil
}

func getURL(year string, section string, name string) string {
	return baseURL + fmt.Sprintf("%s/%s/%s.zip", year, section, name)
}

func buildErrorResult(section string, year string, source string, err error) ImportResult {
	result := ImportResult{}
	result.Message = err.Error()
	result.Success = false
	result.Status = "failed"
	result.Source = source
	result.Section = section
	result.Year = year
	return result
}

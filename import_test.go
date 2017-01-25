package main

// import (
// 	"fmt"
// 	"testing"
// )

// func TestCreateTable(t *testing.T) {
// 	// Testing:
// 	// - downloading a file
// 	// - reading the file contents to determine columns to build a table
// 	// - build sql statements to drop if exists, and create table
// 	// - print description to console

// 	var expectation string
// 	name := "f_sch_c_part1_item1_%s_%s"
// 	section := "latest"
// 	year := "2013"
// 	tableName := fmt.Sprintf(name, year, section)
// 	sqlRunner := createTable(tableName, year, section)
// 	// expectation = utils.ReadAsset("test_examples/test_create_table.txt")

// 	if sqlRunner.Statement != expectation {
// 		expectation = fmt.Sprintf("Expected: %s, Got:", expectation)
// 		t.Error(expectation, sqlRunner.Statement)
// 	}

// 	expectation = "Creating table: f_sch_c_part1_item1_2013_latest"
// 	if sqlRunner.Description != expectation {
// 		expectation = fmt.Sprintf("Expected %s, Got:", expectation)
// 		t.Error(expectation, sqlRunner.Description)
// 	}

// }

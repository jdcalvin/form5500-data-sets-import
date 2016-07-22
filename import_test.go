package main

import (
  "testing" 
  "fmt"
  utils "github.com/jdcalvin/form5500-data-sets-import/internal/utils"
)

func TestCreateTable(t *testing.T) {
  // Testing:
  // - downloading a file
  // - reading the file contents to determine columns to build a table
  // - build sql statements to drop if exists, and create table
  // - print description to console

  var expectation string
  name := "f_sch_c_part1_item1_%s_%s"
  section := "latest"
  year := "2013"
  tableName := fmt.Sprintf(name, year, section)
  statement := createTable(tableName, year, section)
  expectation = utils.ReadFile("test_examples/test_create_table.txt")
    
  if statement.Sql != expectation {
    expectation = fmt.Sprintf("Expected: %s, Got:", expectation)
    t.Error(expectation, statement.Sql)
  }

  expectation = "Creating table: f_sch_c_part1_item1_2013_latest"
  if statement.Description != expectation {
    expectation = fmt.Sprintf("Expected %s, Got:", expectation)
    t.Error(expectation, statement.Description)
  }


}
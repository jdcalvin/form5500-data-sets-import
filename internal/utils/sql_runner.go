package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
)

var db *sql.DB
var dbConnection string
var err error

// SetDBConnection set connection string to be used by sql.DB Open()
func SetDBConnection(connection string) {
	dbConnection = connection
}

// OpenDBConnection opens db connection - call at top of function
func OpenDBConnection() error {
	if dbConnection == "" {
		log.Fatal("SQLRunner.Open() Must set connection with SetConnection(connection string)")
	}
	db, err = sql.Open("postgres", dbConnection)

	if err != nil {
		return err
	}
	return nil
}

// CloseDBConnection closes opened db connection - defer at top of function
func CloseDBConnection() {
	db.Close()
}

// SQLRunner struct to assign and call sql statements throughout the form5500 package
type SQLRunner struct {
	Statement   string
	Description string
}

// Exec runs #sql statement and prints description to command line
func (s SQLRunner) Exec() error {
	s.Print()
	_, err := db.Exec(s.Statement)
	if err != nil {
		fmt.Println(s)
		return err
	}
	return nil
}

func (s SQLRunner) Query() (*sql.Rows, error) {
	s.Print()
	rows, err := db.Query(s.Statement)
	if err != nil {
		fmt.Println(s)
		return rows, err
	}
	return rows, nil
}

// ExecCLI uses psql command line tool to copy data from a csv file
// Cannot use Exec due to permissions error on aws box
func (s SQLRunner) ExecCLI() error {
	s.Print()
	cmd := exec.Command("psql", dbConnection, "-c", s.Statement)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("psql \"" + dbConnection + "\" -c \"" + s.Statement + "\"")
		return err
	}
	return nil
}

// Print print formatted message to console
func (s SQLRunner) Print() {
	fmt.Println(fmt.Sprintf(" - %s", s.Description))
}

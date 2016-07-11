package form5500

import (
	"flag"
	"fmt"
	"strings"
)

const baseUrl string = "http://askebsa.dol.gov/FOIA%20Files/"

var hostFlag = flag.String("host", "localhost", "connection host")
var portFlag = flag.String("port", "5432", "connection port")
var userFlag = flag.String("user", "", "username (optional)")
var passwordFlag = flag.String("password", "", "password (optional)")
var dbFlag = flag.String("db", "", "specify existing database to store data")
var sslFlag = flag.Bool("ssl", false, "connection sslmode enabled")

var yearsFlag = flag.String("years", "", "Specify form5500 year")
var sectionFlag = flag.String("section", "", "Specify form5500 section ('all' or 'latest')")

var isImportFlag = flag.Bool("import", false, "Download csvs into database")
var isBuildFlag = flag.Bool("build", false, "Builds form5500_search table from all long form and short form form5500 tables")
var isExtensionFlag = flag.String("extension", "", "Add extensions to ")

func main() {
	flag.Parse()

	if (!*isImportFlag && !*isBuildFlag && (*isExtensionFlag == "") || (*isImportFlag && *isBuildFlag && (*isExtensionFlag != ""))) {
		fmt.Println("Must specify import, build, or an extension")
	}

	var connectionPartial string
	fmt.Println(*userFlag)
	if *userFlag != "" && *passwordFlag != "" {
		connectionPartial = fmt.Sprintf("user=%s password=%s ", *userFlag, *passwordFlag)
	} else {
		connectionPartial = ""
	}

	var sslMode string

	if *sslFlag {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}

	years := strings.Split(*yearsFlag, ",")

	section := *sectionFlag

	connection := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s %s", *hostFlag, *portFlag, *dbFlag, sslMode, connectionPartial)

	if *isImportFlag {
		for _, year := range years {
			runImport(connection, year, section)
		}
	}

	if *isBuildFlag {
		buildTable(connection, section, years)
	}

	if *isExtensionFlag != "" {
		callExtension(connection, *isExtensionFlag)
	}

}

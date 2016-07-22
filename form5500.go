package main

import (
	"flag"
	"fmt"
	"strings"
	utils "github.com/jdcalvin/form5500/internal/utils"
)

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

	form5500Flag := new(utils.Form5500Flags)

	form5500Flag.Section = *sectionFlag
	form5500Flag.Years = strings.Split(*yearsFlag, ",")
	form5500Flag.SetConnection(*hostFlag, *portFlag, *dbFlag, *sslFlag, *userFlag, *passwordFlag)

	utils.SetDBConnection(form5500Flag.Connection)
	utils.OpenDBConnection()
	defer utils.CloseDBConnection()

	if *isImportFlag {
		runImport(form5500Flag.Section, form5500Flag.Years)
	}

	if *isBuildFlag {
		buildTable(form5500Flag.Section, form5500Flag.Years)
	}

	if *isExtensionFlag != "" {
		callExtension(*isExtensionFlag)
	}

}

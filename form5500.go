package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/fiduciary-benchmarks/form5500/internal/utils"
)

var hostFlag = flag.String("host", "localhost", "connection host")
var portFlag = flag.String("port", "5432", "connection port")
var userFlag = flag.String("user", "", "username (optional)")
var passwordFlag = flag.String("password", "", "password (optional)")
var dbFlag = flag.String("db", "", "specify existing database to store data")
var sslFlag = flag.String("ssl", "disable", "connection sslmode: require, verify-full, verify-ca, disable")

var yearsFlag = flag.String("years", "", "Specify form5500 year")
var sectionFlag = flag.String("section", "", "Specify form5500 section ('all' or 'latest')")

var jiraCreator = flag.String("jira-user", "", "jira user to auto-create issues")
var jiraToken = flag.String("jira-token", "", "api token for specified jira issue creation user")
var jiraAssignee = flag.String("jira-assignee", "", "user to assign auto-created issues to")

var isImportFlag = flag.Bool("import", false, "Download csvs into database")
var isBuildFlag = flag.Bool("build", false, "Builds form5500_search table from all long form and short form form5500 tables")
var isExtensionFlag = flag.String("extension", "", "Add extensions to ")
var outputJSONFlag = flag.Bool("output-json", false, "Return result of import process as JSON. Defaults to false.")
var debugMode = flag.Bool("debug", false, "Output verbose log of import. Defaults to false.")

func main() {
	flag.Parse()

	form5500Flag := new(utils.Form5500Flags)

	if !*isImportFlag && !*isBuildFlag && (*isExtensionFlag == "") {
		fmt.Println("Must specify import, build, or an extension")
	}

	form5500Flag.Section = *sectionFlag
	form5500Flag.Years = strings.Split(*yearsFlag, ",")
	form5500Flag.SetConnection(*hostFlag, *portFlag, *dbFlag, *sslFlag, *userFlag, *passwordFlag)

	utils.SetDBConnection(form5500Flag.Connection)
	utils.OpenDBConnection()
	defer utils.CloseDBConnection()

	var results []ImportResult
	var err error

	if *isImportFlag {
		results, err = runImport(form5500Flag.Section, form5500Flag.Years)

		if err != nil {
			log.Println(err.Error())
			utils.CloseDBConnection()
			results = append(results, buildErrorResult("", "", "", err))
		}
	}

	if *isBuildFlag {
		rebuildSearchTable(form5500Flag.Section, form5500Flag.Years)
		if *jiraCreator != "" && *jiraToken != "" && *jiraAssignee != "" {
			findUnmatchedRks(*jiraCreator, *jiraToken, *jiraAssignee)
		}
	}

	if *isExtensionFlag != "" {
		callExtension(*isExtensionFlag)
	}

	// Return a JSON string so that it can be parsed by a receiving application.
	// Idea was borrowed from AWS CLI commands that can generate a JSON skeleton and
	// other output to be piped to other applications.
	if results != nil && *outputJSONFlag {
		jByte, _ := json.Marshal(results)
		fmt.Println(string(jByte))
	}
}

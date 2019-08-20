package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/andygrunwald/go-jira"
)

//jiraCreator: lgabriel@fiduciarybenchmarks.com
//jiraAssignee: lgabriel
func CreateJiraIssue(jiraCreator string, jiraToken string, jiraAssignee string) {

	tp := jira.BasicAuthTransport{
		Username: jiraCreator,
		Password: jiraToken,
	}

	client, err := jira.NewClient(tp.Client(), "https://fbidev.atlassian.net")
	if err != nil {
		fmt.Println("error getting client: ", err)
	}

	t := time.Now()
	year, month, day := t.Date()

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: jiraAssignee,
			},
			Description: "Please provide an updated mapping file that includes mappings for the missing rks.",
			Type: jira.IssueType{
				Name: "Task",
			},
			Project: jira.Project{
				Key: "DL",
			},
			Summary: fmt.Sprintf("Unmatched RK records from 5500 import from %d-%d-%d", month, day, year),
		},
	}
	issue, _, err := client.Issue.Create(&i)
	if err != nil {
		fmt.Println(err)
	}
	if issue == nil || issue.Key == "" {
		fmt.Println("Unable to create the jira issue.")
		return
	}
	fmt.Printf("key: %v", issue.Key)

	r, err := os.Open("unmatched_rks.txt")
	if err != nil {
		fmt.Println("Failed to open unmatched_rks.txt.  The issue was created but the file was not attached.: ", err)
		return
	}
	_, _, err = client.Issue.PostAttachment(issue.Key, r, "unmatched_rks.txt")
	if err != nil {
		fmt.Println("Failed to post attachment: ", err)
	}
}

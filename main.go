package main

import (
	"context"
	"os"
	"regexp"

	"github.com/google/go-github/v40/github"
	"github.com/posener/goaction"
	log "github.com/sirupsen/logrus"
)

var token = os.Getenv("GITHUB_TOKEN")

func main() {
	titleRegex := os.Getenv("INPUT_TITLE-REGEX")
	errorMessage := os.Getenv("INPUT_ERROR-MESSAGE")
	excludeFilesRegex := os.Getenv("INPUT_EXCLUDE-FILES-REGEX")
	//titleRegex := "^(\\[[a-z ]*\\]|Bump) "
	//errorMessage := ""
	//excludeFilesRegex := ".*(\\.gradle|\\.java|\\.xml).*"

	ctx := context.Background()
	client := github.NewClient(nil)

	pr, err := goaction.GetPullRequest()
	if err != nil {
		log.Fatal(err)
	}

	doTitleCheck := true
	if excludeFilesRegex != "" {
		doTitleCheck = false
		files, resp, err := client.PullRequests.ListFiles(ctx, pr.GetOrganization().GetName(), pr.GetRepo().GetName(), pr.GetNumber(), nil)
		if err != nil || resp.StatusCode != 200 {
			log.Fatalf("Cannot get files for PR. %v", err)
		}
		for _, file := range files {
			match, err := regexp.MatchString(excludeFilesRegex, file.GetFilename())
			if err != nil {
				log.Fatal(err)
			}
			if !match {
				doTitleCheck = true
				break
			}
		}
	}
	if doTitleCheck {
		match, err := regexp.MatchString(titleRegex, pr.GetPullRequest().GetTitle())
		if err != nil {
			log.Fatal(err)
		}
		if !match {
			if errorMessage == "" {
				errorMessage = "Please include JIRA issue ID at the beginning of PR title"
			}
			log.Fatalf(errorMessage)
		}
	}
}

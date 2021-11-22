package main

import (
	"context"
	"golang.org/x/oauth2"
	"os"
	"regexp"

	"github.com/google/go-github/v40/github"
	"github.com/posener/goaction"
	log "github.com/sirupsen/logrus"
)

func main() {
	token := os.Getenv("secrets.GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token = os.Getenv("INPUT_GITHUB-TOKEN")
	}
	titleRegex := os.Getenv("INPUT_TITLE-REGEX")
	errorMessage := os.Getenv("INPUT_ERROR-MESSAGE")
	excludeFilesRegex := os.Getenv("INPUT_EXCLUDE-FILES-REGEX")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	if token == "" {
		log.Infof("Empty GitHub token")
	}

	pr, err := goaction.GetPullRequest()
	if err != nil {
		log.Fatal(err)
	}
	doTitleCheck := true
	if excludeFilesRegex != "" {
		doTitleCheck = false
		files, resp, err := client.PullRequests.ListFiles(ctx, pr.GetRepo().GetOwner().GetLogin(), pr.GetRepo().GetName(), pr.GetNumber(), nil)
		if err != nil || resp.StatusCode != 200 {
			log.Fatalf("Cannot get files for PR. %v", err)
		}
		for _, file := range files {
			match, err := regexp.MatchString(excludeFilesRegex, file.GetFilename())
			if err != nil {
				log.Fatal(err)
			}
			log.Infof("Filename: %v, matches: %v", file.GetFilename(), match)
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

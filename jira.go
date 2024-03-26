package main

import (
	"fmt"
	"github.com/joakimen/goji/pkg/jira"
	"os"
)

func GetJiraIssues() (jira.Issues, error) {
	user, err := RequireEnv("JIRA_API_USER")
	if err != nil {
		return nil, fmt.Errorf("error fetching jira issues: %v", err)
	}

	token, err := RequireEnv("JIRA_API_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("error fetching jira issues: %v", err)
	}

	host, err := RequireEnv("JIRA_HOST")
	if err != nil {
		return nil, fmt.Errorf("error fetching jira issues: %v", err)
	}

	apiCredentials := jira.APICredentials{
		User:  user,
		Token: token,
		Host:  host,
	}
	issues, err := jira.Search(apiCredentials)
	if err != nil {
		return nil, fmt.Errorf("error fetching jira issues: %v", err)
	}
	return issues, nil
}

func RequireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return value, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/joakimen/goji/pkg/jira"
	fz "github.com/ktr0731/go-fuzzyfinder"
)

type NoSelection struct{}

func (m *NoSelection) Error() string {
	return "no selection made"
}

func SelectIssue(issues jira.Issues) (jira.Issue, error) {
	renderFunc := func(selectedIndex int) string {
		return fmt.Sprintf("%s: %s", issues[selectedIndex].Key, issues[selectedIndex].Fields.Summary)
	}

	previewFunc := func(selectedIndex, width, height int) string {
		if selectedIndex == -1 {
			return ""
		}
		issue := issues[selectedIndex]
		issueJson, err := json.MarshalIndent(issue, "", "  ")
		if err != nil {
			return fmt.Sprintf("error marshalling issue to json: %v", err)
		}
		return string(issueJson)
	}

	idx, err := fz.Find(issues, renderFunc, fz.WithPreviewWindow(previewFunc))
	if err != nil {
		return jira.Issue{}, &NoSelection{}
	}

	return issues[idx], nil
}

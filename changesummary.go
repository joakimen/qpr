package main

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"
)

type ChangeSummary string

func (c ChangeSummary) String() string {
	return string(c) // Convert ChangeSummary to string and return.
}

func (c ChangeSummary) CommitSummary() string {
	if len(c) == 0 {
		return errors.New("commit summary is empty").Error()
	}
	if len(c) > 50 {
		return errors.New("commit summary too long (max 50 characters)").Error()
	}
	return string(c)
}

func (c ChangeSummary) BranchName(prefixes []string) string {
	branchName := string(c)
	branchName = removeConventionalCommitPattern(branchName)
	branchName = strings.ReplaceAll(branchName, " ", "-")

	slog.Debug("branch name after sanitizing: " + branchName)

	// exit if changeSummary cont
	pattern := `^[a-zA-Z0-9._/-]+$`
	matched, err := regexp.MatchString(pattern, branchName)
	if err != nil {
		panic(err)
	}

	if !matched {
		panic("branch name contains invalid characters")
	}

	return strings.Join(append(prefixes, branchName), "/")
}

func (c ChangeSummary) PullRequestTitle() string {

	prTitle := string(c)

	// Remove any conventional commit prefixes
	prTitle = removeConventionalCommitPattern(prTitle)

	// Capitalize the first letter
	r := rune(prTitle[0])
	if unicode.IsLetter(r) {
		prTitle = string(unicode.ToUpper(r)) + prTitle[1:]
	}

	return prTitle
}

func removeConventionalCommitPattern(s string) string {
	// Remove any conventional commit prefixes
	conventionalCommitPrefixes := []string{"feat", "fix", "docs", "build", "style", "refactor", "perf", "test", "ci", "chore"}
	prefixPattern := strings.Join(conventionalCommitPrefixes, "|")
	regexPattern := fmt.Sprintf(`(%s)(\([a-zA-Z]*\))?:\s*`, prefixPattern)
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(re.ReplaceAllString(s, ""))
}

type ChangeSummaryConfig struct {
	CommitSummary    string
	BranchName       string
	PullRequestTitle string
}

func NewChangeSummaryConfig(changeSummary ChangeSummary, branchPrefixes []string) ChangeSummaryConfig {
	return ChangeSummaryConfig{
		CommitSummary:    changeSummary.CommitSummary(),
		BranchName:       changeSummary.BranchName(branchPrefixes),
		PullRequestTitle: changeSummary.PullRequestTitle(),
	}
}

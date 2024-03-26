package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Flags struct {
	Verbose  bool
	DryRun   bool
	SkipJira bool
}

func main() {

	flags := parseFlags()

	if flags.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	doPreliminaryChecks()

	// 1. read input string from user (MSG) summarizing changes (max-length: 52)
	changeSummary, err := ReadChangeSummary()
	if err != nil {
		panic(err)
	}
	slog.Debug("change summary: " + changeSummary.String())

	// 2. Create a new branch from the message, adding additional prefixes if available
	var branchPrefixes []string

	// 2.1 add git user prefix is present
	gitUserPrefix := os.Getenv("GIT_USER_PREFIX")
	if gitUserPrefix != "" {
		slog.Debug("using GIT_USER_PREFIX: " + gitUserPrefix)
		branchPrefixes = append(branchPrefixes, gitUserPrefix)
	}

	// 2.2 add Jira issue prefix if available
	if !flags.SkipJira {
		jiraIssues, err := GetJiraIssues()
		if err != nil {
			panic(err)
		}
		jiraIssue, err := SelectIssue(jiraIssues)
		switch {
		case err == nil:
			branchPrefixes = append(branchPrefixes, strings.ToLower(jiraIssue.Key))
		case errors.Is(err, &NoSelection{}):
			break
		default:
			panic(err)
		}
	}

	changeSummaryConfig := NewChangeSummaryConfig(changeSummary, branchPrefixes)

	if flags.DryRun {

		fmt.Println("Flags:")
		PrintAsJson(flags)

		fmt.Println("ChangeSummaryConfig:")
		PrintAsJson(changeSummaryConfig)
		return
	}

	slog.Debug("branch name: " + changeSummaryConfig.BranchName)
	// 3. checkout branch
	err = exec.Command("git", "checkout", "-b", changeSummaryConfig.BranchName).Run()

	// 4. stage all changes
	err = exec.Command("git", "add", "--all").Run()
	if err != nil {
		panic(err)
	}

	// 5. commit all changes with MSG as message
	slog.Debug("commitSummary: " + changeSummary.CommitSummary())
	err = exec.Command("git", "commit", "-m", changeSummaryConfig.CommitSummary).Run()
	if err != nil {
		panic(err)
	}

	// 6. push changes to remote
	err = exec.Command("git", "push", "origin", changeSummaryConfig.BranchName).Run()
	if err != nil {
		panic(err)
	}

	// 7. open a new PR with the commit message as title using gh, with the option "continue in browser" in git
	err = exec.Command("gh", "pr", "create", "--title", changeSummaryConfig.PullRequestTitle, "--web").Run()
}

func parseFlags() Flags {
	var flags Flags
	flag.BoolVar(&flags.Verbose, "v", false, "enable verbose output")
	flag.BoolVar(&flags.DryRun, "dry-run", false, "summarize changes without creating a commit, branch or PR")
	flag.BoolVar(&flags.SkipJira, "skip-jira", false, "disable jira integration")
	flag.Parse()
	return flags
}

func doPreliminaryChecks() {
	// quit immediately if not in a git directory
	err := exec.Command("git", "rev-parse").Run()
	if err != nil {
		panic("not in a git directory")
	}

	// quit if not on trunk branch
	branch, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		panic(err)
	}
	trimmedBranch := strings.TrimSpace(string(branch))
	if !slices.Contains([]string{"main", "master"}, trimmedBranch) {
		panic("not on main or master branch")
	}

	// quit if there are no changes
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		panic("error checking repo status")
	}

	if len(output) == 0 {
		panic("repo is clean, no changes to commit")
	}
}

func ReadChangeSummary() (ChangeSummary, error) {

	userPrompt := "Enter a summary of your changes:"
	p := tea.NewProgram(initialModel(userPrompt))

	bubblesModel, err := p.Run()
	if err != nil {
		return "", err
	}

	textInput := bubblesModel.(model).GetTextInput()
	return ChangeSummary(textInput), nil
}

func PrintAsJson(v interface{}) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}

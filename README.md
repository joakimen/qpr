# qpr

Quick PR

## Description

Automates the process from uncommitted changes to a new PR.

## Install

```bash
go install github.com/joakimen/qpr@latest
```

## Usage

```bash
$ qpr
```

## Features

### Preliminary checks

Exit if:
- not in a git repo
- not on master/main branch
- repo is clean

### Main steps

1. Prompt user for summary of changes (`MSG`)
2. Create branch name based on `MSG` and switch to it
3. Stage everything, commit using sanitized commit message from `MSG`
4. Push branch to origin
5. Open web browser to finalize a new PR with title based on `MSG`

## Flags

- `-v`: print some information
- `--dry-run`: dry run, print generated commit message, branch name and PR title
- `--skip-jira`: skip Jira issue API lookup and selection

## Example

With the following circumstance:
- `GIT_USER_PREFIX` env var set with value `baconator`
- Jira credentials (user, token, host) configured, and issue `ABC-123` selected after invoking `qpr`

```bash
$ qpr --dry-run
Enter a summary of your changes:

> refactor(api): reduce log volume of auth handler

(ctrl-c or esc to quit)
Flags:
{
  "Verbose": false,
  "DryRun": true,
  "SkipJira": false
}

ChangeSummaryConfig:
{
  "CommitSummary": "refactor(api): reduce log volume of auth handler",
  "BranchName": "baconator/abc-123/reduce-log-volume-of-auth-handler",
  "PullRequestTitle": "Reduce log volume of auth handler"
}
```

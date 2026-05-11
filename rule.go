package main

import (
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

type Rule struct {
	Description string
	Severity    Severity
	CheckStage func(*instructions.Stage) []Issue
	CheckRun   func(*instructions.RunCommand) []Issue
	CheckCopy  func(*instructions.CopyCommand) []Issue
	CheckUser  func(*instructions.UserCommand) []Issue
}

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Issue struct {
	Severity Severity
	Message  string
	Line     *int
	Fix      string
}

var NoLatestTag = Rule{
	Description: "Checks that the user doesn't use latest tag when using the FROM command on an image:tag",
	Severity: SeverityWarning,
	CheckStage: func(cmd *instructions.Stage) []Issue{
		if strings.Contains(cmd.BaseName, ":latest") {
			line := stageLine(cmd)
			return []Issue{
				{
					Severity: SeverityWarning,
					Message: "Using latest tag may introduce malicious dependencies if the targeted image is compromised.",
					Line: &line,
					Fix: "Use versioned image.",
				},
			}
		}
		return nil
	},
}

var NoUserDefined = Rule{
	Description: "Checks that the user doesn't use the default user (root) in the image",
	Severity: SeverityWarning,
	CheckStage: func(stage *instructions.Stage) []Issue {
		root_flag := true
		for _, cmd := range stage.Commands {
			switch c := cmd.(type) {
			case *instructions.UserCommand:
				user := strings.ToLower(c.User)
				if user == "root" {
					root_flag = true
				} else {
					root_flag = false
				}
			}
		}
		if root_flag {
			line := stageLine(stage)
			return []Issue{
				{
					Severity: SeverityWarning,
					Message: "Stage on line " + strconv.Itoa(line) + " ends with root user. Define a specific user for the stage.",
				},
			}
		}
		return nil
	},
}

var Rules = []Rule{
	NoLatestTag,
	NoUserDefined,	
} 


// Helper functions
func lineOf(cmd instructions.Command) int {
	if loc := cmd.Location(); len(loc) > 0 {
		return loc[0].Start.Line
	}
	return 0
}

func stageLine(stage *instructions.Stage) int {
	if len(stage.Location) > 0 {
		return stage.Location[0].Start.Line
	}
	return 0
}

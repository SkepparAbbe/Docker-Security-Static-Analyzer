package main

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

type Rule struct {
	Description string
	Severity    Severity
	CheckRun   func(*instructions.RunCommand) []Issue
	CheckCopy  func(*instructions.CopyCommand) []Issue
	CheckUser  func(*instructions.UserCommand) []Issue
	CheckStage func(*instructions.Stage) []Issue
}

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Issue struct {
	Severity Severity
	Message  string
	Line     int
	Fix      string
}

var CheckUser = Rule{
	Description: "Checks that a user is set for the final stage",
}

// TODO
// is found in stage

var noLatestTag = Rule{
	Description: "Checks that the user doesn't use latest tag when using the FROM command on an image:tag",
	Severity: SeverityWarning,
	CheckStage: func(cmd *instructions.Stage) []Issue{
		if strings.Contains(cmd.BaseName, ":latest") {
			return []Issue{
				{
					Severity: SeverityWarning,
					Message: "Using latest tag may introduce malicious dependencies if the targeted image is compromised.",
					Line: stageLine(cmd),
					Fix: "Use versioned image.",
				},
			}
		}
		return nil
	},
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

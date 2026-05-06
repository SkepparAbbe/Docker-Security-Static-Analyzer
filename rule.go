package main

import {
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
}

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
	Description: "Checks that a user",
}

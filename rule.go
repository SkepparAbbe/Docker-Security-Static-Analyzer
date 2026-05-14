package main

import (
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

type Rule struct {
	Description string
	Severity    Severity
	CheckStage  func(*instructions.Stage) []Issue
	CheckRun    func(*instructions.RunCommand) []Issue
	CheckCopy   func(*instructions.CopyCommand) []Issue
	CheckUser   func(*instructions.UserCommand) []Issue
}

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Issue struct {
	Severity    Severity
	Problem     string
	Description string
	Line        *int
	Fix         string
}

var NoLatestTag = Rule{
	Description: "Checks that the user doesn't use latest tag when using the FROM command on an image:tag",
	Severity:    SeverityWarning,
	CheckStage: func(cmd *instructions.Stage) []Issue {
		if strings.Contains(cmd.BaseName, ":latest") {
			line := stageLine(cmd)
			return []Issue{
				{
					Severity:    SeverityWarning,
					Problem:     "Using latest tag in FROM command.",
					Description: "Using latest tag may lead to unexpected updates and potential security vulnerabilities if the image is compromised.",
					Line:        &line,
					Fix:         "Use versioned image.",
				},
			}
		}
		return nil
	},
}

var NoUserDefined = Rule{
	Description: "Checks that the user doesn't use the default user (root) in the image",
	Severity:    SeverityWarning,
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
					Problem:  "Stage on line " + strconv.Itoa(line) + " ends with root user.",
					Description: "If this stage is the runtime of the container, the main process of the container may have unecessary privileges which violates principle of least privilege",
					Fix: "Define a specific user for the stage.",
				},
			}
		}
		return nil
	},
}

var NoHashTagImage = Rule{
	Description: "Checks that an image hash a tag of format \"sha256:(64 characters)\"",
	Severity:    SeverityInfo,
	CheckStage: func(stage *instructions.Stage) []Issue {
		from := strings.Split(stage.OrigCmd, ":")
		if (len(from) == 3 && from[1] == "sha256") {
			if len(from[2]) == 64 {
				return nil
			}
		}
		line := stageLine(stage)
		return []Issue{
			{
				Severity:    SeverityInfo,
				Problem:     "Image on line does not have a hashed tag, or illegal hash (does only support sha256).",
				Description: "Image may have been updated without your knowledge leading to potential malicious content.",
				Line:        &line,
				Fix:         "Look up the hash of the image and use the tag \":sha256:hash\"",
			},
		}
	},
}

var Rules = []Rule{
	NoLatestTag,
	NoUserDefined,
	NoHashTagImage,
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

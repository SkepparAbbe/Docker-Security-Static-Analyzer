package main

import (
	"fmt"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Severity string

func main() {
	f, err := os.Open("Dockerfile")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Parse the Dockerfile into an AST
	result, err := parser.Parse(f)
	if err != nil {
		panic(err) // syntax errors land here
	}

	// Convert the AST into a list of instructions
	stages, _, err := instructions.Parse(result.AST, nil)
	if err != nil {
		panic(err)
	}

	rules := []Rule{}
	issues := Analyze(stages, rules)

	for _, i := range issues {
		fmt.Printf("L%d  %-7s  %s\n", i.Line, i.Severity, i.Message)
	}
}

func Analyze(stages []instructions.Stage, rules []Rule) []Issue {
	var out []Issue
	for _, stage := range stages {

		for _, rule := range rules {
			if rule.CheckStage != nil {
				out = append(out, rule.CheckStage(stage)...)
			}
		}

		for _, cmd := range stage.Commands {
			for _, rule := range rules {
				switch c := cmd.(type) {
				case *instructions.RunCommand:
					if rule.CheckRun != nil {
						out = append(out, rule.CheckRun(c)...)
					}
				case *instructions.CopyCommand:
					if rule.CheckCopy != nil {
						out = append(out, rule.CheckCopy(c)...)
					}
				case *instructions.UserCommand:
					if rule.CheckUser != nil {
						out = append(out, rule.CheckUser(c)...)
					}
				}
			}
		}
	}
	return out
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Severity string

func main() {
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

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

	if *debugFlag {
		debug(stages)
	}

	issues := CheckConfig()
	issues = append(issues, Analyze(stages, Rules)...)

	for _, i := range issues {
		if i.Line != nil {
			fmt.Printf("L%d  %-7s  %s\n%s\n%s\n\n", *i.Line, i.Severity, i.Problem, i.Description, i.Fix)

		} else {
			fmt.Printf("config  %-7s  %s\n%s\n%s\n\n", i.Severity, i.Problem, i.Description, i.Fix)
		}
	}
}

// Iterates over each stage and each instruction (command) in said stage.
// Calls rule.go/Check...() for each stage instruction.
func Analyze(stages []instructions.Stage, rules []Rule) []Issue {
	var out []Issue
	for i := range stages {
		stage := &stages[i]

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
				case *instructions.AddCommand:
					if rule.CheckUser != nil {
						out = append(out, rule.CheckAdd(c)...)
					}
				}
			}
		}
	}
	return out
}

func debug(stages []instructions.Stage) {
	fmt.Printf("Parsed %d stage(s)\n\n", len(stages))
	for i, stage := range stages {
		fmt.Printf("Stage %d\n", i)
		fmt.Printf("├─ BaseName: %s\n", stage.BaseName)
		fmt.Printf("├─ As:       %s\n", dash(stage.Name))
		fmt.Printf("├─ Platform: %s\n", dash(stage.Platform))
		fmt.Printf("└─ Instructions (%d):\n", len(stage.Commands)+1)

		// FROM — synthesized from Stage fields, not from Commands
		fromLine := 0
		if len(stage.Location) > 0 {
			fromLine = stage.Location[0].Start.Line
		}
		fmt.Printf("   L%-3d %-7s %s", fromLine, "FROM", stage.BaseName)
		if stage.Name != "" {
			fmt.Printf(" AS %s", stage.Name)
		}
		fmt.Println()

		// Everything else
		for _, cmd := range stage.Commands {
			line := 0
			if loc := cmd.Location(); len(loc) > 0 {
				line = loc[0].Start.Line
			}
			fmt.Printf("   L%-3d %-7s %s\n", line, cmd.Name(), cmd)
		}
	}
}

func dash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

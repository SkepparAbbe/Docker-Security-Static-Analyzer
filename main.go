package main

import (
	"fmt"
	"os"
	"flag"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Severity string

func main() {
	debugFlag := flag.Bool("debug" , false, "Enable debug mode");
	fixFlag := flag.Bool("fix", false, "Show fixes of potential problems")
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

	//rules := []Rule{}
	//rules = append(rules, NoLatestTag)

	issues := CheckConfig()
	issues = append(issues, Analyze(stages, Rules)...)

	for _, i := range issues {
		if i.Line != nil {
			fmt.Printf("L%d  %-7s  %s\n", *i.Line, i.Severity, i.Message)

		}else {
			fmt.Printf("config  %-7s  %s\n", i.Severity, i.Message)
		}
		if *fixFlag {
			fmt.Printf("Fix:  %s\n", i.Fix)
		}
	}
}


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
		fmt.Println("Docker Configuration:")
		fmt.Println(exportInfo())
	}
}

func dash(s string) string {
    if s == "" {
        return "—"
    }
    return s
}
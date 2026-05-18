package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/charmbracelet/lipgloss"
)

type Severity string

func main() {
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	var f *os.File

	if flag.NArg() >= 1 {
		path := flag.Arg(0)
		f = openFile(path)
	} else {
		f = openFile("Dockerfile")
	}

	// Parse the Dockerfile into an AST
	result, err := parser.Parse(f)
	if err != nil {
		panic(err)
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

	printIssues(issues)
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
				case *instructions.UserCommand:
					if rule.CheckUser != nil {
						out = append(out, rule.CheckUser(c)...)
					}
				case *instructions.AddCommand:
					if rule.CheckAdd != nil {
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
		fmt.Printf("BaseName: %s\n", stage.BaseName)
		fmt.Printf("OrigCmd: %s\n", stage.OrigCmd)
		fmt.Printf("As:       %s\n", stage.Name)
		fmt.Printf("Platform: %s\n", stage.Platform)
		fmt.Printf("Instructions (%d):\n", len(stage.Commands)+1)

		for _, cmd := range stage.Commands {
			line := 0
			if loc := cmd.Location(); len(loc) > 0 {
				line = loc[0].Start.Line
			}
			fmt.Printf("   L%-3d %-7s %s\n", line, cmd.Name(), cmd)
		}
	}
}


func openFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}


var (
	labelStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("244"))
	fieldLabelStyle = labelStyle.Width(13)
	prefixStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Width(7)
	sepStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) 

	severityColors = map[string]lipgloss.Color{
		"warning":   lipgloss.Color("3"), // yellow
		"info":      lipgloss.Color("4"), // blue
	}
)

func severityColor(s string) lipgloss.Color {
	c, ok := severityColors[s];
	if ok {
		return c
	}
	return lipgloss.Color("7") // fallback gray
}

func printIssues(issues []Issue) {
	for _, i := range issues {
		prefix := "config"
		if i.Line != nil {
			prefix = fmt.Sprintf("Line: %d", *i.Line)
		}

		color := severityColor(string(i.Severity))
		sevStyle := lipgloss.NewStyle().Bold(true).Foreground(color)

		lines := []string{
			fmt.Sprintf("%s %s %s %s",
				prefixStyle.Render(prefix),
				sepStyle.Render("|"),
				labelStyle.Render("Severity:"),
				sevStyle.Render(string(i.Severity)),
			),
			fieldLabelStyle.Render("Problem:") + i.Problem,
			fieldLabelStyle.Render("Description:") + i.Description,
			fieldLabelStyle.Render("Fix:") + i.Fix,
		}

		block := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(color).
			PaddingLeft(1).
			MarginBottom(1)

		fmt.Println(block.Render(strings.Join(lines, "\n")))
	}
}
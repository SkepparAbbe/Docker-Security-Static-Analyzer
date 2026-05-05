package main

import (
	"fmt"
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func main() {
	f, err := os.Open("Dockerfile")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	result, err := parser.Parse(f)
	if err != nil {
		panic(err) // syntax errors land here
	}

	// Each child of AST is one instruction
	for _, node := range result.AST.Children {
		fmt.Printf("L%d  %s", node.StartLine, node.Value)

		// Walk the linked-list of arguments
		for n := node.Next; n != nil; n = n.Next {
			fmt.Printf(" %q", n.Value)
		}
		if len(node.Flags) > 0 {
			fmt.Printf("  flags=%v", node.Flags)
		}
		fmt.Println()
	}
}
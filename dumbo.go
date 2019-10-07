package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments")
		fmt.Println("Usage: dumbo <INPUT_DIR> <OUTPUT_DIR>")
		os.Exit(1)
	}

	baseDir := args[0]
	destDir := args[1]

	startStep("reading input files")

	inputFiles, err := scan(baseDir)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startStep("loading templates")

	templates, err := loadTemplates(inputFiles, baseDir)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startStep("building site")
	err = build(inputFiles, templates, destDir)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startStep(step string) {
	fmt.Printf("\n=== %v ===\n\n", strings.ToUpper(step))
}

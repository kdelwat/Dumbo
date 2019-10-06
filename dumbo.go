package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments")
		os.Exit(1)
	}

	baseDir := args[0]
	destDir := args[1]

	inputFiles, err := scan(baseDir)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("LOADING TEMPLATES")
	templates, err := loadTemplates(inputFiles, baseDir)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	build(inputFiles, templates, destDir)
}

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
	// destDir := args[1]

	scan(baseDir)
}

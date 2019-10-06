package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type inputFile struct {
	source    string
	dest      string
	template  string
	extension string
}

func scan(basedir string) ([]inputFile, error) {
	fmt.Printf("Building site at %s\n", basedir)

	inputFiles := []inputFile{}

	err := filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path: %q: %v\n", path, err)
			return err
		}

		if info.IsDir() && info.Name() == "_templates" {
			return filepath.SkipDir
		}

		fmt.Printf("Visited file or dir: %q\n", path)

		extension := filepath.Ext(path)
		dest := getDest(path, basedir)

		fmt.Printf("\tExtension: %v\n", extension)
		fmt.Printf("\tDestination: %v\n", dest)
		if extension == ".html" {
			inputFiles = append(inputFiles, inputFile{source: path, dest: dest, extension: "html"})
		}

		if extension == ".md" {
			template, err := getTemplate(path)

			if err != nil {
				return err
			}

			fmt.Printf("\tTemplate: %v\n", template)
			inputFiles = append(inputFiles, inputFile{source: path, dest: dest, template: template, extension: "md"})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error scanning site directory: %v\n", err)
	}

	return inputFiles, nil
}

func getDest(path string, basedir string) string {
	baseDirRe := regexp.MustCompile(basedir + "/")

	withoutBaseDir := baseDirRe.ReplaceAllString(path, "")

	suffixRe := regexp.MustCompile(`\..*$`)

	return suffixRe.ReplaceAllString(withoutBaseDir, "")
}

func getTemplate(path string) (string, error) {
	re := regexp.MustCompile(`\.(\w*)\.\w*$`)

	matches := re.FindAllString(path, 1)

	if len(matches) == 0 {
		return "", fmt.Errorf("ERROR: could not extract template name from path: %v", path)
	}

	match := matches[0]

	template := strings.Split(match, ".")[1]

	return template, nil
}

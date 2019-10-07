package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// inputFile contains the metadata needed to read and build a HTML
// page based on an input file (either HTML or Markdown)
type inputFile struct {
	source    string
	dest      string
	template  string
	extension string
}

// scan searches a given directory for all HTML and Markdown files and
// readies them for building
func scan(basedir string) ([]inputFile, error) {
	inputFiles := []inputFile{}

	// Walk the directory recursively. For each file, if it's HTML or Markdown,
	// construct metadata that will be used when building
	err := filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path: %q: %v\n", path, err)
			return err
		}

		// Skip templates directory
		if info.IsDir() && info.Name() == "_templates" {
			return filepath.SkipDir
		}

		// Get the file's corresponding destination path in the build output
		// e.g. for file `baseDir/a/b/index.html`, this will be `a/b/index`
		dest := getDest(path, basedir)

		// Get the file extension
		extension := filepath.Ext(path)

		if extension == ".html" {
			inputFiles = append(inputFiles, inputFile{source: path, dest: dest, extension: "html"})
		}

		if extension == ".md" {
			// If the file is Markdown, it is rendered via a template
			template, err := getTemplate(path)

			if err != nil {
				return err
			}

			inputFiles = append(inputFiles, inputFile{source: path, dest: dest, template: template, extension: "md"})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error scanning site directory: %v", err)
	}

	// Print information about input files
	for _, inputFile := range inputFiles {
		fmt.Printf("[VISITED] %v\n          Extension: %v\n          Template: %v\n", inputFile.source, inputFile.extension, inputFile.template)
	}

	return inputFiles, nil
}

// getDest converts an input path to an output path
// e.g. baseDir/a/b/index.html -> a/b/index
func getDest(path string, basedir string) string {
	baseDirRe := regexp.MustCompile(basedir + "/")

	withoutBaseDir := baseDirRe.ReplaceAllString(path, "")

	suffixRe := regexp.MustCompile(`\..*$`)

	return suffixRe.ReplaceAllString(withoutBaseDir, "")
}

// getTemplate extracts the name of the template used for a Markdown file
// e.g. baseDir/a/b/test.page.md -> page
func getTemplate(path string) (string, error) {
	re := regexp.MustCompile(`\.(\w+)\.\w+$`)

	matches := re.FindAllString(path, 1)

	if len(matches) == 0 {
		return "", fmt.Errorf("could not extract template name from path: %v", path)
	}

	match := matches[0]

	template := strings.Split(match, ".")[1]

	return template, nil
}

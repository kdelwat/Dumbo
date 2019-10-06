package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type html struct {
	source string
	dest   string
}

type md struct {
	source   string
	dest     string
	template string
}

func build(basedir string) {
	fmt.Printf("Building site at %s\n", basedir)

	htmlFiles := []html{}
	mdFiles := []md{}

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
			htmlFiles = append(htmlFiles, html{source: path, dest: dest})
		}

		if extension == ".md" {
			template, err := getTemplate(path)

			if err != nil {
				return err
			}

			fmt.Printf("\tTemplate: %v\n", template)
			mdFiles = append(mdFiles, md{source: path, dest: dest, template: template})
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning site directory: %v\n", err)
	}

	fmt.Println(htmlFiles)
	fmt.Println(mdFiles)
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

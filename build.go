package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/russross/blackfriday"
)

type templateInput struct {
	Content string
	Title   string
}

// build takes input file metadata and templates as input and builds the site in the desired
// output directory
func build(inputFiles []inputFile, templates map[string]*template.Template, destDir string) error {
	// If the output directory exists, delete it ready for rebuild
	doesOutputDirExist, err := exists(destDir)

	if err != nil {
		return fmt.Errorf("Could not check existence of output directory %q: %v", destDir, err)
	}

	if doesOutputDirExist {
		err := os.RemoveAll(destDir)

		if err != nil {
			return fmt.Errorf("Could not remove output directory %q: %v", destDir, err)
		}

	}

	// Create the output directory
	err = os.MkdirAll(destDir, 0755)

	if err != nil {
		return fmt.Errorf("Could not create destination directory %q: %v", destDir, err)
	}

	// Build files
	for _, input := range inputFiles {
		err = buildFile(input, templates, destDir)

		if err != nil {
			return err
		}
	}

	return nil
}

// buildFile renders an input file as HTML
func buildFile(input inputFile, templates map[string]*template.Template, destDir string) error {
	// Get the output filename
	destFileName := makeDestFilename(input.dest, destDir)

	// Create the path to the filename if needed
	err := os.MkdirAll(filepath.Dir(destFileName), 0755)

	if err != nil {
		return fmt.Errorf("Could not create directory %q: %v", filepath.Dir(destFileName), err)
	}

	if input.extension == "html" {
		return buildHTML(input, destFileName)
	} else if input.extension == "md" {
		return buildMD(input, templates, destFileName)
	} else {
		return fmt.Errorf("Unsupported file type: %v", input.extension)
	}
}

// buildHTML builds a HTML file by copying it from input to output unmodified
func buildHTML(input inputFile, dest string) error {
	contents, err := ioutil.ReadFile(input.source)

	if err != nil {
		return fmt.Errorf("Could not read file %q: %v", input.source, err)
	}

	err = ioutil.WriteFile(dest, contents, 0644)

	if err != nil {
		return fmt.Errorf("Could not write file %q: %v", dest, err)
	}

	fmt.Printf("[BUILT]   %q\n", dest)

	return nil
}

// buildMD builds a Markdown file by rendering it as HTML and including it in the relevant template
func buildMD(input inputFile, templates map[string]*template.Template, dest string) error {
	// Read the Markdown file
	contents, err := ioutil.ReadFile(input.source)

	if err != nil {
		return fmt.Errorf("Could not read file %q: %v", input.source, err)
	}

	// Render as HTML
	html := blackfriday.Run(contents)

	title := getTitle(html)

	// Open the destination file for writing
	outFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("Could not open for writing file %q: %v", dest, err)
	}

	w := bufio.NewWriter(outFile)

	// Render the template and write it to the destination
	templates[input.template].Execute(w, templateInput{Content: string(html), Title: title})

	err = w.Flush()

	if err != nil {
		return fmt.Errorf("Could not write to file %q: %v", dest, err)
	}

	err = outFile.Close()

	if err != nil {
		return fmt.Errorf("Could not close file %q: %v", dest, err)
	}

	fmt.Printf("[BUILT]   %q (%v)\n", dest, title)

	return nil
}

func makeDestFilename(dest string, destDir string) string {
	return path.Join(destDir, dest) + ".html"
}

func hasSelectedTemplate(templateNames []string, target string) bool {
	for _, name := range templateNames {
		if name == target {
			return true
		}
	}

	return false
}

// https://stackoverflow.com/a/10510783
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getTitle(html []byte) string {
	re := regexp.MustCompile(`<h1>(.*)</h1>`)

	submatches := re.FindStringSubmatch(string(html))

	if len(submatches) == 0 {
		return "Unknown"
	}

	return string(submatches[1])
}

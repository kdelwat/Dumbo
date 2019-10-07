package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/russross/blackfriday"
)

type templateInput struct {
	Content string
	Title   string
}

func build(inputFiles []inputFile, templates map[string]*template.Template, destDir string) error {
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

	err = os.MkdirAll(destDir, 0755)

	if err != nil {
		return fmt.Errorf("Could not create destination directory %q: %v", destDir, err)
	}

	// Build files
	for _, input := range inputFiles {
		err = buildFile(input, templates, destDir)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func buildFile(input inputFile, templates map[string]*template.Template, destDir string) error {
	destFileName := makeDestFilename(input.dest, destDir)

	if input.extension == "html" {
		return buildHTML(input, destFileName)
	} else if input.extension == "md" {
		return buildMD(input, templates, destFileName)
	} else {
		return fmt.Errorf("Unsupported file type: %v", input.extension)
	}
}

func buildHTML(input inputFile, dest string) error {
	fmt.Printf("Building HTML file\n\tSource: %q\n\tDest: %q\n", input.source, dest)

	contents, err := ioutil.ReadFile(input.source)

	if err != nil {
		return fmt.Errorf("Could not read file %q: %v", input.source, err)
	}

	err = ioutil.WriteFile(dest, contents, 0644)

	if err != nil {
		return fmt.Errorf("Could not write file %q: %v", dest, err)
	}

	return nil
}

func buildMD(input inputFile, templates map[string]*template.Template, dest string) error {
	fmt.Printf("Building MD file\n\tSource: %q\n\tDest: %q\n\tTemplate: %v\n", input.source, dest, input.template)

	contents, err := ioutil.ReadFile(input.source)

	if err != nil {
		return fmt.Errorf("Could not read file %q: %v", input.source, err)
	}

	html := blackfriday.Run(contents)

	err = os.MkdirAll(filepath.Dir(dest), 0755)

	if err != nil {
		return fmt.Errorf("Could not create directory %q: %v", filepath.Dir(dest), err)
	}

	outFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("Could not open for writing file %q: %v", dest, err)
	}

	w := bufio.NewWriter(outFile)

	templates[input.template].Execute(w, templateInput{Content: string(html), Title: "A title"})

	err = w.Flush()

	if err != nil {
		return fmt.Errorf("Could not write to file %q: %v", dest, err)
	}

	err = outFile.Close()

	if err != nil {
		return fmt.Errorf("Could not close file %q: %v", dest, err)
	}

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

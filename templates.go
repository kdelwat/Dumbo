package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"text/template"
)

func loadTemplates(inputFiles []inputFile, baseDir string) (map[string]*template.Template, error) {
	templateNames := []string{}

	// Load templates
	for _, input := range inputFiles {
		if input.template == "" || hasSelectedTemplate(templateNames, input.template) {
			continue
		}

		templateNames = append(templateNames, input.template)
	}

	fmt.Println(templateNames)

	templates := make(map[string]*template.Template)
	for _, name := range templateNames {
		templatePath := path.Join(baseDir, "_templates", name) + ".html"

		fmt.Println(templatePath)
		templateContents, err := ioutil.ReadFile(templatePath)

		fmt.Println(string(templateContents))
		if err != nil {
			return nil, fmt.Errorf("Could not read template %v: %v", name, err)
		}

		templates[name], err = template.New(name).Parse(string(templateContents))

		if err != nil {
			return nil, fmt.Errorf("Could not parse template %v: %v", name, err)
		}
	}

	return templates, nil
}

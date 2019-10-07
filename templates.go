package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"text/template"
)

// loadTemplates loads all HTML template files that will be used to render Markdown
func loadTemplates(inputFiles []inputFile, baseDir string) (map[string]*template.Template, error) {
	// Create array of template names that will be used in building
	templateNames := []string{}
	for _, input := range inputFiles {
		if input.template == "" || hasSelectedTemplate(templateNames, input.template) {
			continue
		}

		templateNames = append(templateNames, input.template)
		fmt.Printf("[LOADED]  %v\n", input.template)
	}

	// For each template name, read it from the _templates directory in the input
	templates := make(map[string]*template.Template)
	for _, name := range templateNames {
		templatePath := path.Join(baseDir, "_templates", name) + ".html"

		templateContents, err := ioutil.ReadFile(templatePath)

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

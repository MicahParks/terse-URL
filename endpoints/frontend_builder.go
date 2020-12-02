package endpoints

import (
	"bytes"
	"html/template"

	"github.com/MicahParks/terse-URL/restapi/operations"
)

var (

	// templateMap is a hardcoded mapping of HTML file paths to their appropriate template building function.
	templateMap = map[string]templateBuilder{
		"form.html":  buildForm,
		"table.html": buildTable,
	}
)

// templateBuilder is a function signature used to describe a function that can handle an template.
type templateBuilder func(fileData []byte, params operations.FrontendParams) (err error)

// templateData is the data structure passed to templates during execution.
type templateData struct {
	Shortened *string
	Original  *string
	Visits    *string
}

// buildForm builds the HTML form with the templated data.
func buildForm(fileData []byte, params operations.FrontendParams) (err error) {

	// Create the template.
	var tmpl *template.Template
	if tmpl, err = template.New("form").Parse(string(fileData)); err != nil {
		return err
	}

	// Create the templated data.
	data := templateData{
		Original:  params.Original,
		Shortened: params.Shortened,
	}

	// Execute the template on the file data.
	if err = tmpl.Execute(bytes.NewBuffer(fileData), data); err != nil {
		return err
	}

	return nil
}

// buildHTML builds the appropriate HTML file from a template. It does nothing if the HTML file has no template.
func buildHTML(fileData []byte, params operations.FrontendParams) (err error) {

	// If the file is recognized, send it to the appropriate templateBuilder.
	//
	// Find a better way to do this?
	builder, ok := templateMap[params.Path]
	if !ok {
		return nil
	}

	// Build from the appropriate template builder.
	return builder(fileData, params)
}

// buildTable builds the HTML table with the templated data.
func buildTable(fileData []byte, params operations.FrontendParams) (err error) {

	// Create the template.
	var tmpl *template.Template
	if tmpl, err = template.New("table").Parse(string(fileData)); err != nil {
		return err
	}

	// TODO Get all matching substrings.

	// Execute the template on the file data.
	if err = tmpl.Execute(bytes.NewBuffer(fileData), data); err != nil {
		return err
	}

	return nil
}

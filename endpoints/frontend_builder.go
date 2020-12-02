package endpoints

import (
	"github.com/MicahParks/terse-URL/restapi/operations"
)

type TableData struct {
	Shortened string
	Original  string
	Visits    string
}

type FormData struct {
	Shortened string
	Original  string
}

func buildHTML(fileData []byte, params operations.FrontendParams) (err error) {

	//
	switch params.Path {
	case "table.html":
		buildTable(fileData, params)
	case "form.html":
		buildForm(fileData, params)
	}

	return nil
}

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
		if err = buildTable(fileData, params); err != nil {
			return err
		}
	case "form.html":
		if err = buildForm(fileData, params); err != nil {
			return err
		}
	}

	return nil
}

func buildForm(fileData []byte, params operations.FrontendParams) (err error) {

}

func buildTable(fileData []byte, params operations.FrontendParams) (err error) {

}

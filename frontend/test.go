package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type TableData struct {
	Shortened string
	Original  string
	Visits    string
}

func main() {
	fileData, err := ioutil.ReadFile("table.gohtml")
	if err != nil {
		log.Fatalf("Can't read file.\nError: %s", err.Error())
	}
	tmpl, err := template.New("table").Parse(string(fileData))
	if err != nil {
		log.Fatalf("Failed to parse template.\nError: %s", err.Error())
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data := []TableData{
			{Shortened: "shorty", Original: "originally", Visits: "10"},
			{Shortened: "shorty23", Original: "originally23", Visits: "23"},
			{Shortened: "aa", Original: "aa", Visits: "11"},
			{Shortened: "zz", Original: "zz", Visits: "99"},
		}
		if err := tmpl.Execute(writer, data); err != nil {
			log.Fatalf("Failed to handle request.\nError: %s", err.Error())
		}
	})
	http.HandleFunc("/table.css", func(writer http.ResponseWriter, request *http.Request) {
		fileData, err := ioutil.ReadFile("table.css")
		if err != nil {
			log.Fatalf("css failed\nError: %s", err.Error())
		}
		if _, err = io.Copy(writer, bytes.NewReader(fileData)); err != nil {
			log.Fatalf("lol%s", err.Error())
		}
	})
	http.HandleFunc("/table.js", func(writer http.ResponseWriter, request *http.Request) {
		fileData, err := ioutil.ReadFile("table.js")
		if err != nil {
			log.Fatalf("css failed\nError: %s", err.Error())
		}
		if _, err = io.Copy(writer, bytes.NewReader(fileData)); err != nil {
			log.Fatalf("lol%s", err.Error())
		}
	})
	http.HandleFunc("/search.css", func(writer http.ResponseWriter, request *http.Request) {
		fileData, err := ioutil.ReadFile("search.css")
		if err != nil {
			log.Fatalf("css failed\nError: %s", err.Error())
		}
		if _, err = io.Copy(writer, bytes.NewReader(fileData)); err != nil {
			log.Fatalf("lol%s", err.Error())
		}
	})
	http.HandleFunc("/form.css", func(writer http.ResponseWriter, request *http.Request) {
		fileData, err := ioutil.ReadFile("form.css")
		if err != nil {
			log.Fatalf("css failed\nError: %s", err.Error())
		}
		if _, err = io.Copy(writer, bytes.NewReader(fileData)); err != nil {
			log.Fatalf("lol%s", err.Error())
		}
	})
	http.HandleFunc("/form.html", func(writer http.ResponseWriter, request *http.Request) {
		fileData, err := ioutil.ReadFile("form.html")
		if err != nil {
			log.Fatalf("css failed\nError: %s", err.Error())
		}
		if _, err = io.Copy(writer, bytes.NewReader(fileData)); err != nil {
			log.Fatalf("lol%s", err.Error())
		}
	})
	http.HandleFunc("/form.js", func(writer http.ResponseWriter, request *http.Request) {
		fileData, err := ioutil.ReadFile("form.js")
		if err != nil {
			log.Fatalf("css failed\nError: %s", err.Error())
		}
		if _, err = io.Copy(writer, bytes.NewReader(fileData)); err != nil {
			log.Fatalf("lol%s", err.Error())
		}
	})
	http.ListenAndServe(":9000", nil)
}

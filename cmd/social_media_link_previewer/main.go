package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/MicahParks/terse-URL/models"
)

func main() {

	// Create the metadata.
	meta := &models.MediaPreview{
		AudioURL:     "",
		CanonicalURL: "https://socialmedialinkpreview.micahparks.com/preview",
		Description:  "This is a social media link preview description.",
		Determiner:   "",
		ImageURL:     "https://wallpapercave.com/wp/wp314971.jpg",
		Locale:       "",
		LocaleAlt:    "",
		SiteName:     "",
		Title:        "The Terse URL social media link preview prototype.",
		Twitter: &models.Twitter{
			Card:        "summary_large_image",
			Creator:     "",
			Description: "This is a social media link preview description for Twitter.",
			ImageURL:    "https://i2.wp.com/entertainmentmesh.com/wp-content/uploads/2017/02/cute-parrot-birds.jpg",
			Site:        "",
			SiteID:      "",
			StreamURL:   "",
			Title:       "The Terse URL social media link preview prototype. Twitter version.",
		},
		Type:     "website",
		VideoURL: "",
	}

	// Read the template file.
	fileData, err := ioutil.ReadFile("socialMediaLinkPreview.gohtml")
	if err != nil {
		log.Fatalf("Failed to read file.\nError: %s", err.Error())
	}

	// Create the template.
	tmpl := template.Must(template.New("linkPreview").Parse(string(fileData)))

	// Put the metadata into the template.
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, meta); err != nil {
		log.Fatalf("Failed to execute template.\nError: %s", err.Error())
	}

	// Server up the templated HTML with metadata on the /preview endpoint.
	http.HandleFunc("/preview", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write(buf.Bytes()); err != nil {
			log.Fatalf("Failed to write HTTP response.\nError: %s", err.Error())
		}
	})

	// Sanity log.
	log.Println("Starting up.")

	// Serve via HTTP.
	log.Fatal(http.ListenAndServe(":30000", nil))
}

package terse_URL

import (
	"github.com/MicahParks/terse-URL/models"
)

type LinkPreview struct { // TODO Add favicon info if possible.
	models.MediaPreview `json:"mediaPreview"`
	Redirect            string `json:"redirect"`
}

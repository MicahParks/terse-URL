package terse_URL

import (
	"github.com/MicahParks/terse-URL/models"
)

type LinkPreview struct {
	Favicon             string `json:"favicon"` // TODO Figure out if this is possible.
	models.MediaPreview `json:"mediaPreview"`
	Redirect            string `json:"redirect"`
}

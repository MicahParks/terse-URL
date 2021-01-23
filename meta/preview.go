package meta

import (
	"io"
	"strings"

	"golang.org/x/net/html"

	"github.com/MicahParks/terseurl/models"
)

const (

	// metaContent is the HTML attribute key for content.
	metaContent = "content"

	// metaName is the HTML attribute key for name.
	metaName = "name"

	// metaProperty is the HTML attribute key for property.
	metaProperty = "property"
)

// Preview represents the data structure placed into the Go HTML template when doing a redirect through JavaScript or
// HTML meta tags.
type Preview struct { // TODO Add favicon info if possible.
	models.MediaPreview `json:"mediaPreview"`
	Redirect            string `json:"redirect"`
}

// PreviewTagInfo parses the given HTML io.Reader and returns the social media link preview meta tag information.
func PreviewTagInfo(body io.Reader) (tagInfo models.MediaPreview, err error) {

	// Create an HTML tokenizer.
	tokenizer := html.NewTokenizer(body)

	// Make the HTML meta tag info data structure.
	tagInfo = models.MediaPreview{
		Og:      make(map[string]string),
		Twitter: make(map[string]string),
	}

	// Iterate through the HTML tokens and find the meta tags, extract the social media link preview info, and put it
	// into a map.
	for {

		// Get the token type and token.
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		switch tokenType {
		case html.ErrorToken:
			return tagInfo, tokenizer.Err()
		case html.EndTagToken:
			if token.Data == "head" {
				return tagInfo, nil
			}
		}

		// Only look for HTML meta tags.
		if token.Data == "meta" {
			if len(token.Attr) >= 2 {

				// Turn the attributes into a map.
				attributes := make(map[string]string)

				// Get the tag's attribute's key value pairs.
				for _, attr := range token.Attr {
					attributes[attr.Key] = attr.Val
				}

				// Keep track of Open Graph protocol and Twitter social media link preview key value pairs.
				if strings.HasPrefix(attributes[metaProperty], "og:") {
					tagInfo.Og[attributes[metaProperty]] = attributes[metaContent]
				} else if strings.HasPrefix(attributes[metaName], "twitter:") {
					tagInfo.Twitter[attributes[metaName]] = attributes[metaContent]
				}
			}
		}
	}
}

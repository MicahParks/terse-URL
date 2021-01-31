// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "The Terse URL shortener.",
    "title": "Terse URL",
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    },
    "version": "0.0.1"
  },
  "host": "localhost",
  "basePath": "/",
  "paths": {
    "/api/alive": {
      "get": {
        "description": "Any non-200 response means the service is not alive.",
        "tags": [
          "system"
        ],
        "summary": "Used by Caddy or other reverse proxy to determine if the service is alive.",
        "operationId": "alive",
        "responses": {
          "200": {
            "description": "Service is alive."
          }
        }
      }
    },
    "/api/delete": {
      "delete": {
        "description": "All Terse and or Visits data will be deleted according to the deletion information specified.",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Delete all Terse and or Visits data.",
        "operationId": "terseDelete",
        "parameters": [
          {
            "description": "A JSON object containing the deletion information. If Terse or Visits data is marked for deletion, it will all be deleted.",
            "name": "delete",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Delete"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The deletion request was successfully fulfilled."
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/delete/some": {
      "delete": {
        "description": "If only Terse data is deleted, the API user is responsible for cleaning up its Visits data before adding new Terse data under the same shortened URL.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Delete Terse and or Visits data for the given shortened URL.",
        "operationId": "terseDeleteSome",
        "parameters": [
          {
            "description": "Indicate if Terse and or Visits data should be deleted and for which shortened URLs.",
            "name": "info",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "delete": {
                  "description": "Indicate if Terse and or Visits data should be deleted.",
                  "$ref": "#/definitions/Delete"
                },
                "shortenedURLs": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The shortened URL's data was successfully deleted from the backend storage."
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/export": {
      "get": {
        "description": "Depending on the underlying storage and amount of data, this may take a while.",
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Export all Terse and Visits data from the backend.",
        "operationId": "terseExport",
        "responses": {
          "200": {
            "description": "The export was successfully retrieved.",
            "schema": {
              "description": "All of the Terse and Visits data from the backend.",
              "type": "object",
              "additionalProperties": {
                "$ref": "#/definitions/Export"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/export/some": {
      "post": {
        "description": "Export Terse and Visits data for the given shortened URL.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Export Terse and Visits data for the given shortened URL.",
        "operationId": "terseExportSome",
        "parameters": [
          {
            "description": "The shortened URLs to get the export for.",
            "name": "shortenedURLs",
            "in": "body",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The export was successfully retrieved.",
            "schema": {
              "additionalProperties": {
                "$ref": "#/definitions/Export"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/import": {
      "post": {
        "description": "Any imported data will overwrite existing data. Unless deletion information is specified. In that case all Terse or Visits data can be deleted before importing the new data.",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Import existing Terse and Visits data for the given shortened URLs.",
        "operationId": "terseImport",
        "parameters": [
          {
            "description": "A JSON object containing the deletion information. If Terse or Visits data is marked for deletion, it will all be deleted. An object matching shortened URLs to their previously exported data is also required.",
            "name": "importDelete",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "delete": {
                  "$ref": "#/definitions/Delete"
                },
                "import": {
                  "additionalProperties": {
                    "$ref": "#/definitions/Export"
                  }
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The import request was successfully fulfilled."
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/prefix": {
      "get": {
        "description": "Provides the HTTP prefix all shortened URLs have.",
        "tags": [
          "api"
        ],
        "summary": "Client's web browser is requesting what HTTP prefix all shortened URLs have.",
        "operationId": "tersePrefix",
        "responses": {
          "200": {
            "description": "The HTTP prefix all shortened URLs have.",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/summary": {
      "post": {
        "description": "Terse summary data includes the shortened URL, the original URL, the type of redirect, and the number of visits.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Provide Terse summary data for the requested shortened URLs.",
        "operationId": "terseSummary",
        "parameters": [
          {
            "description": "The array of shortened URLs to get Terse summary data for. If none is provided, all will summaries will be returned.",
            "name": "shortened",
            "in": "body",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The map of shortened URLs to Terse summary data.",
            "schema": {
              "additionalProperties": {
                "$ref": "#/definitions/TerseSummary"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/terse/{shortened}": {
      "get": {
        "description": "Read the Terse data for the given shortened URL.",
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Read the Terse data for the given shortened URL.",
        "operationId": "terseTerse",
        "parameters": [
          {
            "type": "string",
            "description": "The shortened URL to get the Terse data for.",
            "name": "shortened",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The Terse data was successfully retrieved.",
            "schema": {
              "$ref": "#/definitions/Terse"
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/visits/{shortened}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Get the Visits data for the given shortened URL.",
        "operationId": "terseVisits",
        "parameters": [
          {
            "type": "string",
            "description": "The shortened URL to get the s data for.",
            "name": "shortened",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The Visits data was successfully retrieved.",
            "schema": {
              "description": "The visit data for the given shortened URL.",
              "type": "array",
              "items": {
                "$ref": "#/definitions/Visit"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/write/{operation}": {
      "post": {
        "description": "\"insert\" will fail if the shortened URL already exists. \"update\" will fail if the shortened URL does not already exist. \"upsert\" will only fail if there is a failure interacting with the underlying storage. If no shortened URL is included in the given Terse data, one will be generated randomly and returned in the response.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Perform a write operation on Terse data for a given shortened URL.",
        "operationId": "terseWrite",
        "parameters": [
          {
            "description": "The Terse data, with an optional shortened URL. If no shortened URL is given, one will be generated randomly and returned in the response. If no redirect type is given, 302 is used.",
            "name": "terse",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TerseInput"
            }
          },
          {
            "enum": [
              "insert",
              "update",
              "upsert"
            ],
            "type": "string",
            "description": "The write operation to perform with the Terse data.",
            "name": "operation",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The write operation was successful.",
            "schema": {
              "description": "The shortened URL affected.",
              "type": "string"
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/{shortened}": {
      "get": {
        "description": "Use the shortened URL. It will redirect to the full URL if it has not expired.",
        "produces": [
          "text/html"
        ],
        "tags": [
          "public"
        ],
        "summary": "Typically a web browser would visit this endpoint, then get redirected.",
        "operationId": "terseRedirect",
        "parameters": [
          {
            "type": "string",
            "name": "shortened",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The HTML document containing a social media link preview and or JavaScript fingerprinting. Any visitor will be automatically redirected to the original link with JavaScript.",
            "schema": {
              "type": "file"
            }
          },
          "301": {
            "description": "An HTTP response that will server as a permanent redirect to the shortened URL's full URL.",
            "headers": {
              "Location": {
                "type": "string",
                "description": "The full URL that the redirect leads to."
              }
            }
          },
          "302": {
            "description": "An HTTP response that will serve as a temporary redirect to the shortened URL's full URL.",
            "headers": {
              "Location": {
                "type": "string",
                "description": "The full URL that the redirect leads to."
              }
            }
          },
          "404": {
            "description": "The shortened URL expired or never existed."
          }
        }
      }
    }
  },
  "definitions": {
    "Delete": {
      "properties": {
        "terse": {
          "type": "boolean",
          "default": true
        },
        "visits": {
          "type": "boolean",
          "default": true
        }
      }
    },
    "Error": {
      "type": "object",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "x-nullable": false
        },
        "message": {
          "type": "string",
          "x-nullable": false
        }
      }
    },
    "Export": {
      "type": "object",
      "properties": {
        "terse": {
          "$ref": "#/definitions/Terse"
        },
        "visits": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Visit"
          }
        }
      }
    },
    "MediaPreview": {
      "properties": {
        "og": {
          "$ref": "#/definitions/OpenGraph"
        },
        "title": {
          "type": "string"
        },
        "twitter": {
          "$ref": "#/definitions/Twitter"
        }
      }
    },
    "OpenGraph": {
      "additionalProperties": {
        "type": "string"
      }
    },
    "RedirectType": {
      "type": "string",
      "enum": [
        "301",
        "302",
        "meta",
        "js"
      ]
    },
    "Terse": {
      "required": [
        "originalURL",
        "shortenedURL"
      ],
      "properties": {
        "javascriptTracking": {
          "type": "boolean"
        },
        "mediaPreview": {
          "$ref": "#/definitions/MediaPreview"
        },
        "originalURL": {
          "type": "string",
          "x-nullable": false
        },
        "redirectType": {
          "$ref": "#/definitions/RedirectType"
        },
        "shortenedURL": {
          "type": "string",
          "x-nullable": false
        }
      }
    },
    "TerseInput": {
      "required": [
        "originalURL"
      ],
      "properties": {
        "javascriptTracking": {
          "type": "boolean"
        },
        "mediaPreview": {
          "$ref": "#/definitions/MediaPreview"
        },
        "originalURL": {
          "type": "string",
          "x-nullable": false
        },
        "redirectType": {
          "$ref": "#/definitions/RedirectType"
        },
        "shortenedURL": {
          "type": "string"
        }
      }
    },
    "TerseSummary": {
      "properties": {
        "originalURL": {
          "type": "string"
        },
        "redirectType": {
          "$ref": "#/definitions/RedirectType"
        },
        "shortenedURL": {
          "type": "string"
        },
        "visitCount": {
          "type": "integer"
        }
      }
    },
    "Twitter": {
      "additionalProperties": {
        "type": "string"
      }
    },
    "Visit": {
      "required": [
        "accessed",
        "ip"
      ],
      "properties": {
        "accessed": {
          "type": "string",
          "format": "date-time"
        },
        "headers": {
          "type": "object",
          "additionalProperties": {
            "type": "array",
            "items": {
              "type": "string"
            }
          }
        },
        "ip": {
          "type": "string"
        }
      }
    }
  },
  "tags": [
    {
      "description": "Endpoints to perform operations on Terse data.",
      "name": "api"
    },
    {
      "description": "Endpoints that are publicly accessible.",
      "name": "public"
    },
    {
      "description": "Endpoints required by the system that are not public facing and do not affect Terse data.",
      "name": "system"
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "The Terse URL shortener.",
    "title": "Terse URL",
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    },
    "version": "0.0.1"
  },
  "host": "localhost",
  "basePath": "/",
  "paths": {
    "/api/alive": {
      "get": {
        "description": "Any non-200 response means the service is not alive.",
        "tags": [
          "system"
        ],
        "summary": "Used by Caddy or other reverse proxy to determine if the service is alive.",
        "operationId": "alive",
        "responses": {
          "200": {
            "description": "Service is alive."
          }
        }
      }
    },
    "/api/delete": {
      "delete": {
        "description": "All Terse and or Visits data will be deleted according to the deletion information specified.",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Delete all Terse and or Visits data.",
        "operationId": "terseDelete",
        "parameters": [
          {
            "description": "A JSON object containing the deletion information. If Terse or Visits data is marked for deletion, it will all be deleted.",
            "name": "delete",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Delete"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The deletion request was successfully fulfilled."
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/delete/some": {
      "delete": {
        "description": "If only Terse data is deleted, the API user is responsible for cleaning up its Visits data before adding new Terse data under the same shortened URL.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Delete Terse and or Visits data for the given shortened URL.",
        "operationId": "terseDeleteSome",
        "parameters": [
          {
            "description": "Indicate if Terse and or Visits data should be deleted and for which shortened URLs.",
            "name": "info",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "delete": {
                  "description": "Indicate if Terse and or Visits data should be deleted.",
                  "$ref": "#/definitions/Delete"
                },
                "shortenedURLs": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The shortened URL's data was successfully deleted from the backend storage."
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/export": {
      "get": {
        "description": "Depending on the underlying storage and amount of data, this may take a while.",
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Export all Terse and Visits data from the backend.",
        "operationId": "terseExport",
        "responses": {
          "200": {
            "description": "The export was successfully retrieved.",
            "schema": {
              "description": "All of the Terse and Visits data from the backend.",
              "type": "object",
              "additionalProperties": {
                "$ref": "#/definitions/Export"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/export/some": {
      "post": {
        "description": "Export Terse and Visits data for the given shortened URL.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Export Terse and Visits data for the given shortened URL.",
        "operationId": "terseExportSome",
        "parameters": [
          {
            "description": "The shortened URLs to get the export for.",
            "name": "shortenedURLs",
            "in": "body",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The export was successfully retrieved.",
            "schema": {
              "additionalProperties": {
                "$ref": "#/definitions/Export"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/import": {
      "post": {
        "description": "Any imported data will overwrite existing data. Unless deletion information is specified. In that case all Terse or Visits data can be deleted before importing the new data.",
        "consumes": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Import existing Terse and Visits data for the given shortened URLs.",
        "operationId": "terseImport",
        "parameters": [
          {
            "description": "A JSON object containing the deletion information. If Terse or Visits data is marked for deletion, it will all be deleted. An object matching shortened URLs to their previously exported data is also required.",
            "name": "importDelete",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "delete": {
                  "$ref": "#/definitions/Delete"
                },
                "import": {
                  "additionalProperties": {
                    "$ref": "#/definitions/Export"
                  }
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The import request was successfully fulfilled."
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/prefix": {
      "get": {
        "description": "Provides the HTTP prefix all shortened URLs have.",
        "tags": [
          "api"
        ],
        "summary": "Client's web browser is requesting what HTTP prefix all shortened URLs have.",
        "operationId": "tersePrefix",
        "responses": {
          "200": {
            "description": "The HTTP prefix all shortened URLs have.",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/summary": {
      "post": {
        "description": "Terse summary data includes the shortened URL, the original URL, the type of redirect, and the number of visits.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Provide Terse summary data for the requested shortened URLs.",
        "operationId": "terseSummary",
        "parameters": [
          {
            "description": "The array of shortened URLs to get Terse summary data for. If none is provided, all will summaries will be returned.",
            "name": "shortened",
            "in": "body",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "The map of shortened URLs to Terse summary data.",
            "schema": {
              "additionalProperties": {
                "$ref": "#/definitions/TerseSummary"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/terse/{shortened}": {
      "get": {
        "description": "Read the Terse data for the given shortened URL.",
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Read the Terse data for the given shortened URL.",
        "operationId": "terseTerse",
        "parameters": [
          {
            "type": "string",
            "description": "The shortened URL to get the Terse data for.",
            "name": "shortened",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The Terse data was successfully retrieved.",
            "schema": {
              "$ref": "#/definitions/Terse"
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/visits/{shortened}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Get the Visits data for the given shortened URL.",
        "operationId": "terseVisits",
        "parameters": [
          {
            "type": "string",
            "description": "The shortened URL to get the s data for.",
            "name": "shortened",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The Visits data was successfully retrieved.",
            "schema": {
              "description": "The visit data for the given shortened URL.",
              "type": "array",
              "items": {
                "$ref": "#/definitions/Visit"
              }
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/api/write/{operation}": {
      "post": {
        "description": "\"insert\" will fail if the shortened URL already exists. \"update\" will fail if the shortened URL does not already exist. \"upsert\" will only fail if there is a failure interacting with the underlying storage. If no shortened URL is included in the given Terse data, one will be generated randomly and returned in the response.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "api"
        ],
        "summary": "Perform a write operation on Terse data for a given shortened URL.",
        "operationId": "terseWrite",
        "parameters": [
          {
            "description": "The Terse data, with an optional shortened URL. If no shortened URL is given, one will be generated randomly and returned in the response. If no redirect type is given, 302 is used.",
            "name": "terse",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TerseInput"
            }
          },
          {
            "enum": [
              "insert",
              "update",
              "upsert"
            ],
            "type": "string",
            "description": "The write operation to perform with the Terse data.",
            "name": "operation",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The write operation was successful.",
            "schema": {
              "description": "The shortened URL affected.",
              "type": "string"
            }
          },
          "default": {
            "description": "Unexpected error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/{shortened}": {
      "get": {
        "description": "Use the shortened URL. It will redirect to the full URL if it has not expired.",
        "produces": [
          "text/html"
        ],
        "tags": [
          "public"
        ],
        "summary": "Typically a web browser would visit this endpoint, then get redirected.",
        "operationId": "terseRedirect",
        "parameters": [
          {
            "type": "string",
            "name": "shortened",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The HTML document containing a social media link preview and or JavaScript fingerprinting. Any visitor will be automatically redirected to the original link with JavaScript.",
            "schema": {
              "type": "file"
            }
          },
          "301": {
            "description": "An HTTP response that will server as a permanent redirect to the shortened URL's full URL.",
            "headers": {
              "Location": {
                "type": "string",
                "description": "The full URL that the redirect leads to."
              }
            }
          },
          "302": {
            "description": "An HTTP response that will serve as a temporary redirect to the shortened URL's full URL.",
            "headers": {
              "Location": {
                "type": "string",
                "description": "The full URL that the redirect leads to."
              }
            }
          },
          "404": {
            "description": "The shortened URL expired or never existed."
          }
        }
      }
    }
  },
  "definitions": {
    "Delete": {
      "properties": {
        "terse": {
          "type": "boolean",
          "default": true
        },
        "visits": {
          "type": "boolean",
          "default": true
        }
      }
    },
    "Error": {
      "type": "object",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "x-nullable": false
        },
        "message": {
          "type": "string",
          "x-nullable": false
        }
      }
    },
    "Export": {
      "type": "object",
      "properties": {
        "terse": {
          "$ref": "#/definitions/Terse"
        },
        "visits": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Visit"
          }
        }
      }
    },
    "MediaPreview": {
      "properties": {
        "og": {
          "$ref": "#/definitions/OpenGraph"
        },
        "title": {
          "type": "string"
        },
        "twitter": {
          "$ref": "#/definitions/Twitter"
        }
      }
    },
    "OpenGraph": {
      "additionalProperties": {
        "type": "string"
      }
    },
    "RedirectType": {
      "type": "string",
      "enum": [
        "301",
        "302",
        "meta",
        "js"
      ]
    },
    "Terse": {
      "required": [
        "originalURL",
        "shortenedURL"
      ],
      "properties": {
        "javascriptTracking": {
          "type": "boolean"
        },
        "mediaPreview": {
          "$ref": "#/definitions/MediaPreview"
        },
        "originalURL": {
          "type": "string",
          "x-nullable": false
        },
        "redirectType": {
          "$ref": "#/definitions/RedirectType"
        },
        "shortenedURL": {
          "type": "string",
          "x-nullable": false
        }
      }
    },
    "TerseInput": {
      "required": [
        "originalURL"
      ],
      "properties": {
        "javascriptTracking": {
          "type": "boolean"
        },
        "mediaPreview": {
          "$ref": "#/definitions/MediaPreview"
        },
        "originalURL": {
          "type": "string",
          "x-nullable": false
        },
        "redirectType": {
          "$ref": "#/definitions/RedirectType"
        },
        "shortenedURL": {
          "type": "string"
        }
      }
    },
    "TerseSummary": {
      "properties": {
        "originalURL": {
          "type": "string"
        },
        "redirectType": {
          "$ref": "#/definitions/RedirectType"
        },
        "shortenedURL": {
          "type": "string"
        },
        "visitCount": {
          "type": "integer"
        }
      }
    },
    "Twitter": {
      "additionalProperties": {
        "type": "string"
      }
    },
    "Visit": {
      "required": [
        "accessed",
        "ip"
      ],
      "properties": {
        "accessed": {
          "type": "string",
          "format": "date-time"
        },
        "headers": {
          "type": "object",
          "additionalProperties": {
            "type": "array",
            "items": {
              "type": "string"
            }
          }
        },
        "ip": {
          "type": "string"
        }
      }
    }
  },
  "tags": [
    {
      "description": "Endpoints to perform operations on Terse data.",
      "name": "api"
    },
    {
      "description": "Endpoints that are publicly accessible.",
      "name": "public"
    },
    {
      "description": "Endpoints required by the system that are not public facing and do not affect Terse data.",
      "name": "system"
    }
  ]
}`))
}

// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Retrieves a paginated list of songs with optional filtering based on group, song name, release date, text, and link",
                "tags": [
                    "songs"
                ],
                "summary": "Get songs list with filtering and pagination",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Group name for filtering",
                        "name": "groupName",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Song name for filtering",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Release date for filtering",
                        "name": "releaseDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Text for filtering",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Link for filtering",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Number of songs per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Songs retrieved successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid page or limit number",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to connect to database or retrieve songs",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new song with details like group, song name, release date, text, and link",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Add a new song to the library",
                "parameters": [
                    {
                        "description": "Song details",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Song created successfully",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    },
                    "400": {
                        "description": "Invalid input data",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to call external API or insert data into database",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/songs/{id}": {
            "put": {
                "description": "Updates the song information by its ID. Only provided fields will be updated.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Update song details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated song data",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song updated successfully",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    },
                    "400": {
                        "description": "Invalid JSON data",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Song not found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to connect to database or update song",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a song from the database by its ID",
                "tags": [
                    "songs"
                ],
                "summary": "Delete a song by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song deleted successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Song not found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to delete song or connect to database",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/songs/{id}/text": {
            "get": {
                "description": "Retrieves the song's text, paginated by verses, based on the song's ID",
                "tags": [
                    "songs"
                ],
                "summary": "Get song text by verses with pagination",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 2,
                        "description": "Number of verses per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song text retrieved successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid page or limit number",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Song not found or no verses on this page",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to connect to database or retrieve song text",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.Song": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Music Service API",
	Description:      "This is a service to manage songs in a library.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
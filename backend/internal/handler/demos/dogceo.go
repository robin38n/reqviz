package demos

// DogCEOSpec is an OpenAPI 3.0 spec for dog.ceo/api.
var DogCEOSpec = map[string]any{
	"openapi": "3.0.3",
	"info": map[string]any{
		"title":       "Dog CEO",
		"version":     "1.0.0",
		"description": "The internet's biggest collection of open-source dog pictures. Free, no auth required.",
	},
	"servers": []any{
		map[string]any{
			"url":         "https://dog.ceo/api",
			"description": "Production",
		},
	},
	"tags": []any{
		map[string]any{"name": "breeds", "description": "Breed listing"},
		map[string]any{"name": "images", "description": "Dog images"},
	},
	"paths": map[string]any{
		"/breeds/list/all": map[string]any{
			"get": map[string]any{
				"operationId": "listAllBreeds",
				"summary":     "List all breeds",
				"tags":        []any{"breeds"},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "All breeds with sub-breeds",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/BreedList"},
							},
						},
					},
				},
			},
		},
		"/breeds/image/random": map[string]any{
			"get": map[string]any{
				"operationId": "randomImage",
				"summary":     "Get a random dog image",
				"tags":        []any{"images"},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A random dog image URL",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/SingleImage"},
							},
						},
					},
				},
			},
		},
		"/breed/{breed}/images": map[string]any{
			"get": map[string]any{
				"operationId": "breedImages",
				"summary":     "Get all images for a breed",
				"tags":        []any{"images"},
				"parameters": []any{
					map[string]any{
						"name": "breed", "in": "path", "required": true,
						"schema": map[string]any{"type": "string"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "List of image URLs for the breed",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/ImageList"},
							},
						},
					},
					"404": map[string]any{"description": "Breed not found"},
				},
			},
		},
		"/breed/{breed}/images/random": map[string]any{
			"get": map[string]any{
				"operationId": "breedRandomImage",
				"summary":     "Get a random image for a breed",
				"tags":        []any{"images"},
				"parameters": []any{
					map[string]any{
						"name": "breed", "in": "path", "required": true,
						"schema": map[string]any{"type": "string"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A random image URL for the breed",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/SingleImage"},
							},
						},
					},
					"404": map[string]any{"description": "Breed not found"},
				},
			},
		},
	},
	"components": map[string]any{
		"schemas": map[string]any{
			"BreedList": map[string]any{
				"type":     "object",
				"required": []any{"status", "message"},
				"properties": map[string]any{
					"status": map[string]any{"type": "string"},
					"message": map[string]any{
						"type": "object",
						"additionalProperties": map[string]any{
							"type":  "array",
							"items": map[string]any{"type": "string"},
						},
					},
				},
			},
			"SingleImage": map[string]any{
				"type":     "object",
				"required": []any{"status", "message"},
				"properties": map[string]any{
					"status":  map[string]any{"type": "string"},
					"message": map[string]any{"type": "string", "format": "uri"},
				},
			},
			"ImageList": map[string]any{
				"type":     "object",
				"required": []any{"status", "message"},
				"properties": map[string]any{
					"status": map[string]any{"type": "string"},
					"message": map[string]any{
						"type":  "array",
						"items": map[string]any{"type": "string", "format": "uri"},
					},
				},
			},
		},
	},
}

package demos

// PokeAPISpec is an OpenAPI 3.0 spec for pokeapi.co/api/v2.
var PokeAPISpec = map[string]any{
	"openapi": "3.0.3",
	"info": map[string]any{
		"title":       "PokéAPI",
		"version":     "2.0.0",
		"description": "All the Pokémon data you'll ever need in one place. Read-only RESTful API.",
	},
	"servers": []any{
		map[string]any{
			"url":         "https://pokeapi.co/api/v2",
			"description": "Production",
		},
	},
	"tags": []any{
		map[string]any{"name": "pokemon", "description": "Pokémon data"},
		map[string]any{"name": "types", "description": "Pokémon types"},
		map[string]any{"name": "abilities", "description": "Pokémon abilities"},
		map[string]any{"name": "generations", "description": "Game generations"},
	},
	"paths": map[string]any{
		"/pokemon": map[string]any{
			"get": map[string]any{
				"operationId": "listPokemon",
				"summary":     "List pokémon",
				"tags":        []any{"pokemon"},
				"parameters": []any{
					map[string]any{
						"name": "limit", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer", "default": 20},
					},
					map[string]any{
						"name": "offset", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer", "default": 0},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Paginated list of pokémon",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/PaginatedList"},
							},
						},
					},
				},
			},
		},
		"/pokemon/{idOrName}": map[string]any{
			"get": map[string]any{
				"operationId": "getPokemon",
				"summary":     "Get a pokémon by ID or name",
				"tags":        []any{"pokemon"},
				"parameters": []any{
					map[string]any{
						"name": "idOrName", "in": "path", "required": true,
						"schema": map[string]any{"type": "string"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A single pokémon",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Pokemon"},
							},
						},
					},
					"404": map[string]any{"description": "Pokémon not found"},
				},
			},
		},
		"/type": map[string]any{
			"get": map[string]any{
				"operationId": "listTypes",
				"summary":     "List all types",
				"tags":        []any{"types"},
				"parameters": []any{
					map[string]any{
						"name": "limit", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer", "default": 20},
					},
					map[string]any{
						"name": "offset", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer", "default": 0},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Paginated list of types",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/PaginatedList"},
							},
						},
					},
				},
			},
		},
		"/type/{idOrName}": map[string]any{
			"get": map[string]any{
				"operationId": "getType",
				"summary":     "Get a type by ID or name",
				"tags":        []any{"types"},
				"parameters": []any{
					map[string]any{
						"name": "idOrName", "in": "path", "required": true,
						"schema": map[string]any{"type": "string"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A single type",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/PokemonType"},
							},
						},
					},
					"404": map[string]any{"description": "Type not found"},
				},
			},
		},
		"/ability": map[string]any{
			"get": map[string]any{
				"operationId": "listAbilities",
				"summary":     "List all abilities",
				"tags":        []any{"abilities"},
				"parameters": []any{
					map[string]any{
						"name": "limit", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer", "default": 20},
					},
					map[string]any{
						"name": "offset", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer", "default": 0},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Paginated list of abilities",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/PaginatedList"},
							},
						},
					},
				},
			},
		},
		"/ability/{idOrName}": map[string]any{
			"get": map[string]any{
				"operationId": "getAbility",
				"summary":     "Get an ability by ID or name",
				"tags":        []any{"abilities"},
				"parameters": []any{
					map[string]any{
						"name": "idOrName", "in": "path", "required": true,
						"schema": map[string]any{"type": "string"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A single ability",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Ability"},
							},
						},
					},
					"404": map[string]any{"description": "Ability not found"},
				},
			},
		},
		"/generation/{id}": map[string]any{
			"get": map[string]any{
				"operationId": "getGeneration",
				"summary":     "Get a generation by ID",
				"tags":        []any{"generations"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A single generation",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Generation"},
							},
						},
					},
					"404": map[string]any{"description": "Generation not found"},
				},
			},
		},
	},
	"components": map[string]any{
		"schemas": map[string]any{
			"NamedAPIResource": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{"type": "string"},
					"url":  map[string]any{"type": "string", "format": "uri"},
				},
			},
			"PaginatedList": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"count":    map[string]any{"type": "integer"},
					"next":     map[string]any{"type": "string", "format": "uri", "nullable": true},
					"previous": map[string]any{"type": "string", "format": "uri", "nullable": true},
					"results": map[string]any{
						"type":  "array",
						"items": map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
					},
				},
			},
			"Pokemon": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":              map[string]any{"type": "integer"},
					"name":            map[string]any{"type": "string"},
					"base_experience": map[string]any{"type": "integer"},
					"height":          map[string]any{"type": "integer"},
					"weight":          map[string]any{"type": "integer"},
					"is_default":      map[string]any{"type": "boolean"},
					"order":           map[string]any{"type": "integer"},
					"sprites": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"front_default": map[string]any{"type": "string", "format": "uri"},
							"front_shiny":   map[string]any{"type": "string", "format": "uri"},
							"back_default":  map[string]any{"type": "string", "format": "uri"},
							"back_shiny":    map[string]any{"type": "string", "format": "uri"},
						},
					},
					"types": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"slot": map[string]any{"type": "integer"},
								"type": map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
							},
						},
					},
					"abilities": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"is_hidden": map[string]any{"type": "boolean"},
								"slot":      map[string]any{"type": "integer"},
								"ability":   map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
							},
						},
					},
					"stats": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"base_stat": map[string]any{"type": "integer"},
								"effort":    map[string]any{"type": "integer"},
								"stat":      map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
							},
						},
					},
				},
			},
			"PokemonType": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":   map[string]any{"type": "integer"},
					"name": map[string]any{"type": "string"},
					"pokemon": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"slot":    map[string]any{"type": "integer"},
								"pokemon": map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
							},
						},
					},
				},
			},
			"Ability": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":   map[string]any{"type": "integer"},
					"name": map[string]any{"type": "string"},
					"effect_entries": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"effect":       map[string]any{"type": "string"},
								"short_effect": map[string]any{"type": "string"},
								"language":     map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
							},
						},
					},
					"pokemon": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"is_hidden": map[string]any{"type": "boolean"},
								"pokemon":   map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
							},
						},
					},
				},
			},
			"Generation": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":   map[string]any{"type": "integer"},
					"name": map[string]any{"type": "string"},
					"main_region": map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
					"pokemon_species": map[string]any{
						"type":  "array",
						"items": map[string]any{"$ref": "#/components/schemas/NamedAPIResource"},
					},
				},
			},
		},
	},
}

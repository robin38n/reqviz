package demos

// JSONPlaceholderSpec is an OpenAPI 3.0 spec for jsonplaceholder.typicode.com.
var JSONPlaceholderSpec = map[string]any{
	"openapi": "3.0.3",
	"info": map[string]any{
		"title":       "JSONPlaceholder",
		"version":     "1.0.0",
		"description": "Free fake REST API for testing and prototyping. Powered by jsonplaceholder.typicode.com.",
	},
	"servers": []any{
		map[string]any{
			"url":         "https://jsonplaceholder.typicode.com",
			"description": "Production",
		},
	},
	"tags": []any{
		map[string]any{"name": "posts", "description": "Blog post operations"},
		map[string]any{"name": "comments", "description": "Comment operations"},
		map[string]any{"name": "users", "description": "User operations"},
		map[string]any{"name": "todos", "description": "Todo operations"},
	},
	"paths": map[string]any{
		// --- Posts ---
		"/posts": map[string]any{
			"get": map[string]any{
				"operationId": "listPosts",
				"summary":     "List all posts",
				"tags":        []any{"posts"},
				"parameters": []any{
					map[string]any{
						"name": "userId", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A list of posts",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"type":  "array",
									"items": map[string]any{"$ref": "#/components/schemas/Post"},
								},
							},
						},
					},
				},
			},
			"post": map[string]any{
				"operationId": "createPost",
				"summary":     "Create a new post",
				"tags":        []any{"posts"},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{"$ref": "#/components/schemas/NewPost"},
						},
					},
				},
				"responses": map[string]any{
					"201": map[string]any{
						"description": "Post created",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Post"},
							},
						},
					},
				},
			},
		},
		"/posts/{id}": map[string]any{
			"get": map[string]any{
				"operationId": "getPost",
				"summary":     "Get a post by ID",
				"tags":        []any{"posts"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A single post",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Post"},
							},
						},
					},
					"404": map[string]any{"description": "Post not found"},
				},
			},
			"put": map[string]any{
				"operationId": "updatePost",
				"summary":     "Replace a post",
				"tags":        []any{"posts"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{"$ref": "#/components/schemas/NewPost"},
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Updated post",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Post"},
							},
						},
					},
				},
			},
			"patch": map[string]any{
				"operationId": "patchPost",
				"summary":     "Partially update a post",
				"tags":        []any{"posts"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{"$ref": "#/components/schemas/PostPatch"},
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Patched post",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Post"},
							},
						},
					},
				},
			},
			"delete": map[string]any{
				"operationId": "deletePost",
				"summary":     "Delete a post",
				"tags":        []any{"posts"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Post deleted"},
				},
			},
		},
		"/posts/{id}/comments": map[string]any{
			"get": map[string]any{
				"operationId": "getPostComments",
				"summary":     "Get comments for a post",
				"tags":        []any{"posts", "comments"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A list of comments",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"type":  "array",
									"items": map[string]any{"$ref": "#/components/schemas/Comment"},
								},
							},
						},
					},
				},
			},
		},
		// --- Comments ---
		"/comments": map[string]any{
			"get": map[string]any{
				"operationId": "listComments",
				"summary":     "List all comments",
				"tags":        []any{"comments"},
				"parameters": []any{
					map[string]any{
						"name": "postId", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A list of comments",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"type":  "array",
									"items": map[string]any{"$ref": "#/components/schemas/Comment"},
								},
							},
						},
					},
				},
			},
		},
		// --- Users ---
		"/users": map[string]any{
			"get": map[string]any{
				"operationId": "listUsers",
				"summary":     "List all users",
				"tags":        []any{"users"},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A list of users",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"type":  "array",
									"items": map[string]any{"$ref": "#/components/schemas/User"},
								},
							},
						},
					},
				},
			},
		},
		"/users/{id}": map[string]any{
			"get": map[string]any{
				"operationId": "getUser",
				"summary":     "Get a user by ID",
				"tags":        []any{"users"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A single user",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/User"},
							},
						},
					},
					"404": map[string]any{"description": "User not found"},
				},
			},
		},
		// --- Todos ---
		"/todos": map[string]any{
			"get": map[string]any{
				"operationId": "listTodos",
				"summary":     "List all todos",
				"tags":        []any{"todos"},
				"parameters": []any{
					map[string]any{
						"name": "userId", "in": "query", "required": false,
						"schema": map[string]any{"type": "integer"},
					},
					map[string]any{
						"name": "completed", "in": "query", "required": false,
						"schema": map[string]any{"type": "boolean"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A list of todos",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"type":  "array",
									"items": map[string]any{"$ref": "#/components/schemas/Todo"},
								},
							},
						},
					},
				},
			},
		},
		"/todos/{id}": map[string]any{
			"get": map[string]any{
				"operationId": "getTodo",
				"summary":     "Get a todo by ID",
				"tags":        []any{"todos"},
				"parameters": []any{
					map[string]any{
						"name": "id", "in": "path", "required": true,
						"schema": map[string]any{"type": "integer"},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "A single todo",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/Todo"},
							},
						},
					},
					"404": map[string]any{"description": "Todo not found"},
				},
			},
		},
	},
	"components": map[string]any{
		"schemas": map[string]any{
			"Post": map[string]any{
				"type":     "object",
				"required": []any{"id", "userId", "title", "body"},
				"properties": map[string]any{
					"id":     map[string]any{"type": "integer"},
					"userId": map[string]any{"type": "integer"},
					"title":  map[string]any{"type": "string"},
					"body":   map[string]any{"type": "string"},
				},
			},
			"NewPost": map[string]any{
				"type":     "object",
				"required": []any{"userId", "title", "body"},
				"properties": map[string]any{
					"userId": map[string]any{"type": "integer"},
					"title":  map[string]any{"type": "string"},
					"body":   map[string]any{"type": "string"},
				},
			},
			"PostPatch": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{"type": "string"},
					"body":  map[string]any{"type": "string"},
				},
			},
			"Comment": map[string]any{
				"type":     "object",
				"required": []any{"id", "postId", "name", "email", "body"},
				"properties": map[string]any{
					"id":     map[string]any{"type": "integer"},
					"postId": map[string]any{"type": "integer"},
					"name":   map[string]any{"type": "string"},
					"email":  map[string]any{"type": "string", "format": "email"},
					"body":   map[string]any{"type": "string"},
				},
			},
			"User": map[string]any{
				"type":     "object",
				"required": []any{"id", "name", "username", "email"},
				"properties": map[string]any{
					"id":       map[string]any{"type": "integer"},
					"name":     map[string]any{"type": "string"},
					"username": map[string]any{"type": "string"},
					"email":    map[string]any{"type": "string", "format": "email"},
					"phone":    map[string]any{"type": "string"},
					"website":  map[string]any{"type": "string"},
					"address":  map[string]any{"$ref": "#/components/schemas/Address"},
					"company":  map[string]any{"$ref": "#/components/schemas/Company"},
				},
			},
			"Address": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"street":  map[string]any{"type": "string"},
					"suite":   map[string]any{"type": "string"},
					"city":    map[string]any{"type": "string"},
					"zipcode": map[string]any{"type": "string"},
					"geo":     map[string]any{"$ref": "#/components/schemas/Geo"},
				},
			},
			"Geo": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"lat": map[string]any{"type": "string"},
					"lng": map[string]any{"type": "string"},
				},
			},
			"Company": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name":        map[string]any{"type": "string"},
					"catchPhrase": map[string]any{"type": "string"},
					"bs":          map[string]any{"type": "string"},
				},
			},
			"Todo": map[string]any{
				"type":     "object",
				"required": []any{"id", "userId", "title", "completed"},
				"properties": map[string]any{
					"id":        map[string]any{"type": "integer"},
					"userId":    map[string]any{"type": "integer"},
					"title":     map[string]any{"type": "string"},
					"completed": map[string]any{"type": "boolean"},
				},
			},
		},
	},
}

package parser_test

import (
	"testing"

	"github.com/robin38n/reqviz/backend/internal/parser"
)

// minimal valid spec: 1 path with GET+POST (2 operations), 2 component schemas,
// 1 tag. Used by both the JSON and YAML happy-path tests.
const validJSONSpec = `{
  "openapi": "3.0.3",
  "info": { "title": "Pet API", "version": "2.1.0" },
  "tags": [ { "name": "pets" } ],
  "paths": {
    "/pets": {
      "get": { "responses": { "200": { "description": "ok" } } },
      "post": { "responses": { "201": { "description": "created" } } }
    }
  },
  "components": {
    "schemas": {
      "Pet": { "type": "object" },
      "Error": { "type": "object" }
    }
  }
}`

const validYAMLSpec = `
openapi: 3.0.3
info:
  title: Pet API
  version: 2.1.0
tags:
  - name: pets
paths:
  /pets:
    get:
      responses:
        "200": { description: ok }
    post:
      responses:
        "201": { description: created }
components:
  schemas:
    Pet:
      type: object
    Error:
      type: object
`

func assertSummary(t *testing.T, r *parser.ParseResult) {
	t.Helper()
	if r.Title != "Pet API" {
		t.Errorf("Title = %q, want %q", r.Title, "Pet API")
	}
	if r.Version != "2.1.0" {
		t.Errorf("Version = %q, want %q", r.Version, "2.1.0")
	}
	if r.EndpointCount != 2 {
		t.Errorf("EndpointCount = %d, want %d", r.EndpointCount, 2)
	}
	if r.SchemaCount != 2 {
		t.Errorf("SchemaCount = %d, want %d", r.SchemaCount, 2)
	}
	if len(r.Tags) != 1 || r.Tags[0] != "pets" {
		t.Errorf("Tags = %v, want [pets]", r.Tags)
	}
}

func TestFromJSON_Valid(t *testing.T) {
	r, err := parser.FromJSON([]byte(validJSONSpec))
	if err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}
	assertSummary(t, r)
}

func TestFromYAML_Valid(t *testing.T) {
	r, err := parser.FromYAML([]byte(validYAMLSpec))
	if err != nil {
		t.Fatalf("FromYAML returned error: %v", err)
	}
	assertSummary(t, r)
}

func TestParse_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		parse func([]byte) (*parser.ParseResult, error)
	}{
		{"malformed JSON", "{", parser.FromJSON},
		{"malformed YAML", "key: : :\n  - [", parser.FromYAML},
		{"JSON wrong type for info", `{"openapi":"3.0.3","info":"not-an-object"}`, parser.FromJSON},
		{"YAML wrong type for paths", "openapi: 3.0.3\ninfo:\n  title: x\n  version: \"1\"\npaths: not-an-object\n", parser.FromYAML},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.parse([]byte(tt.input)); err == nil {
				t.Errorf("parse(%q) = nil error, want error", tt.input)
			}
		})
	}
}

func TestFromJSON_CountsAllOperations(t *testing.T) {
	// /a: get, post, put, patch, delete (5); /b: head, options (2) => 7
	spec := `{
      "openapi": "3.0.3",
      "info": { "title": "t", "version": "1" },
      "paths": {
        "/a": {
          "get": { "responses": { "200": { "description": "ok" } } },
          "post": { "responses": { "200": { "description": "ok" } } },
          "put": { "responses": { "200": { "description": "ok" } } },
          "patch": { "responses": { "200": { "description": "ok" } } },
          "delete": { "responses": { "200": { "description": "ok" } } }
        },
        "/b": {
          "head": { "responses": { "200": { "description": "ok" } } },
          "options": { "responses": { "200": { "description": "ok" } } }
        }
      }
    }`
	r, err := parser.FromJSON([]byte(spec))
	if err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}
	if r.EndpointCount != 7 {
		t.Errorf("EndpointCount = %d, want %d", r.EndpointCount, 7)
	}
}

func TestFromJSON_DefaultsTitleAndVersion(t *testing.T) {
	// info present but empty title/version -> parser falls back to defaults.
	spec := `{ "openapi": "3.0.3", "info": { "title": "", "version": "" }, "paths": {} }`
	r, err := parser.FromJSON([]byte(spec))
	if err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}
	if r.Title != "Untitled" {
		t.Errorf("Title = %q, want %q", r.Title, "Untitled")
	}
	if r.Version != "unknown" {
		t.Errorf("Version = %q, want %q", r.Version, "unknown")
	}
}

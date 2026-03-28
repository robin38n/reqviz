package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/robin38n/restatlas/backend/internal/handler/demos"
	"github.com/robin38n/restatlas/backend/internal/store"
)

type demoEntry struct {
	Slug        string
	Title       string
	Description string
	Spec        map[string]any
}

var demoSpecs = []demoEntry{
	{
		Slug:        "jsonplaceholder",
		Title:       "JSONPlaceholder",
		Description: "Fake REST API for testing — posts, comments, users, todos. Supports all HTTP methods.",
		Spec:        demos.JSONPlaceholderSpec,
	},
	{
		Slug:        "pokeapi",
		Title:       "PokéAPI",
		Description: "Pokémon data API — browse pokemon, types, abilities. Read-only.",
		Spec:        demos.PokeAPISpec,
	},
	{
		Slug:        "dogceo",
		Title:       "Dog CEO",
		Description: "Random dog images by breed. Simple and fun.",
		Spec:        demos.DogCEOSpec,
	},
}

func demoBySlug(slug string) *demoEntry {
	for i := range demoSpecs {
		if demoSpecs[i].Slug == slug {
			return &demoSpecs[i]
		}
	}
	return nil
}

// HandleDemoList returns the list of available demo specs.
func (s *Server) HandleDemoList(w http.ResponseWriter, r *http.Request) {
	type demoInfo struct {
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	list := make([]demoInfo, len(demoSpecs))
	for i, d := range demoSpecs {
		list[i] = demoInfo{Slug: d.Slug, Title: d.Title, Description: d.Description}
	}
	writeJSON(w, http.StatusOK, list)
}

// HandleDemoSpec returns the raw spec JSON for a demo by slug.
func (s *Server) HandleDemoSpec(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	entry := demoBySlug(slug)
	if entry == nil {
		writeError(w, http.StatusNotFound, "unknown demo: "+slug)
		return
	}
	writeJSON(w, http.StatusOK, entry.Spec)
}

// HandleDemoUpload loads a demo spec by slug through the normal pipeline and returns a SpecSummary.
func (s *Server) HandleDemoUpload(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	entry := demoBySlug(slug)
	if entry == nil {
		writeError(w, http.StatusNotFound, "unknown demo: "+slug)
		return
	}

	rawBytes, _ := json.Marshal(entry.Spec)

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(rawBytes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "demo spec failed validation: "+err.Error())
		return
	}

	title := doc.Info.Title
	version := doc.Info.Version

	endpointCount := 0
	if doc.Paths != nil {
		for _, pathItem := range doc.Paths.Map() {
			if pathItem.Get != nil {
				endpointCount++
			}
			if pathItem.Post != nil {
				endpointCount++
			}
			if pathItem.Put != nil {
				endpointCount++
			}
			if pathItem.Patch != nil {
				endpointCount++
			}
			if pathItem.Delete != nil {
				endpointCount++
			}
			if pathItem.Head != nil {
				endpointCount++
			}
			if pathItem.Options != nil {
				endpointCount++
			}
		}
	}

	schemaCount := 0
	if doc.Components != nil {
		schemaCount = len(doc.Components.Schemas)
	}

	var tags []string
	for _, tag := range doc.Tags {
		tags = append(tags, tag.Name)
	}

	stored := &store.StoredSpec{
		Title:         title,
		Version:       version,
		EndpointCount: endpointCount,
		SchemaCount:   schemaCount,
		Tags:          tags,
		Raw:           entry.Spec,
	}
	id := s.store.Save(stored)

	now := time.Now()
	writeJSON(w, http.StatusCreated, SpecSummary{
		Id:            openapi_types.UUID(id),
		Title:         title,
		Version:       version,
		EndpointCount: endpointCount,
		SchemaCount:   schemaCount,
		Tags:          &tags,
		CreatedAt:     &now,
	})
}

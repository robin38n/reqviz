package parser

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ParseResult holds the validated metadata and raw content extracted from an
// OpenAPI specification.
type ParseResult struct {
	Title         string
	Version       string
	EndpointCount int
	SchemaCount   int
	Tags          []string
	Raw           map[string]any
}

// extract validates a decoded spec map against the OpenAPI 3.x schema and
// extracts summary metadata. It is the shared core called by both parsers.
func extract(data []byte, raw map[string]any) (*ParseResult, error) {
	doc, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAPI spec: %w", err)
	}

	title := "Untitled"
	version := "unknown"
	if doc.Info != nil {
		if doc.Info.Title != "" {
			title = doc.Info.Title
		}
		if doc.Info.Version != "" {
			version = doc.Info.Version
		}
	}

	var tags []string
	for _, tag := range doc.Tags {
		tags = append(tags, tag.Name)
	}

	schemaCount := 0
	if doc.Components != nil {
		schemaCount = len(doc.Components.Schemas)
	}

	return &ParseResult{
		Title:         title,
		Version:       version,
		EndpointCount: countEndpoints(doc),
		SchemaCount:   schemaCount,
		Tags:          tags,
		Raw:           raw,
	}, nil
}

func countEndpoints(doc *openapi3.T) int {
	if doc.Paths == nil {
		return 0
	}
	count := 0
	for _, item := range doc.Paths.Map() {
		for _, op := range []*openapi3.Operation{
			item.Get, item.Post, item.Put, item.Patch,
			item.Delete, item.Head, item.Options,
		} {
			if op != nil {
				count++
			}
		}
	}
	return count
}

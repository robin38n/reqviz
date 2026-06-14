package parser

import (
	"encoding/json"
	"fmt"
)

func FromJSON(data []byte) (*ParseResult, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return extract(data, raw)
}

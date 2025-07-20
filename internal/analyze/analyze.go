package analyze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Entry represents a single JSON input document.
type Entry struct {
	Input struct {
		Payload string `json:"payload"`
	} `json:"input"`

	Scan map[string]interface{} `json:"scan"`
}

// ResultGroup holds a merged scan result with all payloads that produced it.
type ResultGroup struct {
	Scan     map[string]interface{} `json:"scan"`
	Payloads []string               `json:"payloads"`
}

// GroupByScan takes a slice of raw JSON objects and groups entries by identical scan results.
// It returns a slice of ResultGroup with scan results and associated payloads.
func GroupByScan(rawEntries [][]byte) ([]ResultGroup, error) {
	groups := make(map[string]*ResultGroup)

	for _, raw := range rawEntries {
		var entry Entry
		if err := json.Unmarshal(raw, &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entry: %w", err)
		}

		scanKey, err := canonicalJSON(entry.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize scan: %w", err)
		}

		if group, exists := groups[scanKey]; exists {
			group.Payloads = append(group.Payloads, entry.Input.Payload)
		} else {
			groups[scanKey] = &ResultGroup{
				Scan:     entry.Scan,
				Payloads: []string{entry.Input.Payload},
			}
		}
	}

	// Finalize grouped results
	var result []ResultGroup
	for _, group := range groups {
		sort.Strings(group.Payloads)
		result = append(result, *group)
	}

	return result, nil
}

// canonicalJSON returns a stable string representation of a JSON object (used for grouping keys).
func canonicalJSON(v interface{}) (string, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	return buf.String(), nil
}

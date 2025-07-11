package httpnode

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONNode holds counts of JSON keys, values, and paths in a document.
type JSONNode struct {
	KeyCount   map[string]int
	ValueCount map[string]int
	PathCount  map[string]int
}

// JSONNodeMerge holds slices of counts for merging multiple JSONNode results.
type JSONNodeMerge struct {
	KeyCount   map[string][]int `json:"KeyCount"`
	ValueCount map[string][]int `json:"ValueCount"`
	PathCount  map[string][]int `json:"PathCount"`
}

// NewJSONNode creates an initialized JSONNode.
func NewJSONNode() JSONNode {
	return JSONNode{
		KeyCount:   make(map[string]int),
		ValueCount: make(map[string]int),
		PathCount:  make(map[string]int),
	}
}

// NewMergeJSONNode creates an initialized JSONNodeMerge.
func NewMergeJSONNode() JSONNodeMerge {
	return JSONNodeMerge{
		KeyCount:   make(map[string][]int),
		ValueCount: make(map[string][]int),
		PathCount:  make(map[string][]int),
	}
}

// GetJSONNode parses the JSON body and returns a JSONNode with counts.
func GetJSONNode(body string) (JSONNode, error) {
	node := NewJSONNode()
	var data interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return node, fmt.Errorf("invalid JSON: %w", err)
	}
	traverseJSON(data, "", &node)
	return node, nil
}

// traverseJSON recursively walks the JSON structure, recording keys, values, and paths.
func traverseJSON(val interface{}, path string, node *JSONNode) {
	switch v := val.(type) {
	case map[string]interface{}:
		for key, child := range v {
			// record key
			node.KeyCount[key]++
			// record path
			p := joinPath(path, key)
			node.PathCount[p]++
			// recurse
			traverseJSON(child, p, node)
		}
	case []interface{}:
		for idx, child := range v {
			// record path index
			p := fmt.Sprintf("%s[%d]", path, idx)
			node.PathCount[p]++
			traverseJSON(child, p, node)
		}
	case string:
		node.ValueCount[v]++
		p := joinPath(path, v)
		node.PathCount[p]++
	case float64:
		s := fmt.Sprint(v)
		node.ValueCount[s]++
		p := joinPath(path, s)
		node.PathCount[p]++
	case bool:
		s := fmt.Sprint(v)
		node.ValueCount[s]++
		p := joinPath(path, s)
		node.PathCount[p]++
	case nil:
		node.ValueCount["null"]++
		p := joinPath(path, "null")
		node.PathCount[p]++
	}
}

// joinPath concatenates parent path and key/value into a dotted JSON path.
func joinPath(parent, elem string) string {
	if parent == "" {
		return elem
	}
	return strings.Join([]string{parent, elem}, ".")
}

// Merge merges a JSONNode into the JSONNodeMerge, avoiding duplicate counts.
func (m *JSONNodeMerge) Merge(node JSONNode) {
	for key, cnt := range node.KeyCount {
		if !contains(m.KeyCount[key], cnt) {
			m.KeyCount[key] = append(m.KeyCount[key], cnt)
		}
	}
	for val, cnt := range node.ValueCount {
		if !contains(m.ValueCount[val], cnt) {
			m.ValueCount[val] = append(m.ValueCount[val], cnt)
		}
	}
	for path, cnt := range node.PathCount {
		if !contains(m.PathCount[path], cnt) {
			m.PathCount[path] = append(m.PathCount[path], cnt)
		}
	}
}

// contains checks if slice s contains v.
func contains(s []int, v int) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

// JSONNodeToJSON serializes a JSONNode to JSON bytes.
func JSONNodeToJSON(node JSONNode) ([]byte, error) {
	return json.Marshal(node)
}

// JSONMergeToJSON serializes a JSONNodeMerge to JSON bytes.
func JSONMergeToJSON(m JSONNodeMerge) ([]byte, error) {
	return json.Marshal(m)
}

package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Brum3ns/firefly/pkg/httpnode"
)

func Test_httpnode_json(t *testing.T) {
	// Example JSON documents
	json1 := `{
		"user": {"id": 1, "name": "Bob"},
		"user2": {"id": 1, "name": "Bob"}
	}`

	json2 := `{
		"user": {"id": 1, "name": "Bob"}
	}`

	// Parse each JSON into a JSONNode
	node1, err := httpnode.GetJSONNode(json1)
	if err != nil {
		log.Fatalf("failed to parse json1: %v", err)
	}
	node2, err := httpnode.GetJSONNode(json2)
	if err != nil {
		log.Fatalf("failed to parse json2: %v", err)
	}

	// Print per-document counts
	fmt.Println("--- JSON Node 1 ---")
	fmt.Printf("Keys: %+v\n", node1.KeyCount)
	fmt.Printf("Values: %+v\n", node1.ValueCount)
	fmt.Printf("Paths: %+v\n", node1.PathCount)

	fmt.Println("--- JSON Node 2 ---")
	fmt.Printf("Keys: %+v\n", node2.KeyCount)
	fmt.Printf("Values: %+v\n", node2.ValueCount)
	fmt.Printf("Paths: %+v\n", node2.PathCount)

	// Merge results
	merge := httpnode.NewMergeJSONNode()
	merge.Merge(node1)
	merge.Merge(node2)

	// Serialize merged result to JSON
	mergedJSON, err := httpnode.JSONMergeToJSON(merge)
	if err != nil {
		log.Fatalf("failed to marshal merged node: %v", err)
	}

	fmt.Println("--- Merged JSONNode ---")
	fmt.Println(string(mergedJSON))
}

package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Brum3ns/firefly/pkg/httpdiff"
	"github.com/Brum3ns/firefly/pkg/httpnode"
)

func Test_httpdiff_json(t *testing.T) {
	// 1. Prepare multiple JSON docs to form the merged baseline
	docs := []string{
		`{"user":{"id":1,"name":"Alice"},"tags":["go","json"]}`,
		`{"user":{"id":2,"name":"Bob"},"tags":["json","http"]}`,
	}

	// 3. New JSON document to diff
	newDoc := `{"user":{"id":2,"name":"Alice"},"tags":["go","json"]}`

	// Initialize a new MergeJSONNode
	mergeNode := httpnode.NewMergeJSONNode()

	// Parse each doc and merge into baseline
	for _, doc := range docs {
		node, err := httpnode.GetJSONNode(doc)
		if err != nil {
			log.Fatalf("failed to parse JSON doc: %v", err)
		}
		mergeNode.Merge(node)
	}

	// (Optional) Serialize and view merged baseline
	baselineJSON, _ := httpnode.JSONMergeToJSON(mergeNode)
	fmt.Println("Baseline (merged JSONNodeMerge):")
	fmt.Println(string(baselineJSON))

	// 2. Configure HttpDiff to use our merged baseline
	diffConfig := httpdiff.Config{
		Merge: httpdiff.Merge{
			JSONMergeNode: mergeNode,
		},
	}
	hdiff := httpdiff.NewHttpDiff(diffConfig)

	newNode, err := httpnode.GetJSONNode(newDoc)
	if err != nil {
		log.Fatalf("failed to parse new JSON doc: %v", err)
	}

	// 4. Compute diff against baseline
	result := hdiff.GetJSONNodeDiff(newNode)

	// 5. Print diff results
	fmt.Printf("Diff OK: %v\n", result.OK)
	fmt.Println("--- Appear ---")
	fmt.Printf("KeyHits: %d, ValueHits: %d, PathHits: %d\n", result.Appear.KeyHits, result.Appear.ValueHits, result.Appear.PathHits)
	fmt.Printf("Appear details: %+v\n", result.Appear.JSONNode)

	fmt.Println("--- Disappear ---")
	fmt.Printf("KeyHits: %d, ValueHits: %d, PathHits: %d\n", result.Disappear.KeyHits, result.Disappear.ValueHits, result.Disappear.PathHits)
	fmt.Printf("Disappear details: %+v\n", result.Disappear.JSONNode)
}

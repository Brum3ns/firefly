package httpdiff

import "github.com/Brum3ns/firefly/pkg/httpnode"

// JSONNodeDiff holds counts of appearing/disappearing JSON metrics and the raw JSONNode data.
type JSONNodeDiff struct {
	KeyHits   int               `json:"keyhits"`
	ValueHits int               `json:"valuehits"`
	PathHits  int               `json:"pathhits"`
	JSONNode  httpnode.JSONNode `json:"jsonnode"`
}

// JSONResult encapsulates the diff results for JSONNodes.
type JSONResult struct {
	OK        bool         `json:"ok"`
	Malformed bool         `json:"malformed"`
	Appear    JSONNodeDiff `json:"appear"`
	Disappear JSONNodeDiff `json:"disappear"`
}

// GetJSONNodeDiff compares a single JSONNode against the merged baseline and returns diffs.
func (hdiff *HttpDiff) GetJSONNodeDiff(jsonNode httpnode.JSONNode) JSONResult {
	// Prepare current metrics
	current := [3]diffNode{
		{data: jsonNode.KeyCount},
		{data: jsonNode.ValueCount},
		{data: jsonNode.PathCount},
	}
	// Known merged baseline metrics
	known := [3]map[string][]int{
		hdiff.config.Merge.JSONMergeNode.KeyCount,
		hdiff.config.Merge.JSONMergeNode.ValueCount,
		hdiff.config.Merge.JSONMergeNode.PathCount,
	}

	// Compute only "appear" diffs: items in new JSON not seen in baseline
	appearDiffs := make([]diffNode, 3)
	totalHits := 0
	for i := 0; i < len(current); i++ {
		appear, _ := hdiff.nodeDiff(current[i], known[i], hdiff.config.Payload)
		appearDiffs[i] = appear
		totalHits += appear.hit
	}

	// Build result: only Appear, Disappear always empty
	return JSONResult{
		OK: totalHits > 0,
		Appear: JSONNodeDiff{
			KeyHits:   appearDiffs[0].hit,
			ValueHits: appearDiffs[1].hit,
			PathHits:  appearDiffs[2].hit,
			JSONNode: httpnode.JSONNode{
				KeyCount:   appearDiffs[0].data,
				ValueCount: appearDiffs[1].data,
				PathCount:  appearDiffs[2].data,
			},
		},
		Disappear: JSONNodeDiff{
			KeyHits:   0,
			ValueHits: 0,
			PathHits:  0,
			JSONNode: httpnode.JSONNode{
				KeyCount:   make(map[string]int),
				ValueCount: make(map[string]int),
				PathCount:  make(map[string]int),
			},
		},
	}
}

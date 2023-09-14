package verify

import (
	"github.com/Brum3ns/firefly/pkg/output"
	"github.com/Brum3ns/firefly/pkg/prepare"
)

// Store the verified behavior from the target.
// Note : (This struct is best to use to store a single target behavior related to the response and scan results. It's recommended to make a map that holds this struct as a list for each target behavior result)
type TargetKnowledge struct {
	Payload  string
	HTMLNode prepare.HTMLNode
	Response output.Response
	Request  output.Request
	Scanner  output.Scanner
}

func NewVerifiedStorage() *TargetKnowledge {
	return &TargetKnowledge{}
}

// Check if the map already contains verified storage, if not then make a new map and return it.
func Prepare(m map[string][]TargetKnowledge) map[string][]TargetKnowledge {
	if m != nil {
		return m
	}
	return make(map[string][]TargetKnowledge)
}

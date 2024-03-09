package scan

import (
	"github.com/Brum3ns/firefly/internal/knowledge"
	"github.com/Brum3ns/firefly/pkg/request"
)

type behavior struct {
	status bool
}

func NewBehavior() *behavior {
	return &behavior{}
}

// Quick detection for unkown behavior
func (b *behavior) QuickDetect(http request.Response, known knowledge.Knowledge) bool {
	count := 0
	countMax := len(known.Responses) * 3 //Note : (Number represents the number of if statements in the loop)
	for _, resp := range known.Responses {
		if resp.StatusCode == http.StatusCode {
			count++
		}
		if resp.Title == http.Title {
			count++
		}
		if resp.ContentType == http.ContentType {
			count++
		}
	}
	//If no test did hit and the count is the same length as the list of known data. No unexpected behavior was discovered:
	return count != countMax
}

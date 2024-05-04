package scan

type behavior struct {
	status bool
}

func NewBehavior() *behavior {
	return &behavior{}
}

// Quick detection for unkown behavior
func (b *behavior) QuickDetect(job Job) bool {
	count := 0
	countMax := len(job.Knowledge.Responses) * 3 //Note : (Number represents the number of if statements in the loop)
	for _, resp := range job.Knowledge.Responses {
		if resp.StatusCode == job.Http.Response.StatusCode {
			count++
		}
		if resp.Title == job.Http.Response.Title {
			count++
		}
		if resp.ContentType == job.Http.Response.ContentType {
			count++
		}
	}
	//If no test did hit and the count is the same length as the list of known data. No unexpected behavior was discovered:
	return count != countMax
}

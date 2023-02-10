package runner

import (
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/functions"
	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
	"github.com/Brum3ns/firefly/pkg/storage"
)

/** Define jobs to a channel. It also insert the payload to all sources ( Adjusted for the fuzz/verify process )*/
func Jobs(r *Runner, c chan<- storage.Target) {
	for _, url := range r.options.Target["urls"] {
		for _, method := range r.options.Target["methods"] {
			for _, attackMethod := range r.wordlist.UseTechniques {

				//Select the wordlist(s) to use for the process:
				tag, wordlist := functions.WordlistToUse(attackMethod, r.wordlist)

				for id, payload := range wordlist {
					//Payload suffix & prefix setup: (if any)
					payload = (r.options.PayloadPrefix + (payload) + r.options.PayloadSuffix)

					//Display payloads and exit
					if r.options.ShowPayload {
						design.DisplayPayload(id, len(wordlist), payload)

						//Start adding jobs to the channels:
					} else {
						//Add all finish payload(s) to all finished URL(s) and it's request raw data:
						payload = functions.PayloadPattern(payload)

						c <- storage.Target{
							ID: id,

							TAG:     tag,
							URL:     url,
							METHOD:  method,
							RAW:     []string{}, //DELETE?? - Unsure if needed headers are handled in request.go
							PAYLOAD: payload,
						}
					}
				}
			}
		}
	}
	//}
}

func JobsVerification(r *Runner, c chan<- storage.Target) {
	var (
		vPayload = (G.PayloadPattern + (r.options.VerifyPayload) + G.PayloadPattern)
		id       = -1 //-> 0 'id++' used first
	)

	for _, url := range r.options.Target["urls"] {
		for _, method := range r.options.Target["vmethods"] {

			//Special char modification/filter detection jobs:
			for i := 0; i < r.options.Verify; i++ {
				id++
				c <- storage.Target{
					ID: id,

					TAG:     "verifyPayload",
					URL:     url,
					METHOD:  method,
					PAYLOAD: vPayload,
				}
			}
			//Default response behavior jobs:
			for _, payload := range r.wordlist.VerifyPayloadChars {
				id++
				c <- storage.Target{
					ID: id,

					TAG:     "verifyChar",
					URL:     url,
					METHOD:  method,
					PAYLOAD: payload,
				}
			}
		}
	}
}

/**[Jobs][Filter] - Give job to channel runner.MatchData*/
func JobsFilter(m map[string][]string, c chan<- FData) {
	for k, l := range m {
		c <- FData{
			typ: k,
			lst: l,
		}
	}
}

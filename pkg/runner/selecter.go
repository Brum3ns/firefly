package runner

import (
	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
)

func (r *Runner) SelectTotal(v bool) (int, error) {
	var total int
	//Verify runner process
	if v {
		total = (len(r.options.Target["vurls"]) * len(r.options.Target["vmethods"]) * r.options.Verify) + len(r.wordlist.VerifyPayloadChars)

		//Fuzz runner process
	} else {
		total = (len(r.options.Target["urls"]) * len(r.options.Target["methods"])) * G.Amount_Item
	}

	return total, nil
}

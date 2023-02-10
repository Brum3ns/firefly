package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Brum3ns/firefly/pkg/design"
	fc "github.com/Brum3ns/firefly/pkg/functions"
	"github.com/Brum3ns/firefly/pkg/parse"
	"github.com/Brum3ns/firefly/pkg/runner"
	"github.com/Brum3ns/firefly/pkg/storage"
)

func main() {
	var (
		opt = parse.UserArguments()
		wl  = storage.ConfWordlist()
	)

	//Disclaimer text + Configure all data before starting the process:
	design.Disclaimer()
	design.InfoBanner(opt.ShowConfig)

	//Configuration setup from given user input:
	state, msg := parse.Configure(opt, wl)
	if !state {
		fc.IFFail(msg)
	}
	fmt.Println(design.OK, msg)

	//Timer to track the process time:
	timer := time.Now()

	//Check for "CTRL +C" to stop the process and it's threads:
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n\r " + design.Warning + "CTRL+C pressed - ")
			os.Exit(1)
		}
	}()

	//Start Runners: Verification / Fuzz
	VerifyResp, err := runner.New(opt, wl, true)
	if err != nil {
		fc.IFFail("vrunner")
	}

	Runner, err := runner.New(opt, wl, false)
	if err != nil {
		fc.IFFail("runner")
	}

	//[TODO] Add features (ATM just to ignore Golang variable "unused" error)
	if VerifyResp == nil || Runner == nil {
	}
	/*====================================================================*/

	fmt.Printf(design.OK+" Process finished in [%v], Success "+design.Green+"%v"+design.White+"/\033[2;32m%v"+design.White+"]\n", time.Since(timer), storage.Count_valid, storage.Count)

	//Graph how FireFly adjusted itself under the process and what the target are likely to run in the clientside/backend:
	/*//[GRAPH PRINTING] - [TODO make passive recon analyze and testing]
	fmt.Println(">>", vRes.LstG_Tag)
	adjust.Adjuster(vRes)*/
}

package main

import (
	"fmt"
	"time"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/fail"
	"github.com/Brum3ns/firefly/pkg/firefly/config"
	"github.com/Brum3ns/firefly/pkg/firefly/keypress"
	"github.com/Brum3ns/firefly/pkg/firefly/precheck"
	"github.com/Brum3ns/firefly/pkg/option"
	"github.com/Brum3ns/firefly/pkg/runner"
)

func main() {

	//Check resources before starting (first time use):
	precheck.Setup()

	//Configure needed resources
	opt := option.NewOptions()
	conf := config.NewConfigure(opt)

	design.Disclaimer()

	//Listen for user keypress input:
	keypress.CTRL_C()

	timer := time.Now()
	//Run the runner in verifyication process mode to detect normal behavior and patterns within the target:
	KnowledgeStorage, err := runner.Run(conf, nil)
	if err != nil {
		fail.IFFail(1009)
	}
	//Run the black-box enumiration process:
	runner.Run(conf, KnowledgeStorage)

	//Display summary of the process:
	fmt.Printf(design.STATUS.OK+" Process finished in [%v], Success "+design.COLOR.GREEN+"%v"+design.COLOR.WHITE+"/\033[2;32m%v"+design.COLOR.WHITE+"]\n", time.Since(timer), 1337, 1337)

}

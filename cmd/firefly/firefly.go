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
	KnowledgeStorage, _, err := runner.Run(conf, nil)
	if err != nil {
		fail.IFFail(1009)
	}
	//Run the black-box enumiration process:
	_, Statistic, _ := runner.Run(conf, KnowledgeStorage)

	//Display summary of the process:
	fmt.Printf(
		":: Process finished: Requests/Responses:[%d/%d], Completed:[\033[1;32m%d\033[0m], Filtered:[\033[1;36m%d\033[0m], Failed:[\033[31m%d\033[0m], Time:[%v]\n",
		Statistic.Request.Count,
		Statistic.Response,
		Statistic.Completed,
		Statistic.Request.Filtered,
		Statistic.Failed,
		time.Since(timer),
	)
}

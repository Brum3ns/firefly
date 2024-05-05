package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Brum3ns/firefly/internal/banner"
	"github.com/Brum3ns/firefly/internal/config"
	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/internal/option"
	"github.com/Brum3ns/firefly/internal/runner"
	"github.com/Brum3ns/firefly/internal/setup"
	"github.com/Brum3ns/firefly/pkg/design"
)

func main() {
	//Check resources before starting (first time use):
	if _, err := setup.Setup(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//Configure needed resources
	opt := option.NewOptions()
	conf := config.NewConfigure(opt)

	if !conf.TerminalUI {
		banner.Banner()
		banner.Disclaimer()
	}

	//Listen for user keypress (CTRL + C):
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n\r"+design.STATUS.WARNING, "CTRL+C pressed - Exiting")
			os.Exit(130)
		}
	}()

	timer := time.Now()

	//Run the runner in verifyication process mode to detect normal behavior and patterns within the target:
	VerifyRunner := runner.NewRunner(conf, nil)
	KnowledgeStorage, _, err := VerifyRunner.Run()
	if err != nil {
		log.Fatal(err)
	}

	//Run the black-box enumiration process:
	AttackRunner := runner.NewRunner(conf, KnowledgeStorage)
	_, Statistic, err := AttackRunner.Run()
	if err != nil {
		log.Fatal(err)
	}

	//Display summary of the process:
	fmt.Printf(
		"%s\033[1;32m\u2713\033[0m Process finished: Requests/Responses:[%d/%d], Scanned:[\033[1;32m%d\033[0m], Behavior:[\033[1;33m%d\033[0m], Filtered:[\033[1;36m%d\033[0m], Error:[\033[31m%d\033[0m], Time:[%v]\n",
		global.TERMINAL_CLEAR,
		Statistic.Request.GetCount(),
		Statistic.Response.GetCount(),
		Statistic.Scanner.GetCount(),
		Statistic.Behavior.GetCount(),
		Statistic.Request.GetFilterCount(),
		Statistic.Request.GetErrorCount(),
		time.Since(timer),
	)
}

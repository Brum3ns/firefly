package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Brum3ns/firefly/internal/option"
	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/internal/runner"
	"github.com/Brum3ns/firefly/internal/setup"
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/faker"
)

func main() {
	//Check resources before starting (first time use):
	if _, err := setup.Setup(); err != nil {
		log.Fatalln(err)
	}

	//Setup user arguments (options)
	opt, err := option.NewOption()
	if err != nil {
		log.Fatalln(err)
	}

	// Check output options
	if output.FileExist(opt.Output.OutputFile) {
		if opt.Output.Overwrite {
			// Remove old output file
			output.RemoveFile(opt.Output.OutputFile)
		} else {
			log.Fatalln("the given output file already exists. To overwrite, use the flag: -oW")
		}
	}

	if opt.Placeholder.List {
		faker.PrintList()
		return
	}

	if opt.Placeholder.Test != "" {
		result, err := faker.Generate(opt.Placeholder.Test)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)
		return
	}

	// Configure user arguments (options)
	/* conf, err := config.NewConfigure(opt)
	if err != nil {
		log.Fatal(design.STATUS.ERROR, err)
	}

	if !conf.Option.TerminalUI {
		banner.Banner()
		banner.Disclaimer()
	} */

	//Listen for user keypress (CTRL + C):
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n\r"+design.STATUS.WARNING, "CTRL+C pressed - Exiting")
			os.Exit(130)
		}
	}()

	// Run the runner
	r, err := runner.NewRunner(opt)
	if err != nil {
		log.Fatalf("failed to create runner instance, error : %v", err)
	}

	if err := r.FetchKnowledge(); err != nil {
		log.Fatalln(err)
	}

	//if !r.HasKnowledge() {
	//	fmt.Println("No knowledge could be gather from the target")
	//	return
	//}

	stats, err := r.RunFuzz()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(stats)

	//Run the runner in verifyication process mode to detect normal behavior and patterns within the target:
	/* VerifyRunner := runner.NewRunner(conf, nil)
	KnowledgeStorage, _, err := VerifyRunner.Run()
	if err != nil {
		log.Fatal(err)
	}

	//Run the black-box enumiration process:
	r := runner.NewRunner(conf)
	_, Statistic, err := r.Run()
	if err != nil {
		log.Fatal(err)
	}*/

	//Display summary of the process:
	/* fmt.Printf(
		"%s\033[1;32m\u2713\033[0m Process finished: Requests/Responses:[%d/%d], Scanned:[\033[1;32m%d\033[0m], Behavior:[\033[1;33m%d\033[0m], Filtered:[\033[1;36m%d\033[0m], Error:[\033[31m%d\033[0m], Time:[%v]\n",
		global.TERMINAL_CLEAR,
		Statistic.Request.GetCount(),
		Statistic.Response.GetCount(),
		Statistic.Scanner.GetCount(),
		Statistic.Behavior.GetCount(),
		Statistic.Request.GetFilterCount(),
		Statistic.Request.GetErrorCount(),
		time.Since(timer),
	) */
}

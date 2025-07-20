package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Brum3ns/firefly/internal/option"
	"github.com/Brum3ns/firefly/internal/runner"
	"github.com/Brum3ns/firefly/internal/setup"
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/faker"
)

func main() {
	//Setup user arguments (options)
	opt, err := option.NewOption()
	if err != nil {
		log.Fatalln(err)
	}

	if err := setup.Setup(opt); err != nil {
		log.Fatalln(err)
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

	_, err = r.RunFuzz()
	if err != nil {
		log.Fatalln(err)
	}
}

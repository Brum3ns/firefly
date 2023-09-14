package keypress

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Brum3ns/firefly/pkg/design"
)

func CTRL_C() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n\r"+design.STATUS.WARNING, "CTRL+C pressed - Exiting")
			os.Exit(130)
		}
	}()
}

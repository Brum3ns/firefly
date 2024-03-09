package verbose

import (
	"log"

	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/pkg/design"
)

// Show verbose on screen (if verbose is enable by the user)
func Show(msg any) {
	if global.VERBOSE {
		log.Printf(global.TERMINAL_CLEAR, "%v", msg)
	}
}

// Output the error messages using the 'log' package (Function used for easy customization and better error output)
func Error(err error) {
	if err != nil {
		log.Println(global.TERMINAL_CLEAR, design.STATUS.FAIL, err)
	}
}

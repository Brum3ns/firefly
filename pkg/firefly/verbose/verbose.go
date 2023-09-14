package verbose

import (
	"log"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/firefly/global"
)

// Show verbose on screen (if verbose is enable by the user)
func Show(msg string) {
	if global.VERBOSE {
		log.Println(msg)
	}
}

// Output the error messages using the 'log' package (Function used for easy customization and better error output)
func Error(err error) {
	if err != nil {
		log.Println(design.STATUS.FAIL, err)
	}
}

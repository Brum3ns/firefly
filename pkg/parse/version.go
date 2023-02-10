package parse

import (
	"fmt"
	"os"
)

/**Show the version of FireFly and the latest*/
func VersionFireFly() {
	//[TODO] - Func to request github version @latest and get -> (v1.x) + show "now version" and newest.
	v := "1.0"
	fmt.Printf("Version %s \n", v)
	os.Exit(0)
}

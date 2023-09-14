package update

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/firefly/global"
)

// Update all resources in the ".config/firefly/db/*" and exit
func Resources() {
	cmd := exec.Command("svn", "checkout", "-r", "HEAD", "https://github.com/Brum3ns/firefly/trunk/db", global.DIR_DB)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(design.STATUS.FAIL, err.Error())
		fmt.Println(design.STATUS.WARNING, "Make sure that the tool \"svn\" is installed on your system.")
		os.Exit(1)
	}

	//Install verbose:
	fmt.Println(string(stdout))
	os.Exit(0)
}

package precheck

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/files"
	"github.com/Brum3ns/firefly/pkg/firefly/global"
)

// Struct only if new development in the future will happen (inc none used variables)
type resources struct {
	folder_DB     string
	folder_HOME   string
	folder_CONFIG string
}

/**Setup the .config folder to store data (Only run first time the tool is executed)*/
func Setup() (bool, error) {
	rsource := &resources{
		folder_DB:     global.DIR_DB,
		folder_HOME:   global.DIR_HOME,
		folder_CONFIG: global.DIR_CONFIG,
	}

	//If the "".config/firefly/db" do not exist, create it:
	if !files.ExistAny(rsource.folder_DB) {
		design.Banner()
		if rsource.svg_getGithubFiles() {
			fmt.Println(design.STATUS.SUCCESS, "Installation finish!")
		} else {
			fmt.Println(design.STATUS.FAIL, "The installation failed, make sure to update your distro and golang to the latest versions.")
		}
		os.Exit(0)
	}

	return true, nil
}

/**Execute the OS command: "svn checkout 'https://github.com/Brum3ns/firefly/trunk/db'"*/
func (rs *resources) svg_getGithubFiles() bool {
	if rs.svm_confirm() {
		fmt.Println(design.STATUS.INFO, "Installing resources...")

		//Execute the command to get the files
		cmd := exec.Command("svn", "checkout", "-r", "HEAD", "https://github.com/Brum3ns/firefly/trunk/db", string(rs.folder_DB))
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(design.STATUS.FAIL, err.Error())
			fmt.Println(design.STATUS.WARNING, "Make sure that the tool \"svn\" is installed on your system.")
			os.Exit(1)
		}

		//Install verbose:
		fmt.Println(string(stdout))

	} else {
		return false
	}

	return true
}

func (rs *resources) svm_confirm() bool {
	var confirm string
	fmt.Println(design.STATUS.INFO+"  The 'db' folder containing all the resources for Firefly to work will be stored at:"+design.COLOR.ORANGE, rs.folder_DB, design.COLOR.WHITE)
	fmt.Print(design.DEBUG.INPUT + " Write \"" + design.COLOR.ORANGE + "ok" + design.COLOR.WHITE + "\" to confirm: ")
	fmt.Scanln(&confirm)

	if confirm == "ok" {
		return true
	} else {
		rs.exit(0)
	}

	return false
}

func (rs *resources) exit(n int) {
	fmt.Println(design.STATUS.FAIL, "Setup aborted")
	os.Exit(n)
}

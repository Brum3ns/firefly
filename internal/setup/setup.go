package setup

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/Brum3ns/firefly/internal/banner"
	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/pkg/design"
)

// Struct only if new development in the future will happen (inc none used variables)
type resource struct {
	folder_DB     string
	folder_HOME   string
	folder_CONFIG string
}

// Setup the .config folder and install needed resources (Only run first time the tool is executed)
// If a new installation is being processed, exit when finished
func Setup() (bool, error) {
	rsource := &resource{
		folder_DB:     global.DIR_DB,
		folder_HOME:   global.DIR_HOME,
		folder_CONFIG: global.DIR_CONFIG,
	}

	// Create config folder and ignore error in case it already exists
	os.MkdirAll(rsource.folder_CONFIG, os.ModePerm)

	// Check if the db git repository is installed already, otherwise install it
	if _, err := os.Stat(rsource.folder_DB); err != nil && os.IsNotExist(err) {
		banner.Banner()
		if rsource.approve() {
			fmt.Println("Installing DB from the firefly-db Github repository")
			InstallDB()
			fmt.Printf("Firefly's database is installed and can be find in the folder: %s\n", rsource.folder_DB)
			os.Exit(0)

		} else {
			return false, errors.New("Firefly's needs to have its database (resources) installed to work properly")
		}
	}
	return true, nil
}

// Approve the firefly db download and installation
func (rs *resource) approve() bool {
	var confirm string
	fmt.Println(design.STATUS.INFO+"  The Github repository - \"https://github.com/Brum3ns/firefly-db\" contains all the resources that Firefly use, and is needed to be installed. It will be installed into the folder:"+design.COLOR.ORANGE, rs.folder_DB, design.COLOR.WHITE)
	fmt.Print(design.DEBUG.INPUT + " Write \"" + design.COLOR.ORANGE + "ok" + design.COLOR.WHITE + "\" to confirm: ")
	fmt.Scanln(&confirm)

	if confirm == "ok" {
		return true
	} else {
		fmt.Println(design.STATUS.FAIL + " Firefly needs its resources to run")
		rs.exit(0)
	}
	return false
}

// Exit the setup process
func (rs *resource) exit(n int) {
	fmt.Println(design.STATUS.FAIL, "Setup process aborted")
	os.Exit(n)
}

// Update all resources in the ".config/firefly/db/*" from the firefly-db Github repository
// Note : Repository - https://github.com/Brum3ns/firefly-db.git
func InstallDB() (string, error) {
	const gitURL = "https://github.com/Brum3ns/firefly-db.git"

	cmd := exec.Command("git", "clone", gitURL, global.DIR_CONFIG)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

func UpdateDB() (string, error) {
	cmd := exec.Command("git", "-C", global.DIR_CONFIG, "pull")

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

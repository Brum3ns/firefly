package files

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Check if the file(s) or folder(s) already exists.
// Return true if it exist and false if it do not.
func ExistAny(f ...string) bool {
	l := []string{}

	//Extract all the files/folders:
	for _, i := range f {
		if _, err := os.Stat(i); !os.IsNotExist(err) {
			l = append(l, i)
		}
	}
	return len(l) == len(f)
}

func FileExist(filename string) bool {
	//File dose not exist:
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	//File exists:
	return true
}

// Read all content in a directory. Return a list of: Files | Folders
func InDir(folder string) ([]string, []string) {
	var (
		l_files   []string
		l_folders []string
	)
	//Read the files from the directory:
	filesFolders, _ := ioutil.ReadDir(folder)

	//If there is atleast one *file* found start adding the names to a list:
	for _, f := range filesFolders {
		if !f.IsDir() {
			l_files = append(l_files, f.Name())

		} else if f.IsDir() {
			l_folders = append(l_folders, f.Name())
		}
	}
	return l_files, l_folders
}

// Create a folder
func CreateFolder(name string) {
	if err := os.Mkdir(name, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

// Take a file and convert it's data to a given map
func FileToMap(file string) map[string]int {
	m := make(map[string]int)

	//Open the file and check for errors:
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	//Read the file items and add it to map:
	for scanner.Scan() {
		item := scanner.Text()
		if item != "" {
			m[item] += 1
		}
	}
	return m
}

// Check file size
func FileSize(f string) (int, error) {
	fInfo, err := os.Stat(f)
	if err != nil {
		return 0, err
	}
	return int(fInfo.Size()), nil
}

func FileOrFolder(f string) (string, error) {
	finfo, err := os.Stat(f)
	switch {
	case err != nil:
		return "", err
	case finfo.IsDir():
		return "folder", nil
	default:
		return "file", nil
	}
}

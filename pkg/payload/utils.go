package payload

import (
	"bufio"
	"fmt"
	"os"
)

// Read file and return file content as a list
func readFile(file string) ([]string, error) {
	var lst []string
	f, err := os.Open(file)
	if err != nil {
		return lst, fmt.Errorf("could not open file, error : %v", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		item := scanner.Text()
		if len(item) > 0 {
			lst = append(lst, item)
		}
	}
	f.Close()
	return lst, nil
}

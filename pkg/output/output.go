package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var mutex sync.Mutex

var (
	prefix   = []byte("[\r\n")
	suffix   = []byte("\r\n]")
	filesize = 0
)

func WriteJSON(count int, f *os.File, result ResultFinal) error {
	if !result.OK {
		return nil
	}
	if count == 0 {
		f.Write(prefix)
		filesize += len(prefix)
	}

	dataJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	f.Write(dataJSON)
	if err != nil {
		return err
	}

	//Replace last binary characters and append the new JSON data to output file:
	f.Write(suffix)

	filesize += len(dataJSON)
	filesize += len(suffix)

	info, _ := f.Stat()
	fmt.Println("size:", info.Size(), "|", filesize)

	return nil
}

package output

import (
	"encoding/json"
	"os"

	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpdiff"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

type Output struct {
	//Id    string `json:"id"`
	Date  string `json:"date"`
	Input Input  `json:"input"`
	Scan  Scan   `json:"scan"`
	Http  Http   `json:"http"`
}

type Input struct {
	Payload string `json:"payload"`
}

type Scan struct {
	Httpdiff httpdiff.Result `json:"httpdiff"`
	Extract  extract.Result  `json:"extract"`
}

type Http struct {
	RawRequest string         `json:"raw_request"`
	Response   rhttp.Response `json:"response"`
}

func MakeOutputJSON(output Output) ([]byte, error) {
	return json.Marshal(output)
}

func AppendOutputToFile(outputJson []byte, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(append(outputJson, '\n')); err != nil {
		return err
	}
	return nil
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func RemoveFile(filename string) error {
	return os.Remove(filename)
}

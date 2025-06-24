package runner

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func makeLog(logFile string) (*os.File, error) {
	// Check if logfile do exists, otherwise create it
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		dirLog := filepath.Dir(logFile)
		if err := os.MkdirAll(dirLog, os.ModePerm); err != nil {
			return nil, fmt.Errorf("could not create log directory, directory:[%s], Error : %v", dirLog, err)
		}
		if _, err := os.Create(logFile); err != nil {
			return nil, fmt.Errorf("could not create log file, file:[%s], Error : %v", logFile, err)
		}
	}
	// Open log file and set it as the logging output
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return file, fmt.Errorf("failed to open log file, Error : %v", err)
	}
	log.SetOutput(file)
	return file, nil
}

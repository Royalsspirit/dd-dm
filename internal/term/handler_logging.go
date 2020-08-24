package term

import (
	"fmt"
	"log"
	"os"
)

var (
	logFile = "errors.log"
)

func setupLogfile() (*os.File, error) {
	// open the log file
	dir, _ := os.Getwd()
	logfile, err := os.OpenFile(dir+"/"+logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// log time, filename, and line number
	log.SetFlags(log.Ltime | log.Lshortfile)
	// log to file
	log.SetOutput(logfile)

	return logfile, nil
}

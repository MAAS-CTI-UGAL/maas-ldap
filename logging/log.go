package logging

import (
	"io"
	"log"
	"os"
)

// Configure sends standard log output to stderr and the configured log file.
func Configure(logFile string) (*os.File, error) {
	if logFile == "" {
		log.SetOutput(os.Stderr)
		return nil, nil
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return nil, err
	}

	log.SetOutput(io.MultiWriter(os.Stderr, file))
	return file, nil
}

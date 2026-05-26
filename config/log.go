package config

import (
	"io"
	"log"
	"os"
)

// LogSettings contains logging resources configured at startup.
type LogSettings struct {
	FilePath string
	File     *os.File
}

func loadLogSettings() (LogSettings, error) {
	logFilePath := envOrDefault(envLogPath, "")
	if logFilePath == "" {
		log.SetOutput(os.Stderr)
		return LogSettings{}, nil
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return LogSettings{}, err
	}

	log.SetOutput(io.MultiWriter(os.Stderr, file))
	return LogSettings{
		FilePath: logFilePath,
		File:     file,
	}, nil
}

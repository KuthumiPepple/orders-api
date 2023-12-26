package logger

import (
	"log"
	"os"
	"path/filepath"
)

func SetupLog() (*os.File, error) {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err = os.Mkdir(logsDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// set up log file
	logFilePath := filepath.Join(logsDir, "error.log")
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)

	return file, nil
}

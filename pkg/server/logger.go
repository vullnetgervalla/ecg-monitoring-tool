package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Loggers struct {
	General     *log.Logger
	Alert       *log.Logger
	generalFile *os.File
	alertFile   *os.File
}

func SetupLoggers(generalLogPath, alertLogPath string) (*Loggers, error) {
	logsDir := filepath.Dir(generalLogPath)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, err
	}

	generalFile, err := os.OpenFile(generalLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	alertFile, err := os.OpenFile(alertLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		generalFile.Close()
		return nil, err
	}

	generalMultiWriter := io.MultiWriter(os.Stdout, generalFile)
	generalLogger := log.New(generalMultiWriter, "", log.LstdFlags)

	alertLogger := log.New(alertFile, "", log.LstdFlags)

	fmt.Printf("General logs will be written to %s\n", generalLogPath)
	fmt.Printf("Alert logs will be written to %s\n", alertLogPath)

	return &Loggers{
		General:     generalLogger,
		Alert:       alertLogger,
		generalFile: generalFile,
		alertFile:   alertFile,
	}, nil
}

func (l *Loggers) Close() error {
	var errs []error

	if l.generalFile != nil {
		if err := l.generalFile.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing general log file: %w", err))
		}
		l.generalFile = nil
	}

	if l.alertFile != nil {
		if err := l.alertFile.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing alert log file: %w", err))
		}
		l.alertFile = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing log files: %v", errs)
	}

	return nil
}

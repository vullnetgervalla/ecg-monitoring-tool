package server_test

import (
	"os"
	"path/filepath"
	"testing"

	"arhm/ecg-monitoring/pkg/server"
)

func TestSetupLoggers(t *testing.T) {
	tempDir := t.TempDir()

	generalLogPath := filepath.Join(tempDir, "logs", "test_general.log")
	alertLogPath := filepath.Join(tempDir, "logs", "test_alert.log")

	loggers, err := server.SetupLoggers(generalLogPath, alertLogPath)
	if err != nil {
		t.Fatalf("Failed to setup loggers: %v", err)
	}
	defer loggers.Close()

	if loggers.General == nil {
		t.Error("General logger is nil")
	}

	if loggers.Alert == nil {
		t.Error("Alert logger is nil")
	}

	testMessage := "Test log message"
	loggers.General.Println(testMessage)
	loggers.Alert.Println(testMessage)

	if err := loggers.Close(); err != nil {
		t.Fatalf("Failed to close loggers: %v", err)
	}

	generalContent, err := os.ReadFile(generalLogPath)
	if err != nil {
		t.Fatalf("Failed to read general log file: %v", err)
	}

	if len(generalContent) == 0 {
		t.Error("General log file is empty")
	}

	alertContent, err := os.ReadFile(alertLogPath)
	if err != nil {
		t.Fatalf("Failed to read alert log file: %v", err)
	}

	if len(alertContent) == 0 {
		t.Error("Alert log file is empty")
	}
}

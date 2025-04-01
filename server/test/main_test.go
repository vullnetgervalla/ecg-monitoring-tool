package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"arhm/ecg-monitoring/pkg/server"
)

func TestServerHandlers(t *testing.T) {
	// Setup a test logger
	tempDir := t.TempDir()
	loggers, err := server.SetupLoggers(tempDir+"/test.log", tempDir+"/alert.log")
	if err != nil {
		t.Fatalf("Failed to setup test loggers: %v", err)
	}
	defer loggers.Close()

	// Create the handler that would be used in main()
	ecgHandler := server.NewECGHandler(loggers)

	// Create a test HTTP server
	mux := http.NewServeMux()
	mux.Handle("/ecg", ecgHandler)
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	// Test a regular HTTP request to WebSocket endpoint
	// It should fail since we need a WebSocket handshake
	resp, err := http.Get(testServer.URL + "/ecg")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should get a 400 Bad Request for non-WebSocket request
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for non-WebSocket request, got %d", resp.StatusCode)
	}
}

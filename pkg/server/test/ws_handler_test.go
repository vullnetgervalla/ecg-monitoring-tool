package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"arhm/ecg-monitoring/pkg/server"

	"github.com/gorilla/websocket"
)

func setupTestECGHandler(t *testing.T) (*server.ECGHandler, *httptest.Server, *server.Loggers) {
	tempDir := t.TempDir()
	loggers, err := server.SetupLoggers(tempDir+"/test.log", tempDir+"/alerts.log")
	if err != nil {
		t.Fatalf("Failed to setup test loggers: %v", err)
	}

	handler := server.NewECGHandler(loggers)

	testServer := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))

	return handler, testServer, loggers
}

func TestECGHandlerWebSocketConnection(t *testing.T) {
	_, testServer, loggers := setupTestECGHandler(t)
	defer testServer.Close()
	defer loggers.Close()

	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Could not open websocket connection: %v", err)
	}
	defer ws.Close()

	timeout := time.After(2 * time.Second)
	received := false

	for !received {
		select {
		case <-timeout:
			t.Fatal("Timed out waiting for ECG reading")
		default:
			_, message, err := ws.ReadMessage()
			if err != nil {
				t.Fatalf("Failed to read message: %v", err)
			}

			if len(message) > 0 {
				received = true

				if !strings.Contains(string(message), "heart_rate") {
					t.Errorf("Message doesn't contain heart_rate: %s", string(message))
				}

				if !strings.Contains(string(message), "rr_interval") {
					t.Errorf("Message doesn't contain rr_interval: %s", string(message))
				}
			}
		}
	}
}

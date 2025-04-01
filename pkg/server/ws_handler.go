package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"arhm/ecg-monitoring/pkg/ecg"
	"arhm/ecg-monitoring/pkg/simulation"

	"github.com/gorilla/websocket"
)

type ECGHandler struct {
	Loggers   *Loggers
	Upgrader  websocket.Upgrader
	Simulator *simulation.Controller
}

func NewECGHandler(loggers *Loggers) *ECGHandler {
	return &ECGHandler{
		Loggers: loggers,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Simulator: simulation.NewController(),
	}
}

func (h *ECGHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Loggers.General.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer c.Close()

	h.Loggers.General.Printf("New client connected from %s", c.RemoteAddr())

	ticker := h.Simulator.RunWithCallback(1*time.Second, func(reading ecg.ECGReading, condition simulation.Condition) {
		switch condition {
		case simulation.ConditionTachycardia:
			h.Loggers.General.Printf("Tachycardia - HR=%d, RR=%0.2f", reading.HeartRate, reading.RRInterval)
			alertMsg := fmt.Sprintf("ALERT: TACHYCARDIA detected - High heart rate (HR=%d BPM, RR=%0.2f s)",
				reading.HeartRate, reading.RRInterval)
			h.Loggers.Alert.Println(alertMsg)
			h.Loggers.General.Println(alertMsg)
		case simulation.ConditionBradycardia:
			h.Loggers.General.Printf("Bradycardia - HR=%d, RR=%0.2f", reading.HeartRate, reading.RRInterval)
			alertMsg := fmt.Sprintf("ALERT: BRADYCARDIA detected - Low heart rate (HR=%d BPM, RR=%0.2f s)",
				reading.HeartRate, reading.RRInterval)
			h.Loggers.Alert.Println(alertMsg)
			h.Loggers.General.Println(alertMsg)
		case simulation.ConditionArrhythmia:
			h.Loggers.General.Printf("Arrhythmia - HR=%d, RR=%0.2f", reading.HeartRate, reading.RRInterval)
			alertMsg := fmt.Sprintf("ALERT: ARRHYTHMIA detected - Irregular heartbeat (HR=%d BPM, RR=%0.2f s)",
				reading.HeartRate, reading.RRInterval)
			h.Loggers.Alert.Println(alertMsg)
			h.Loggers.General.Println(alertMsg)
		default:
			h.Loggers.General.Printf("Normal - HR=%d, RR=%0.2f", reading.HeartRate, reading.RRInterval)
		}

		data, err := json.Marshal(reading)
		if err != nil {
			h.Loggers.General.Printf("Marshal error: %v", err)
			return
		}

		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			h.Loggers.General.Printf("Write error: %v", err)
			return
		}

		h.Loggers.General.Printf("Sent reading: HR=%d, RR=%0.2f", reading.HeartRate, reading.RRInterval)
	})
	defer ticker.Stop()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			h.Loggers.General.Printf("Read error: %v", err)
			break
		}
	}
}

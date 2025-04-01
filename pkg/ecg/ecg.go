package ecg

import (
	"fmt"
	"time"
)

type ECGReading struct {
	Timestamp  time.Time `json:"timestamp"`
	HeartRate  int       `json:"heart_rate"`
	RRInterval float64   `json:"rr_interval"`
}

type HeartCondition struct {
	Type        string
	Description string
	Reading     ECGReading
	Severity    string // "normal", "warning", "critical"
}

const (
	MinNormalHeartRate = 60
	MaxNormalHeartRate = 100

	MinNormalRRInterval = 0.6 // 100 BPM
	MaxNormalRRInterval = 1.0 // 60 BPM

	ConditionNormal      = "NORMAL"
	ConditionTachycardia = "TACHYCARDIA"
	ConditionBradycardia = "BRADYCARDIA"
	ConditionArrhythmia  = "ARRHYTHMIA"
)

func AnalyzeReading(reading ECGReading) HeartCondition {
	condition := HeartCondition{
		Type:        ConditionNormal,
		Description: "Normal heart activity",
		Reading:     reading,
		Severity:    "normal",
	}

	if reading.HeartRate > MaxNormalHeartRate {
		condition = HeartCondition{
			Type:        ConditionTachycardia,
			Description: fmt.Sprintf("High heart rate: %d BPM", reading.HeartRate),
			Reading:     reading,
			Severity:    "warning",
		}

		if reading.HeartRate > 120 {
			condition.Severity = "critical"
		}

		return condition
	} else if reading.HeartRate < MinNormalHeartRate {
		condition = HeartCondition{
			Type:        ConditionBradycardia,
			Description: fmt.Sprintf("Low heart rate: %d BPM", reading.HeartRate),
			Reading:     reading,
			Severity:    "warning",
		}

		if reading.HeartRate < 45 {
			condition.Severity = "critical"
		}

		return condition
	}

	if reading.RRInterval < MinNormalRRInterval || reading.RRInterval > MaxNormalRRInterval {
		condition = HeartCondition{
			Type:        ConditionArrhythmia,
			Description: fmt.Sprintf("Irregular heartbeat: RR interval %0.2f s", reading.RRInterval),
			Reading:     reading,
			Severity:    "warning",
		}

		if reading.RRInterval > 1.5 || reading.RRInterval < 0.4 {
			condition.Severity = "critical"
		}
	}

	return condition
}

func FormatAlert(condition HeartCondition) string {
	timestamp := condition.Reading.Timestamp.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("ALERT: %s - %s at %s", condition.Type, condition.Description, timestamp)
}

package ecg_test

import (
	"testing"
	"time"

	"arhm/ecg-monitoring/pkg/ecg"
)

func TestAnalyzeReading(t *testing.T) {
	testCase := func(name string, heartRate int, rrInterval float64, expectedType string, expectedSeverity string) {
		t.Run(name, func(t *testing.T) {
			reading := ecg.ECGReading{
				Timestamp:  time.Now(),
				HeartRate:  heartRate,
				RRInterval: rrInterval,
			}

			condition := ecg.AnalyzeReading(reading)

			if condition.Type != expectedType {
				t.Errorf("Expected condition type %s, got %s", expectedType, condition.Type)
			}

			if condition.Severity != expectedSeverity {
				t.Errorf("Expected severity %s, got %s", expectedSeverity, condition.Severity)
			}
		})
	}


	// Normal heart rates
	testCase("Normal lower bound", ecg.MinNormalHeartRate, 1.0, ecg.ConditionNormal, "normal")
	testCase("Normal mid range", 80, 0.75, ecg.ConditionNormal, "normal")
	testCase("Normal upper bound", ecg.MaxNormalHeartRate, 0.6, ecg.ConditionNormal, "normal")

	// Tachycardia (high heart rate)
	testCase("Mild tachycardia", ecg.MaxNormalHeartRate+1, 0.59, ecg.ConditionTachycardia, "warning")
	testCase("Severe tachycardia", 130, 0.45, ecg.ConditionTachycardia, "critical")

	// Bradycardia (low heart rate)
	testCase("Mild bradycardia", ecg.MinNormalHeartRate-1, 1.01, ecg.ConditionBradycardia, "warning")
	testCase("Severe bradycardia", 40, 1.5, ecg.ConditionBradycardia, "critical")

	// Arrhythmia (irregular RR interval)
	testCase("Arrhythmia high interval", 70, ecg.MaxNormalRRInterval+0.2, ecg.ConditionArrhythmia, "warning")
	testCase("Arrhythmia low interval", 70, ecg.MinNormalRRInterval-0.2, ecg.ConditionArrhythmia, "warning")
}

func TestFormatAlert(t *testing.T) {
	reading := ecg.ECGReading{
		Timestamp:  time.Date(2025, 4, 1, 14, 30, 5, 0, time.UTC),
		HeartRate:  120,
		RRInterval: 0.5,
	}

	condition := ecg.HeartCondition{
		Type:        ecg.ConditionTachycardia,
		Description: "High heart rate: 120 BPM",
		Reading:     reading,
		Severity:    "warning",
	}

	alert := ecg.FormatAlert(condition)

	expected := "ALERT: TACHYCARDIA - High heart rate: 120 BPM at 2025-04-01 14:30:05"
	if alert != expected {
		t.Errorf("Expected alert format: %s, got: %s", expected, alert)
	}
}

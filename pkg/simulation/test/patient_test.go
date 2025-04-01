package simulation_test

import (
	"arhm/ecg-monitoring/pkg/ecg"
	"arhm/ecg-monitoring/pkg/simulation"
	"testing"
	"time"
)

func init() {
	time.Local = time.UTC
}

func TestNewDefaultPatient(t *testing.T) {
	patient := simulation.NewDefaultPatient()

	if patient.ID != "PATIENT" {
		t.Errorf("Expected ID 'PATIENT', got '%s'", patient.ID)
	}

	if patient.BaseHeartRate != 80 {
		t.Errorf("Expected BaseHeartRate 80, got %d", patient.BaseHeartRate)
	}

	if patient.SimulateTachycardia || patient.SimulateBradycardia || patient.SimulateArrhythmia {
		t.Errorf("Expected default patient to have no simulated conditions")
	}
}

func TestGenerateECGReading(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		patient := simulation.NewDefaultPatient()

		for i := 0; i < 10; i++ {
			reading := simulation.GenerateECGReading(patient)

			if reading.HeartRate < patient.BaseHeartRate-patient.Variability ||
				reading.HeartRate > patient.BaseHeartRate+patient.Variability {
				t.Errorf("Heart rate %d outside expected range %d±%d",
					reading.HeartRate, patient.BaseHeartRate, patient.Variability)
			}

			expectedRR := 60.0 / float64(reading.HeartRate)
			rrDiff := reading.RRInterval - expectedRR
			if rrDiff < -patient.RRVariability || rrDiff > patient.RRVariability {
				t.Errorf("RR interval %f outside expected range of base %f±%f",
					reading.RRInterval, expectedRR, patient.RRVariability)
			}
		}
	})

	t.Run("Tachycardia", func(t *testing.T) {
		patient := simulation.NewDefaultPatient()
		patient.SimulateTachycardia = true

		for i := 0; i < 10; i++ {
			reading := simulation.GenerateECGReading(patient)

			if reading.HeartRate <= ecg.MaxNormalHeartRate {
				t.Errorf("Expected tachycardia heart rate > %d, got %d",
					ecg.MaxNormalHeartRate, reading.HeartRate)
			}

			expectedRR := 60.0 / float64(reading.HeartRate)
			if reading.RRInterval != expectedRR {
				t.Errorf("Expected RR interval %f, got %f", expectedRR, reading.RRInterval)
			}
		}
	})

	t.Run("Bradycardia", func(t *testing.T) {
		patient := simulation.NewDefaultPatient()
		patient.SimulateBradycardia = true

		for i := 0; i < 10; i++ {
			reading := simulation.GenerateECGReading(patient)

			if reading.HeartRate >= ecg.MinNormalHeartRate {
				t.Errorf("Expected bradycardia heart rate < %d, got %d",
					ecg.MinNormalHeartRate, reading.HeartRate)
			}

			expectedRR := 60.0 / float64(reading.HeartRate)
			if reading.RRInterval != expectedRR {
				t.Errorf("Expected RR interval %f, got %f", expectedRR, reading.RRInterval)
			}
		}
	})

	t.Run("Arrhythmia", func(t *testing.T) {
		patient := simulation.NewDefaultPatient()
		patient.SimulateArrhythmia = true

		for i := 0; i < 10; i++ {
			reading := simulation.GenerateECGReading(patient)

			baseRR := 60.0 / float64(reading.HeartRate)
			maxDiff := baseRR * patient.ArrhythmiaIntensity / 2.0

			diff := reading.RRInterval - baseRR
			if diff < -maxDiff || diff > maxDiff {
				t.Errorf("Arrhythmia RR interval %f outside expected range %f±%f",
					reading.RRInterval, baseRR, maxDiff)
			}
		}
	})
}

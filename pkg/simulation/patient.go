package simulation

import (
	"math/rand"
	"time"

	"arhm/ecg-monitoring/pkg/ecg"
)

type SimulatedPatient struct {
	ID            string
	BaseHeartRate int
	Variability   int
	RRVariability float64

	SimulateArrhythmia  bool
	SimulateTachycardia bool
	SimulateBradycardia bool

	ArrhythmiaIntensity float64
}

func NewDefaultPatient() SimulatedPatient {
	return SimulatedPatient{
		ID:                  "PATIENT",
		BaseHeartRate:       80,
		Variability:         5,
		RRVariability:       0.05,
		SimulateArrhythmia:  false,
		SimulateTachycardia: false,
		SimulateBradycardia: false,
		ArrhythmiaIntensity: 0.7,
	}
}

func GenerateECGReading(patient SimulatedPatient) ecg.ECGReading {
	heartRate := patient.BaseHeartRate
	var rrInterval float64

	if patient.SimulateTachycardia {
		heartRate = ecg.MaxNormalHeartRate + 1 + rand.Intn(29)
		rrInterval = 60.0 / float64(heartRate)
	} else if patient.SimulateBradycardia {
		heartRate = ecg.MinNormalHeartRate - 1 - rand.Intn(19)
		rrInterval = 60.0 / float64(heartRate)
	} else if patient.SimulateArrhythmia {
		heartRate = patient.BaseHeartRate + rand.Intn(20) - 10

		intensity := patient.ArrhythmiaIntensity
		baseRR := 60.0 / float64(heartRate)
		rrInterval = baseRR + (rand.Float64()*intensity-intensity/2.0)*baseRR
	} else {
		heartRate += rand.Intn(patient.Variability*2) - patient.Variability
		rrVariation := (rand.Float64() * patient.RRVariability * 2) - patient.RRVariability
		rrInterval = 60.0/float64(heartRate) + rrVariation
	}

	return ecg.ECGReading{
		Timestamp:  time.Now(),
		HeartRate:  heartRate,
		RRInterval: rrInterval,
	}
}

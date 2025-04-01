package simulation

import (
	"arhm/ecg-monitoring/pkg/ecg"
	"time"
)

type Condition string

const (
	ConditionNormal      Condition = "normal"
	ConditionTachycardia Condition = "tachycardia"
	ConditionBradycardia Condition = "bradycardia"
	ConditionArrhythmia  Condition = "arrhythmia"
)

type Controller struct {
	Patient         SimulatedPatient
	SimulationCycle []Condition
	CycleIndex      int
	CycleTime       int
	CycleLength     int
}

func NewController() *Controller {
	return &Controller{
		Patient:         NewDefaultPatient(),
		SimulationCycle: []Condition{ConditionNormal, ConditionTachycardia, ConditionNormal, ConditionBradycardia, ConditionNormal, ConditionArrhythmia},
		CycleIndex:      0,
		CycleTime:       0,
		CycleLength:     3,
	}
}

func (c *Controller) NextReading() (ecg.ECGReading, Condition) {
	currentCondition := c.SimulationCycle[c.CycleIndex]

	c.Patient.SimulateTachycardia = false
	c.Patient.SimulateBradycardia = false
	c.Patient.SimulateArrhythmia = false

	switch currentCondition {
	case ConditionTachycardia:
		c.Patient.SimulateTachycardia = true
	case ConditionBradycardia:
		c.Patient.SimulateBradycardia = true
	case ConditionArrhythmia:
		c.Patient.SimulateArrhythmia = true
	}

	reading := GenerateECGReading(c.Patient)

	c.advanceCycle()

	return reading, currentCondition
}

func (c *Controller) advanceCycle() {
	c.CycleTime++
	if c.CycleTime >= c.CycleLength {
		c.CycleTime = 0
		c.CycleIndex = (c.CycleIndex + 1) % len(c.SimulationCycle)
	}
}

func (c *Controller) RunWithCallback(interval time.Duration, callback func(reading ecg.ECGReading, condition Condition)) *time.Ticker {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			reading, condition := c.NextReading()
			callback(reading, condition)
		}
	}()

	return ticker
}

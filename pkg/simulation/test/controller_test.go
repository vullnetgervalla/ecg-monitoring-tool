package simulation_test

import (
	"arhm/ecg-monitoring/pkg/ecg"
	"arhm/ecg-monitoring/pkg/simulation"
	"testing"
	"time"
)

func TestNewController(t *testing.T) {
	controller := simulation.NewController()

	if controller.CycleIndex != 0 {
		t.Errorf("Expected initial CycleIndex to be 0, got %d", controller.CycleIndex)
	}

	if controller.CycleTime != 0 {
		t.Errorf("Expected initial CycleTime to be 0, got %d", controller.CycleTime)
	}

	if len(controller.SimulationCycle) != 6 {
		t.Errorf("Expected SimulationCycle to have 6 elements, got %d", len(controller.SimulationCycle))
	}

	if controller.SimulationCycle[0] != simulation.ConditionNormal {
		t.Errorf("Expected first condition to be 'normal', got '%s'", controller.SimulationCycle[0])
	}
}

func TestNextReading(t *testing.T) {
	controller := simulation.NewController()

	expectedConditions := []simulation.Condition{
		simulation.ConditionNormal,
		simulation.ConditionNormal,
		simulation.ConditionNormal,
		simulation.ConditionTachycardia,
		simulation.ConditionTachycardia,
		simulation.ConditionTachycardia,
	}

	for i, expectedCondition := range expectedConditions {
		reading, condition := controller.NextReading()

		if condition != expectedCondition {
			t.Errorf("Reading %d: Expected condition %s, got %s", i, expectedCondition, condition)
		}

		if condition == simulation.ConditionTachycardia && reading.HeartRate <= 100 {
			t.Errorf("Reading %d: Expected tachycardia heart rate > 100, got %d", i, reading.HeartRate)
		}

		if condition == simulation.ConditionBradycardia && reading.HeartRate >= 60 {
			t.Errorf("Reading %d: Expected bradycardia heart rate < 60, got %d", i, reading.HeartRate)
		}
	}
}

func TestAdvanceCycle(t *testing.T) {
	controller := simulation.NewController()

	if controller.CycleIndex != 0 || controller.CycleTime != 0 {
		t.Fatalf("Initial state incorrect: CycleIndex=%d, CycleTime=%d",
			controller.CycleIndex, controller.CycleTime)
	}

	controller.NextReading()
	if controller.CycleTime != 1 || controller.CycleIndex != 0 {
		t.Errorf("After 1 advance: Expected CycleTime=1, CycleIndex=0, got CycleTime=%d, CycleIndex=%d",
			controller.CycleTime, controller.CycleIndex)
	}

	controller.NextReading()
	if controller.CycleTime != 2 || controller.CycleIndex != 0 {
		t.Errorf("After 2 advances: Expected CycleTime=2, CycleIndex=0, got CycleTime=%d, CycleIndex=%d",
			controller.CycleTime, controller.CycleIndex)
	}

	controller.NextReading()
	if controller.CycleTime != 0 || controller.CycleIndex != 1 {
		t.Errorf("After 3 advances: Expected CycleTime=0, CycleIndex=1, got CycleTime=%d, CycleIndex=%d",
			controller.CycleTime, controller.CycleIndex)
	}
}

func TestRunWithCallback(t *testing.T) {
	controller := simulation.NewController()
	callCount := 0
	conditionsSeen := make(map[simulation.Condition]bool)

	callback := func(reading ecg.ECGReading, condition simulation.Condition) {
		callCount++
		conditionsSeen[condition] = true
	}

	ticker := controller.RunWithCallback(20*time.Millisecond, callback)
	defer ticker.Stop()

	time.Sleep(500 * time.Millisecond)

	if callCount < 5 {
		t.Errorf("Expected at least 5 calls to callback, got %d", callCount)
	}

	if len(conditionsSeen) < 2 {
		t.Errorf("Expected to see at least 2 different conditions, saw %d", len(conditionsSeen))
	}
}

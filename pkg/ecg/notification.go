package ecg

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
)

type NotificationType int

const (
	NotificationLog NotificationType = iota
	NotificationBeep
	NotificationBoth
)

type Notifier interface {
	Notify(condition HeartCondition)
}

type LogNotifier struct {
	Logger *log.Logger
}

func NewLogNotifier(output *os.File) *LogNotifier {
	return &LogNotifier{
		Logger: log.New(output, "ECG MONITOR: ", log.LstdFlags),
	}
}

func (n *LogNotifier) Notify(condition HeartCondition) {
	if condition.Type != ConditionNormal {
		n.Logger.Println(FormatAlert(condition))
	} else {
		n.Logger.Println("Normal heart activity detected")
	}
}

type BeepNotifier struct {
	MinSeverity string // Minimum severity to trigger a beep (normal, warning, critical)
	initialized bool
}

func NewBeepNotifier(minSeverity string) *BeepNotifier {
	notifier := &BeepNotifier{
		MinSeverity: minSeverity,
		initialized: false,
	}

	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		fmt.Printf("Warning: Audio initialization failed: %v", err)
		return notifier
	}

	notifier.initialized = true
	return notifier
}

func (n *BeepNotifier) Notify(condition HeartCondition) {
	if n.shouldBeep(condition) {
		n.triggerBeep(condition)
	}
}

func (n *BeepNotifier) shouldBeep(condition HeartCondition) bool {
	if condition.Type == ConditionNormal {
		return false
	}

	switch n.MinSeverity {
	case "normal":
		return true
	case "warning":
		return condition.Severity == "warning" || condition.Severity == "critical"
	case "critical":
		return condition.Severity == "critical"
	default:
		return false
	}
}

func (n *BeepNotifier) triggerBeep(condition HeartCondition) {
	if !n.initialized {
		fmt.Printf("BEEP ALERT: %s (no audio)\n", condition.Type)
		return
	}

	frequency := 800
	duration := time.Millisecond * 300

	if condition.Severity == "critical" {
		frequency = 1200
	}

	playBeep(frequency, duration)
}

func playBeep(frequency int, duration time.Duration) {
	sr := beep.SampleRate(44100)
	sine, err := generators.SinTone(sr, frequency)
	if err != nil {
		fmt.Printf("Error generating tone: %v\n", err)
		return
	}

	samples := sr.N(duration)
	limited := beep.Take(samples, sine)
	speaker.Play(limited)

	time.Sleep(duration)
}

type CompositeNotifier struct {
	Notifiers []Notifier
}

func NewCompositeNotifier(notifiers ...Notifier) *CompositeNotifier {
	return &CompositeNotifier{
		Notifiers: notifiers,
	}
}

func (n *CompositeNotifier) Notify(condition HeartCondition) {
	for _, notifier := range n.Notifiers {
		notifier.Notify(condition)
	}
}

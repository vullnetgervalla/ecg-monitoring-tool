package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"arhm/ecg-monitoring/pkg/ecg"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var minSeverity = flag.String("minseverity", "warning", "minimum severity for beep alerts (normal, warning, critical)")
var noColor = flag.Bool("no-color", false, "disable colored output")

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

const (
	timestampWidth  = 19 // YYYY-MM-DD HH:MM:SS
	heartRateWidth  = 10
	rrIntervalWidth = 11
	statusWidth     = 25
)

func formatWithColor(text string, condition ecg.HeartCondition) string {
	if *noColor {
		return text
	}

	var color string
	switch condition.Type {
	case ecg.ConditionNormal:
		color = colorGreen
	case ecg.ConditionTachycardia:
		if condition.Severity == "critical" {
			color = colorRed
		} else {
			color = colorYellow
		}
	case ecg.ConditionBradycardia:
		if condition.Severity == "critical" {
			color = colorRed
		} else {
			color = colorYellow
		}
	case ecg.ConditionArrhythmia:
		color = colorPurple
	default:
		color = colorWhite
	}

	return color + text + colorReset
}

type ConsoleNotifier struct {
	UseColor       bool
	SuppressOutput bool
}

func NewConsoleNotifier(useColor bool, suppressOutput bool) *ConsoleNotifier {
	return &ConsoleNotifier{
		UseColor:       !*noColor,
		SuppressOutput: suppressOutput,
	}
}

func (n *ConsoleNotifier) Notify(condition ecg.HeartCondition) {
	if n.SuppressOutput {
		return
	}

	if condition.Type == ecg.ConditionNormal {
		msg := "Normal heart activity detected"
		if n.UseColor {
			fmt.Println(formatWithColor(msg, condition))
		} else {
			fmt.Println(msg)
		}
	} else {
		alert := ecg.FormatAlert(condition)
		if n.UseColor {
			fmt.Println(formatWithColor(alert, condition))
		} else {
			fmt.Println(alert)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	fmt.Println("╔═══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                           ║")
	fmt.Println("║                      ECG Monitoring Tool - Client                         ║")
	fmt.Println("║                                                                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════════════════╝")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ecg"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	consoleNotifier := NewConsoleNotifier(!*noColor, true)
	beepNotifier := ecg.NewBeepNotifier(*minSeverity)
	notifier := ecg.NewCompositeNotifier(consoleNotifier, beepNotifier)

	done := make(chan struct{})
	readingCh := make(chan ecg.ECGReading)

	fmt.Println(string(colorCyan) + "\nMonitoring started.\n" + string(colorReset))

	tableWidth := timestampWidth + heartRateWidth + rrIntervalWidth + statusWidth + 11
	headerBorder := "╔"
	for i := 0; i < tableWidth; i++ {
		headerBorder += "═"
	}
	headerBorder += "╗"

	fmt.Println(headerBorder)
	fmt.Printf("║ %-*s │ %-*s │ %-*s │ %-*s ║\n",
		timestampWidth, "Timestamp",
		heartRateWidth, "Heart Rate",
		rrIntervalWidth, "RR Interval",
		statusWidth, "Status")

	separatorRow := "╟"
	for i := 0; i < timestampWidth+2; i++ {
		separatorRow += "─"
	}
	separatorRow += "┼"
	for i := 0; i < heartRateWidth+2; i++ {
		separatorRow += "─"
	}
	separatorRow += "┼"
	for i := 0; i < rrIntervalWidth+2; i++ {
		separatorRow += "─"
	}
	separatorRow += "┼"
	for i := 0; i < statusWidth+2; i++ {
		separatorRow += "─"
	}
	separatorRow += "╢"
	fmt.Println(separatorRow)


	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			var reading ecg.ECGReading
			err = json.Unmarshal(message, &reading)
			if err != nil {
				log.Println("Unmarshal error:", err)
				continue
			}

			readingCh <- reading
		}
	}()

	go func() {
		for reading := range readingCh {
			condition := ecg.AnalyzeReading(reading)

			timestamp := reading.Timestamp.Format("2006-01-02 15:04:05")
			status := condition.Type
			if condition.Type != ecg.ConditionNormal {
				status = fmt.Sprintf("%s (%s)", condition.Type, condition.Severity)
			}

			statusText := fmt.Sprintf("║ %-*s │ %*d │ %*.2f │ %-*s ║",
				timestampWidth, timestamp,
				heartRateWidth, reading.HeartRate,
				rrIntervalWidth, reading.RRInterval,
				statusWidth, status)

			fmt.Println(formatWithColor(statusText, condition))

			notifier.Notify(condition)
		}
	}()

	footerBorder := "╚"
	for i := 0; i < tableWidth; i++ {
		footerBorder += "═"
	}
	footerBorder += "╝"

	for {
		select {
		case <-done:
			fmt.Println(footerBorder)
			fmt.Println("Connection closed")
			return
		case <-interrupt:
			fmt.Println(footerBorder)
			log.Println("Interrupt received, closing connection...")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

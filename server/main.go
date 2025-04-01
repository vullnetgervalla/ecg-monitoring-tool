package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"arhm/ecg-monitoring/pkg/server"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var logFile = flag.String("logfile", "server/logs/ecg.log", "general log file")
var alertLogFile = flag.String("alertlog", "server/logs/alerts.log", "alerts-only log file")

func main() {
	flag.Parse()

	loggers, err := server.SetupLoggers(*logFile, *alertLogFile)
	if err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}
	defer loggers.Close()

	ecgHandler := server.NewECGHandler(loggers)
	http.Handle("/ecg", ecgHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		loggers.General.Printf("Starting ECG Simulation Server on %s", *addr)
		if err := http.ListenAndServe(*addr, nil); err != nil {
			loggers.General.Printf("HTTP server error: %v", err)
			stop <- os.Interrupt
		}
	}()

	<-stop
	loggers.General.Println("Shutting down server...")
}

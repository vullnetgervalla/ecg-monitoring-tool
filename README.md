# ECG Monitoring Tool

A real-time ECG monitoring system that simulates and analyzes heart activity.

## Running the Application

### Server
```bash
go run ./server
```

### Client
```bash
go run ./client
```

## Testing

Run tests:
```bash
go test ./pkg/ecg/test/... ./pkg/server/test/... ./pkg/simulation/test/... ./server/test/...
```

Or run with verbose output:
```bash
go test ./... -v
```

## Simulation Overview

The ECG Monitoring Tool simulates different heart conditions:

- **Normal**: Heart rate between 60-100 BPM with regular intervals
- **Tachycardia**: Heart rate >100 BPM
- **Bradycardia**: Heart rate <60 BPM
- **Arrhythmia**: Normal heart rate with irregular RR intervals

The simulation automatically cycles through these conditions to demonstrate the monitoring system's detection capabilities. Heart rates and RR intervals are generated based on the simulated condition, with appropriate randomization to create realistic variations.

The server sends readings to the client via WebSocket, where they are analyzed and displayed. Alerts are generated for abnormal conditions and logged to separate files. 
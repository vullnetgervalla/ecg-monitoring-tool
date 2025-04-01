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

### Note
For linux, `libasound2-dev` is needed for the beep sounds, you can install it by running the following:
```bash
sudo apt install libasound2-dev
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

The server sends readings to the client via WebSocket, where they are analyzed and displayed. Alerts are generated for abnormal conditions and logged to separate files in `server/logs/`. 

## Project Architecture

The application is divided into several packages, each with a specific responsibility:

### Client
Located in `./client/main.go`, the client application:
- Connects to the server via WebSockets
- Receives and processes ECG readings
- Displays real-time ECG data in a formatted table
- Provides visual alerts for abnormal heart conditions
- Produces audio alerts for critical conditions

### Server
Located in `./server/main.go`, the server application:
- Hosts the WebSocket endpoint for ECG data
- Manages logging to both general and alert-specific log files
- Controls the ECG simulation through the simulation package

### Packages

#### pkg/ecg
Core ECG data structures and analysis:
- `ecg.go`: Defines ECG readings and heart conditions
  - `ECGReading`: Data structure for heart rate and RR interval
  - `HeartCondition`: Classification of readings with severity
  - `AnalyzeReading()`: Analyzes readings to detect abnormal conditions
- `notification.go`: Alert mechanisms for abnormal conditions
  - Supports both console and audio notifications

#### pkg/server
Server-side components:
- `logger.go`: Logging infrastructure for general and alert logs
- `ws_handler.go`: WebSocket handler that:
  - Establishes connections with clients
  - Generates simulated ECG readings
  - Sends readings to connected clients
  - Logs alerts for abnormal conditions

#### pkg/simulation
ECG simulation components:
- `controller.go`: Controls the simulation cycle and parameters
  - Cycles through different heart conditions
  - Generates realistic ECG readings based on condition
- `patient.go`: Simulates a patient with configurable heart conditions
  - Generates realistic variations in heart rate and RR intervals
  - Supports simulation of tachycardia, bradycardia, and arrhythmia

### Data Flow
1. The server initiates the simulation controller
2. The controller generates ECG readings based on the simulated heart condition
3. Readings are sent to connected clients via WebSocket
4. The client analyzes the readings and provides visual/audio alerts
5. All abnormal conditions are logged on the server for record-keeping 

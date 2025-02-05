package simulation

import (
	"time"
)

// RobotStatus represents the current operational status of a robot
type RobotStatus string

const (
	StatusOperational RobotStatus = "OPERATIONAL"
	StatusWarning     RobotStatus = "WARNING"
	StatusError       RobotStatus = "ERROR"
	StatusOffline     RobotStatus = "OFFLINE"
)

// HeartbeatMessage represents the periodic connectivity check
type HeartbeatMessage struct {
	// RobotID uniquely identifies the robot sending the heartbeat
	RobotID    string      `json:"robot_id"`
	// Timestamp indicates when the heartbeat was generated
	Timestamp  time.Time   `json:"timestamp"`
	// Status represents the current operational state of the robot
	Status     RobotStatus `json:"status"`
	// BatteryPct indicates the remaining battery charge (0-100)
	BatteryPct float64     `json:"battery_percentage"`
}

// HealthMessage represents detailed hardware diagnostics
type HealthMessage struct {
	RobotID      string    `json:"robot_id"`
	Timestamp    time.Time `json:"timestamp"`
	CPUTemp      float64   `json:"cpu_temperature"`
	MotorTemps   []float64 `json:"motor_temperatures"`
	VoltageLevel float64   `json:"voltage_level"`
	CurrentDraw  float64   `json:"current_draw"`
	ErrorCodes   []string  `json:"error_codes,omitempty"`
}

// NavigationMessage represents path-related information
type NavigationMessage struct {
	RobotID    string     `json:"robot_id"`
	Timestamp  time.Time  `json:"timestamp"`
	Position   Position   `json:"position"`
	Heading    float64    `json:"heading"`
	Velocity   float64    `json:"velocity"`
	Obstacles  []Obstacle `json:"obstacles,omitempty"`
	PathStatus string     `json:"path_status"`
}

// Position represents 3D coordinates
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Obstacle represents detected obstacles
type Obstacle struct {
	Position Position `json:"position"`
	Size     float64  `json:"size"`
	Type     string   `json:"type"`
	Severity string   `json:"severity"`
}

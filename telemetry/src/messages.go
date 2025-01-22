package telemetry

import (
	"context"
	"encoding/json"
	"time"

	"telemetry/include/logger"
)

// SensorData represents the basic structure for all sensor readings
type SensorData struct {
	Timestamp time.Time   `json:"timestamp"`
	SensorID  string      `json:"sensor_id"`
	DataType  string      `json:"data_type"`
	Value     interface{} `json:"value"`
}

// Sensor defines the interface that all sensors must implement
type Sensor interface {
	ID() string
	Read(ctx context.Context) (SensorData, error)
	Initialize() error
	Shutdown() error
}

// IMUData represents the 6 DoF sensor readings
type IMUData struct {
	AccelX float64 `json:"accel_x"`
	AccelY float64 `json:"accel_y"`
	AccelZ float64 `json:"accel_z"`
	GyroX  float64 `json:"gyro_x"`
	GyroY  float64 `json:"gyro_y"`
	GyroZ  float64 `json:"gyro_z"`
}

// UltrasonicData represents distance measurements
type UltrasonicData struct {
	Distance float64 `json:"distance"` // Distance in centimeters
}

// MotorData represents motor telemetry
type MotorData struct {
	Speed     float64 `json:"speed"`     // Current speed
	Direction int     `json:"direction"` // 1 forward, -1 backward, 0 stopped
	Current   float64 `json:"current"`   // Current draw in amps
}

// Storage interface for persisting telemetry data
type Storage interface {
	Store(data SensorData) error
	Retrieve(sensorID string, startTime, endTime time.Time) ([]SensorData, error)
}

// SDCardStorage implements Storage interface for microSD card
type SDCardStorage struct {
	FilePath string
}

// TelemetryManager handles collection and storage of sensor data
type TelemetryManager struct {
	sensors  []Sensor
	storage  Storage
	interval time.Duration
	log      *logger.Logger
}

// NewTelemetryManager creates a new telemetry manager instance
func NewTelemetryManager(storage Storage, interval time.Duration) *TelemetryManager {
	return &TelemetryManager{
		sensors:  make([]Sensor, 0),
		storage:  storage,
		interval: interval,
		log:      logger.New(logger.INFO),
	}
}

// AddSensor registers a new sensor with the telemetry manager
func (tm *TelemetryManager) AddSensor(s Sensor) error {
	if err := s.Initialize(); err != nil {
		return err
	}
	tm.sensors = append(tm.sensors, s)
	return nil
}

// Start begins collecting telemetry data from all sensors
func (tm *TelemetryManager) Start(ctx context.Context) error {
	ticker := time.NewTicker(tm.interval)
	defer ticker.Stop()

	tm.log.Info("Starting telemetry collection")

	for {
		select {
		case <-ctx.Done():
			tm.log.Warn("Telemetry collection stopped: %v", ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			for _, sensor := range tm.sensors {
				data, err := sensor.Read(ctx)
				if err != nil {
					tm.log.Error("Error reading sensor %s: %v", sensor.ID(), err)
					continue
				}
				if err := tm.storage.Store(data); err != nil {
					tm.log.Error("Error storing data from sensor %s: %v", sensor.ID(), err)
					continue
				}
				tm.log.Debug("Collected data from sensor %s", sensor.ID())
			}
		}
	}
}

// Example IMU sensor implementation
type IMUSensor struct {
	id string
}

func (s *IMUSensor) ID() string {
	return s.id
}

func (s *IMUSensor) Read(ctx context.Context) (SensorData, error) {
	// Implement actual IMU reading logic here
	imuData := IMUData{
		// Fill with actual sensor readings
	}

	return SensorData{
		Timestamp: time.Now(),
		SensorID:  s.id,
		DataType:  "imu",
		Value:     imuData,
	}, nil
}

func (s *IMUSensor) Initialize() error {
	// Implement IMU initialization
	return nil
}

func (s *IMUSensor) Shutdown() error {
	// Implement IMU shutdown
	return nil
}

// Store implements Storage interface for SDCardStorage
func (s *SDCardStorage) Store(data SensorData) error {
	log := logger.New(logger.INFO)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("Failed to marshal sensor data: %v", err)
		return err
	}
	// Implement actual file writing logic
	log.Debug("Stored data for sensor %s", data.SensorID)
	return nil
}

func (s *SDCardStorage) Retrieve(sensorID string, startTime, endTime time.Time) ([]SensorData, error) {
	log := logger.New(logger.INFO)
	log.Debug("Retrieving data for sensor %s between %v and %v", sensorID, startTime, endTime)
	// Implementation for reading from SD card
	return nil, nil
}

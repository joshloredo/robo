package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// MockSensor implements Sensor interface for testing
type MockSensor struct {
	id        string
	readCount int
	mockData  SensorData
	shouldErr bool
}

func (s *MockSensor) ID() string {
	return s.id
}

func (s *MockSensor) Read(ctx context.Context) (SensorData, error) {
	s.readCount++
	if s.shouldErr {
		return SensorData{}, fmt.Errorf("mock sensor read error")
	}
	return s.mockData, nil
}

func (s *MockSensor) Initialize() error {
	if s.shouldErr {
		return fmt.Errorf("mock sensor initialization error")
	}
	return nil
}

func (s *MockSensor) Shutdown() error {
	if s.shouldErr {
		return fmt.Errorf("mock sensor shutdown error")
	}
	return nil
}

// MockStorage implements Storage interface for testing
type MockStorage struct {
	storedData []SensorData
	shouldErr  bool
}

func (s *MockStorage) Store(data SensorData) error {
	if s.shouldErr {
		return fmt.Errorf("mock storage error")
	}
	s.storedData = append(s.storedData, data)
	return nil
}

func (s *MockStorage) Retrieve(sensorID string, startTime, endTime time.Time) ([]SensorData, error) {
	if s.shouldErr {
		return nil, fmt.Errorf("mock retrieve error")
	}
	var result []SensorData
	for _, data := range s.storedData {
		if data.SensorID == sensorID && data.Timestamp.After(startTime) && data.Timestamp.Before(endTime) {
			result = append(result, data)
		}
	}
	return result, nil
}

// TestNewTelemetryManager tests the creation of a new telemetry manager
func TestNewTelemetryManager(t *testing.T) {
	storage := &MockStorage{}
	interval := time.Second
	tm := NewTelemetryManager(storage, interval)

	if tm == nil {
		t.Fatal("Expected non-nil TelemetryManager")
	}
	if tm.storage != storage {
		t.Error("Storage not properly set")
	}
	if tm.interval != interval {
		t.Error("Interval not properly set")
	}
}

// TestAddSensor tests adding a sensor to the telemetry manager
func TestAddSensor(t *testing.T) {
	tm := NewTelemetryManager(&MockStorage{}, time.Second)

	// Test successful sensor addition
	mockSensor := &MockSensor{id: "test-sensor"}
	err := tm.AddSensor(mockSensor)
	if err != nil {
		t.Errorf("Failed to add sensor: %v", err)
	}
	if len(tm.sensors) != 1 {
		t.Error("Sensor not added to manager")
	}

	// Test adding sensor that fails initialization
	failingSensor := &MockSensor{id: "failing-sensor", shouldErr: true}
	err = tm.AddSensor(failingSensor)
	if err == nil {
		t.Error("Expected error when adding failing sensor")
	}
}

// TestTelemetryManagerStart tests the Start method of TelemetryManager
func TestTelemetryManagerStart(t *testing.T) {
	storage := &MockStorage{}
	tm := NewTelemetryManager(storage, 100*time.Millisecond)

	mockData := SensorData{
		Timestamp: time.Now(),
		SensorID:  "test-sensor",
		DataType:  "test",
		Value:     "test-value",
	}

	mockSensor := &MockSensor{
		id:       "test-sensor",
		mockData: mockData,
	}

	err := tm.AddSensor(mockSensor)
	if err != nil {
		t.Fatalf("Failed to add sensor: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	// Start telemetry collection in a goroutine
	go func() {
		if err := tm.Start(ctx); err != nil && err != context.DeadlineExceeded {
			t.Errorf("Unexpected error during telemetry collection: %v", err)
		}
	}()

	// Wait for some data to be collected
	time.Sleep(200 * time.Millisecond)

	// Verify that data was collected
	if len(storage.storedData) == 0 {
		t.Error("No data collected")
	}
}

// TestIMUSensor tests the IMU sensor implementation
func TestIMUSensor(t *testing.T) {
	imu := &IMUSensor{id: "imu-1"}

	// Test ID method
	if imu.ID() != "imu-1" {
		t.Error("Incorrect sensor ID")
	}

	// Test Read method
	ctx := context.Background()
	data, err := imu.Read(ctx)
	if err != nil {
		t.Errorf("Failed to read IMU data: %v", err)
	}

	// Verify data structure
	if data.SensorID != "imu-1" {
		t.Error("Incorrect sensor ID in data")
	}
	if data.DataType != "imu" {
		t.Error("Incorrect data type")
	}

	// Test data conversion
	imuData, ok := data.Value.(IMUData)
	if !ok {
		t.Error("Failed to convert sensor data to IMUData")
	}
}

// TestSDCardStorage tests the SD card storage implementation
func TestSDCardStorage(t *testing.T) {
	storage := &SDCardStorage{FilePath: "/tmp/test.json"}

	testData := SensorData{
		Timestamp: time.Now(),
		SensorID:  "test-sensor",
		DataType:  "test",
		Value:     "test-value",
	}

	// Test storing data
	err := storage.Store(testData)
	if err != nil {
		t.Errorf("Failed to store data: %v", err)
	}

	// Test retrieving data
	startTime := time.Now().Add(-time.Hour)
	endTime := time.Now().Add(time.Hour)
	data, err := storage.Retrieve("test-sensor", startTime, endTime)
	if err != nil {
		t.Errorf("Failed to retrieve data: %v", err)
	}

	// Note: In a real implementation, you would verify the retrieved data
	// matches the stored data
}

// TestSensorDataSerialization tests JSON serialization of sensor data
func TestSensorDataSerialization(t *testing.T) {
	originalData := SensorData{
		Timestamp: time.Now(),
		SensorID:  "test-sensor",
		DataType:  "imu",
		Value: IMUData{
			AccelX: 1.0,
			AccelY: 2.0,
			AccelZ: 3.0,
			GyroX:  4.0,
			GyroY:  5.0,
			GyroZ:  6.0,
		},
	}

	// Test marshaling
	jsonData, err := json.Marshal(originalData)
	if err != nil {
		t.Errorf("Failed to marshal sensor data: %v", err)
	}

	// Test unmarshaling
	var decodedData SensorData
	err = json.Unmarshal(jsonData, &decodedData)
	if err != nil {
		t.Errorf("Failed to unmarshal sensor data: %v", err)
	}

	// Compare timestamps and basic fields
	if !originalData.Timestamp.Equal(decodedData.Timestamp) {
		t.Error("Timestamp mismatch after serialization")
	}
	if originalData.SensorID != decodedData.SensorID {
		t.Error("SensorID mismatch after serialization")
	}
	if originalData.DataType != decodedData.DataType {
		t.Error("DataType mismatch after serialization")
	}
}

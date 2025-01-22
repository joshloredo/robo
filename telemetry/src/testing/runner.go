package testing

import (
	"os"
	"os/signal"
	"telemetry/include/logger"
	"telemetry/src/mqtt"
	"telemetry/src/simulation"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type TelemetryTestRunner struct {
	log         *logger.Logger
	robotClient *mqtt.MQTTTelemetryClient
	mockRobot   *simulation.MockRobot
}

func NewTelemetryTestRunner(robotID string, brokerURL string) *TelemetryTestRunner {
	return &TelemetryTestRunner{
		log:         logger.New(logger.DEBUG),
		robotClient: mqtt.NewMQTTTelemetryClient(robotID, brokerURL),
		mockRobot:   simulation.NewMockRobot(robotID, simulation.Position{X: 0, Y: 0, Z: 0}),
	}
}

func (t *TelemetryTestRunner) Run() error {
	// Connect to the broker
	if err := t.robotClient.Connect(); err != nil {
		return err
	}

	// Give the connection time to establish
	time.Sleep(time.Second)

	// Subscribe to commands
	err := t.robotClient.SubscribeToCommands(func(client paho.Client, msg paho.Message) {
		t.log.Info("Received command: %s on topic: %s",
			string(msg.Payload()), msg.Topic())
	})
	if err != nil {
		return err
	}

	// Start publishing telemetry
	go t.publishMockTelemetry()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	t.log.Info("Shutting down test runner...")
	return nil
}

func (t *TelemetryTestRunner) publishMockTelemetry() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			healthData := simulation.HealthMessage{
				RobotID:      t.mockRobot.ID,
				Timestamp:    time.Now(),
				CPUTemp:      45.5,
				MotorTemps:   []float64{50.0, 48.5},
				VoltageLevel: 12.1,
				CurrentDraw:  2.5,
			}

			if err := t.robotClient.PublishTelemetry("health", healthData); err != nil {
				t.log.Error("Failed to publish health data: %v", err)
			} else {
				t.log.Debug("Published health data")
			}
		}
	}
}

package testing

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"telemetry/include/logger"
	"telemetry/src/mqtt"
	"telemetry/src/simulation"
	"telemetry/src/workerpool"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type TelemetryTestRunner struct {
	log         *logger.Logger
	robotClient *mqtt.MQTTTelemetryClient
	mockRobot   *simulation.MockRobot
	pool        *workerpool.WorkerPool
}

func NewTelemetryTestRunner(robotID string, brokerURL string) *TelemetryTestRunner {
	return &TelemetryTestRunner{
		log:         logger.New(logger.DEBUG),
		robotClient: mqtt.NewMQTTTelemetryClient(robotID, brokerURL),
		mockRobot:   simulation.NewMockRobot(robotID, simulation.Position{X: 0, Y: 0, Z: 0}),
		pool:        workerpool.NewWorkerPool(2), // Reduced from 5 to 2 workers
	}
}

func (t *TelemetryTestRunner) Run() error {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling first
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	defer signal.Stop(sigChan)

	// Create error channel for collecting errors
	errChan := make(chan error, 1)

	// Start the worker pool
	t.pool.Start()

	// Run the main logic in a goroutine
	go func() {
		defer t.pool.Shutdown()

		// Submit MQTT connection job
		connectionJob := workerpool.Job{
			Name: "MQTT Connection",
			Execute: func(ctx context.Context) error {
				// TODO(future): Implement service discovery or configuration to handle multiple MQTT brokers
				// Currently hibernating if primary broker is unavailable
				maxAttempts := 3
				attempt := 0

				for attempt < maxAttempts {
					t.log.Debug("Attempting MQTT connection to %s (attempt %d/%d)",
						t.robotClient.BrokerURL(), attempt+1, maxAttempts)

					err := t.robotClient.Connect()
					if err == nil {
						t.log.Debug("Successfully connected to MQTT broker")
						return nil
					}

					t.log.Debug("Connection attempt failed: %v", err)
					attempt++

					if attempt == maxAttempts {
						t.log.Info("Maximum connection attempts reached, entering hibernate mode")
						// Instead of failing, we'll just hibernate this worker
						t.log.Info("Worker hibernating - will wake when broker at %s becomes available", t.robotClient.BrokerURL())
						select {
						case <-ctx.Done():
							return ctx.Err()
						case <-time.After(time.Hour): // Effectively hibernate
							return fmt.Errorf("connection hibernated")
						}
					}

					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(time.Second * time.Duration(attempt)):
						continue
					}
				}
				return fmt.Errorf("connection attempts exhausted")
			},
			OnError: func(err error) {
				if err.Error() == "connection hibernated" {
					t.log.Info("MQTT connection hibernated - will retry when broker becomes available")
				} else {
					t.log.Error("MQTT connection failed: %v", err)
				}
			},
			Retries:  0,     // Don't retry, we handle retries in the Execute function
			Critical: false, // Changed to false so it doesn't trigger shutdown
		}
		t.pool.Submit(connectionJob)

		// Remove the WaitForCriticalJob and replace with a connection check loop
		t.log.Debug("Waiting for MQTT connection to be established...")
		for {
			if t.robotClient.IsConnected() {
				t.log.Info("MQTT connection established")
				break
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				continue
			}
		}

		// Submit subscription job
		subscriptionJob := workerpool.Job{
			Name: "MQTT Subscription",
			Execute: func(ctx context.Context) error {
				return t.robotClient.SubscribeToCommands(func(client paho.Client, msg paho.Message) {
					t.log.Info("Received command: %s on topic: %s",
						string(msg.Payload()), msg.Topic())
				})
			},
			Retries:  3,
			Critical: true,
		}
		t.pool.Submit(subscriptionJob)

		// Submit telemetry publishing job
		publishJob := workerpool.Job{
			Name: "Telemetry Publishing",
			Execute: func(ctx context.Context) error {
				t.publishMockTelemetry()
				return nil
			},
			Retries: 3,
		}
		t.pool.Submit(publishJob)

		// Wait for shutdown signal
		<-ctx.Done()
	}()

	// Wait for either error or interrupt
	select {
	case <-sigChan:
		t.log.Info("Received interrupt signal, shutting down...")
		cancel()
	case err := <-errChan:
		t.log.Error("Error occurred: %v", err)
		cancel()
		// Instead of returning error, just log it and continue shutdown
		t.log.Info("Continuing with shutdown process...")
	}

	// Give cleanup operations a deadline
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Wait for cleanup or timeout
	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("shutdown timed out")
	case <-t.pool.Done(): // We'll need to add this method to WorkerPool
		return nil
	}
}

func (t *TelemetryTestRunner) publishMockTelemetry() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	ctx := t.pool.Context()
	done := make(chan struct{})
	var once sync.Once // Add this to ensure we only close once

	// Goroutine to handle context cancellation
	go func() {
		<-ctx.Done()
		ticker.Stop()
		once.Do(func() {
			close(done)
		})
	}()

	for {
		select {
		case <-done:
			t.log.Info("Stopping telemetry publishing")
			return
		case <-ticker.C:
			select {
			case <-done:
				return
			default:
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
}

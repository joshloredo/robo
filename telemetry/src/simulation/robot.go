package simulation

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"telemetry/include/logger"
)

// MockRobot represents a simulated robot
type MockRobot struct {
	ID            string
	status        RobotStatus
	position      Position
	batteryLevel  float64
	log           *logger.Logger
	msgChan       chan interface{}
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

// NewMockRobot creates a new simulated robot
func NewMockRobot(id string, initialPosition Position) *MockRobot {
	return &MockRobot{
		ID:           id,
		status:       StatusOperational,
		position:     initialPosition,
		batteryLevel: 100.0,
		log:          logger.New(logger.DEBUG),
		msgChan:      make(chan interface{}, 100),
		stopChan:     make(chan struct{}),
	}
}

// Start begins the robot simulation
func (r *MockRobot) Start(ctx context.Context) {
	r.wg.Add(3) // One for each message type routine

	// Heartbeat routine
	go r.heartbeatRoutine(ctx)
	
	// Health monitoring routine
	go r.healthRoutine(ctx)
	
	// Navigation routine
	go r.navigationRoutine(ctx)

	r.log.Info("Robot %s started", r.ID)
}

func (r *MockRobot) heartbeatRoutine(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msg := HeartbeatMessage{
				RobotID:    r.ID,
				Timestamp:  time.Now(),
				Status:     r.status,
				BatteryPct: r.batteryLevel,
			}
			r.msgChan <- msg
		}
	}
}

func (r *MockRobot) healthRoutine(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msg := HealthMessage{
				RobotID:      r.ID,
				Timestamp:    time.Now(),
				CPUTemp:      35.0 + rand.Float64()*15.0,
				MotorTemps:   []float64{40.0 + rand.Float64()*20.0, 41.0 + rand.Float64()*20.0},
				VoltageLevel: 11.5 + rand.Float64()*1.0,
				CurrentDraw:  2.0 + rand.Float64()*3.0,
			}
			r.msgChan <- msg
		}
	}
}

func (r *MockRobot) navigationRoutine(ctx context.Context) {
	defer r.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Simulate movement
			r.position.X += (rand.Float64() - 0.5) * 0.5
			r.position.Y += (rand.Float64() - 0.5) * 0.5
			
			msg := NavigationMessage{
				RobotID:    r.ID,
				Timestamp:  time.Now(),
				Position:   r.position,
				Heading:    rand.Float64() * 360.0,
				Velocity:   rand.Float64() * 2.0,
				PathStatus: "CLEAR",
			}
			
			// Randomly generate obstacles
			if rand.Float64() < 0.1 { // 10% chance of obstacle
				msg.Obstacles = []Obstacle{{
					Position: Position{
						X: r.position.X + rand.Float64()*5.0,
						Y: r.position.Y + rand.Float64()*5.0,
						Z: 0,
					},
					Size:     rand.Float64() * 2.0,
					Type:     "STATIC",
					Severity: "MEDIUM",
				}}
			}
			
			r.msgChan <- msg
		}
	}
}

// GetMessageChannel returns the channel for receiving messages from the robot
func (r *MockRobot) GetMessageChannel() <-chan interface{} {
	return r.msgChan
}

// Stop gracefully stops the robot simulation
func (r *MockRobot) Stop() {
	close(r.stopChan)
	r.wg.Wait()
	close(r.msgChan)
	r.log.Info("Robot %s stopped", r.ID)
} 
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"telemetry/include/logger"
	"telemetry/src/simulation"
)

func main() {
	log := logger.New(logger.DEBUG)
	fleet := simulation.NewFleetManager()

	// Add some mock robots
	initialPositions := []simulation.Position{
		{X: 0, Y: 0, Z: 0},
		{X: 10, Y: 10, Z: 0},
		{X: -10, Y: -10, Z: 0},
	}

	for i, pos := range initialPositions {
		robotID := fmt.Sprintf("ROBOT_%03d", i+1)
		if err := fleet.AddRobot(robotID, pos); err != nil {
			log.Error("Failed to add robot: %v", err)
			return
		}
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start all robots
	fleet.StartAll(ctx)

	// Handle message processing
	channels := fleet.GetRobotChannels()
	go func() {
		for id, ch := range channels {
			go func(robotID string, msgChan <-chan interface{}) {
				for msg := range msgChan {
					// Process messages as needed
					log.Debug("Received message from %s: %+v", robotID, msg)
				}
			}(id, ch)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Info("Shutting down simulation...")
	cancel()
	fleet.StopAll()
}

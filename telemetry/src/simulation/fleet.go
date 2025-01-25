package simulation

import (
	"context"
	"fmt"
	"sync"
	"telemetry/include/logger"
)

// FleetManager manages a collection of mock robots
type FleetManager struct {
	robots map[string]*MockRobot
	log    *logger.Logger
	mu     sync.RWMutex
}

// NewFleetManager creates a new fleet manager
func NewFleetManager() *FleetManager {
	return &FleetManager{
		robots: make(map[string]*MockRobot),
		log:    logger.New(logger.DEBUG),
	}
}

// AddRobot adds a new robot to the fleet
func (fm *FleetManager) AddRobot(id string, initialPosition Position) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, exists := fm.robots[id]; exists {
		fm.log.Warn("Attempted to add duplicate robot with ID: %s", id)
		return ErrDuplicateRobot
	}

	fm.robots[id] = NewMockRobot(id, initialPosition)
	fm.log.Info("Added robot %s to fleet", id)
	return nil
}

// StartAll starts all robots in the fleet
func (fm *FleetManager) StartAll(ctx context.Context) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	for _, robot := range fm.robots {
		robot.Start(ctx)
	}
	fm.log.Info("Started all robots in fleet")
}

// StopAll stops all robots in the fleet
func (fm *FleetManager) StopAll() {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	for _, robot := range fm.robots {
		robot.Stop()
	}
	fm.log.Info("Stopped all robots in fleet")
}

// GetRobotChannels returns all message channels for the fleet
func (fm *FleetManager) GetRobotChannels() map[string]<-chan interface{} {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	channels := make(map[string]<-chan interface{})
	for id, robot := range fm.robots {
		channels[id] = robot.GetMessageChannel()
	}
	return channels
}

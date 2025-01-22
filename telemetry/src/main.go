package main

import (
	"telemetry/include/logger"
	"telemetry/src/testing"
)

func main() {
	log := logger.New(logger.DEBUG)

	runner := testing.NewTelemetryTestRunner(
		"TEST_ROBOT_001",
		"tcp://localhost:1883",
	)

	if err := runner.Run(); err != nil {
		log.Fatal("Test runner failed: %v", err)
	}
}

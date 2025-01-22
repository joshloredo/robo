package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"telemetry/include/logger"
	"time"
)

func printMenu() {
	fmt.Println("\n=== Test Menu ===")
	fmt.Println("1. Print Hello World")
	fmt.Println("2. Test Logger Levels")
	fmt.Println("3. Print Current Time")
	fmt.Println("0. Exit")
	fmt.Print("Enter your choice: ")
}

func testLoggerLevels(log *logger.Logger) {
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
}

func main() {
	// Initialize logger with DEBUG level
	log := logger.New(logger.DEBUG)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()

		// Read user input
		scanner.Scan()
		choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			log.Error("Invalid input: %v", err)
			continue
		}

		// Process user choice
		switch choice {
		case 0:
			log.Info("Exiting program...")
			return
		case 1:
			log.Info("Hello, World!")
		case 2:
			testLoggerLevels(log)
		case 3:
			log.Info("Current time: %v", time.Now().Format("2006-01-02 15:04:05"))
		default:
			log.Warn("Invalid option: %d", choice)
		}
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	// Default values
	defaultInterval := 2 * time.Second
	defaultLogPath := "./spammer.log"
	defaultLogFormat := "JSON"

	// Handle INTERVAL
	intervalStr := os.Getenv("INTERVAL")
	interval := defaultInterval
	if intervalStr != "" {
		parsedInterval, err := time.ParseDuration(intervalStr)
		if err != nil {
			log.Printf("Invalid INTERVAL '%s', using default %s", intervalStr, defaultInterval)
		} else {
			interval = parsedInterval
		}
	} else {
		log.Printf("INTERVAL not set, using default %s", defaultInterval)
	}

	// Handle LOG_PATH
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		log.Printf("LOG_PATH not set, using default: %s", defaultLogPath)
		logPath = defaultLogPath
	} else {
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Printf("Failed to create log directory %s: %v", logDir, err)
			log.Printf("Using default log path: %s", defaultLogPath)
			logPath = defaultLogPath
		}
	}

	// Attempt to open the log file
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file at %s: %v", logPath, err)
		log.Printf("Using default log path: %s", defaultLogPath)
		logFile, err = os.OpenFile(defaultLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open default log file: %v", err)
		}
	}
	defer logFile.Close()

	// Handle LOG_FORMAT
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat != "JSON" && logFormat != "PLAIN" {
		log.Printf("Invalid or missing LOG_FORMAT '%s', using default: %s", logFormat, defaultLogFormat)
		logFormat = defaultLogFormat
	}

	// Configure logger based on LOG_FORMAT
	var logger func(string)
	if logFormat == "JSON" {
		zlogger := zerolog.New(logFile).With().Timestamp().Logger()
		logger = func(message string) {
			zlogger.Info().Msg(message)
		}
	} else {
		plainLogger := log.New(logFile, "", log.LstdFlags)
		logger = func(message string) {
			plainLogger.Println(message)
		}
	}

	// Log messages at fixed intervals
	fmt.Printf("Logging every %s to %s in %s format\n", interval, logPath, logFormat)
	for {
		currentTime := time.Now().Format(time.RFC3339)
		message := fmt.Sprintf("Log message at %s", currentTime)
		logger(message)
		time.Sleep(interval)
	}
}

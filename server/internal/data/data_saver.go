package data

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logFile *os.File
	logMu   sync.Mutex
)

func init() {
	var err error
	logFile, err = os.OpenFile("snap.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}
}

// LogCommand writes an operation to the log with thread safety and timestamps.
func LogCommand(op, key, value string) {
	logMu.Lock()
	defer logMu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	line := fmt.Sprintf("%s | %s %s %s\n", timestamp, op, key, value)

	_, err := logFile.WriteString(line)
	if err != nil {
		// Also write to stderr so we don't silently lose logs
		log.Printf("failed to write to log file: %v", err)
	}

	// Flush to disk to prevent data loss on crash
	err = logFile.Sync()
	if err != nil {
		log.Printf("failed to flush log file: %v", err)
	}
}

// CloseLog should be called on program exit to release resources.
func CloseLog() {
	if logFile != nil {
		logFile.Close()
	}
}

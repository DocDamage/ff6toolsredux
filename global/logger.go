package global

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile *os.File
	logPath string
)

// InitLogger initializes the log file
func InitLogger() error {
	logPath = filepath.Join(os.TempDir(), fmt.Sprintf("ff6_editor_%d.log", time.Now().Unix()))
	f, err := os.Create(logPath)
	if err != nil {
		return err
	}
	logFile = f
	Log("[INIT] Log file: %s", logPath)
	return nil
}

// Log writes a message to the log file
func Log(format string, args ...interface{}) {
	if logFile != nil {
		timestamp := time.Now().Format("15:04:05.000")
		fmt.Fprintf(logFile, "[%s] %s\n", timestamp, fmt.Sprintf(format, args...))
		logFile.Sync() // Flush immediately
	}
}

// LogPath returns the path to the log file
func LogPath() string {
	return logPath
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// WriteCrashReport writes a crash report to a separate file
func WriteCrashReport(err interface{}) string {
	crashFile := filepath.Join(os.TempDir(), fmt.Sprintf("ff6_crash_%d.txt", time.Now().Unix()))
	if f, fileErr := os.Create(crashFile); fileErr == nil {
		fmt.Fprintf(f, "Crash time: %s\nError: %v\n", time.Now().Format(time.RFC3339), err)
		f.Close()
		return crashFile
	}
	return ""
}

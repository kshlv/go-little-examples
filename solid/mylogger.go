package main

import (
	"fmt"
	"os"
	"time"
)

// MyLogger ...
type MyLogger struct {
	format string
}

// New ...
func New() *MyLogger {
	return &MyLogger{
		format: "%d: %s\n",
	}
}

// Log ...
func (logger *MyLogger) Log(msg ...interface{}) {
	fmt.Fprintf(os.Stdout, logger.format, time.Now().Unix(), msg)
}

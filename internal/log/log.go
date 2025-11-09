package log

import "log"

// Simple wrapper to centralize logging (placeholder).
func Info(msg string)  { log.Println("INFO:", msg) }
func Error(msg string) { log.Println("ERROR:", msg) }

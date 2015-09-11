package logging

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	StandardTimeFormat = "2006-01-02 15:04:05.000Z07:00"
	tolerance          = 25 * time.Millisecond
)

// Formatter is interface for defining new custom log messages formats.
type Formatter interface {
	// Should return string that represents log message based on provided Record.
	Format(record *Record) string
}

// StandardFormatter adds time of logging with message to be logged.
type StandardFormatter struct {
	TimeFormat string
}

// Format constructs string for logging. Time of logging is added to log message.
// Also, if time of logging is more then 25 miliseconds in the passt, both
// times will be added to message (time when application sent log and time when
// message was processed). Otherwise, only time of processing will be written.
func (formatter *StandardFormatter) Format(record *Record) string {
	var message string
	now := time.Now()
	if now.Sub(record.Time) <= tolerance {
		message = record.Message
	} else {
		message = fmt.Sprintf("[%v] %v", record.Time.Format(formatter.TimeFormat), record.Message)
	}
	return fmt.Sprintf("[%v] %v %v", now.Format(formatter.TimeFormat), record.Level, message)
}

// JsonFormatter creates JSON struct with provided record.
type JsonFormatter struct{}

// Format creates JSON struct from provided record and returns it.
func (formatter *JsonFormatter) Format(record *Record) string {
	data, _ := json.Marshal(record)
	return string(data)
}

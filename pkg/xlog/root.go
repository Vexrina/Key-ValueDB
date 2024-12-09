package xlog

import (
	"fmt"
	"log"
	"os"
)

type LogField struct {
	Key   string
	Value any
}

func createMsg(msg string, fields ...LogField) string {
	s := msg + " "
	for _, field := range fields {
		logField := field
		s = s + fmt.Sprintf("%s: %v; ", logField.Key, logField.Value)
	}
	return s
}

func Debug(msg string, fields ...LogField) {
	s := createMsg(msg, fields...)
	log.Printf("[DEBUG] %s\n", s)
}

func Error(msg string, fields ...LogField) {
	s := createMsg(msg, fields...)
	log.Printf("[ERROR] %s\n", s)
}

func Info(msg string, fields ...LogField) {
	s := createMsg(msg, fields...)
	log.Printf("[INFO] %s\n", s)
}

func SetupLog() {
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)
}

func Field(key string, value any) LogField {
	return LogField{Key: key, Value: value}
}

func ErrorField(err error) LogField {
	return LogField{Key: "error", Value: err.Error()}
}

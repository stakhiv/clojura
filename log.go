package main

import (
	slog "log"
)

const (
	Fatal int = iota
	Error
	Info
	Debug
)

type Logger struct {
	level int
}

func NewLogger(level int) *Logger {
	return &Logger{
		level: level,
	}
}

func (l Logger) Debug(args ...interface{}) {
	l.print(Debug, args...)
}
func (l Logger) Info(args ...interface{}) {
	l.print(Info, args...)
}
func (l Logger) Error(args ...interface{}) {
	l.print(Error, args...)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.printf(Debug, format, args...)
}
func (l Logger) Infof(format string, args ...interface{}) {
	l.printf(Info, format, args...)
}
func (l Logger) Errorf(format string, args ...interface{}) {
	l.printf(Error, format, args...)
}

func (l *Logger) print(level int, args ...interface{}) {
	if l.level >= level {
		slog.Println(args...)
	}
}

func (l *Logger) printf(level int, format string, args ...interface{}) {
	if l.level >= level {
		slog.Printf(format, args...)
	}
}

package logger

import "go.uber.org/zap"

type Logger struct {
	zap.SugaredLogger
}

func New() Logger {
	return Logger{
		*zap.L().Sugar(),
	}
}

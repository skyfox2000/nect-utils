package utils

import (
	logger "github.com/skyfox2000/nect-utils/logger"
)

type UtilsTool struct {
	Name   string
	Debug  []string
	Logger *logger.LoggerEntry
}

type CustomError struct {
	Errno int
	Msg   string
	Data  interface{}
}

func NewError(errno int, msg string) *CustomError {
	return &CustomError{
		Errno: errno,
		Msg:   msg,
		Data:  nil,
	}
}

func (e *CustomError) Error() string {
	return e.Msg
}

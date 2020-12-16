package main

import (
	"context"
	"github.com/KeKe-Li/log"
)

type LogFunc func(ctx context.Context, msg string, fields ...interface{})

func GetLogFunc(err error) LogFunc {
	var logFunc LogFunc = log.ErrorContext
	return logFunc
}

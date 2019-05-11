package main

import (
	"os"
	"runtime"
	"time"

	"github.com/go-kit/kit/log"
)

type LogMiddleware struct {
	logger log.Logger
	next   KeyValueService
}

func NewLogMiddleware(next KeyValueService) *LogMiddleware {
	rst := LogMiddleware{
		logger: log.NewLogfmtLogger(os.Stderr),
		next:   next,
	}
	return &rst
}

func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func (logMiddleware *LogMiddleware) Set(key string, value string) (param string, err error) {
	defer func(funcName string) {
		logMiddleware.logger.Log(
			"time", time.Now(),
			"method", funcName,
			"input", param,
			"err", err,
		)
	}(GetFuncName())

	param, err = logMiddleware.next.Set(key, value)
	return
}

func (logMiddleware *LogMiddleware) Get(key string) (param string, err error) {
	defer func(funcName string) {
		logMiddleware.logger.Log(
			"time", time.Now(),
			"method", funcName,
			"input", param,
			"err", err,
		)
	}(GetFuncName())

	return logMiddleware.next.Get(key)
}

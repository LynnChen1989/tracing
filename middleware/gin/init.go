package middleware

import (
	"go.uber.org/zap"
	"tracing/lib"
)

var (
	Log *zap.Logger
)

func init() {
	Log = lib.GetLogger()
}

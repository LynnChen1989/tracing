package main

import (
	"tracing/lib"
	"go.uber.org/zap"
	"tracing/example"
)

var (
	Log *zap.Logger
)

func init() {
	Log = lib.GetLogger()
}

func main() {
	//api.Router()
	//example.TracerEntry()
	example.GinMain()
}

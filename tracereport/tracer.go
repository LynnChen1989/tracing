package tracereport

import (
	"go.uber.org/zap"
	"github.com/opentracing/opentracing-go"
	"io"
	"tracing/conf"
	"fmt"
	"tracing/lib"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"time"
)

var (
	Log *zap.Logger
)

func init() {
	Log = lib.GetLogger()
}

// 初始化追踪器
func InitTracer() (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: conf.JaegerAgent,
		},
		ServiceName: conf.SysName,
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

// NewTracer创建追踪器
func NewTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: conf.SysName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(conf.JaegerAgent, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
	)

	return tracer, closer, err
}

func GinTracer() (opentracing.Tracer, io.Closer) {

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const", //固定采样
			Param: 1,       //1=全采样、0=不采样
		},

		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: conf.JaegerAgent,
		},

		ServiceName: conf.SysName,
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))

	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

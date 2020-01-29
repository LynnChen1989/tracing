package middleware

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"tracing/tracereport"
)

func TraceSetUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var parentSpan opentracing.Span

		//tracer, closer := tracereport.InitTracer()
		tracer, closer := tracereport.GinTracer()
		defer closer.Close()

		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		spCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
		// TraceContextHeaderName: uber-trace-id
		// 必须使用这个
		if err != nil {
			Log.Warn("ctx waning", zap.String("waning", err.Error()))
			parentSpan = tracer.StartSpan(c.Request.URL.Path)
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer parentSpan.Finish()
		}
		c.Set("Tracer", tracer)
		c.Set("ParentSpanContext", parentSpan.Context())

		c.Next()
	}
}

package gogin

import (
	"testing"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/opentracing/opentracing-go"
	"net/http/httptest"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

func TestExample(t *testing.T) {

	tracer, closer := jaeger.NewTracer(
		"serviceName",
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
	)
	defer closer.Close()

	fn := func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		if span == nil {
			t.Error("Span is nil")
		}
	}

	r := gin.New()
	r.Use(Middleware(tracer))
	group := r.Group("")
	group.GET("", fn)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Error non-nil %v", err)
	}
	r.ServeHTTP(w, req)
}

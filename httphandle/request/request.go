package request

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"crypto/tls"
	"time"
	"go.uber.org/zap"
)

func HttpGet(url string, c *gin.Context) (string, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 5, //默认5秒超时时间
		Transport: tr,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	tracer, _ := c.Get("Tracer")
	parentSpanContext, _ := c.Get("ParentSpanContext")

	span := opentracing.StartSpan(
		"HTTP GET",
		opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)),
		// set tag
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		opentracing.Tag{Key: "http.url", Value: url},
		opentracing.Tag{Key: "http.method", Value: "GET"},
		ext.SpanKindRPCClient,
	)

	defer span.Finish()

	injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	if injectErr != nil {
		Log.Error("couldn't inject headers", zap.String("error", injectErr.Error()))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	return string(content), err
}

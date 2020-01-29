// 基于net.http的调用链追踪例子
package example

import (
	"io"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"fmt"
	"tracing/nethttp"
	"github.com/opentracing/opentracing-go/ext"
	"context"
	"io/ioutil"
	"tracing/tracereport"
	otLog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
)

var tracer opentracing.Tracer

func TracerEntry() {
	var (
		err      error
		ioCloser io.Closer
	)
	//创建tracer对象
	tracer, ioCloser, err = tracereport.NewTracer()
	if err != nil {
		Log.Error("tracer.NewTracer error")
	}
	defer ioCloser.Close()
	opentracing.SetGlobalTracer(tracer)

	//server
	http.HandleFunc("/ip", getIP)
	Log.Info("Starting server on port :8002")
	err = http.ListenAndServe(
		fmt.Sprintf(":%d", 8002),
		// use nethttp.Middleware to enable OpenTracing for server
		nethttp.Middleware(tracer, http.DefaultServeMux))
	if err != nil {
		Log.Error("Cannot start server")
	}
}

// 获取自己的ip地址
func getIP(w http.ResponseWriter, r *http.Request) {
	Log.Info("Received getIP request")

	//client
	client := &http.Client{Transport: &nethttp.Transport{}}
	span := tracer.StartSpan("getIP")
	span.SetTag(string(ext.Component), "getIP")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	req, err := http.NewRequest(
		"GET",
		"http://icanhazip.com",
		nil,
	)
	if err != nil {
		Log.Error("request error", zap.String("error", err.Error()))
	}

	req = req.WithContext(ctx)
	// wrap the request in net http.TraceRequest
	req, ht := nethttp.TraceRequest(tracer, req)
	defer ht.Finish()

	res, err := client.Do(req)
	if err != nil {
		onError(span, err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		onError(span, err)
		return
	}
	Log.Info("Received result", zap.String("result", string(body)))
	io.WriteString(w, fmt.Sprintf("ip %s", body))
}

func onError(span opentracing.Span, err error) {
	span.SetTag(string(ext.Error), true)
	span.LogKV(otLog.Error(err))
	Log.Error("client error", zap.String("error", err.Error()))
}

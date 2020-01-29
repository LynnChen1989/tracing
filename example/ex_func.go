// 基于函数的调用链追踪例子

package example

import (
	"context"
	"time"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"tracing/tracereport"
)

func foo3(req string, ctx context.Context) (reply string) {
	//1.创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "span_foo3")
	defer func() {
		//4.接口调用完，在tag中设置request和reply
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	Log.Info(req)
	//2.模拟处理耗时
	time.Sleep(time.Second / 2)
	//3.返回reply
	reply = "foo3Reply"
	Log.Info("foo3 replay", zap.String("content", reply))
	return
}

//跟foo3一样逻辑
func foo4(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_foo4")
	defer func() {
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	Log.Info(req)
	time.Sleep(time.Second / 2)
	reply = "foo4Reply"
	Log.Info("foo4 replay", zap.String("content", reply))
	return
}

// 链路追踪
func ExFuncTracing() {
	tracer, closer := tracereport.InitTracer()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer) //Start span From Context创建新span时会用到
	span := tracer.StartSpan("span_root")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	foo3("Hello foo3", ctx)
	foo4("Hello foo4", ctx)
	span.Finish()
}

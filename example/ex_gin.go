package example

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"encoding/json"
	"time"
	"fmt"
	"os"
	"tracing/httphandle/response"
	"go.uber.org/zap"
	ginMiddleWare "tracing/middleware/gin"
	"os/signal"
	"context"

	"tracing/httphandle/request"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func GinMain() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	GinRouter(engine)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8002),
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	fmt.Println("|-----------------------------------|")
	fmt.Println("|            jaeger ex gin          |")
	fmt.Println("|-----------------------------------|")
	fmt.Println("|  Go Http Server Start Successful  |")
	fmt.Println("|    Port: 8002" + "    Pid:" + fmt.Sprintf("%d", os.Getpid()) + "        |")
	fmt.Println("|-----------------------------------|")
	fmt.Println("")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log.Error("HTTP server listen Failed")
		}
	}()

	//
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	Log.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		Log.Error("Server Shutdown error")
	}
	Log.Info("server exiting")
}

func GinRouter(engine *gin.Engine) {

	logger, _ := zap.NewProduction()

	//加载Gin中间件
	engine.Use(
		ginMiddleWare.HeaderSetUp(),
		ginMiddleWare.TraceSetUp(),                       // 调用链追踪中间件
		ginMiddleWare.GinZap(logger, time.RFC3339, true), // zap日志中间件
	)

	//404
	engine.NoRoute(func(c *gin.Context) {
		utilGin := response.Gin{Ctx: c}
		utilGin.Response(404, "请求方法或路径不存在", nil)
	})

	// 路由配置
	engine.GET("/api/v1/ip", GetIpHandler)
}

func GetIpHandler(c *gin.Context) {
	utilGin := response.Gin{Ctx: c}
	//c.Get("")

	parentSpanContext, _ := c.Get("ParentSpanContext")
	sp := opentracing.StartSpan(
		"GetIpHandler",
		opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)),
		opentracing.Tag{Key: "action", Value: GetIpHandler},
		ext.SpanKindRPCServer,

	)
	defer sp.Finish()

	c.Set("ParentSpanContext", sp.Context())
	_, err := request.HttpGet("http://icanhazip.com", c)
	if err == nil {
		utilGin.Response(200, "ok", nil)
	} else {
		utilGin.Response(500, "", nil)
	}

}

// -----** 无用 -----
// 登录接口
func LoginHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "username|password is missed",
		})
		return
	}

	if username == "chenlin" && password == "chenlin@2020" {
		Log.Info("user: " + username + " login success")
		c.JSON(http.StatusOK, gin.H{
			"data": "login success",
		})
		return
	} else {
		Log.Error("user: " + username + " login failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": "login failed",
		})
		return
	}
}

// 注册接口
func RegisterHandler(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	var postMap map[string]interface{}

	if err := json.Unmarshal(buf[0:n], &postMap); err == nil {
		username := postMap["username"]
		password := postMap["password"]
		if username != "" && password != "" {
			Log.Info("user:" + username.(string) + " register success")
			c.JSON(http.StatusOK, gin.H{
				"data": "register success",
			})
			return

		} else {
			Log.Error("user:" + username.(string) + " register failed")
			c.JSON(http.StatusInternalServerError, gin.H{
				"data": "register failed",
			})
			return
		}
	}
}

func CreditHandler(c *gin.Context) {

}

func RealNameHandler(c *gin.Context) {

}

func WithdrawHandler(c *gin.Context) {

}

// -----** end --------

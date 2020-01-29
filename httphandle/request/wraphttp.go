package request

import (
	"net/http"
	"bytes"
	"time"
	"io/ioutil"
	"go.uber.org/zap"
	"github.com/opentracing/opentracing-go"
	"context"
	"tracing/nethttp"
	"encoding/json"
	"errors"
)

type HTTPClient struct {
	Tracer opentracing.Tracer
	Client *http.Client
}

// 调用链封装GET, 只能用于原生http请求，不能用在框架
func (c *HTTPClient) GetJSON(ctx context.Context, endpoint string, url string, out interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(c.Tracer, req)
	defer ht.Finish()

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

// 常规封装POST
func HttpPost(url string, postData string) {
	var jsonStr = []byte(postData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("http request error", zap.String("error", err.Error()))
	}
	defer resp.Body.Close()

	Log.Info("response status", zap.String("code", resp.Status))
	body, _ := ioutil.ReadAll(resp.Body)
	Log.Info("response Body:", zap.String("body", string(body)))
}

// 常规封装GET
func (c *HTTPClient) HttpGet(url string) {

	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: 5 * time.Second}
	//client := &c.Client
	//client.Timeout = 5 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("http request error", zap.String("error", err.Error()))
	}
	defer resp.Body.Close()

	Log.Info("response status", zap.String("code", resp.Status))
	body, _ := ioutil.ReadAll(resp.Body)
	Log.Info("response Body:", zap.String("body", string(body)))
}

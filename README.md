# 链路追踪

## 模式

+ gRPC拦截
+ gRPC注入
+ HTTP拦截，被请求时
+ HTTp注入,请求时


## go实践

```
（1）net/http原生服务，通过https://github.com/opentracing-contrib/go-stdlib.git中间件
（2）gin框架通过编写自定义中间件
```

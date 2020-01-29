package response

import "gopkg.in/gin-gonic/gin.v1"

type Gin struct {
	Ctx *gin.Context
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (g *Gin) Response(code int, msg string, data interface{}) {
	g.Ctx.JSON(200, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
	return
}

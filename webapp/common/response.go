package common

import "github.com/gin-gonic/gin"

type MyGinContext struct {
	gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (c *MyGinContext) Response(httpCode int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  msg,
		Data: data,
	})
}

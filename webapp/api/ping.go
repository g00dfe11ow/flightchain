package api

import "github.com/gin-gonic/gin"

// 测试服务器是否成功的功能
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

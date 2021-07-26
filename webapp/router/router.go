package router

import (
	"flightchain/api"

	"github.com/gin-gonic/gin"
)

// InitRouter初始化路由信息
func InitRouter() *gin.Engine {
	router := gin.Default()
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/ping", api.Ping)                  // 测试网络连通性
		apiGroup.POST("/createplan", api.CreatePlan)     // 创建计划
		apiGroup.POST("/readplan", api.ReadPlan)         // 读取已有计划
		apiGroup.POST("/modifyplan", api.ModifyPlan)     // 修改计划
		apiGroup.POST("/approvalplan", api.ApprovalPlan) // 审批计划
		apiGroup.POST("/addreader", api.AddReader)       // 增加读者
		apiGroup.POST("/removereader", api.RemoveReader) // 移除读者
		apiGroup.POST("/getallplans", api.GetAllPlans)   // 获取所有计划
		// apiGroup.GET("/getallplans", api.GetAllPlans) // 获取所有计划
	}
	return router
}

package api

import (
	"flightchain/blockchain"
	"flightchain/common"
	"flightchain/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreatePlan(c *gin.Context) {
	body := new(model.NewPlan)
	// 解析body参数
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "解析body参数失败",
			"errorlog": err.Error(),
		})
		return
	}
	// 检查参数
	if body.ID == "" || body.UnitOrIndividual == "" || !common.VerifyMobileFormat(body.Contact) || body.Aircraft == (model.Aircraft{}) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "参数检查失败",
			"errorlog": "参数均不能为空",
		})
		return
	}
	// 准备传参
	var bodyBytes []string
	bodyBytes = append(bodyBytes, string(body.ID))
	bodyBytes = append(bodyBytes, string(body.UnitOrIndividual))
	bodyBytes = append(bodyBytes, string(body.Contact))
	stringed := common.StructToString(body.Aircraft)
	lower := strings.ToLower(stringed)
	bodyBytes = append(bodyBytes, lower)
	// 调用智能合约
	resp, err := blockchain.ChannelExecute("CreatePlan", bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":  "chaincode调用失败",
			"errorlog": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "chaincode调用成功",
		"info":    resp,
	})
}

func ReadPlan(c *gin.Context) {
	body := new(model.QueryID)
	// 解析body参数
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "解析body参数失败",
			"errorlog": err.Error(),
		})
		return
	}
	// 检查参数
	if body.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "参数检查失败",
			"errorlog": "参数均不能为空",
		})
		return
	}
	// 准备传参
	var bodyBytes []string
	bodyBytes = append(bodyBytes, string(body.ID))
	// 调用智能合约
	resp, err := blockchain.ChannelExecute("ReadPlan", bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":  "chaincode调用失败",
			"errorlog": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "chaincode调用成功",
		"info":    resp,
	})
}

func ModifyPlan(c *gin.Context) {
	body := new(model.NewPlan)
	// 解析body参数
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "解析body参数失败",
			"errorlog": err.Error(),
		})
		return
	}
	// 检查参数
	if body.ID == "" || body.UnitOrIndividual == "" || !common.VerifyMobileFormat(body.Contact) || body.Aircraft == (model.Aircraft{}) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "参数检查失败",
			"errorlog": "参数均不能为空",
		})
		return
	}
	// 准备传参
	var bodyBytes []string
	bodyBytes = append(bodyBytes, string(body.ID))
	bodyBytes = append(bodyBytes, string(body.UnitOrIndividual))
	bodyBytes = append(bodyBytes, string(body.Contact))
	// Aircraft为结构体，需要多处理一下才可以使用
	stringed := common.StructToString(body.Aircraft)
	lower := strings.ToLower(stringed)
	bodyBytes = append(bodyBytes, lower)
	// 调用智能合约
	resp, err := blockchain.ChannelExecute("ModifyPlan", bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":  "chaincode调用失败",
			"errorlog": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "chaincode调用成功",
		"info":    resp,
	})
}

func ApprovalPlan(c *gin.Context) {
	body := new(model.ApprovalType)
	// 解析body参数
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "解析body参数失败",
			"errorlog": err.Error(),
		})
		return
	}
	// 检查参数
	if body.ID == "" || body.State == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "参数检查失败",
			"errorlog": "参数均不能为空",
		})
		return
	}
	// 准备传参
	var bodyBytes []string
	bodyBytes = append(bodyBytes, string(body.ID))
	bodyBytes = append(bodyBytes, string(body.State))

	// 调用智能合约
	resp, err := blockchain.ChannelExecute("ApprovalPlan", bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":  "chaincode调用失败",
			"errorlog": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "chaincode调用成功",
		"info":    resp,
	})
}

func AddReader(c *gin.Context) {
	body := new(model.Readers)
	// 解析body参数
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "解析body参数失败",
			"errorlog": err.Error(),
		})
		return
	}
	// 检查参数
	if body.ID == "" || body.Reader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "参数检查失败",
			"errorlog": "参数均不能为空",
		})
		return
	}
	// 准备传参
	var bodyBytes []string
	bodyBytes = append(bodyBytes, string(body.ID))
	bodyBytes = append(bodyBytes, body.Reader)

	// 调用智能合约
	resp, err := blockchain.ChannelExecute("AddReader", bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":  "chaincode调用失败",
			"errorlog": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "chaincode调用成功",
		"info":    resp,
	})
}

func RemoveReader(c *gin.Context) {
	body := new(model.Readers)
	// 解析body参数
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "解析body参数失败",
			"errorlog": err.Error(),
		})
		return
	}
	// 检查参数
	if body.ID == "" || body.Reader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "参数检查失败",
			"errorlog": "参数均不能为空",
		})
		return
	}
	// 准备传参
	var bodyBytes []string
	bodyBytes = append(bodyBytes, string(body.ID))
	bodyBytes = append(bodyBytes, body.Reader)

	// 调用智能合约
	resp, err := blockchain.ChannelExecute("RemoveReader", bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":  "chaincode调用失败",
			"errorlog": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "chaincode调用成功",
		"info":    resp,
	})
}

func GetAllPlans(c *gin.Context) {
	// 准备传参
	var bodyBytes []string

	// 调用智能合约
	resp, err := blockchain.ChannelExecute("GetAllPlans", bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":  "chaincode调用失败",
			"errorlog": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "chaincode调用成功",
		"info":    resp,
	})
}

// func GetAllPlans(c *gin.Context) {
// 	// 调用智能合约
// 	resp, err := blockchain.GetAllPlans()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message":  "chaincode调用失败",
// 			"errorlog": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "chaincode调用成功",
// 		"info":    resp,
// 	})
// }

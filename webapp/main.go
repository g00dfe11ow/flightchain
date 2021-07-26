package main

import (
	"flightchain/blockchain"
	"flightchain/router"
)

func main() {
	blockchain.Init()                // 初始化区块链链接配置
	newrouter := router.InitRouter() // 初始化路由器
	newrouter.Run()                  // 运行服务器
}

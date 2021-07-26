package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 配置信息
var (
	SDK           *fabsdk.FabricSDK                                                     // Fabric提供的SDK
	ChannelName   = "mychannel"                                                         // 通道名称
	ChainCodeName = "plan"                                                              // 链码名称
	Org           = "applicant"                                                         // 组织名称
	User          = "Admin"                                                             // 用户
	ConfigPath    = "/home/pillow/Desktop/flightchain/webapp/connection-applicant.yaml" // 配置文件路径
)

// Init 初始化
func Init_old() {
	var err error
	// 通过配置文件初始化SDK
	SDK, err = fabsdk.New(config.FromFile(ConfigPath))
	if err != nil {
		panic(err)
	}
}

// ChannelExecute 区块链交互，会改变账本的操作
func ChannelExecute_old(fcn string, args [][]byte) (channel.Response, error) {
	// 创建客户端，表明在通道的身份
	ctx := SDK.ChannelContext(ChannelName, fabsdk.WithOrg(Org), fabsdk.WithUser(User))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}
	// 对区块链增删改的操作（调用了链码的invoke）
	// 向区块连网络请求的内容
	req := channel.Request{
		ChaincodeID: ChainCodeName, // 调用的链码
		Fcn:         fcn,           // 调用的函数
		Args:        args,          // 函数的参数
	}
	// 请求的peer节点
	reqPeers := channel.WithTargetEndpoints("peer0.applicant.flight.com")
	resp, err := cli.Execute(req, reqPeers)
	if err != nil {
		return channel.Response{}, err
	}
	//返回链码执行后的结果
	return resp, nil
}

// ChannelQuery 区块链查询，不会改变账本的操作
func ChannelQuery_old(fcn string, args [][]byte) (channel.Response, error) {
	// 创建客户端，表明在通道的身份
	ctx := SDK.ChannelContext(ChannelName, fabsdk.WithOrg(Org), fabsdk.WithUser(User))
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}
	// 对区块链查询的操作（调用了链码的invoke），将结果返回
	// 向区块连网络请求的内容
	req := channel.Request{
		ChaincodeID: ChainCodeName, // 调用的链码
		Fcn:         fcn,           // 调用的函数
		Args:        args,          // 函数的参数
	}
	// 请求的peer节点
	reqPeers := channel.WithTargetEndpoints("peer0.applicant.flight.com")
	resp, err := cli.Execute(req, reqPeers)
	if err != nil {
		return channel.Response{}, err
	}
	//返回链码执行后的结果
	return resp, nil
}

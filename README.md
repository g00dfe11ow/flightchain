# ~~网络启动方法~~
1. ~~进入`test-network`文件夹;~~  
2. ~~运行`./network.sh up createChannel -c mychannel -ca`~~  
   
    > ~~该命令会在创建网络的同时创建一个名为`mychannel`的channel，以及ca~~  
3. ~~`./network.sh deployCC -ccn plan -ccp ../asset-transfer-basic-modified/chaincode-go/ -ccl go`~~  
   
    > ~~该命令会在上一步中创建好的通道上部署一个名为`plan`的chaincode~~  
4. ~~`cd asset-transfer-basic-modified/application-go`~~  
   
    > ~~进入`application-go`文件夹，分别有`applicant`和`approver`的文件夹，里面有各自的应用程序~~
5. ~~进入相应的文件夹运行`go run <identity>.go`运行程序~~  
***
> ~~**注意：**每次重启网络时，要将`application-go`文件夹下个应用的`keystore`和`wallet`都删掉~~

## TL；DR

1. 进入`test-network`文件夹;  

2. 运行`./network.sh up createChannel -c mychannel -ca`  

   > 该命令会在创建网络的同时创建一个名为`mychannel`的channel，以及ca  

3. `./network.sh deployCC -ccn plan -ccp ../asset-transfer-basic-modified/chaincode-go/ -ccl go`  

   > 该命令会在上一步中创建好的通道上部署一个名为`plan`的chaincode  

4. 在`test-network`中的网络创建完成后，cd到`webapp`目录，即在`main.go`文件同级目录下，运行`go run main.go`启动后端，往相应的api接口发送请求即可使用。

5. `test.json`中为一测试样例，可用于测试`CreatePlan`，其他的功能可以根据对应功能接收的参数进行修改。

# 项目结构

```bash
┌─[pillow@ubuntu] - [~/Desktop/flightchain] - [2021-09-07 01:42:18]
└─[0] <git:(main ac88fda✱) > tree -L 1
.
├── asset-transfer-basic-modified ---> 根据官方实例修改的chaincode 
├── bin                           ---> 项目需要的二进制文件
├── config                        ---> 项目所需的配置文件
├── test-network                  ---> 根据官方示例修改的网络
├── webapp                        ---> 项目的后端
├── test.json                     ---> 测试样例
└── README.md                     ---> 项目说明文档，即本文档

5 directories, 2 file
```

以下进行具体说明：

## `assert-transfer-basic-modified`

```bash
─[pillow@ubuntu] - [~/Desktop/flightchain/asset-transfer-basic-modified] - [2021-09-07 01:48:43]
└─[0] <git:(main ac88fda✱) > tree -L 2
.
├── application-go        ---> 根据官方示例写的测试用application，有了`webapp`后可以说没用了
│   ├── applicant
│   ├── approver
│   ├── go.mod
│   └── go.sum
└── chaincode-go          ---> 使用的chaincode，修改了官方示例，同时加入了访问控制
    ├── assetTransfer.go
    ├── chaincode
    ├── go.mod
    ├── go.sum
    └── vendor
    
6 directories, 5 files
```

chaincode的主要逻辑在`chaincode-go/chaincode/smartcontract.go`中，其他的文件是官方示例就有的，应该没有改动。

`smartcontract.go`实现的功能在代码中有相应的注释，具体功能说明如下：

```go
InitLedger()                   ---> 初始化账本，主要用于测试
CreatePlan()                   ---> 创建飞行计划，只有ApplicantMSP有权限创建飞行计划
ReadPlan()                     ---> 读计划，如果是applicant和approver可以直接读取，如果是user，则判断是否在ReadSet中
ModifyPlan()                   ---> 用参数来修改已有的飞行计划，只有ApplicantMSP有权限修改飞行计划，且是能修改自己提交的飞行计划
DeletePlan()                   ---> 删除已有的飞行计划（访问控制没做，为了溯源可能不需要这个功能）
PlanExists()                   ---> 查询对应ID的飞行计划是否存在
ApprovalPlan()                 ---> approver审批提交上来的飞行计划
AddReader()                    ---> 添加新的用户到ReadSet中
RemoveReader()                 ---> 将某用户从ReadSet中删除
GetAllPlans()                  ---> 获取所有飞行计划
GetSubmittingClientIdentity()  ---> 获取用户身份
InString()                     ---> 判断字符串是否在字符数组中
Remove()                       ---> 删除index位置的字符串
```

> 注意：以上说明时只保留了函数名，没有参数列表和返回值。

## `test-network`

`test-network`中内容，大体与官方实例中相同，改动的地方为结合本项目的需求对`network.sh`和相应的配置文件进行了修改。

具体改动没有进行记录，可以通过下载官方原版后与改动后的通过`diff`指令查看改动处。

本文件夹主要用于启动网络：

> 1. 进入`test-network`文件夹;  
>
> 2. 运行`./network.sh up createChannel -c mychannel -ca`  
>
>    > 该命令会在创建网络的同时创建一个名为`mychannel`的channel，以及ca  
>
> 3. `./network.sh deployCC -ccn plan -ccp ../asset-transfer-basic-modified/chaincode-go/ -ccl go`  
>
>    > 该命令会在上一步中创建好的通道上部署一个名为`plan`的chaincode  

## `webapp`

后端的实现上使用了`gin`这个框架，具体说明如下：

```bash
┌─[pillow@ubuntu] - [~/Desktop/flightchain/webapp] - [2021-09-07 02:13:13]
└─[0] <git:(main ac88fda✱) > tree
.
├── api
│   ├── applicant.go  ---> 主要功能的实现，与chaincode中的功能基本一一对应，对chaincode中的功能进行封装以方便使用
│   └── ping.go       ---> 测试网络是否搭建成功的功能
├── blockchain
│   ├── connect.go    ---> 与区块链相关的功能，如初始化网络，生成钱包和与区块链网络进行交互
│   └── sdk.go        ---> coding中间产物，好像没用了
├── common
│   ├── response.go   ---> 对常用功能的封装，为了使用更方便
│   └── util.go       ---> 同样也是coding中间产物，应该也没用了，本来打算用这里的代码解决的问题用别的方式解决了
├── router
    └── router.go     ---> 路由文件
├── model
│   └── model.go      ---> 模型，项目中用的数据结构
├── go.mod
├── go.sum
└── main.go           ---> 项目入口

5 directories, 11 files
```

在`test-network`中的网络创建完成后，cd到`webapp`目录，即在`main.go`文件同级目录下，运行`go run main.go`启动后端，往相应的api接口发送请求即可使用；`test.json`中为一测试样例，可用于测试`CreatePlan`，其他的功能可以根据对应功能接收的参数进行修改。


# 网络启动方法
1. 进入`test-network`文件夹;  
2. 运行`./network.sh up createChannel -c mychannel -ca`  
    > 该命令会在创建网络的同时创建一个名为`mychannel`的channel，以及ca  
3. `./network.sh deployCC -ccn plan -ccp ../asset-transfer-basic-modified/chaincode-go/ -ccl go`  
    > 该命令会在上一步中创建好的通道上部署一个名为`plan`的chaincode  
4. `cd asset-transfer-basic-modified/application-go`  
    > 进入`application-go`文件夹，分别有`applicant`和`approver`的文件夹，里面有各自的应用程序
5. 进入相应的文件夹运行`go run <identity>.go`运行程序  
***
> **注意：**每次重启网络时，要将`application-go`文件夹下个应用的`keystore`和`wallet`都删掉
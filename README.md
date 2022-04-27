fabric-samples-gm
============================

基于`fabric-samples`的master分支(27ac653)的国密改造。

# 使用说明
这里说明一下`go1.17.5`环境下国密改造后最简单的智能合约运行示例。

- fabric-gm : 参考 "https://gitee.com/zhaochuninhefei/fabric-gm", 编译出对应的本地二进制文件与docker镜像。
- 智能合约: `asset-transfer-basic/chaincode-go`

## 示例运行
将`fabric-gm`的编译结果拷贝到`bin`目录下，并cd到工作目录:
```sh
cd fabric-samples-gm/bin
cp ../../fabric-gm/release/linux-amd64/bin/* ./

cd ../test-network
# 此时工作目录: fabric-samples-gm/test-network
dir_test_network=${PWD}
```

启动fabric网络并创建通道:
```sh
# 启动网络
./network.sh up

# 创建通道
./network.sh createChannel
```

发布合约:
```sh
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
```

切换peer节点:peer0.org1.example.com，并初始化合约资产:
```sh
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"InitLedger","Args":[]}'

```

查看资产:
```sh
# 查看合约资产初始数据
peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'

# 查看指定资产
peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","asset6"]}'
```

转移指定资产，将asset6的owner改为Christopher:
```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"TransferAsset","Args":["asset6","Christopher"]}'
```

查看指定资产:
```sh
# 切换peer节点:peer0.org2.example.com(可选)
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# 查看指定资产
peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","asset6"]}'
```

至此，可以看到国密改造后的fabric能够正常运行智能合约了。

接下来检查当前通道有没有使用国密算法:
```sh
# 拉取当前通道配置，会在当前目录下生成文件: mychannel_config.block
peer channel fetch config -c mychannel
# 将mychannel_config.block转为json，可以在json中查看到hash算法使用的是SM3(搜索关键字"hash_function")
configtxlator proto_decode  --type common.Block --input mychannel_config.block > mychannel_config.json

# 检查相关证书
# 随便找一个证书，如:"test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp/signcerts/peer0.org1.example.com-cert.pem"
# 通过openssl x509命令可以查看证书内容，看到签名算法为 "SM2-with-SM3" ，证书公钥算法为 "sm2"
# 注意，需要本地openssl支持国密算法。
dir_test_network=${PWD}
cd organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp/signcerts/
openssl x509 -text -in peer0.org1.example.com-cert.pem

```

最后，关闭fabric网络:
```sh
# 关闭fabric网络
cd ${dir_test_network}
./network.sh down
```


# 版权声明
本项目采取木兰宽松许可证, 第2版，具体参见`LICENSE`文件。

本项目基于`github.com/hyperledger/fabric-samples`进行了二次开发，对应版权声明文件:`thrid_licenses/github.com/hyperledger/fabric-samples/LICENSE`。

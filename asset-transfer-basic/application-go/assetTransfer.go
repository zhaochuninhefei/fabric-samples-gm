/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gitee.com/zhaochuninhefei/fabric-sdk-go-gm/pkg/core/config"
	"gitee.com/zhaochuninhefei/fabric-sdk-go-gm/pkg/gateway"
	"gitee.com/zhaochuninhefei/zcgolog/zclog"
)

func main() {
	zclog.Infoln("============ asset-transfer-basic/application-go/assetTransfer.go 开始执行 ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		zclog.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		zclog.Fatalf("Failed to create wallet: %v", err)
	}
	// 清理钱包，确保获取最新的客户端信息
	wallet.Remove("appUser")
	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			zclog.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		zclog.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		zclog.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic")

	zclog.Infoln("--> Submit Transaction: InitLedger, 初始化合约数据")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		zclog.Errorf("Failed to Submit transaction: %v", err)
	} else {
		zclog.Infof("成功初始化合约数据 %s", string(result))
	}

	zclog.Infoln("--> Evaluate Transaction: GetAllAssets, 查看合约所有数据")
	result, err = contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		zclog.Fatalf("Failed to evaluate transaction: %v", err)
	}
	zclog.Infoln(string(result))

	zclog.Infoln("--> Evaluate Transaction: AssetExists, 检查数据asset13是否存在")
	result, err = contract.EvaluateTransaction("AssetExists", "asset13")
	if err != nil {
		zclog.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	zclog.Infoln(string(result))
	if string(result) != "true" {
		zclog.Infoln("--> Submit Transaction: CreateAsset, 创建新数据asset13")
		result, err = contract.SubmitTransaction("CreateAsset", "asset13", "yellow", "5", "Tom", "1300")
		if err != nil {
			zclog.Errorf("Failed to Submit transaction: %v", err)
		} else {
			zclog.Infoln(string(result))
		}
	}

	zclog.Infoln("--> Evaluate Transaction: ReadAsset, 读取数据asset13")
	result, err = contract.EvaluateTransaction("ReadAsset", "asset13")
	if err != nil {
		zclog.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	zclog.Infoln(string(result))

	if string(result) == "{\"ID\":\"asset13\",\"color\":\"yellow\",\"size\":5,\"owner\":\"Tom\",\"appraisedValue\":1300}" {
		zclog.Infoln("--> Submit Transaction: TransferAsset asset13, 将数据asset13的所有者改为zhaochun")
		_, err = contract.SubmitTransaction("TransferAsset", "asset13", "zhaochun")
		if err != nil {
			zclog.Fatalf("Failed to Submit transaction: %v", err)
		}
	} else {
		zclog.Infoln("--> Submit Transaction: TransferAsset asset13, 将数据asset13的所有者改为Tom")
		_, err = contract.SubmitTransaction("TransferAsset", "asset13", "Tom")
		if err != nil {
			zclog.Fatalf("Failed to Submit transaction: %v", err)
		}
	}

	zclog.Infoln("--> Evaluate Transaction: ReadAsset, 读取数据asset13")
	result, err = contract.EvaluateTransaction("ReadAsset", "asset13")
	if err != nil {
		zclog.Fatalf("Failed to evaluate transaction: %v", err)
	}
	zclog.Infoln(string(result))

	zclog.Infoln("============ asset-transfer-basic/application-go/assetTransfer.go 执行结束 ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	zclog.Infoln("============ 获取客户端钱包数据:User1@org1.example.com ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}

/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"gitee.com/zhaochuninhefei/fabric-contract-api-go-gm/contractapi"
	"gitee.com/zhaochuninhefei/fabric-samples-gm/token-erc-20/chaincode-go/chaincode"
)

func main() {
	tokenChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating token-erc-20 chaincode: %v", err)
	}

	if err := tokenChaincode.Start(); err != nil {
		log.Panicf("Error starting token-erc-20 chaincode: %v", err)
	}
}

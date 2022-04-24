/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"gitee.com/zhaochuninhefei/fabric-contract-api-go-gm/contractapi"
	"gitee.com/zhaochuninhefei/fabric-samples-gm/token-utxo/chaincode-go/chaincode"
)

func main() {
	tokenChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating token-utxo chaincode: %v", err)
	}

	if err := tokenChaincode.Start(); err != nil {
		log.Panicf("Error starting token-utxo chaincode: %v", err)
	}
}

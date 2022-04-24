/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"gitee.com/zhaochuninhefei/fabric-contract-api-go-gm/contractapi"
	auction "gitee.com/zhaochuninhefei/fabric-samples-gm/auction/chaincode-go/smart-contract"
)

func main() {
	auctionSmartContract, err := contractapi.NewChaincode(&auction.SmartContract{})
	if err != nil {
		log.Panicf("Error creating auction chaincode: %v", err)
	}

	if err := auctionSmartContract.Start(); err != nil {
		log.Panicf("Error starting auction chaincode: %v", err)
	}
}

/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/consensus"
	"chainmaker.org/chainmaker/pb-go/v2/discovery"
	"chainmaker.org/chainmaker/pb-go/v2/store"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/hokaccha/go-prettyjson"
)

const (
	createContractTimeout       = 5
	sdkConfigOrg1Client1Path    = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkPwkConfigOrg1Client1Path = "../sdk_configs/sdk_config_pwk_org1_admin1.yml"
	sdkPkConfigUser1Path        = "../sdk_configs/sdk_config_pk_user1.yml"
)

func main() {
	//testMain(sdkConfigOrg1Client1Path)
	testNativeContract()
	//testMain(sdkPwkConfigOrg1Client1Path)
	//testMain(sdkPkConfigUser1Path)
}

func testMain(sdkPath string) {
	testSystemContract(sdkPath)
	testSystemContractArchive(sdkPath)
	testGetMerklePathByTxId(sdkPath)
}

func testNativeContract() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	toAddContractList := []string{syscontract.SystemContract_DPOS_ERC20.String(),
		syscontract.SystemContract_DPOS_STAKE.String()}
	testNativeContractAccessGrant(client, false, toAddContractList, usernames)
	time.Sleep(time.Second * 3)
	testGetNativeContractInfo(client, toAddContractList[0])
	testGetNativeContractList(client)
	testNativeContractAccessRevoke(client, false, toAddContractList, usernames)
	time.Sleep(time.Second * 3)
	testGetDisabledNativeContractList(client)
}

// [系统合约]
func testSystemContract(sdkPath string) {
	client, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	genesisBlockInfo := testSystemContractGetBlockByHeight(client, 1)
	testSystemContractGetTxByTxId(client, genesisBlockInfo.Block.Txs[0].Payload.TxId)
	testSystemContractGetBlockByHash(client, hex.EncodeToString(genesisBlockInfo.Block.Header.BlockHash))
	testSystemContractGetBlockByTxId(client, genesisBlockInfo.Block.Txs[0].Payload.TxId)
	blockInfo := testSystemContractGetLastConfigBlock(client)
	fmt.Printf("GetLastConfigBlock BlockType=%s, BlockHeight=%d\n",
		blockInfo.Block.Header.BlockType, blockInfo.Block.Header.BlockHeight)
	testSystemContractGetLastBlock(client)
	testSystemContractGetChainInfo(client)

	systemChainClient, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	testSystemContractGetNodeChainList(systemChainClient)
}

func testSystemContractArchive(sdkPath string) {
	client, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	var blockHeight uint64 = 4
	fullBlock := testSystemContractGetFullBlockByHeight(client, blockHeight)
	heightByTxId := testSystemContractGetBlockHeightByTxId(client, fullBlock.Block.Txs[0].Payload.TxId)
	if blockHeight != heightByTxId {
		log.Fatalln("blockHeight != heightByTxId")
	}
	heightByHash := testSystemContractGetBlockHeightByHash(client, hex.EncodeToString(fullBlock.Block.Header.BlockHash))
	if blockHeight != heightByHash {
		log.Fatalln("blockHeight != heightByHash")
	}

	testSystemContractGetCurrentBlockHeight(client)
	testSystemContractGetArchivedBlockHeight(client)
	testSystemContractGetBlockHeaderByHeight(client)
}

func testSystemContractGetTxByTxId(client *sdk.ChainClient, txId string) *common.TransactionInfo {
	transactionInfo, err := client.GetTxByTxId(txId)
	if err != nil {
		log.Fatalln(err)
	}
	return transactionInfo
}

func testSystemContractGetBlockByHeight(client *sdk.ChainClient, blockHeight uint64) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHeight(blockHeight, true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetBlockByHash(client *sdk.ChainClient, blockHash string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHash(blockHash, true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetBlockByTxId(client *sdk.ChainClient, txId string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByTxId(txId, true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetLastConfigBlock(client *sdk.ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastConfigBlock(true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetLastBlock(client *sdk.ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastBlock(true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("last block height: %d\n", blockInfo.Block.Header.BlockHeight)
	marshal, err := prettyjson.Marshal(blockInfo)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("blockInfo: %s\n", marshal)
	return blockInfo
}

func testSystemContractGetCurrentBlockHeight(client *sdk.ChainClient) uint64 {
	height, err := client.GetCurrentBlockHeight()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("current block height: %d\n", height)
	return height
}

func testSystemContractGetArchivedBlockHeight(client *sdk.ChainClient) uint64 {
	height, err := client.GetArchivedBlockHeight()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("archived block height: %d\n", height)
	return height
}

func testSystemContractGetBlockHeightByTxId(client *sdk.ChainClient, txId string) uint64 {
	height, err := client.GetBlockHeightByTxId(txId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("txId [%s] => block height: %d\n", txId, height)
	return height
}

func testSystemContractGetBlockHeightByHash(client *sdk.ChainClient, blockHash string) uint64 {
	height, err := client.GetBlockHeightByHash(blockHash)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("blockHash [%s] => block height: %d\n", blockHash, height)
	return height
}

func testSystemContractGetChainInfo(client *sdk.ChainClient) *discovery.ChainInfo {
	chainConfig, err := client.GetChainConfig()
	if err != nil {
		log.Fatalln(err)
	}
	chainInfo := &discovery.ChainInfo{}
	if chainConfig.Consensus.Type != consensus.ConsensusType_SOLO {
		var err error
		chainInfo, err = client.GetChainInfo()
		if err != nil {
			log.Fatalln(err)
		}
	}
	return chainInfo
}

func testSystemContractGetNodeChainList(client *sdk.ChainClient) *discovery.ChainList {
	chainList, err := client.GetNodeChainList()
	if err != nil {
		log.Fatalln(err)
	}
	return chainList
}

func testSystemContractGetFullBlockByHeight(client *sdk.ChainClient, blockHeight uint64) *store.BlockWithRWSet {
	fullBlockInfo, err := client.GetFullBlockByHeight(blockHeight)
	if err != nil {
		if sdkutils.IsArchivedString(err.Error()) {
			fmt.Println("Is archived...")
		}
	}
	if err != nil {
		log.Fatalln(err)
	}
	marshal, err := prettyjson.Marshal(fullBlockInfo)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("fullBlockInfo: %s\n", marshal)
	return fullBlockInfo
}

func testSystemContractGetBlockHeaderByHeight(client *sdk.ChainClient) {
	_, err := client.GetBlockHeaderByHeight(0)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = client.GetBlockHeaderByHeight(5)
	if err != nil {
		log.Fatalln(err)
	}
}

func testGetMerklePathByTxId(sdkPath string) {
	client, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	genesisBlockInfo := testSystemContractGetBlockByHeight(client, 1)
	txId := genesisBlockInfo.Block.Txs[0].Payload.TxId

	merklePath, err := client.GetMerklePathByTxId(txId)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("GetMerklePathByTxId: ", merklePath)
}

func testNativeContractAccessGrant(client *sdk.ChainClient, withSyncResult bool,
	grantContractList []string, usernames []string) {
	payload, err := client.CreateNativeContractAccessGrantPayload(grantContractList)
	if err != nil {
		log.Fatalln(err)
	}

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	fmt.Printf("resp: %+v\n", resp)
	if err != nil {
		log.Fatalln(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("testNativeContractAccessGrant resp: %+v\n", resp)
}

func testNativeContractAccessRevoke(client *sdk.ChainClient, withSyncResult bool,
	revokeContractList []string, usernames []string) {
	payload, err := client.CreateNativeContractAccessRevokePayload(revokeContractList)
	if err != nil {
		log.Fatalln(err)
	}

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	fmt.Printf("resp: %+v\n", resp)
	if err != nil {
		log.Fatalln(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("testNativeContractAccessRevoke resp: %+v\n", resp)
}

func testGetNativeContractInfo(client *sdk.ChainClient, contractName string) {
	resp, err := client.GetContractInfo(contractName)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("testGetNativeContractInfo resp: %+v\n", resp)
}

func testGetNativeContractList(client *sdk.ChainClient) {
	resp, err := client.GetContractList()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("testGetNativeContractList resp: %+v\n", resp)
}

func testGetDisabledNativeContractList(client *sdk.ChainClient) {
	resp, err := client.GetDisabledNativeContractList()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("testGetDisabledNativeContractList resp: %+v\n", resp)
}

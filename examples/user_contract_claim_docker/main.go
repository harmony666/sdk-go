/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	createContractTimeout = 5
	claimVersion          = "1.0.0"

	claimContractName = "claim_docker_001"
	//claimContractName     = "claim_zxl_002"

	claimByteCodePath = "../../testdata/claim-docker-demo/claim_docker_001.7z"
	//claimByteCodePath     = "../../testdata/claim-docker-demo/claim_zxl_002.7z"

	sdkConfigPKUser1Path     = "../sdk_configs/sdk_config_pk_user1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testUserContractClaim()
}

func testUserContractClaim() {
	fmt.Println("====================== create client ======================")
	//client, err := examples.CreateChainClientWithSDKConf(sdkConfigPKUser1Path)
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 创建合约 ======================")
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	//usernames := []string{examples.UserNameOrg1Admin1}
	txId := testUserContractClaimCreate(client, true, true, usernames...)

	tx := testContractGetTxByTxId(client, txId)
	fmt.Printf("result code:%d, msg:%s\n", tx.Transaction.Result.Code, tx.Transaction.Result.Code.String())
	fmt.Printf("contract result code:%d, msg:%s\n",
		tx.Transaction.Result.ContractResult.Code, tx.Transaction.Result.ContractResult.Message)

	fmt.Println("====================== 调用合约 ======================")
	fileHash, err := testUserContractClaimInvoke(client, "invoke_contract", true)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 执行合约查询接口 ======================")
	//txId := "1cbdbe6106cc4132b464185ea8275d0a53c0261b7b1a470fb0c3f10bd4a57ba6"
	//fileHash = txId[len(txId)/2:]
	kvs := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("findByFileHash"),
		},
		{
			Key:   "file_hash",
			Value: []byte(fileHash),
		},
	}
	testUserContractClaimQuery(client, "invoke_contract", kvs)

	// ====================== create client ======================
	//====================== 创建合约 ======================
	//CREATE claim contract resp: message:"OK" contract_result:<result:"\n\020claim_docker_001\022\0051.0.0\030\006*<\n\026wx-org1.chainmaker.org\020\001\032 \226n\021\033rX\325\3153\260\213\003Vb\310 ,\300\037\364\276\240\376\300\315=u\303#\251l/" message:"OK" > tx_id:"3cee2d4e5b1a4f13a919c0e8a977b4730225be0ad48646d8827b4fcd77628338"
	//====================== 调用合约 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:result:"file_1631947501c51875808b1f42099a19067b7f11e526" message:"Success" contract_event:<topic:"topic_vx" tx_id:"1f69aaf2396145359888e3d658681fa2125c66e46f7c41c7b6c0eb5cdca6093d" contract_name:"claim_docker_001" contract_version:"1.0.0" event_data:"c51875808b1f42099a19067b7f11e526" event_data:"file_1631947501" > ]
	//====================== 执行合约查询接口 ======================
	//QUERY claim contract resp: message:"SUCCESS" contract_result:<result:"{\"FileHash\":\"c51875808b1f42099a19067b7f11e526\",\"FileName\":\"file_1631947501\",\"Time\":1631947501}" message:"Success" > tx_id:"1e002ba323a14f48bd1c90ddac402f5cf29c6722aa724521bdedd747652ce9ed"
}

func testUserContractClaimCreate(client *sdk.ChainClient, withSyncResult bool, isIgnoreSameContract bool, usernames ...string) string {

	resp, err := createUserContract(client, claimContractName, claimVersion, claimByteCodePath,
		common.RuntimeType_DOCKER_GO, []*common.KeyValuePair{}, withSyncResult, usernames...)
	if err != nil {
		if !isIgnoreSameContract {
			log.Fatalln(err)
		} else {
			fmt.Printf("CREATE claim contract failed, err: %s, resp: %+v\n", err, resp)
		}
	} else {
		fmt.Printf("CREATE claim contract success, resp: %+v\n", resp)
	}

	if resp != nil {
		return resp.TxId
	}

	return ""
}

func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string, runtime common.RuntimeType,
	kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	payload = client.AttachGasLimit(payload, &common.Limit{
		GasLimit: 60000000,
	})

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func testUserContractClaimInvoke(client *sdk.ChainClient,
	method string, withSyncResult bool) (string, error) {

	curTime := strconv.FormatInt(time.Now().Unix(), 10)

	fileHash := uuid.GetUUID()
	kvs := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("save"),
		},
		{
			Key:   "time",
			Value: []byte(curTime),
		},
		{
			Key:   "file_hash",
			Value: []byte(fileHash),
		},
		{
			Key:   "file_name",
			Value: []byte(fmt.Sprintf("file_%s", curTime)),
		},
	}

	err := invokeUserContract(client, claimContractName, method, "", kvs, withSyncResult, &common.Limit{GasLimit: 200000})
	if err != nil {
		return "", err
	}

	return fileHash, nil
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string,
	kvs []*common.KeyValuePair, withSyncResult bool, limit *common.Limit) error {

	resp, err := client.InvokeContractWithLimit(contractName, method, txId, kvs, -1, withSyncResult, limit)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}

func testUserContractClaimQuery(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) {
	resp, err := client.QueryContract(claimContractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("QUERY claim contract resp: %+v\n", resp)
}

func testContractGetTxByTxId(client *sdk.ChainClient, txId string) *common.TransactionInfo {
	transactionInfo, err := client.GetTxByTxId(txId)
	if err != nil {
		log.Fatalln(err)
	}
	return transactionInfo
}

/*
Copyright (C) BABEC. All rights reserved.
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
	// "chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"

	"crypto/sha1"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const (
	createContractTimeout = 5
	claimContractName     = "claimnew"
	claimVersion          = "2.0.0"
	claimByteCodePath     = "../../testdata/claim-wasm-demo/chainmaker_contract2.wasm"

	sdkConfigOrg1Client1Path = "../../examples/sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	engine := gin.Default()

	engine.GET("/save", func(context *gin.Context) {
		file_content := context.Query("file_content")

		fmt.Println("====================== create client ======================")
		client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("====================== 调用合约 ======================")
		fileHash, err := testUserContractClaimInvoke(client, "save", true, file_content)
		fmt.Println(fileHash)
		context.Writer.Write([]byte("filehash:" + fileHash))
		if err != nil {
			log.Fatalln(err)
		}

	})

	engine.GET("/find", func(context *gin.Context) {
		file_hash := context.Query("file_hash")

		fmt.Println("====================== create client ======================")
		client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("====================== 执行合约查询接口 ======================")
		//txId := "1cbdbe6106cc4132b464185ea8275d0a53c0261b7b1a470fb0c3f10bd4a57ba6"
		//fileHash = txId[len(txId)/2:]
		kvs := []*common.KeyValuePair{
			{
				Key:   "file_hash",
				Value: []byte(file_hash),
			},
		}
		resp := testUserContractClaimQuery(client, "find_by_file_hash", kvs)

		context.Writer.Write([]byte((resp.ContractResult.Result)))
	})
	engine.Run()
	// testUserContractClaim()
}

func testUserContractClaim() {
	fmt.Println("====================== create client ======================")
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println("====================== 创建合约 ======================")
	// usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	// testUserContractClaimCreate(client, true, usernames...)

	fmt.Println("====================== 调用合约 ======================")
	fileHash, err := testUserContractClaimInvoke(client, "save", true, "123456")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 执行合约查询接口 ======================")
	//txId := "1cbdbe6106cc4132b464185ea8275d0a53c0261b7b1a470fb0c3f10bd4a57ba6"
	//fileHash = txId[len(txId)/2:]
	kvs := []*common.KeyValuePair{
		{
			Key:   "file_hash",
			Value: []byte(fileHash),
		},
	}
	testUserContractClaimQuery(client, "find_by_file_hash", kvs)

	//====================== 创建合约 ======================
	//CREATE claim contract resp: message:"OK" contract_result:<result:"\n\010claim001\022\0052.0.0\030\002*<\n\026wx-org1.chainmaker.org\020\001\032 $p^\215Q\366\236\2120\007\233eW\210\220\3746\250\027\331h\212\024\253\370Ecl\214J'\322" message:"OK" > tx_id:"e40e126cf093472bbb1c80cbd9e6c18ef64e0f8e276046a38f7cc98df1d0cba7"
	//====================== 调用合约 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:gas_used:14538222 ]
	//====================== 执行合约查询接口 ======================
	//QUERY claim contract resp: message:"SUCCESS" contract_result:<result:"{\"file_hash\":\"8f4c3500833040919ea63bfe1059e117\",\"file_name\":\"file_2021-07-20 19:47:24\",\"time\":\"2021-07-20 19:47:24\"}" gas_used:24597022 > tx_id:"154d3f1bb53d432098de1664b5dbdbfa1e1420cdb4634bd3ba92431ce037ca29"
}

func testUserContractClaimCreate(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {

	resp, err := createUserContract(client, claimContractName, claimVersion, claimByteCodePath,
		common.RuntimeType_WASMER, []*common.KeyValuePair{}, withSyncResult, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE claim contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string, runtime common.RuntimeType,
	kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func testUserContractClaimInvoke(client *sdk.ChainClient,
	method string, withSyncResult bool, file_content string) (string, error) {

	curTime := strconv.FormatInt(time.Now().Unix(), 10)

	// fileHash := uuid.GetUUID()

	fileHash := Sha1(file_content, "file_"+curTime)
	fmt.Printf(fileHash)

	kvs := []*common.KeyValuePair{
		{
			Key:   "file_content",
			Value: []byte(file_content),
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

	err := invokeUserContract(client, claimContractName, method, "", kvs, withSyncResult)
	if err != nil {
		return "", err
	}

	return fileHash, nil
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string,
	kvs []*common.KeyValuePair, withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
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

func testUserContractClaimQuery(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) *common.TxResponse {
	resp, err := client.QueryContract(claimContractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("QUERY claim contract resp: %+v\n", resp)
	return resp
}

func Sha1(data1 string, data2 string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data1))
	sha1.Write([]byte(data2))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

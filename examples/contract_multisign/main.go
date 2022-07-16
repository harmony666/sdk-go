/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	contractByteCodePath = "../../testdata/claim-wasm-demo/rust-fact-2.0.0.wasm"
	contractName         = "claim010"
	contractVersion      = "v1.0.0"

	sdkConfigPKUser1Path = "../sdk_configs/sdk_config_pk_user1.yml"

	needEndorserCount = 1
)

var (
	contractRuntimeType = common.RuntimeType_WASMER.String()
)

func main() {
	cc, err := examples.CreateChainClientWithSDKConf(sdkConfigPKUser1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 发送线上多签部署合约的交易 ======================")
	kvs := newContractInitPairs() //构造交易发起 kv pairs
	payload := cc.CreateMultiSignReqPayload(kvs)
	resp, err := cc.MultiSignContractReq(payload)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("send MultiSignContractReq resp: %+v\n", resp)

	fmt.Println("====================== 各链管理员开始投票 ======================")
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[cc.GetHashType()],
		cc.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatalln(err)
	}

	for i, e := range endorsers {
		fmt.Printf("====================== 链管理员 %d 投票 ======================\n", i)
		time.Sleep(3 * time.Second)

		resp, err = cc.MultiSignContractVote(payload, e, true)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("send MultiSignContractVote resp: %+v\n", resp)

		time.Sleep(3 * time.Second)
		fmt.Println("====================== 查询本多签交易的投票情况 ======================")

		//resp, err = cc.MultiSignContractQuery(payload.TxId)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//fmt.Printf("query MultiSignContractQuery resp: %+v\n", resp)

		//查看详细投票信息
		//multiSignInfo := &syscontract.MultiSignInfo{}
		//err = proto.Unmarshal(resp.ContractResult.Result, multiSignInfo)
		//if err != nil {
		//	fmt.Printf(" multiSignInfoDB Unmarshal error: %s", err)
		//}
		//fmt.Printf("multi sign query finished, vote info: %v", multiSignInfo.VoteInfos)

		result, err := cc.GetSyncResult(resp.TxId)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("result code:%d, msg:%s\n", result.Code, result.Code.String())
		fmt.Printf("contract result code:%d, msg:%s\n",
			result.ContractResult.Code, result.ContractResult.Message)

		if i+1 == needEndorserCount {
			break
		}
	}

	fmt.Println("====================== 调用上边创建的合约 ======================")
	time.Sleep(time.Second * 3)
	fileHash := invokeUserContract(cc, true)

	fmt.Println("====================== 执行合约查询接口 ======================")
	kvs = []*common.KeyValuePair{
		{
			Key:   "file_hash",
			Value: []byte(fileHash),
		},
	}
	testUserContractClaimQuery(cc, "find_by_file_hash", kvs)
}

func newContractInitPairs() []*common.KeyValuePair {
	wasmBin, err := ioutil.ReadFile(contractByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}
	return []*common.KeyValuePair{
		{
			Key:   syscontract.MultiReq_SYS_CONTRACT_NAME.String(),
			Value: []byte(syscontract.SystemContract_CONTRACT_MANAGE.String()),
		},
		{
			Key:   syscontract.MultiReq_SYS_METHOD.String(),
			Value: []byte(syscontract.ContractManageFunction_INIT_CONTRACT.String()),
		},
		{
			Key:   syscontract.InitContract_CONTRACT_NAME.String(),
			Value: []byte(contractName),
		},
		{
			Key:   syscontract.InitContract_CONTRACT_VERSION.String(),
			Value: []byte(contractVersion),
		},
		{
			Key:   syscontract.InitContract_CONTRACT_BYTECODE.String(),
			Value: wasmBin,
		},
		{
			Key:   syscontract.InitContract_CONTRACT_RUNTIME_TYPE.String(),
			Value: []byte(contractRuntimeType),
		},
	}
}

func invokeUserContract(client *sdk.ChainClient, withSyncResult bool) string {

	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	fileHash := uuid.GetUUID()
	kvs := []*common.KeyValuePair{
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

	resp, err := client.InvokeContractWithLimit(contractName, "save", "", kvs, -1, withSyncResult,
		&common.Limit{GasLimit: 100000000})
	if err != nil {
		log.Fatalln(err)
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		log.Fatalln(fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message))
	}

	fmt.Printf("query invokeUserContract resp: %+v\n", resp)

	//if !withSyncResult {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	//} else {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	//}

	return fileHash
}

func testUserContractClaimQuery(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) {
	resp, err := client.QueryContract(contractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("QUERY claim contract resp: %+v\n", resp)
}

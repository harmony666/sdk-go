/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	createContractTimeout = 5
	contractName          = "counter001"
	version               = "1.0.0"
	upgradeVersion        = "2.0.0"
	byteCodePath          = "../../testdata/counter-go-demo/rust-counter-1.0.0.wasm"
	upgradeByteCodePath   = "../../testdata/counter-go-demo/rust-counter-2.0.0.wasm"

	sdkConfigOrg1Client1Path    = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkPwkConfigOrg1Client1Path = "../sdk_configs/sdk_config_pwk_org1_admin1.yml"
	sdkPkConfigUser1Path        = "../sdk_configs/sdk_config_pk_user1.yml"
)

func main() {
	testUserContractCounterGo(sdkConfigOrg1Client1Path)
	//testUserContractCounterGo(sdkPwkConfigOrg1Client1Path)
	//testUserContractCounterGo(sdkPkConfigUser1Path)
}

func testUserContractCounterGo(sdkPath string) {
	client, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 创建合约（异步）======================")
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	testUserContractCounterGoCreate(client, false, usernames...)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用合约（异步）======================")
	testUserContractCounterGoInvoke(client, "increase", nil, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约查询接口1 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	fmt.Println("====================== 冻结合约 ======================")
	testUserContractCounterGoFreeze(client, false, usernames...)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口2 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	fmt.Println("====================== 解冻合约 ======================")
	testUserContractCounterGoUnfreeze(client, false, usernames...)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口3 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	//fmt.Println("====================== 吊销合约 ======================")
	//testUserContractCounterGoRevoke(client, admin1, admin2, admin3, admin4, false)
	//time.Sleep(5 * time.Second)
	//fmt.Println("====================== 执行合约查询接口 ======================")
	//testUserContractCounterGoQuery(client, "query", nil)

	fmt.Println("====================== 调用合约（同步）======================")
	testUserContractCounterGoInvoke(client, "increase", nil, true)

	fmt.Println("====================== 执行合约查询接口 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	fmt.Println("====================== 调用升级合约的接口(报错) ======================")
	testUserContractCounterGoInvoke(client, "add_two", nil, true)

	fmt.Println("====================== 升级合约 ======================")
	testUserContractCounterGoUpgrade(client, false, usernames...)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用新接口 ======================")
	testUserContractCounterGoInvoke(client, "add_two", nil, true)

	fmt.Println("====================== 执行合约查询接口 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	//====================== 创建合约（异步）======================
	//CREATE counter-go contract resp: message:"OK" tx_id:"383cf7ca18f6459d8cc643f021db16f61fd0df87ee3b46d787634f5a28a06e62"
	//====================== 调用合约（异步）======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[txId:5452383de99a471999c1530c9c777377ff27a705bafb40a5ba9b72688ab88eec]
	//====================== 执行合约查询接口1 ======================
	//QUERY counter-go contract resp: message:"SUCCESS" contract_result:<result:"1" gas_used:11541125 > tx_id:"fa8cddc2cd3b4b81a465e7d01ce4804ad6acc3336a17403ea6c755f8c7e54916"
	//====================== 冻结合约 ======================
	//resp: message:"OK" tx_id:"2e94e47369ad45c588ac712c18598f28f976e54a716941e38f03cfdb35c5a0ce"
	//Freeze counter-go contract resp: message:"OK" tx_id:"2e94e47369ad45c588ac712c18598f28f976e54a716941e38f03cfdb35c5a0ce"
	//====================== 执行合约查询接口2 ======================
	//QUERY counter-go contract resp: code:CONTRACT_FREEZE_FAILED message:"txStatusCode:36, resultCode:1, contractName[counter001] method[query] txType[QUERY_CONTRACT], failed to run user contract, counter001 has been frozen." contract_result:<code:1 message:"failed to run user contract, counter001 has been frozen." > tx_id:"f2c749f260b140af9e3c0a1cab2c4e966c7e7f823ddd4aaab88f355e7f457e01"
	//====================== 解冻合约 ======================
	//unfreeze resp: message:"OK" tx_id:"137c39c5e3dc4e96a5847c9536c9c0c292e38bf459f74d7eb65202a5e08ebb9d"
	//Unfreeze counter-go contract resp: message:"OK" tx_id:"137c39c5e3dc4e96a5847c9536c9c0c292e38bf459f74d7eb65202a5e08ebb9d"
	//====================== 执行合约查询接口3 ======================
	//QUERY counter-go contract resp: message:"SUCCESS" contract_result:<result:"1" gas_used:11541125 > tx_id:"aec5a93436f84ef69c949e03f6e57b8038490ae5d509460c91fae7f4f3d2c4bf"
	//====================== 调用合约（同步）======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:result:"++ stone success count=2" gas_used:14081908 ]
	//====================== 执行合约查询接口 ======================
	//QUERY counter-go contract resp: message:"SUCCESS" contract_result:<result:"2" gas_used:11543981 > tx_id:"5f141e2ef34b4047892c0addecdd0d81cf8948d7bef4404ea50a89f933af25af"
	//====================== 调用升级合约的接口(报错) ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:code:1 message:"contract[counter001] invoke failed, method [add_two] not export" gas_used:71537 ]
	//====================== 升级合约 ======================
	//UPGRADE counter-go contract resp: message:"OK" tx_id:"66177733a6944147b30bde9510c1ae11b8f2e803bce947e3877dae1459d4ce65"
	//====================== 调用新接口 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:result:"++ stone success count=4" gas_used:14065288 ]
	//====================== 执行合约查询接口 ======================
	//QUERY counter-go contract resp: message:"SUCCESS" contract_result:<result:"4" gas_used:11488459 > tx_id:"549c1b3e4cf3447785a6185e81d15ef51bf689aa52d3485792427257d429d1cc"
}

// [用户合约]
func testUserContractCounterGoCreate(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {

	resp, err := createUserContract(client, contractName, version, byteCodePath, common.RuntimeType_WASMER,
		[]*common.KeyValuePair{}, withSyncResult, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

// 更新合约
func testUserContractCounterGoUpgrade(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {
	payload, err := client.CreateContractUpgradePayload(contractName, upgradeVersion, upgradeByteCodePath,
		common.RuntimeType_WASMER, []*common.KeyValuePair{})
	if err != nil {
		log.Fatalln(err)
	}

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, -1, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("UPGRADE counter-go contract resp: %+v\n", resp)
}

// 冻结合约
func testUserContractCounterGoFreeze(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {
	payload, err := client.CreateContractFreezePayload(contractName)
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

	fmt.Printf("Freeze counter-go contract resp: %+v\n", resp)
}

// 解冻合约
func testUserContractCounterGoUnfreeze(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {
	payload, err := client.CreateContractUnfreezePayload(contractName)
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
	fmt.Printf("unfreeze resp: %+v\n", resp)
	if err != nil {
		log.Fatalln(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Unfreeze counter-go contract resp: %+v\n", resp)
}

// 吊销合约
func testUserContractCounterGoRevoke(client *sdk.ChainClient, withSyncResult bool, usernames ...string) {
	payload, err := client.CreateContractRevokePayload(contractName)
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
	fmt.Printf("revoke resp: %+v\n", resp)
	if err != nil {
		log.Fatalln(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("revoke counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoInvoke(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair,
	withSyncResult bool) {

	err := invokeUserContract(client, contractName, method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserContractCounterGoQuery(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) {
	resp, err := client.QueryContract(contractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string,
	runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

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
	if err != nil {
		return nil, err
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, kvs []*common.KeyValuePair, withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.TxId)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%+v]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}

func invokeUserContractStepByStep(client *sdk.ChainClient, contractName, method, txId string,
	kvs []*common.KeyValuePair, withSyncResult bool) error {
	req, err := client.GetTxRequest(contractName, method, "", kvs)
	if err != nil {
		return err
	}

	resp, err := client.SendTxRequest(req, -1, withSyncResult)
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

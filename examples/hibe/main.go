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
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

// test contract functionName
const (
	// save Hibe Message
	saveHibeMsg = "save_hibe_msg"

	// save params
	saveHibeParams = "save_hibe_params"

	// find params by ogrId
	findParamsByOrgId = "find_params_by_org_id"

	createContractTimeout = 5

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	//sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	//sdkConfigOrg2Admin1Path  = "../sdk_configs/sdk_config_org2_admin1.yml"
	//sdkConfigOrg3Admin1Path  = "../sdk_configs/sdk_config_org3_admin1.yml"
	//sdkConfigOrg4Admin1Path  = "../sdk_configs/sdk_config_org4_admin1.yml"
	//sdkConfigOrg5Admin1Path  = "../sdk_configs/sdk_config_org5_admin1.yml"
)

// test data
const (
	hibeContractByteCodePath = "../../testdata/hibe-wasm-demo/contract-hibe.wasm"

	hibeContractName = "contracthibe10000005"

	// 本地 hibe params 文件路径
	localHibeParamsFilePath = "../../testdata/hibe-data/wx-org1.chainmaker.org/wx-org1.chainmaker.org.params"

	// 测试源消息
	msg = "这是一条HIBE测试存证 ✔✔✔"

	// hibe_msg 的消息 Id
	bizId2 = "1234567890123452"

	// Id 和 对应私钥文件路径 这里测试3组
	localTopLevelId                 = "wx-topL"
	localTopLevelHibePrvKeyFilePath = "../../testdata/hibe-data/wx-org1.chainmaker.org/privateKeys/wx-topL.privateKey"

	localSecondLevelId                 = "wx-topL/secondL"
	localSecondLevelHibePrvKeyFilePath = "../../testdata/hibe-data/wx-org1.chainmaker.org/privateKeys/wx-topL#secondL.privateKey"

	localThirdLevelId                 = "wx-topL/secondL/thirdL"
	localThirdLevelHibePrvKeyFilePath = "../../testdata/hibe-data/wx-org1.chainmaker.org/privateKeys/wx-topL#secondL#thirdL.privateKey"
)

var txId = ""

func main() {
	testHibeContractCounterGo()
}

func testHibeContractCounterGo() {

	txId = utils.GetTimestampTxId()
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 创建合约（异步）======================")
	testUserHibeContractCounterGoCreate(client, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用合约 params 上链 （异步）======================")
	testUserHibeContractParamsGoInvoke(client, saveHibeParams, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约 params 查询接口 ======================")
	testUserHibeContractParamsGoQuery(client, findParamsByOrgId, nil)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用合约 加密数据上链（异步）======================")
	testUserHibeContractMsgGoInvoke(client, saveHibeMsg, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约 加密数据查询接口 ======================")
	testUserHibeContractMsgGoQuery(client)
}

// 创建Hibe合约
func testUserHibeContractCounterGoCreate(client *sdk.ChainClient, admin1, admin2, admin3,
	admin4 string, withSyncResult bool) {
	resp, err := createUserHibeContract(client, admin1, admin2, admin3, admin4,
		hibeContractName, examples.Version, hibeContractByteCodePath, common.RuntimeType_GASM, []*common.KeyValuePair{}, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE contract-hibe-1 contract resp: %+v\n", resp)
}

// 调用Hibe合约
// params 上链
func testUserHibeContractParamsGoInvoke(client *sdk.ChainClient, method string, withSyncResult bool) {
	err := invokeUserHibeContractParams(client, hibeContractName, method, "", withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

// params 查询
func testUserHibeContractParamsGoQuery(client *sdk.ChainClient, method string, params map[string]string) {
	hibeParams, err := client.QueryHibeParamsWithOrgId(hibeContractName, findParamsByOrgId, examples.OrgId1, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY %s contract resp -> hibeParams:%s\n", hibeContractName, hibeParams)
}

// 加密数据上链
func testUserHibeContractMsgGoInvoke(client *sdk.ChainClient, method string, withSyncResult bool) {
	err := invokeUserHibeContractMsg(client, hibeContractName, method, txId, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

// 获取加密数据
func testUserHibeContractMsgGoQuery(client *sdk.ChainClient) {
	//keyType := crypto.AES
	keyType := crypto.SM4

	localParams, err := utils.ReadHibeParamsWithFilePath(localHibeParamsFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	topHibePrvKey, err := utils.ReadHibePrvKeysWithFilePath(localTopLevelHibePrvKeyFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	secondHibePrvKey, err := utils.ReadHibePrvKeysWithFilePath(localSecondLevelHibePrvKeyFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	thirdHibePrvKey, err := utils.ReadHibePrvKeysWithFilePath(localThirdLevelHibePrvKeyFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	msgBytes1, err := client.DecryptHibeTxByTxId(localTopLevelId, localParams, topHibePrvKey, txId, keyType)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY hibe-contract-go-1 contract resp DecryptHibeTxByBizId [Decrypt Msg By TopLevel privateKey] message: %s\n", string(msgBytes1))

	msgBytes2, err := client.DecryptHibeTxByTxId(localSecondLevelId, localParams, secondHibePrvKey, txId, keyType)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY hibe-contract-go-1 contract resp DecryptHibeTxByBizId [Decrypt Msg By SecondLevel privateKey] message: %s\n", string(msgBytes2))

	msgBytes3, err := client.DecryptHibeTxByTxId(localThirdLevelId, localParams, thirdHibePrvKey, txId, keyType)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY hibe-contract-go-1 contract resp DecryptHibeTxByBizId [Decrypt Msg By ThirdLevel privateKey] message: %s\n", string(msgBytes3))

}

func createUserHibeContract(client *sdk.ChainClient, admin1, admin2, admin3, admin4 string,
	contractName, version, byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	//endorsers, err := examples.GetEndorsers(payload, admin1, admin2, admin3, admin4)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, admin1, admin2, admin3, admin4)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	// TODO: ??
	//err = examples.CheckProposalRequestResp(resp, true)
	//if err != nil {
	//	return nil, err
	//}

	return resp, nil
}

func invokeUserHibeContractParams(client *sdk.ChainClient, contractName, method, txId string,
	withSyncResult bool) error {
	localParams, err := utils.ReadHibeParamsWithFilePath(localHibeParamsFilePath)
	if err != nil {
		return err
	}
	payloadParams, err := client.CreateHibeInitParamsTxPayloadParams(examples.OrgId1, localParams)

	// resp, err := client.InvokeContract(contractName, method, txId, payloadParams, -1, withSyncResult)
	resp, err := client.InvokeContract(contractName, method, txId, payloadParams, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	// TODO: ??
	//if !withSyncResult {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	//} else {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	//}

	return nil
}

func invokeUserHibeContractMsg(client *sdk.ChainClient, contractName, method, txId string, withSyncResult bool) error {
	receiverId := make([]string, 3)
	receiverId[0] = localSecondLevelId
	receiverId[1] = localThirdLevelId
	receiverId[2] = localTopLevelId

	// fetch orgId []string from receiverId []string
	org := make([]string, len(receiverId))
	org[0] = "wx-org1.chainmaker.org"
	org[1] = "wx-org1.chainmaker.org"
	org[2] = "wx-org1.chainmaker.org"

	// query params
	var paramsBytesList [][]byte
	for _, id := range org {
		hibeParamsBytes, err := client.QueryHibeParamsWithOrgId(hibeContractName, findParamsByOrgId, id, -1)
		if err != nil {
			//t.Logf("QUERY hibe-contract-go-1 contract resp: %+v\n", hibeParams)
			return fmt.Errorf("client.QueryHibeParamsWithOrgId(hibeContractName, id, -1) failed, err: %v\n", err)
		}

		if len(hibeParamsBytes) == 0 {
			return fmt.Errorf("no souch params of %s's org, please check it", id)
		}

		paramsBytesList = append(paramsBytesList, hibeParamsBytes)
	}

	//keyType := crypto.AES
	keyType := crypto.SM4
	params, err := client.CreateHibeTxPayloadParamsWithHibeParams([]byte(msg), receiverId, paramsBytesList, txId, keyType)
	if err != nil {
		return err
	}

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	// TODO: ??
	//if !withSyncResult {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	//} else {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	//}

	return nil
}

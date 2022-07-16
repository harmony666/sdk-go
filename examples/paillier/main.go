/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"testing"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/paillier"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
	"github.com/stretchr/testify/require"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"

	createContractTimeout = 5
)

const (

	// go 合约
	paillierContractName = "pailliergo100001"
	paillierByteCodePath = "../../testdata/paillier-wasm-demo/contract-paillier.wasm"
	runtime              = common.RuntimeType_GASM

	// rust 合约
	//paillierContractName = "paillier-rust-10001"
	//paillierByteCodePath = "./testdata/counter-go-demo/chainmaker_contract.wasm"
	//runtime              = common.RuntimeType_WASMER

	paillierPubKeyFilePath = "../../testdata/paillier-key/test1.pubKey"
	paillierPrvKeyFilePath = "../../testdata/paillier-key/test1.prvKey"
)

func main() {
	TestPaillierContractCounterGo()
}

func TestPaillierContractCounterGo() {
	t := new(testing.T)
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	require.Nil(t, err)

	fmt.Println("======================================= 创建合约（异步）=======================================")
	testPaillierCreate(client, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1, false)
	time.Sleep(5 * time.Second)

	fmt.Println("======================================= 调用合约运算（异步）=======================================")
	testPaillierOperation(client, "paillier_test_set", false)
	time.Sleep(5 * time.Second)

	fmt.Println("======================================= 查询结果并解密（异步）=======================================")
	testPaillierQueryResult(t, client, "paillier_test_get")
}

func testPaillierCreate(client *sdk.ChainClient, admin1, admin2, admin3,
	admin4 string, withSyncResult bool) {
	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		paillierContractName, examples.Version, paillierByteCodePath, runtime, []*common.KeyValuePair{}, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE contract-hibe-1 contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, admin1, admin2, admin3, admin4 string,
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

// 调用合约进行同态运算
func testPaillierOperation(client *sdk.ChainClient, s string, b bool) {
	pubKeyBytes, err := ioutil.ReadFile(paillierPubKeyFilePath)
	//require.Nil(t, err)
	if err != nil {
		log.Fatalln(err)
	}
	payloadParams, err := CreatePaillierTransactionPayloadParams(pubKeyBytes, 1, 1000000)
	resp, err := client.InvokeContract(paillierContractName, s, "", payloadParams, -1, b)
	//require.Nil(t, err)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		fmt.Printf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	//if !b {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	//} else {
	//	fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	//}
}

func CreatePaillierTransactionPayloadParams(pubKeyBytes []byte, plaintext1, plaintext2 int64) ([]*common.KeyValuePair, error) {
	pubKey := new(paillier.PubKey)
	err := pubKey.Unmarshal(pubKeyBytes)
	if err != nil {

	}

	pt1 := new(big.Int).SetInt64(plaintext1)
	ciphertext1, err := pubKey.Encrypt(pt1)
	if err != nil {
		return nil, err
	}
	ct1Bytes, err := ciphertext1.Marshal()
	if err != nil {
		return nil, err
	}

	prv2, _ := paillier.GenKey()

	pub2, _ := prv2.GetPubKey()
	_, _ = pub2.Marshal()

	pt2 := new(big.Int).SetInt64(plaintext2)
	ciphertext2, err := pubKey.Encrypt(pt2)
	if err != nil {
		return nil, err
	}
	ct2Bytes, err := ciphertext2.Marshal()
	if err != nil {
		return nil, err
	}

	ct1Str := base64.StdEncoding.EncodeToString(ct1Bytes)
	ct2Str := base64.StdEncoding.EncodeToString(ct2Bytes)

	payloadParams := []*common.KeyValuePair{
		{
			Key:   "handletype",
			Value: []byte("SubCiphertext"),
		},
		{
			Key:   "para1",
			Value: []byte(ct1Str),
		},
		{
			Key:   "para2",
			Value: []byte(ct2Str),
		},
		{
			Key:   "pubkey",
			Value: pubKeyBytes,
		},
	}
	/*
		old
		payloadParams := make(map[string]string)
		//payloadParams["handletype"] = "AddCiphertext"
		//payloadParams["handletype"] = "AddPlaintext"
		payloadParams["handletype"] = "SubCiphertext"
		//payloadParams["handletype"] = "SubCiphertextStr"
		//payloadParams["handletype"] = "SubPlaintext"
		//payloadParams["handletype"] = "NumMul"

		payloadParams["para1"] = ct1Str
		payloadParams["para2"] = ct2Str
		payloadParams["pubkey"] = string(pubKeyBytes)
	*/

	return payloadParams, nil
}

// 查询同态执行结果并解密
func testPaillierQueryResult(t *testing.T, c *sdk.ChainClient, s string) {
	//paillierMethod := "AddCiphertext"
	//paillierMethod := "AddPlaintext"
	paillierMethod := "SubCiphertext"
	//paillierMethod := "SubCiphertextStr"
	//paillierMethod := "SubPlaintext"
	params1, err := QueryPaillierResult(c, paillierContractName, s, paillierMethod, -1, paillierPrvKeyFilePath)
	require.Nil(t, err)
	fmt.Printf("QUERY %s contract resp -> encrypt(cipher 10): %d\n", paillierContractName, params1)
}

func QueryPaillierResult(c *sdk.ChainClient, contractName, method, paillierDataItemId string, timeout int64, paillierPrvKeyPath string) (int64, error) {

	resultStr, err := QueryPaillierResultById(c, contractName, method, paillierDataItemId, timeout)
	if err != nil {
		return 0, err
	}

	ct := new(paillier.Ciphertext)
	resultBytes, err := base64.StdEncoding.DecodeString(string(resultStr))
	if err != nil {
		return 0, err
	}
	err = ct.Unmarshal(resultBytes)
	if err != nil {
		return 0, err
	}

	prvKey := new(paillier.PrvKey)

	prvKeyBytes, err := ioutil.ReadFile(paillierPrvKeyPath)
	if err != nil {
		return 0, fmt.Errorf("open paillierKey file failed, [err:%s]", err)
	}

	err = prvKey.Unmarshal(prvKeyBytes)
	if err != nil {
		return 0, err
	}

	decrypt, err := prvKey.Decrypt(ct)
	if err != nil {
		return 0, err
	}

	return decrypt.Int64(), nil
}

func QueryPaillierResultById(c *sdk.ChainClient, contractName, method, paillierMethod string, timeout int64) ([]byte, error) {
	pairs := []*common.KeyValuePair{
		{Key: "handletype", Value: []byte(paillierMethod)},
	}

	/*
		old
		pairsMap := make(map[string]string)
		pairsMap["handletype"] = paillierMethod
	*/

	resp, err := c.QueryContract(contractName, method, pairs, timeout)
	if err != nil {
		return nil, err
	}

	result := resp.ContractResult.Result

	return result, nil
}

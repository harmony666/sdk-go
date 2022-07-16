/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	createContractTimeout = 5
	tokenContractName     = "token001"
	tokenVersion          = "1.0.0"
	tokenByteCodePath     = "../../testdata/token-evm-demo/token.bin"
	tokenABIPath          = "../../testdata/token-evm-demo/token.abi"
	amount                = 200

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

var client1AddrInt, client2AddrInt, client1AddrSki, client1EthAddr, client2EthAddr string

func init() {
	userClient1, err := examples.GetUser(examples.UserNameOrg1Client1)
	if err != nil {
		log.Fatalln(err)
	}

	client1AddrInt, client1EthAddr, client1AddrSki, err = examples.MakeAddrAndSkiFromCrtFilePath(userClient1.SignCrtPath)
	if err != nil {
		log.Fatalln(err)
	}

	userClient2, err := examples.GetUser(examples.UserNameOrg2Client1)
	if err != nil {
		log.Fatalln(err)
	}

	client2AddrInt, client2EthAddr, _, err = examples.MakeAddrAndSkiFromCrtFilePath(userClient2.SignCrtPath)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("client1AddrInt: %s\nclient1EthAddr: %s\nclient1AddrSki: %s\nclient2AddrInt: %s\nclient2EthAddr: %s\n",
		client1AddrInt, client1EthAddr, client1AddrSki, client2AddrInt, client2EthAddr)
}

func main() {
	testUserContractTokenEVM(sdkConfigOrg1Client1Path)
}

func testUserContractTokenEVM(sdkPath string) {
	fmt.Println("====================== create client ======================")
	client, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 创建Token合约,给client1地址分配初始代币 ======================")
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	testUserContractTokenEVMCreate(client, true, true, usernames...)

	fmt.Println("====================== 查看余额 ======================")
	testUserContractTokenEVMBalanceOf(client, client1AddrInt, true)
	testUserContractTokenEVMBalanceOf(client, client2AddrInt, true)

	fmt.Println("====================== client1给client2地址转账 ======================")
	testUserContractTokenEVMTransfer(client, amount, true)

	fmt.Println("====================== 查看余额 ======================")
	testUserContractTokenEVMBalanceOf(client, client1AddrInt, true)
	testUserContractTokenEVMBalanceOf(client, client2AddrInt, true)

	//====================== 创建Token合约,给client1地址分配初始代币 ======================
	//CREATE EVM token contract resp: message:"OK" contract_result:<result:"\n(a90b890d4b2577b3d8c9bb487fcb5e90167e4842\022\0051.0.0\030\005*<\n\026wx-org1.chainmaker.org\020\001\032 $p^\215Q\366\236\2120\007\233eW\210\220\3746\250\027\331h\212\024\253\370Ecl\214J'\322" message:"OK" > tx_id:"df90acc322724e1e852313f510bbbf81dbb65774ac5c4aa4a83a7258f002d8eb"
	//====================== 查看余额 ======================
	//addr [1018109374098032500766612781247089211099623418384] => [100000000000000000]
	//addr [1317892642413437150535769048733130623036570974971] => [0]
	//====================== client1给client2地址转账 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:result:"\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\001" gas_used:19916 contract_event:<topic:"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" tx_id:"e45f546577994c0eaec0c33a6281cd6219fac062ecea4defb7411772c582c4c7" contract_name:"a90b890d4b2577b3d8c9bb487fcb5e90167e4842" contract_version:"1.0.0" event_data:"000000000000000000000000b2559a706dbab942671d5978ea97a13c82af6a10" event_data:"000000000000000000000000e6d859965e3dd4ef88c99130350b44b2ae9fbafb" event_data:"00000000000000000000000000000000000000000000000000000000000000c8" > ]
	//====================== 查看余额 ======================
	//addr [1018109374098032500766612781247089211099623418384] => [99999999999999800]
	//addr [1317892642413437150535769048733130623036570974971] => [200]
}

func testUserContractTokenEVMCreate(client *sdk.ChainClient, withSyncResult bool, isIgnoreSameContract bool, usernames ...string) {

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	// 方式1: 16进制
	//addr := evmutils.BigToAddress(evmutils.FromHexString(client1Addr[2:]))
	// 方式2: Int
	//addr := evmutils.BigToAddress(evmutils.FromDecimalString(client1AddrInt))
	// 方式3: ski
	addrInt, err := evmutils.MakeAddressFromHex(client1AddrSki)
	if err != nil {
		log.Fatalln(err)
	}
	addr := evmutils.BigToAddress(addrInt)

	dataByte, err := myAbi.Pack("", addr)
	if err != nil {
		log.Fatalln(err)
	}

	data := hex.EncodeToString(dataByte)
	pairs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(data),
		},
	}

	byteCode, err := ioutil.ReadFile(tokenByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}

	//bc := string(byteCode)
	//bc = strings.TrimSpace(bc)

	resp, err := createUserContract(client, examples.CalcContractName(tokenContractName), tokenVersion, string(byteCode),
		common.RuntimeType_EVM, pairs, withSyncResult, usernames...)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM token contract resp: %+v\n", resp)
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

	// 发送创建合约请求
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

func testUserContractTokenEVMTransfer(client *sdk.ChainClient, amount int64, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	//addr := evmutils.StringToAddress(client2Addr)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(client2AddrInt))

	dataByte, err := myAbi.Pack("transfer", addr, big.NewInt(amount))
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, examples.CalcContractName(tokenContractName), method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserContractTokenEVMBalanceOf(client *sdk.ChainClient, address string, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(tokenABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	//addr := evmutils.StringToAddress(address)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))

	methodName := "balanceOf"
	dataByte, err := myAbi.Pack(methodName, addr)
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	result, err := invokeUserContractWithResult(client, examples.CalcContractName(tokenContractName), method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	balance, err := myAbi.Unpack(methodName, result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("addr [%s] => %d\n", address, balance)
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
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}

func invokeUserContractWithResult(client *sdk.ChainClient, contractName, method, txId string,
	kvs []*common.KeyValuePair, withSyncResult bool) ([]byte, error) {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
	if err != nil {
		return nil, err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	return resp.ContractResult.Result, nil
}

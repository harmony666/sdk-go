/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	queryAddr             = "query_address"
	createContractTimeout = 5
	issueLimit            = "100000000"
	issueSupply           = "100000000"

	sdkConfigOrg2Client1Path = "../sdk_configs/sdk_config_org2_client1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

var (
	assetContractName = "asset001"
	assetVersion      = "2.0.0"
	assetByteCodePath = "../../testdata/asset-wasm-demo/rust-asset-2.0.0.wasm"
)

func main() {
	testUserContractAsset()
	testUserContractAssetBalanceOf()
}

func testUserContractAsset() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	client2, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 1)安装钱包合约 ======================")
	pairs := []*common.KeyValuePair{
		{
			Key:   "issue_limit",
			Value: []byte(issueLimit),
		},
		{
			Key:   "total_supply",
			Value: []byte(issueSupply),
		},
	}

	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	testUserContractAssetCreate(client, pairs, true, false, usernames...)

	fmt.Println("====================== 2)注册另一个用户 ======================")
	testUserContractAssetInvokeRegister(client2, "register", true)

	fmt.Println("====================== 3)查询钱包地址 ======================")
	addr1 := testUserContractAssetQuery(client, queryAddr, nil)
	fmt.Printf("client1 address: %s\n", addr1)
	addr2 := testUserContractAssetQuery(client2, queryAddr, nil)
	fmt.Printf("client2 address: %s\n", addr2)

	fmt.Println("====================== 4)给用户分别发币100000 ======================")
	amount := "100000"
	testUserContractAssetInvoke(client, "issue_amount", amount, addr1, true)
	testUserContractAssetInvoke(client, "issue_amount", amount, addr2, true)

	fmt.Println("====================== 5)分别查看余额 ======================")
	getBalance(client, addr1)
	getBalance(client, addr2)

	fmt.Println("====================== 6)A给B转账100 ======================")
	amount = "100"
	testUserContractAssetInvoke(client, "transfer", amount, addr2, true)

	fmt.Println("====================== 7)再次分别查看余额 ======================")
	getBalance(client, addr1)
	getBalance(client, addr2)

	// ====================== 1)安装钱包合约 ======================
	//CREATE asset contract resp: message:"OK" contract_result:<result:"\n\010asset001\022\0052.0.0\030\002*<\n\026wx-org1.chainmaker.org\020\001\032 $p^\215Q\366\236\2120\007\233eW\210\220\3746\250\027\331h\212\024\253\370Ecl\214J'\322" message:"OK" > tx_id:"3174f3876a084c0aa701f07e7b3e78ad058ab83932af4303b113685c17d2a124"
	//====================== 2)注册另一个用户 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:result:"20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f" gas_used:15045866 ]
	//====================== 3)查询钱包地址 ======================
	//QUERY asset contract [query_address] resp: message:"SUCCESS" contract_result:<result:"4d2b2301e06ca9269361fce6105296cc00ee19ffaa6a5f5b37b4c7faf8889697" gas_used:19349195 > tx_id:"2b8fa49241204f92aa9b9a64814d912ba9e91d27f2614ab2859b5816b0a2c82b"
	//client1 address: 4d2b2301e06ca9269361fce6105296cc00ee19ffaa6a5f5b37b4c7faf8889697
	//QUERY asset contract [query_address] resp: message:"SUCCESS" contract_result:<result:"20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f" gas_used:19349195 > tx_id:"f307af3ca50948eca99cab947373d3ca1907b13463284dddb60e4a84d102f54d"
	//client2 address: 20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f
	//====================== 4)给用户分别发币100000 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:gas_used:47322933 ]
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:gas_used:47350108 ]
	//====================== 5)分别查看余额 ======================
	//QUERY asset contract [balance_of] resp: message:"SUCCESS" contract_result:<result:"100000" gas_used:15530495 > tx_id:"acfa53cbf6e241819c0fb7aa4271481f048f7ab5c92643b6ba825645c5ff1027"
	//client [4d2b2301e06ca9269361fce6105296cc00ee19ffaa6a5f5b37b4c7faf8889697] balance: 100000
	//QUERY asset contract [balance_of] resp: message:"SUCCESS" contract_result:<result:"100000" gas_used:15530495 > tx_id:"34fae5285fe242428f5a4bc359fdb8914b15a096e4a047fab9a24957df559222"
	//client [20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f] balance: 100000
	//====================== 6)A给B转账100 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:result:"ok" gas_used:30310634 ]
	//====================== 7)再次分别查看余额 ======================
	//QUERY asset contract [balance_of] resp: message:"SUCCESS" contract_result:<result:"99900" gas_used:15741987 > tx_id:"780123c71508467ca8c5490d886709b1460acda905b14d1d971c41c5ed1a3c40"
	//client [4d2b2301e06ca9269361fce6105296cc00ee19ffaa6a5f5b37b4c7faf8889697] balance: 99900
	//QUERY asset contract [balance_of] resp: message:"SUCCESS" contract_result:<result:"100100" gas_used:15781095 > tx_id:"7405f9895304408d9be642a3e2fe7242cf97ef3e53e847498e5295433b6eca1d"
	//client [20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f] balance: 100100
	//====================== 1)查询钱包地址 ======================
	//QUERY asset contract [query_address] resp: message:"SUCCESS" contract_result:<result:"4d2b2301e06ca9269361fce6105296cc00ee19ffaa6a5f5b37b4c7faf8889697" gas_used:19349528 > tx_id:"e8345c199706453ebb7014fd88da98ce5e3a39b1693b4404bef8f25886adb8e0"
	//client1 address: 4d2b2301e06ca9269361fce6105296cc00ee19ffaa6a5f5b37b4c7faf8889697
	//QUERY asset contract [query_address] resp: message:"SUCCESS" contract_result:<result:"20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f" gas_used:19349528 > tx_id:"e76429e857a049a49057e19f21f5471544802c092fa64ffd9c667e836a4fcfdc"
	//client2 address: 20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f
	//====================== 2)查询钱包余额 ======================
	//QUERY asset contract [balance_of] resp: message:"SUCCESS" contract_result:<result:"99900" gas_used:15508529 > tx_id:"0fc793354de1422bb646aab7a8246d23fa143cb9935e4090a50c41ca68fa9b8f"
	//client [4d2b2301e06ca9269361fce6105296cc00ee19ffaa6a5f5b37b4c7faf8889697] balance: 99900
	//QUERY asset contract [balance_of] resp: message:"SUCCESS" contract_result:<result:"100100" gas_used:15547637 > tx_id:"fb5932fb618e466480a560f3cbba3a33e217eca0e2254a3b91d60e0a6f4c9d43"
	//client [20ea85650709f2e533b531d4a3e992312406b54c5171d431afcc6f555125073f] balance: 100100
}

func testUserContractAssetBalanceOf() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	client2, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 1)查询钱包地址 ======================")
	addr1 := testUserContractAssetQuery(client, queryAddr, nil)
	fmt.Printf("client1 address: %s\n", addr1)
	addr2 := testUserContractAssetQuery(client2, queryAddr, nil)
	fmt.Printf("client2 address: %s\n", addr2)

	fmt.Println("====================== 2)查询钱包余额 ======================")
	getBalance(client, addr1)
	getBalance(client, addr2)
}

func testUserContractAssetCreate(client *sdk.ChainClient, kvs []*common.KeyValuePair, withSyncResult bool,
	isIgnoreSameContract bool, usernames ...string) {

	resp, err := createUserContract(client, assetContractName, assetVersion, assetByteCodePath,
		common.RuntimeType_WASMER, kvs, withSyncResult, usernames...)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE asset contract resp: %+v\n", resp)
}

func testUserContractAssetInvokeRegister(client *sdk.ChainClient, method string, withSyncResult bool) {
	err := invokeUserContract(client, assetContractName, method, "", nil, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserContractAssetQuery(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) string {
	resp, err := client.QueryContract(assetContractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY asset contract [%s] resp: %+v\n", method, resp)

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		log.Fatalln(err)
	}
	return string(resp.ContractResult.Result)
}

func testUserContractAssetInvoke(client *sdk.ChainClient, method string, amount, addr string, withSyncResult bool) {
	kvs := []*common.KeyValuePair{
		{
			Key:   "amount",
			Value: []byte(amount),
		},
		{
			Key:   "to",
			Value: []byte(addr),
		},
	}

	err := invokeUserContract(client, assetContractName, method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func getBalance(client *sdk.ChainClient, addr string) {
	kvs := []*common.KeyValuePair{
		{
			Key:   "owner",
			Value: []byte(addr),
		},
	}

	balance := testUserContractAssetQuery(client, "balance_of", kvs)

	fmt.Printf("client [%s] balance: %s\n", addr, balance)
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

func createUserContract(client *sdk.ChainClient, contractName, version,
	byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair,
	withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

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

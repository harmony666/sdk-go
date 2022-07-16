package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	createContractTimeout = 5
	claimContractName     = "zxlcc"
	claimVersion          = "1.0.2"

	//claimByteCodePath = "/Users/suxiongye/Code/chainmaker-sdk-go/testdata/zxl_demo/zxldemo.7z"
	claimByteCodePath = "../../testdata/zxl_demo/zxlcc.7z"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	fmt.Println("====================== create client ======================")
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}
	createOrUpgradeContract(client, false)
	var pointKey = "testPoint"
	// testGet(client, pointKey)
	startTime := time.Now()
	var concurrentNum = 10000
	var roundNum = 5

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		totalTxNum := 0
		currHeight, err := client.GetCurrentBlockHeight()
		if err != nil {
			panic(err)
		}
		for {
			newHeight, err := client.GetCurrentBlockHeight()
			if err != nil {
				panic(err)
			}
			if newHeight == currHeight+1 {
				currHeight++
				newHeader, err := client.GetBlockHeaderByHeight(newHeight)
				if err != nil {
					panic(err)
				}
				totalTxNum += int(newHeader.TxCount)
				if totalTxNum >= concurrentNum*roundNum {
					duration := time.Since(startTime)
					fmt.Printf("finished time : %+v, tps: %f\n", duration, float64(concurrentNum*roundNum)/(float64(duration)/1e9))
					// testGet(client, pointKey)
					wg.Done()
					return
				}
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	// 第一层控制并发
	for i := 0; i < concurrentNum; i++ {
		wg.Add(1)
		// 第二层控制时间
		go func(i int) {
			for j := 0; j < roundNum; j++ {
				testAddPoint(client, fmt.Sprintf("%s-%d-%d", pointKey, i, j))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// 创建或更新合约
func createOrUpgradeContract(client *sdk.ChainClient, updateTag bool) {
	var resp *common.TxResponse
	var err error
	fmt.Println("====================== 创建合约 ======================")
	//usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1,
	//	examples.UserNameOrg4Admin1}
	usernames := []string{examples.UserNameOrg1Admin1}
	if updateTag {
		resp, err = _upgradeContract(client, claimContractName, claimVersion, claimByteCodePath,
			common.RuntimeType_DOCKER_GO, []*common.KeyValuePair{}, true, usernames...)
	} else {
		resp, err = _createContract(client, claimContractName, claimVersion, claimByteCodePath,
			common.RuntimeType_DOCKER_GO, []*common.KeyValuePair{}, true, usernames...)
	}
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("CREATE claim contract resp: %+v\n", resp)
}

// 创建合约
func _createContract(client *sdk.ChainClient, contractName, version, byteCodePath string,
	runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool,
	usernames ...string) (*common.TxResponse, error) {
	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}
	endorsers, err := examples.GetEndorsers(payload, usernames...)
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

// 更新合约
func _upgradeContract(client *sdk.ChainClient, contractName, version, byteCodePath string,
	runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool,
	usernames ...string) (*common.TxResponse, error) {
	payload, err := client.CreateContractUpgradePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}
	endorsers, err := examples.GetEndorsers(payload, usernames...)
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

// 查询用户合约
func testAuth(client *sdk.ChainClient) {
	kvsAuth := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("auth"),
		},
		{
			Key:   "appId",
			Value: []byte("190815000210001"),
		},
	}
	query(client, "invoke_contract", kvsAuth)
}

// 绑定用户合约
func testBind(client *sdk.ChainClient) error {
	kvsBind := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("bind"),
		},
		{
			Key:   "appId",
			Value: []byte("190815000210001"),
		}, {
			Key:   "pubKey",
			Value: []byte("04f9728bacace1b8572c4814ee2db0422f45e51ff873fa1003e50380157de29a34235c2be1ed15df6776360839a710f3a069e3cf79ac68fc2de1dc2998c5d40a14"),
		}, {
			Key:   "signature",
			Value: []byte("3046022100b7ed41cb7d2501105cf9eeed804c70d7e0f425b66cd10393b60b5472174d4b58022100e1d9aa52ffb16ca15569f0cae4b1399b44e3020cc3f4df9b243468b54e0fac88"),
		},
	}
	err := invoke(client, claimContractName, "invoke_contract", "", kvsBind, true)
	if err != nil {
		return err
	}
	return nil
}

// 增加积分
func testAddPoint(client *sdk.ChainClient, key string) error {
	kvsBind := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("addPoint"),
		},
		{
			Key:   "key",
			Value: []byte(key),
		}, {
			Key:   "value",
			Value: []byte("10"),
		},
	}
	err := invoke(client, claimContractName, "invoke_contract", "", kvsBind, false)
	if err != nil {
		return err
	}
	return nil
}

// 更新
func testUpdate(client *sdk.ChainClient) error {
	var value = make([]byte, 2000)
	rand.Read(value)
	kvsBind := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("updateKey"),
		},
		{
			Key:   "key",
			Value: []byte("nftTest"),
		}, {
			Key:   "value",
			Value: value,
		},
	}
	err := invoke(client, claimContractName, "invoke_contract", "", kvsBind, true)
	if err != nil {
		return err
	}
	return nil
}

// 绑定用户合约
func testAdd(client *sdk.ChainClient) error {
	kvsBind := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("add"),
		},
		{
			Key:   "appId",
			Value: []byte("190815000210001"),
		}, {
			Key:   "evId",
			Value: []byte("evev"),
		}, {
			Key:   "signature",
			Value: []byte("3046022100b7ed41cb7d2501105cf9eeed804c70d7e0f425b66cd10393b60b5472174d4b58022100e1d9aa52ffb16ca15569f0cae4b1399b44e3020cc3f4df9b243468b54e0fac88"),
		},
		{
			Key:   "evHash",
			Value: []byte("190815000210001"),
		}, {
			Key:   "ext",
			Value: []byte("ext"),
		}, {
			Key:   "timestamp",
			Value: []byte("3046022100b7ed41cb7d2501105cf9eeed804c70d7e0f425b66cd10393b60b5472174d4b58022100e1d9aa52ffb16ca15569f0cae4b1399b44e3020cc3f4df9b243468b54e0fac88"),
		}, {
			Key:   "time",
			Value: []byte("3046022100b7ed41cb7d2501105cf9eeed804c70d7e0f425b66cd10393b60b5472174d4b58022100e1d9aa52ffb16ca15569f0cae4b1399b44e3020cc3f4df9b243468b54e0fac88"),
		},
	}
	err := invoke(client, claimContractName, "invoke_contract", "", kvsBind, true)
	if err != nil {
		return err
	}
	return nil
}

// 查询用户合约
func testGet(client *sdk.ChainClient, key string) {
	kvsAuth := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("get"),
		},
		{
			Key:   "key",
			Value: []byte(key),
		},
	}
	query(client, "invoke_contract", kvsAuth)
}

// 查询用户合约
func testGetNftAddrTime(client *sdk.ChainClient) int {
	kvsAuth := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("get"),
		},
		{
			Key:   "key",
			Value: []byte("nftTest"),
		},
	}
	return queryNum(client, "invoke_contract", kvsAuth)
}

func invoke(client *sdk.ChainClient, contractName, method, txId string,
	kvs []*common.KeyValuePair, withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		//fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message,
		//	resp.TxId)
	} else {
		//fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code,
		//resp.Message, resp.ContractResult)
	}

	return nil
}

func query(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) {
	resp, err := client.QueryContract(claimContractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY claim contract resp: %s\n", string(resp.ContractResult.Result))
}

func queryNum(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) int {
	resp, err := client.QueryContract(claimContractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}
	var addrMap map[string]string
	var i = 0
	json.Unmarshal(resp.ContractResult.Result, &addrMap)
	for _, v := range addrMap {
		if strings.Contains(v, "newAddr") {
			i++
		}
	}
	return i
}

// 创建nft
func testNewNft(client *sdk.ChainClient) {
	kvsAuth := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("newNft"),
		},
		{
			Key:   "key",
			Value: []byte("nftTest"),
		},
	}
	invoke(client, claimContractName, "invoke_contract", "", kvsAuth, true)
}

// 更新nft
func testUpdateNft(client *sdk.ChainClient, index, addr []byte) {
	kvsAuth := []*common.KeyValuePair{
		{
			Key:   "method",
			Value: []byte("updateNft"),
		},
		{
			Key:   "key",
			Value: []byte("nftTest"),
		},
		{
			Key:   "index",
			Value: index,
		},
		{
			Key:   "addr",
			Value: addr,
		},
	}
	invoke(client, claimContractName, "invoke_contract", "", kvsAuth, false)
}

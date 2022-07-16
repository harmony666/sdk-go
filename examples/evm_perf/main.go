/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	createContractTimeout = 5
	balanceContractName   = "balance001"
	balanceVersion        = "1.0.0"
	balanceByteCodePath   = "../../testdata/balance-evm-demo/ledger_balance.bin"
	balanceABIPath        = "../../testdata/balance-evm-demo/ledger_balance.abi"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

var client1AddrInt, client2AddrInt, client1EthAddr, client2EthAddr string
var timeSum int64
var successNum int64
var totalNum int64

var randSize = 2

//var randSize = int(^uint(0) >> 1)

func init() {
	userClient1, err := examples.GetUser(examples.UserNameOrg1Client1)
	if err != nil {
		log.Fatalln(err)
	}

	client1AddrInt, client1EthAddr, _, err = examples.MakeAddrAndSkiFromCrtFilePath(userClient1.SignCrtPath)
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

	fmt.Printf("client1AddrInt: %s\nclient1EthAddr: %s\nclient2AddrInt: %s\nclient2EthAddr: %s\n",
		client1AddrInt, client1EthAddr, client2AddrInt, client2EthAddr)

}

func main() {
	flag.Parse()
	fmt.Println("====================== create client ======================")
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}
	usernames := []string{examples.UserNameOrg1Admin1}

	fmt.Println("====================== create contract ======================")
	testUserContractBalanceEVMCreate(client, true, true, usernames...)
	time.Sleep(time.Second * 5)

	fmt.Println("====================== start perf test ======================")
	for i := 0; i < execTimes; i++ {
		// fmt.Printf("rand size: %d\n", randSize)
		timeSum = 0
		successNum = 0
		totalNum = 0
		txsMap = sync.Map{}
		finishChan := make(chan struct{}, 1)
		perf(client, finishChan)
		<-finishChan
		printStats()
		// randSize *= 2
	}
}

func perf(client *sdk.ChainClient, done chan struct{}) {
	switch method {
	case "increaseBalance":
		invokePerfTest(client, "increaseBalance", done)
	//case "query":
	//	fmt.Println("====================== 查询合约perf测试 ======================")
	//queryPerfTest(client, done)
	default:
		panic("invalid contract method, must be 'create|invoke|query'")
	}
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func invokePerfTest(client *sdk.ChainClient, method string, done chan struct{}) {
	startTime := time.Now().UnixNano()
	totalTxNum := 0
	currHeight, err := client.GetCurrentBlockHeight()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		for true {
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
				if totalTxNum >= threadNum*loopNum {
					timeSum = time.Now().UnixNano() - startTime
					wg.Done()
					return
				}
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()
	runParallel(func(threadIndex int) {
		txId := sdkutils.GetTimestampTxId()
		preHandler(txId)
		token := IntToBytes(rand.Intn(randSize))
		ski := hex.EncodeToString(token)
		addrInt, err := evmutils.MakeAddressFromHex(ski)
		if err != nil {
			log.Fatalln(err)
		}
		testUserContractBalanceEVMIncreaseBalance(client, addrInt.String(), method, false)
		postHandler(txId, nil)
	}, done)
	timeSum = time.Now().UnixNano() - startTime

	wg.Wait()
}

func runParallel(f func(threadIndex int), done chan struct{}) {
	var wg sync.WaitGroup
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		index := i
		go func() {
			for j := 0; j < loopNum; j++ {
				//time.Sleep(time.Duration(interval) * time.Second)
				//fmt.Printf("loop %d", j)
				f(index)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	done <- struct{}{}
}

func preHandler(id string) {
	txsMap.Store(id, &transaction{requestId: id, startTime: time.Now().UnixNano()})
}

func postHandler(id string, err error) {
	val, ok := txsMap.Load(id)
	if !ok {
		return
	}
	tx, ok := val.(*transaction)
	if !ok {
		return
	}
	if err == nil {
		tx.success = true
	}
	tx.endTime = time.Now().UnixNano()
}

func printStats() {
	txsMap.Range(func(key, val interface{}) bool {
		tx, ok := val.(*transaction)
		if !ok {
			return false
		}
		if tx.success {
			successNum++
			//timeSum += tx.endTime - tx.startTime
		}
		totalNum++
		return true
	})
	fmt.Printf("total num: %d\tsuccess num: %d\tfaild num: %d\ttps: %f\n", totalNum, successNum,
		totalNum-successNum, float64(successNum)/(float64(timeSum)/1e9))
}

func testUserContractBalanceEVMCreate(client *sdk.ChainClient, withSyncResult bool, isIgnoreSameContract bool, usernames ...string) {

	byteCode, err := ioutil.ReadFile(balanceByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := createUserContract(client, examples.CalcContractName(balanceContractName), balanceVersion,
		string(byteCode), common.RuntimeType_EVM, nil, withSyncResult, usernames...)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM balance contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string,
	runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	endorsers, err := examples.GetEndorsers(payload, usernames...)
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

func testUserContractBalanceEVMIncreaseBalance(client *sdk.ChainClient, address string, method string, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(balanceABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	//addr := evmutils.StringToAddress(address)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))

	dataByte, err := myAbi.Pack(method, addr)
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	methodStr := dataString[0:8]

	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, examples.CalcContractName(balanceContractName), methodStr, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, kvs []*common.KeyValuePair, withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	return nil
}

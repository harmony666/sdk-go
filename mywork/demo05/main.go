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

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"

	"crypto/sha1"
	"encoding/hex"

	"github.com/gin-gonic/gin"

	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

const (
	createContractTimeout = 5
	claimContractName     = "claim003"
	claimVersion          = "2.0.0"
	claimByteCodePath     = "../../testdata/claim-wasm-demo/chainmaker_contract01.wasm"

	sdkConfigOrg1Client1Path = "../../examples/sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	engine := gin.Default()

	engine.POST("/save", func(context *gin.Context) {
		file_content := context.PostForm("file_content")
		key := context.PostForm("key")
		cert := context.GetHeader("Authorization")

		bodyAsByteArray, _ := ioutil.ReadAll(context.Request.Body)
		jsonBody := string(bodyAsByteArray)
		fmt.Println(jsonBody)

		// 提取header,对cert进行解码
		decodeBytes, err := base64.StdEncoding.DecodeString(cert)
		if err != nil {
			log.Fatalln(err)
		}

		certBlock, _ := pem.Decode([]byte(string(decodeBytes)))
		if certBlock == nil {
			fmt.Println("1")
			return
		}
		// fmt.Println(certBlock)
		certBody, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			fmt.Println(err)
			return
		}

		ou := certBody.Subject.OrganizationalUnit[0]
		id := certBody.Subject.CommonName
		org := certBody.Subject.Organization[0]

		time_stamp := strconv.FormatInt(time.Now().Unix(), 10)

		fmt.Println(ou)
		fmt.Println(id)
		fmt.Println(org)
		fmt.Println("====================== create client ======================")
		client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("====================== 调用合约 ======================")
		fileHash, err := testUserContractClaimInvoke(client, "save", true, file_content, key, id, ou, org, time_stamp)
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
}

func testUserContractClaimInvoke(client *sdk.ChainClient,
	method string, withSyncResult bool, file_content string, key string, id string, ou string, org string, time_stamp string) (string, error) {

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
			Key:   "key",
			Value: []byte(key),
		},
		{
			Key:   "ou",
			Value: []byte(ou),
		},
		{
			Key:   "id",
			Value: []byte(id),
		},
		{
			Key:   "org",
			Value: []byte(org),
		},
		{
			Key:   "time_stamp",
			Value: []byte(time_stamp),
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

/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testChainClientCheckNewBlockChainConfig()
}

func testChainClientCheckNewBlockChainConfig() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}
	err = client.CheckNewBlockChainConfig()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("check new block chain config: ok")
}

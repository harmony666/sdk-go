/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

func main() {
	testGetEVMAddressFromCertPath()
}

func testGetEVMAddressFromCertPath() {
	userOrg1Client1, err := examples.GetUser(examples.UserNameOrg1Client1)
	if err != nil {
		log.Fatalln(err)
	}

	addrInt, err := sdk.GetEVMAddressFromCertPath(userOrg1Client1.SignCrtPath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("client1 addrInt: %s\n", addrInt)
}

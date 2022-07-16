/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testCertHash()
	testCertManage()
}

func testCertHash() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 用户证书上链 ======================")
	certHash := testCertAdd(client)
	time.Sleep(3 * time.Second)

	fmt.Println("====================== 用户证书查询 ======================")
	certInfos := testQueryCert(client, []string{certHash})
	if len(certInfos.CertInfos) != 1 {
		log.Fatalln("require len(certInfos.CertInfos) == 1")
	}

	fmt.Println("====================== 用户证书删除 ======================")
	testDeleteCert(client, []string{certHash})
	time.Sleep(3 * time.Second)

	fmt.Println("====================== 再次查询用户证书 ======================")
	certInfos = testQueryCert(client, []string{certHash})
	if len(certInfos.CertInfos) != 1 {
		log.Fatalln("require len(certInfos.CertInfos) == 1")
	}
	if certInfos.CertInfos[0].Cert != nil {
		log.Fatalln("require certInfos.CertInfos[0].Cert == nil")
	}
}

func testCertManage() {
	// org2 client证书
	var certs = []string{
		"-----BEGIN CERTIFICATE-----\nMIICiDCCAi6gAwIBAgIDCuSTMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMi5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcyLmNoYWlubWFrZXIub3JnMB4XDTIwMTExNjA2NDYwNFoXDTI1\nMTExNTA2NDYwNFowgZAxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcyLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjbGllbnQxKzApBgNVBAMTImNsaWVudDEudGxzLnd4LW9yZzIu\nY2hhaW5tYWtlci5vcmcwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQjsmPDqPjx\nikMpRPkmWH8RFgUXwpzwaoMF9OQY6sAty2U8Q6TPlafMbm/xBls//UPZpi5uhwTv\neunkar0HqfvRo3sweTAOBgNVHQ8BAf8EBAMCAaYwDwYDVR0lBAgwBgYEVR0lADAp\nBgNVHQ4EIgQgjqe9Y2WHp+WC/GfKlvwummg3xvKPi9hbDja0QVFKa/EwKwYDVR0j\nBCQwIoAgmZcrtYWpTzN56LDZdqiHah3fG5w0kPaLoEBtyC8GfaEwCgYIKoZIzj0E\nAwIDSAAwRQIgbz8Du0bvtlWVJfBFzUamyfY2OodQDGBbKnr/eFXNeIECIQDnnJs5\nAX2NCT42Be3et+jhwxshehNsYm3WOOdTq/y+yg==\n-----END CERTIFICATE-----\n",
	}

	// org2 client证书的CRL
	var certCrl = "-----BEGIN CRL-----\nMIIBXTCCAQMCAQEwCgYIKoZIzj0EAwIwgYoxCzAJBgNVBAYTAkNOMRAwDgYDVQQI\nEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcyLmNo\nYWlubWFrZXIub3JnMRIwEAYDVQQLEwlyb290LWNlcnQxIjAgBgNVBAMTGWNhLnd4\nLW9yZzIuY2hhaW5tYWtlci5vcmcXDTIxMDEyMTA2NDYwM1oXDTIxMDEyMTEwNDYw\nM1owFjAUAgMK5JMXDTI0MDMyMzE1MDMwNVqgLzAtMCsGA1UdIwQkMCKAIJmXK7WF\nqU8zeeiw2Xaoh2od3xucNJD2i6BAbcgvBn2hMAoGCCqGSM49BAMCA0gAMEUCIEgb\nQsHoMkKAKAurOUUfAJpb++DYyxXS3zhvSWPxIUPWAiEAyLSd4TgB9PbSgHyGzS5D\nU1knUTu/4HKTol6GuzmV0Kg=\n-----END CRL-----"

	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 用户证书冻结 ======================")
	testCertManageFrozen(client, certs)

	fmt.Println("====================== 用户证书解冻 ======================")
	testCertManageUnfrozen(client, certs)

	fmt.Println("====================== 用户证书吊销 ======================")
	testCertManageRevoke(client, certCrl)
}

func testCertManageFrozen(client *sdk.ChainClient, certs []string) {
	payload := client.CreateCertManageFrozenPayload(certs)

	//endorsers, err := examples.GetEndorsers(payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SendCertManageRequest(payload, endorsers, -1, true)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("frozen resp: %+v\n", resp)
}

func testCertManageUnfrozen(client *sdk.ChainClient, certs []string) {
	payload := client.CreateCertManageUnfrozenPayload(certs)

	//endorsers, err := examples.GetEndorsers(payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SendCertManageRequest(payload, endorsers, -1, true)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("unfrozen resp: %+v\n", resp)
}

func testCertManageRevoke(client *sdk.ChainClient, certCrl string) {
	payload := client.CreateCertManageRevocationPayload(certCrl)

	//endorsers, err := examples.GetEndorsers(payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SendCertManageRequest(payload, endorsers, -1, true)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("cert revoke resp: %+v\n", resp)
}

func testCertAdd(client *sdk.ChainClient) string {
	resp, err := client.AddCert()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("add cert resp: %+v\n", resp)
	return hex.EncodeToString(resp.ContractResult.Result)
}

func testQueryCert(client *sdk.ChainClient, certHashes []string) *common.CertInfos {
	certInfos, err := client.QueryCert(certHashes)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("query cert resp: %+v\n", certInfos)
	return certInfos
}

func testDeleteCert(client *sdk.ChainClient, certHashes []string) {
	pairs := []*common.KeyValuePair{
		{
			Key:   "cert_hashes",
			Value: []byte(strings.Join(certHashes, ",")),
		},
	}

	payload := client.CreateCertManagePayload(syscontract.CertManageFunction_CERTS_DELETE.String(), pairs)

	//endorsers, err := examples.GetEndorsers(payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SendCertManageRequest(payload, endorsers, -1, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("delete cert resp: %+v\n", resp)
}

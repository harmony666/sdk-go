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
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	sdkConfigPath = "../sdk_configs/sdk_config_org1_client1.yml"
	certAlias     = "mycertalias"
	newCertPEM    = "-----BEGIN CERTIFICATE-----\nMIICiTCCAi+gAwIBAgIDA+zYMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMi5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcyLmNoYWlubWFrZXIub3JnMB4XDTIwMTIwODA2NTM0M1oXDTI1\nMTIwNzA2NTM0M1owgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcyLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjbGllbnQxLDAqBgNVBAMTI2NsaWVudDEuc2lnbi53eC1vcmcy\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEZd92CJez\nCiOMzLSTrJfX5vIUArCycg05uKru2qFaX0uvZUCwNxbfSuNvkHRXE8qIBUhTbg1Q\nR9rOlfDY1WfgMaN7MHkwDgYDVR0PAQH/BAQDAgGmMA8GA1UdJQQIMAYGBFUdJQAw\nKQYDVR0OBCIEICfLatSyyebzRsLbnkNKZJULB2bZOtG+88NqvAHCsXa3MCsGA1Ud\nIwQkMCKAIPGP1bPT4/Lns2PnYudZ9/qHscm0pGL6Kfy+1CAFWG0hMAoGCCqGSM49\nBAMCA0gAMEUCIQDzHrEHrGNtoNfB8jSJrGJU1qcxhse74wmDgIdoGjvfTwIgabRJ\nJNvZKRpa/VyfYi3TXa5nhHRIn91ioF1dQroHQFc=\n-----END CERTIFICATE-----\n"
)

func main() {
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfigPath),
		sdk.WithChainClientAlias(certAlias),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== add alias ======================")
	resp, err := cc.AddAlias()
	if err != nil {
		log.Fatal(err)
	}
	examples.PrintPrettyJson(resp)

	fmt.Println("====================== query alias ======================")
	aliasInfos, err := cc.QueryCertsAlias([]string{certAlias})
	if err != nil {
		log.Fatal(err)
	}
	examples.PrintPrettyJson(aliasInfos)

	fmt.Println("====================== update cert by alias ======================")
	fmt.Printf("newCertPEM=%s", newCertPEM)
	updateAliasPayload := cc.CreateUpdateCertByAliasPayload(certAlias, newCertPEM)

	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[cc.GetHashType()],
		cc.GetAuthType(), updateAliasPayload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatal(err)
	}
	resp2, err := cc.UpdateCertByAlias(updateAliasPayload, endorsers, -1, true)
	if err != nil {
		log.Fatal(err)
	}
	examples.PrintPrettyJson(resp2)

	fmt.Println("====================== query alias ======================")
	aliasInfos2, err := cc.QueryCertsAlias([]string{certAlias})
	if err != nil {
		log.Fatal(err)
	}
	examples.PrintPrettyJson(aliasInfos2)

	fmt.Println("====================== delete alias ======================")
	deleteAliasPayload := cc.CreateDeleteCertsAliasPayload([]string{certAlias})

	endorsers2, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[cc.GetHashType()],
		cc.GetAuthType(), deleteAliasPayload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	if err != nil {
		log.Fatal(err)
	}
	resp3, err := cc.DeleteCertsAlias(deleteAliasPayload, endorsers2, -1, true)
	if err != nil {
		log.Fatal(err)
	}
	examples.PrintPrettyJson(resp3)

	fmt.Println("====================== query alias ======================")
	aliasInfos3, err := cc.QueryCertsAlias([]string{certAlias})
	if err != nil {
		log.Fatal(err)
	}
	examples.PrintPrettyJson(aliasInfos3)
}

/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	commonCrt "chainmaker.org/chainmaker/common/v2/cert"
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/common/v2/serialize"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

// Deprecated: SignPayload use ./utils.MakeEndorserWithPem
func SignPayload(keyPem, certPem []byte, payload *common.Payload) (*common.EndorsementEntry, error) {
	key, err := asym.PrivateKeyFromPEM(keyPem, nil)
	if err != nil {
		return nil, fmt.Errorf("asym.PrivateKeyFromPEM failed, %s", err)
	}

	blockCrt, _ := pem.Decode(certPem)
	if blockCrt == nil {
		return nil, fmt.Errorf("decode pem failed, invalid certificate")
	}
	crt, err := bcx509.ParseCertificate(blockCrt.Bytes)
	if err != nil {
		return nil, fmt.Errorf("bcx509.ParseCertificate failed, %s", err)
	}

	signature, err := utils.SignPayload(key, crt, payload)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
	}

	var orgId string
	if len(crt.Subject.Organization) != 0 {
		orgId = crt.Subject.Organization[0]
	}

	sender := &accesscontrol.Member{
		OrgId:      orgId,
		MemberInfo: certPem,
		MemberType: accesscontrol.MemberType_CERT,
	}

	entry := &common.EndorsementEntry{
		Signer:    sender,
		Signature: signature,
	}

	return entry, nil
}

// Deprecated: SignPayloadWithPath use ./utils.MakeEndorserWithPath instead.
func SignPayloadWithPath(keyFilePath, crtFilePath string, payload *common.Payload) (*common.EndorsementEntry, error) {
	// 读取私钥
	keyBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file failed, %s", err)
	}

	// 读取证书
	crtBytes, err := ioutil.ReadFile(crtFilePath)
	if err != nil {
		return nil, fmt.Errorf("read crt file failed, %s", err)
	}

	return SignPayload(keyBytes, crtBytes, payload)
}

func GetEVMAddressFromCertPath(certFilePath string) (string, error) {
	certBytes, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return "", fmt.Errorf("read cert file [%s] failed, %s", certFilePath, err)
	}

	return GetEVMAddressFromCertBytes(certBytes)
}

func GetEVMAddressFromPrivateKeyPath(privateKeyFilePath, hashType string) (string, error) {
	keyPem, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("readFile failed, %s", err.Error())
	}

	return GetEVMAddressFromPrivateKeyBytes(keyPem, hashType)
}

func GetEVMAddressFromCertBytes(certBytes []byte) (string, error) {
	block, _ := pem.Decode(certBytes)
	cert, err := bcx509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("ParseCertificate cert failed, %s", err)
	}

	ski := hex.EncodeToString(cert.SubjectKeyId)
	addrInt, err := evmutils.MakeAddressFromHex(ski)
	if err != nil {
		return "", fmt.Errorf("make address from cert SKI failed, %s", err)
	}

	return addrInt.String(), nil
}

func GetEVMAddressFromPrivateKeyBytes(privateKeyBytes []byte, hashType string) (string, error) {
	privateKey, err := asym.PrivateKeyFromPEM(privateKeyBytes, nil)
	if err != nil {
		return "", fmt.Errorf("PrivateKeyFromPEM failed, %s", err.Error())
	}

	publicKey := privateKey.PublicKey()

	ski, err := commonCrt.ComputeSKI(crypto.HashAlgoMap[hashType], publicKey.ToStandardKey())
	if err != nil {
		return "", fmt.Errorf("computeSKI from publickey failed")
	}

	skiStr := hex.EncodeToString(ski)

	addrInt, err := evmutils.MakeAddressFromHex(skiStr)
	if err != nil {
		return "", fmt.Errorf("make address from cert SKI failed, %s", err)
	}

	return addrInt.String(), nil
}

func (cc *ChainClient) EasyCodecItemToParamsMap(items []*serialize.EasyCodecItem) map[string][]byte {
	return serialize.EasyCodecItemToParamsMap(items)
}

func GetZXAddressFromPKHex(pkHex string) (string, error) {
	pk, err := hex.DecodeString(pkHex)
	if err != nil {
		return "", err
	}

	return evmutils.ZXAddressFromPublicKeyDER(pk)
}

func GetZXAddressFromPKPEM(pkPEM string) (string, error) {
	return evmutils.ZXAddressFromPublicKeyPEM([]byte(pkPEM))
}

func GetZXAddressFromCertPEM(certPEM string) (string, error) {
	return evmutils.ZXAddressFromCertificatePEM([]byte(certPEM))
}

func GetZXAddressFromCertPath(certPath string) (string, error) {
	return evmutils.ZXAddressFromCertificatePath(certPath)
}

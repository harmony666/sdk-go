/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"chainmaker.org/chainmaker/common/v2/crypto"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	computeName  = "compute_name"
	computeCode  = "compute_code"
	computeCode2 = "compute_code2"
	ComputeRes   = "private_compute_result"
	//enclaveId    = "enclave_id"
	quoteId                     = "quote_id"
	quote                       = "quote_content"
	orderId                     = "order_id"
	ContractresultcodeOk uint32 = 0 //todo pb create

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

var (
	proof     []byte
	enclaveId string
	caCert    []byte
	report    string
)

func main() {
	//testChainClientSaveData()
	//testChainClientSaveDir()
	//testChainClientGetContract()
	//testChainClientGetData()
	//testChainClientGetDir()
	//testChainClientSaveCACert()
	//testChainClientGetCACert()
	testChainClientSaveEnclaveReport()
	//testChainClientSaveRemoteAttestationProof()
	//testChainClientGetEnclaveEncryptPubKey()
	//testChainClientGetEnclaveVerificationPubKey()
	//testChainClientGetEnclaveReport()
	//testChainClientGetEnclaveChallenge()
	//testChainClientGetEnclaveSignature()
	//testChainClientGetEnclaveProof()
}

func readFileData(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	return data
}

func initCaCert() {
	caCert = readFileData("../../testdata/remote_attestation/enclave_cacert.crt")
}

func initProof() {
	proof = readFileData("../../testdata/remote_attestation/proof.hex")

}

func initEnclaveId() {
	enclaveId = "global_enclave_id"
}

func initReport() {
	reportBytes := readFileData("../../testdata/remote_attestation/report.dat")
	report = hex.EncodeToString(reportBytes)
}

var priDir = &common.StrSlice{
	StrArr: []string{"dir_key1", "dir_key2", "dir_key3"},
}

func testChainClientSaveData() {
	codeHash := sha256.Sum256([]byte(computeCode))
	txid := utils.GetTimestampTxId()
	result := &common.ContractResult{
		Code:    0,
		Result:  nil,
		Message: "",
		GasUsed: 0,
	}
	rwSet := &common.TxRWSet{
		TxReads: []*common.TxRead{
			{Key: []byte("key2"), Value: []byte("value2"), ContractName: computeName},
		},
		TxWrites: []*common.TxWrite{
			{Key: []byte("key1"), Value: []byte("value_1"), ContractName: computeName},
			{Key: []byte("key3"), Value: []byte("value_3"), ContractName: computeName},
			{Key: []byte("key4"), Value: []byte("value_4"), ContractName: computeName},
			{Key: []byte("key5"), Value: []byte("value_5"), ContractName: computeName},
		},
	}
	// todo add reportHash,sign
	var reportHash []byte
	var reportSign []byte

	cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := cc.SaveData(computeName, examples.Version, false, codeHash[:], reportHash, result, []byte(""), txid, rwSet,
		reportSign, nil, []byte(""), false, 1)
	if err != nil {
		log.Fatalln(err)
	}
	if res.ContractResult.Code != ContractresultcodeOk {
		log.Fatalln("res.ContractResult.Code != common.ContractResultCode_OK")
	}
}

func testChainClientSaveDir() {
	txid := utils.GetTimestampTxId()

	cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}
	got, err := cc.SaveDir(orderId, txid, priDir, false, 1)
	if err != nil {
		log.Fatalln(err)
	}
	if got.ContractResult.Code != ContractresultcodeOk {
		log.Fatalln("got.ContractResult.Code != common.ContractResultCode_OK")
	}
}

func testChainClientGetContract() {
	type args struct {
		contractName string
		codeHash     string
	}

	codeHash := sha256.Sum256([]byte(computeCode))
	codeHash2 := sha256.Sum256([]byte(computeCode2))

	tests := []struct {
		name    string
		args    args
		want    *common.PrivateGetContract
		wantErr interface{}
	}{
		{
			name: "test1",
			args: args{
				contractName: computeName,
				codeHash:     string(codeHash[:]),
			},
			want: &common.PrivateGetContract{
				ContractCode: []byte(computeCode),
				Version:      examples.Version,
				GasLimit:     10000000000,
			},
			wantErr: nil,
		},
		{
			name: "test2",
			args: args{
				contractName: computeName,
				codeHash:     string(codeHash2[:]),
			},
			want: &common.PrivateGetContract{
				ContractCode: []byte(computeCode2),
				Version:      examples.UpgradeVersion,
				GasLimit:     10000000000,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetContract(tt.args.contractName, tt.args.codeHash)
		if err != nil {
			log.Fatalln(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			log.Fatalln("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientGetData() {

	type args struct {
		contractName string
		key          string
	}

	dirByte, _ := priDir.Marshal()
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				contractName: computeName,
				key:          "key1",
			},
			want:    []byte("value_1"),
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				contractName: "",
				key:          orderId,
			},
			want:    dirByte,
			wantErr: true,
		},
		{
			name: "test3",
			args: args{
				contractName: "",
				key:          enclaveId,
			},
			want:    caCert,
			wantErr: true,
		},
		{
			name: "test4",
			args: args{
				contractName: "",
				key:          quoteId,
			},
			want:    []byte(quote),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetData(tt.args.contractName, tt.args.key)
		if err != nil {
			log.Fatalln(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			log.Fatalln("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientGetDir() {

	type args struct {
		orderId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				orderId: "orderId",
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetDir(tt.args.orderId)
		if err != nil {
			log.Fatalln(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			log.Fatalln("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientSaveCACert() {

	initCaCert()

	type args struct {
		caCert         string
		txId           string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				caCert:         string(caCert),
				txId:           "",
				withSyncResult: true,
				timeout:        1,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		chainClient, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		resp, err := saveEnclaveCACert(chainClient, tt.args.caCert, tt.args.txId, tt.args.withSyncResult, tt.args.timeout,
			"org1admin1", "org2admin1", "org3admin1", "org4admin1", "org5admin1")
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("testChainClientSaveCACert() got response ==> \n%+v\n", resp)
	}
	log.Println("########## testChainClientSaveCACert() execute successfully.")
}

func saveEnclaveCACert(client *sdk.ChainClient, enclaveCACert string, txId string, withSyncResult bool,
	timeout int64, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateSaveEnclaveCACertPayload(enclaveCACert, txId)
	if err != nil {
		return nil, err
	}

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendMultiSigningRequest(payload, endorsers, timeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func testChainClientGetCACert() {

	initCaCert()

	type args struct {
		txId           string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				txId:           "",
				withSyncResult: true,
				timeout:        1,
			},
			want:    caCert,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetEnclaveCACert()
		if err != nil {
			log.Fatalln(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			log.Fatalln("!reflect.DeepEqual(got, tt.want)")
		}
	}
	log.Println("########## testChainClientGetCACert() execute successfully.")
}

func testChainClientSaveEnclaveReport() {

	initEnclaveId()
	initReport()

	type args struct {
		enclaveId      string
		report         string
		txId           string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId:      enclaveId,
				report:         report,
				txId:           "",
				withSyncResult: true,
				timeout:        1,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		chainClient, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := saveEnclaveReport(chainClient, tt.args.enclaveId, tt.args.report, tt.args.txId, tt.args.withSyncResult, tt.args.timeout,
			"org1admin1", "org2admin1", "org3admin1", "org4admin1")
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("testChainClientSaveEnclaveReport got %+v\n", got)
	}
	log.Println("########## testChainClientSaveEnclaveReport() execute successfully.")
}

func saveEnclaveReport(client *sdk.ChainClient, enclaveId string, report string, txId string, withSyncResult bool,
	timeout int64, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateSaveEnclaveReportPayload(enclaveId, report, txId)
	if err != nil {
		return nil, err
	}

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendMultiSigningRequest(payload, endorsers, timeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func testChainClientSaveRemoteAttestationProof() {

	initProof()

	type args struct {
		proof          string
		txId           string
		withSyncResult bool
		timeout        int64
	}

	tests := []struct {
		name    string
		args    args
		want    *common.TxResponse
		wantErr bool
	}{
		{
			name: "TEST1",
			args: args{
				proof:          string(proof),
				txId:           "",
				withSyncResult: true,
				timeout:        -1,
			},
			want: &common.TxResponse{
				Code:           0,
				Message:        "OK",
				ContractResult: nil,
			},
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.SaveRemoteAttestationProof(tt.args.proof, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("enclaveId = %s \n", string(got.ContractResult.Result))
	}
	log.Println("########## testChainClientSaveRemoteAttestationProof() execute successfully.")
}

func testChainClientGetEnclaveEncryptPubKey() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetEnclaveEncryptPubKey(tt.args.enclaveId)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("encrypt pub key => %s \n", got)
	}
	log.Println("########## testChainClientGetEnclaveEncryptPubKey() execute successfully.")
}

func testChainClientGetEnclaveVerificationPubKey() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetEnclaveVerificationPubKey(tt.args.enclaveId)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("verification pub key => %s \n", got)
	}
	log.Println("########## testChainClientGetEnclaveVerificationPubKey() execute successfully.")
}

func testChainClientGetEnclaveReport() {

	initEnclaveId()
	initReport()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    report,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetEnclaveReport(tt.args.enclaveId)
		if err != nil {
			log.Fatalln(err)
		}

		report, err := hex.DecodeString(string(got))
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("testChainClientGetEnclaveReport got =>\n%s\n", report)
	}
	log.Println("########## testChainClientGetEnclaveReport() execute successfully.")
}

func testChainClientGetEnclaveChallenge() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetEnclaveChallenge(tt.args.enclaveId)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("challenge => %s \n", got)
	}
	log.Println("########## testChainClientGetEnclaveChallenge() execute successfully.")
}

func testChainClientGetEnclaveSignature() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetEnclaveSignature(tt.args.enclaveId)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("signature = %x \n", got)
	}
	log.Println("########## testChainClientGetEnclaveSignature() execute successfully.")
}

func testChainClientGetEnclaveProof() {

	initEnclaveId()
	initProof()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    proof,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			log.Fatalln(err)
		}
		got, err := cc.GetEnclaveProof(tt.args.enclaveId)
		if err != nil {
			log.Fatalln(err)
		}

		if !reflect.DeepEqual(got, tt.want) {
			log.Fatalln("!reflect.DeepEqual(got, tt.want)")
		}
	}
	log.Println("########## testChainClientGetEnclaveProof() execute successfully.")
}

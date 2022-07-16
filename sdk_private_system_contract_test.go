/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"crypto/sha256"
	"errors"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/stretchr/testify/require"
)

const (
	computeCodeGood = "computeCodeGood"
	computeCodeBad  = "computeCodeBad"
	computeName     = "computeName"
	chainId         = "chain1"
	enclaveCACert   = "enclaveCACert"
	enclaveId       = "enclaveId"
	report          = "report"
	orderId         = "orderId"
)

func TestChainClient_CheckCallerCertAuth(t *testing.T) {
	type args struct {
		payload   string
		orgIds    []string
		signPairs []*syscontract.SignInfo
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.TxResponse
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code:           common.TxStatusCode_SUCCESS,
				Message:        "",
				ContractResult: &common.ContractResult{},
			},
			serverTxErr: nil,
			args:        args{},
			want: &common.TxResponse{
				Code:           0,
				Message:        "",
				ContractResult: &common.ContractResult{},
				TxId:           "9bb59a8d04074e9a87255128735d4ec88712bb398a844d0cb91ee8fd21235414",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.CheckCallerCertAuth(tt.args.payload, tt.args.orgIds, tt.args.signPairs)
			got.TxId = tt.want.TxId
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckCallerCertAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckCallerCertAuth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_CreateSaveEnclaveCACertPayload(t *testing.T) {
	type args struct {
		enclaveCACert string
		txId          string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.Payload
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				ContractResult: &common.ContractResult{},
				TxId:           "ba9299ce68854ff0b9968241bba592b0af52d0120fd744f9b443d2a73cabc1d5",
			},
			serverTxErr: nil,
			args: args{
				enclaveCACert: enclaveCACert,
				txId:          "ba9299ce68854ff0b9968241bba592b0af52d0120fd744f9b443d2a73cabc1d5",
			},
			want: &common.Payload{
				ChainId:      chainId,
				ContractName: syscontract.SystemContract_PRIVATE_COMPUTE.String(),
				Method:       syscontract.PrivateComputeFunction_SAVE_CA_CERT.String(),
				Parameters: []*common.KeyValuePair{
					{
						Key:   utils.KeyCaCert,
						Value: []byte(enclaveCACert),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.CreateSaveEnclaveCACertPayload(tt.args.enclaveCACert, tt.args.txId)
			tt.want.Timestamp = got.Timestamp
			tt.want.TxId = got.TxId
			tt.want.Parameters = got.Parameters
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSaveEnclaveCACertPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSaveEnclaveCACertPayload() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_CreateSaveEnclaveReportPayload(t *testing.T) {
	type args struct {
		enclaveId string
		report    string
		txId      string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.Payload
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				ContractResult: &common.ContractResult{},
				TxId:           "19d012d5b15a4a82945eeeca891aac396a6e0d7ee0094cffbce7bb40ec79fc11",
			},
			serverTxErr: nil,
			args: args{
				enclaveId: enclaveId,
				report:    report,
				txId:      "19d012d5b15a4a82945eeeca891aac396a6e0d7ee0094cffbce7bb40ec79fc11",
			},
			want: &common.Payload{
				ChainId:      chainId,
				TxId:         "19d012d5b15a4a82945eeeca891aac396a6e0d7ee0094cffbce7bb40ec79fc11",
				ContractName: syscontract.SystemContract_PRIVATE_COMPUTE.String(),
				Method:       syscontract.PrivateComputeFunction_SAVE_ENCLAVE_REPORT.String(),
				Parameters: []*common.KeyValuePair{
					{
						Key:   utils.KeyEnclaveId,
						Value: []byte(enclaveId),
					},
					{
						Key:   utils.KeyReport,
						Value: []byte(report),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.CreateSaveEnclaveReportPayload(tt.args.enclaveId, tt.args.report, tt.args.txId)
			tt.want.Timestamp = got.Timestamp
			tt.want.TxId = got.TxId
			tt.want.Parameters = got.Parameters
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSaveEnclaveReportPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSaveEnclaveReportPayload() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetContract(t *testing.T) {
	codeHashGood := sha256.Sum256([]byte(computeCodeGood))
	codeHashBad := sha256.Sum256([]byte(computeCodeBad))
	type args struct {
		contractName string
		codeHash     string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.PrivateGetContract
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code: ContractResultCode_OK,
				},
			},
			args: args{
				contractName: computeName,
				codeHash:     string(codeHashBad[:]),
			},
			serverTxErr: nil,
			want:        &common.PrivateGetContract{},
			wantErr:     false,
		},

		{
			name: "badCheckProposalRequestResp",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
			},
			args: args{
				contractName: computeName,
				codeHash:     string(codeHashGood[:]),
			},
			serverTxErr: nil,
			want:        nil,
			wantErr:     true,
		},
		{
			name: "badResp",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_CONTRACT_FAIL,
			},
			args: args{
				contractName: computeName,
				codeHash:     string(codeHashBad[:]),
			},
			serverTxErr: errors.New("get a error "),
			want:        nil,
			wantErr:     true,
		},

		{
			name: "badUnmarshal",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("bad struct slice"),
				},
			},
			args: args{
				contractName: computeName,
				codeHash:     string(codeHashBad[:]),
			},
			serverTxErr: nil,
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetContract(tt.args.contractName, tt.args.codeHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetData(t *testing.T) {
	type args struct {
		contractName string
		key          string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code: ContractResultCode_OK,
				},
			},
			args: args{
				contractName: computeName,
				key:          "key1",
			},
			wantErr: false,
		},
		{
			name: "badRequest",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_CONTRACT_FAIL},
			args: args{
				contractName: computeName,
				key:          "key2",
			},
			serverTxErr: errors.New("send QUERY_CONTRACT failed, client.call failed"),
			want:        nil,
			wantErr:     true,
		},

		{
			name: "badResp",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
			},
			args: args{
				contractName: computeName,
				key:          "key2",
			},
			serverTxErr: nil,
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetData(tt.args.contractName, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetDir(t *testing.T) {
	type args struct {
		orderId string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				orderId: orderId,
			},
			want:    []byte("result"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetDir(tt.args.orderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetEnclaveCACert(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			want:        []byte("result"),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetEnclaveCACert()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnclaveCACert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnclaveCACert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetEnclaveChallenge(t *testing.T) {
	type args struct {
		enclaveId string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				enclaveId: enclaveId,
			},
			want:    []byte("result"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetEnclaveChallenge(tt.args.enclaveId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnclaveChallenge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnclaveChallenge() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetEnclaveEncryptPubKey(t *testing.T) {
	type args struct {
		enclaveId string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{

		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				enclaveId: enclaveId,
			},
			want:    []byte("result"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetEnclaveEncryptPubKey(tt.args.enclaveId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnclaveEncryptPubKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnclaveEncryptPubKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetEnclaveProof(t *testing.T) {
	type args struct {
		enclaveId string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				enclaveId: enclaveId,
			},
			want:    []byte("result"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetEnclaveProof(tt.args.enclaveId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnclaveProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnclaveProof() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetEnclaveReport(t *testing.T) {
	type args struct {
		enclaveId string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				enclaveId: enclaveId,
			},
			want:    []byte("result"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetEnclaveReport(tt.args.enclaveId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnclaveReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnclaveReport() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetEnclaveSignature(t *testing.T) {
	type args struct {
		enclaveId string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				enclaveId: enclaveId,
			},
			want:    []byte("result"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetEnclaveSignature(tt.args.enclaveId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnclaveSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnclaveSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetEnclaveVerificationPubKey(t *testing.T) {
	type args struct {
		enclaveId string
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         []byte
		wantErr      bool
	}{
		{
			name: "good",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				enclaveId: enclaveId,
			},
			want:    []byte("result"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.GetEnclaveVerificationPubKey(tt.args.enclaveId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnclaveVerificationPubKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnclaveVerificationPubKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_SaveData(t *testing.T) {
	type args struct {
		contractName    string
		contractVersion string
		isDeployment    bool
		codeHash        []byte
		reportHash      []byte
		result          *common.ContractResult
		codeHeader      []byte
		txId            string
		rwSet           *common.TxRWSet
		sign            []byte
		events          *common.StrSlice
		privateReq      []byte
		withSyncResult  bool
		timeout         int64
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.TxResponse
		wantErr      bool
	}{
		{
			name: "bad",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},

			serverTxErr: nil,
			args:        args{},

			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.SaveData(tt.args.contractName, tt.args.contractVersion, tt.args.isDeployment, tt.args.codeHash, tt.args.reportHash, tt.args.result, tt.args.codeHeader, tt.args.txId, tt.args.rwSet, tt.args.sign, tt.args.events, tt.args.privateReq, tt.args.withSyncResult, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_SaveDir(t *testing.T) {
	type args struct {
		orderId        string
		txId           string
		privateDir     *common.StrSlice
		withSyncResult bool
		timeout        int64
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.TxResponse
		wantErr      bool
	}{
		{
			name: "bad",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				privateDir: &common.StrSlice{StrArr: nil},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.SaveDir(tt.args.orderId, tt.args.txId, tt.args.privateDir, tt.args.withSyncResult, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_SaveEnclaveCACert(t *testing.T) {
	type args struct {
		enclaveCACert  string
		txId           string
		withSyncResult bool
		timeout        int64
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.TxResponse
		wantErr      bool
	}{
		{
			name: "bad",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args:        args{},
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.SaveEnclaveCACert(tt.args.enclaveCACert, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveEnclaveCACert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveEnclaveCACert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_SaveEnclaveReport(t *testing.T) {
	type args struct {
		enclaveId      string
		report         string
		txId           string
		withSyncResult bool
		timeout        int64
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.TxResponse
		wantErr      bool
	}{
		{
			name: "bad",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args:        args{},
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.SaveEnclaveReport(tt.args.enclaveId, tt.args.report, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveEnclaveReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveEnclaveReport() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_SaveRemoteAttestationProof(t *testing.T) {
	type args struct {
		proof          string
		txId           string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.TxResponse
		wantErr      bool
	}{
		{
			name: "bad",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args:        args{},
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.SaveRemoteAttestationProof(tt.args.proof, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveRemoteAttestationProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveRemoteAttestationProof() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_SendMultiSigningRequest(t *testing.T) {
	type args struct {
		payload        *common.Payload
		endorsers      []*common.EndorsementEntry
		timeout        int64
		withSyncResult bool
	}

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		args         args
		want         *common.TxResponse
		wantErr      bool
	}{
		{
			name: "bad",
			serverTxResp: &common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Code:   ContractResultCode_OK,
					Result: []byte("result"),
				},
			},
			serverTxErr: nil,
			args: args{
				payload: &common.Payload{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			got, err := cli.SendMultiSigningRequest(tt.args.payload, tt.args.endorsers, tt.args.timeout, tt.args.withSyncResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMultiSigningRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendMultiSigningRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

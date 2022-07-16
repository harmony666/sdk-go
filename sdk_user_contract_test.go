/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"errors"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/stretchr/testify/require"
)

func TestSendTxRequest(t *testing.T) {
	var txID = "b374f23e4e6747e4b5fcb3ca975ef1655ad56555adfd4534ae8676cd9f1eb145"

	tests := []struct {
		name         string
		cliTxReq     *common.TxRequest
		serverTxResp *common.TxResponse
		serverTxErr  error
		wantErr      bool
	}{
		{
			"bad",
			&common.TxRequest{Payload: &common.Payload{TxId: txID}},
			&common.TxResponse{
				Code: common.TxStatusCode_CONTRACT_FAIL,
			},
			errors.New("rpc server throw an error"),
			true,
		},
		{
			"good",
			&common.TxRequest{Payload: &common.Payload{TxId: txID}},
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
			},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			txResp, err := cli.SendTxRequest(tt.cliTxReq, -1, false)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				require.Contains(t, txResp.Message, tt.serverTxErr.Error())
			} else {
				require.Equal(t, tt.serverTxResp.Code, txResp.Code)
			}
		})
	}
}

func TestCreateContractCreatePayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateContractCreatePayload("claim", "v1",
				"./testdata/claim-wasm-demo/rust-fact-2.0.0.wasm", common.RuntimeType_WASMER,
				[]*common.KeyValuePair{})
			require.Nil(t, err)
		})
	}
}

func TestCreateContractUpgradePayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateContractUpgradePayload("claim", "v1",
				"./testdata/claim-wasm-demo/rust-fact-2.0.0.wasm", common.RuntimeType_WASMER,
				[]*common.KeyValuePair{})
			require.Nil(t, err)
		})
	}
}

func TestCreateContractFreezePayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateContractFreezePayload("claim")
			require.Nil(t, err)
		})
	}
}

func TestCreateContractUnfreezePayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateContractUnfreezePayload("claim")
			require.Nil(t, err)
		})
	}
}

func TestCreateContractRevokePayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateContractRevokePayload("claim")
			require.Nil(t, err)
		})
	}
}

func TestSignContractManagePayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.SignContractManagePayload(&common.Payload{})
			require.Nil(t, err)
		})
	}
}

func TestSendContractManageRequest(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.SendContractManageRequest(&common.Payload{}, nil, -1, false)
			require.Nil(t, err)
		})
	}
}

func TestInvokeContract(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.InvokeContract("claim", "save", "",
				[]*common.KeyValuePair{}, -1, false)
			require.Nil(t, err)
		})
	}
}

func TestInvokeContractWithLimit(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.InvokeContractWithLimit("claim", "save", "",
				[]*common.KeyValuePair{}, -1, false, &common.Limit{GasLimit: 1000000})
			require.Nil(t, err)
		})
	}
}

func TestQueryContract(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.QueryContract("claim", "save", []*common.KeyValuePair{}, -1)
			require.Nil(t, err)
		})
	}
}

func TestGetTxRequest(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte{},
				Message: "OK",
			}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.GetTxRequest("claim", "save", "", []*common.KeyValuePair{})
			require.Nil(t, err)
		})
	}
}

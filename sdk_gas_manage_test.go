/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"strconv"
	"testing"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"github.com/stretchr/testify/require"
)

func TestChainClient_CreateSetGasAdminPayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateSetGasAdminPayload("ZXaaa6f45415493ffb832ca28faa14bef5c357f5f0")
			require.Nil(t, err)
		})
	}
}

func TestChainClient_GetGasAdmin(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{
				ContractResult: &common.ContractResult{
					Message: "OK",
					Result:  []byte("ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429"),
				},
				Message: "SUCCESS",
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			adminAddr, err := cc.GetGasAdmin()
			require.Nil(t, err)
			require.Equal(t, adminAddr, string(tt.serverTxResp.ContractResult.Result))
		})
	}
}

func TestChainClient_CreateRechargeGasPayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateRechargeGasPayload([]*syscontract.RechargeGas{
				{
					Address:   "ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429",
					GasAmount: 100,
				},
			})
			require.Nil(t, err)
		})
	}
}

func TestChainClient_GetGasBalance(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{
				ContractResult: &common.ContractResult{
					Message: "OK",
					Result:  []byte("100"),
				},
				Message: "SUCCESS",
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			balance, err := cc.GetGasBalance("ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429")
			require.Nil(t, err)
			respBalance, err := strconv.ParseInt(string(tt.serverTxResp.ContractResult.Result), 10, 64)
			require.Nil(t, err)
			require.Equal(t, balance, respBalance)
		})
	}
}

func TestChainClient_CreateRefundGasPayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateRefundGasPayload("ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429", 100)
			require.Nil(t, err)
		})
	}
}

func TestChainClient_CreateFrozenGasAccountPayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateFrozenGasAccountPayload("ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429")
			require.Nil(t, err)
		})
	}
}

func TestChainClient_CreateUnfrozenGasAccountPayload(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.CreateUnfrozenGasAccountPayload("ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429")
			require.Nil(t, err)
		})
	}
}

func TestChainClient_GetGasAccountStatus(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{
				ContractResult: &common.ContractResult{
					Message: "OK",
					Result:  []byte("1"),
				},
				Message: "SUCCESS",
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			status, err := cc.GetGasAccountStatus("ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429")
			require.Nil(t, err)
			require.Equal(t, status, string(tt.serverTxResp.ContractResult.Result) == "0")
		})
	}
}

func TestChainClient_SendGasManageRequest(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			_, err = cc.SendGasManageRequest(&common.Payload{
				ChainId:      "chain1",
				TxId:         "2043b7e1f98a4a8385e0aaad1843531680225e43834d400283e38865cb03691b",
				Timestamp:    time.Now().Unix(),
				ContractName: "ACCOUNT_MANAGER",
				Method:       "RECHARGE_GAS",
				Parameters: []*common.KeyValuePair{
					{
						Key:   "address_key",
						Value: []byte("ZXf9fc99018a7d5d09b8e909e5c7d34c8b4fce9429"),
					},
					{
						Key:   "charge_gas_amount",
						Value: []byte("100"),
					},
				},
			}, nil, -1, false)
			require.Nil(t, err)
		})
	}
}

func TestChainClient_AttachGasLimit(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cc.Stop()

			payload := &common.Payload{}
			limit := &common.Limit{GasLimit: 100}
			payload = cc.AttachGasLimit(payload, limit)
			require.Equal(t, payload.Limit, limit)
		})
	}
}

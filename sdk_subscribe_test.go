/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/stretchr/testify/require"
)

func TestChainClient_SubscribeBlock(t *testing.T) {
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, err = cc.SubscribeBlock(ctx, 0, 10, true, true)
			require.Nil(t, err)
		})
	}
}

func TestChainClient_SubscribeTx(t *testing.T) {
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, err = cc.SubscribeTx(ctx, 0, 10, "claim", nil)
			require.Nil(t, err)
		})
	}
}

func TestChainClient_SubscribeContractEvent(t *testing.T) {
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, err = cc.SubscribeContractEvent(ctx, 0, 10, "claim", "topic1")
			require.Nil(t, err)
		})
	}
}

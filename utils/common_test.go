/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/stretchr/testify/require"
)

func TestGetRandTxId(t *testing.T) {
	txId := GetRandTxId()
	require.Len(t, txId, 64)
}

func TestCheckProposalRequestResp(t *testing.T) {
	tests := []struct {
		name               string
		serverTxResp       *common.TxResponse
		needContractResult bool
		wantErr            bool
	}{
		{
			"good",
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS, ContractResult: &common.ContractResult{
				Code: SUCCESS,
			}},
			true,
			false,
		},
		{
			"bad",
			&common.TxResponse{Code: common.TxStatusCode_CONTRACT_FAIL},
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckProposalRequestResp(tt.serverTxResp, tt.needContractResult)
			require.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestGetTimestampTxId(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "正常流",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimestampTxId()
			t.Log(got)
			require.Len(t, got, 64)
		})
	}
}

func TestGetNanosecondByTxId(t *testing.T) {
	nano := time.Now().UnixNano()
	type args struct {
		nano int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "正常流",
			args: args{nano: nano},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimestampTxIdByNano(tt.args.nano)
			nanosecond, err := GetNanoByTimestampTxId(got)
			if err != nil {
				return
			}

			require.Truef(t, nanosecond == nano, "emmm not ok")
		})
	}
}

func TestHexEB(t *testing.T) {
	b, err := hex.DecodeString("ca")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
}

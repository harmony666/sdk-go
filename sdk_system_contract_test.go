/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"encoding/json"
	"strconv"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/discovery"
	"chainmaker.org/chainmaker/pb-go/v2/store"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestGetTxByTxId(t *testing.T) {
	var txID = "b374f23e4e6747e4b5fcb3ca975ef1655ad56555adfd4534ae8676cd9f1eb145"

	goodTxInfoBz, err := proto.Marshal(&common.TransactionInfo{Transaction: &common.Transaction{Payload: &common.Payload{TxId: txID}}})
	require.Nil(t, err)

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		wantErr      bool
	}{
		{
			"bad",
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: []byte("this is a bad *common.TransactionInfo bytes"),
				},
			},
			nil,
			true,
		},
		{
			"good",
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: goodTxInfoBz,
				},
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

			txInfo, err := cli.GetTxByTxId(txID)
			require.Equal(t, tt.wantErr, err != nil)

			if txInfo != nil {
				bz, err := proto.Marshal(txInfo)
				require.Nil(t, err)
				require.Equal(t, tt.serverTxResp.ContractResult.Result, bz)
			}
		})
	}
}

func TestGetTxWithRWSetByTxId(t *testing.T) {
	var txID = "b374f23e4e6747e4b5fcb3ca975ef1655ad56555adfd4534ae8676cd9f1eb145"

	goodTxInfoBz, err := proto.Marshal(&common.TransactionInfoWithRWSet{
		Transaction: &common.Transaction{
			Payload: &common.Payload{
				TxId: txID,
			},
		},
		RwSet: &common.TxRWSet{
			TxId:     txID,
			TxReads:  []*common.TxRead{{}},
			TxWrites: []*common.TxWrite{{}},
		},
	})
	require.Nil(t, err)

	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
		wantErr      bool
	}{
		{
			"bad",
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: []byte("this is a bad *common.TransactionInfo bytes"),
				},
			},
			nil,
			true,
		},
		{
			"good",
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: goodTxInfoBz,
				},
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

			txInfo, err := cli.GetTxWithRWSetByTxId(txID)
			require.Equal(t, tt.wantErr, err != nil)

			if txInfo != nil {
				bz, err := proto.Marshal(txInfo)
				require.Nil(t, err)
				require.Equal(t, tt.serverTxResp.ContractResult.Result, bz)
			}
		})
	}
}

func TestGetBlockByHeight(t *testing.T) {
	var height uint64 = 1
	rawBlkInfo, err := proto.Marshal(&common.BlockInfo{
		Block: &common.Block{
			Header: &common.BlockHeader{
				BlockHeight: height,
			},
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawBlkInfo,
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

			blk, err := cc.GetBlockByHeight(height, false)
			require.Nil(t, err)
			require.Equal(t, blk.Block.Header.BlockHeight, height)
		})
	}
}

func TestGetBlockByHash(t *testing.T) {
	var blkHash = []byte("block hash")
	rawBlkInfo, err := proto.Marshal(&common.BlockInfo{
		Block: &common.Block{
			Header: &common.BlockHeader{
				BlockHash: blkHash,
			},
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawBlkInfo,
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

			blk, err := cc.GetBlockByHash(string(blkHash), false)
			require.Nil(t, err)
			require.Equal(t, blk.Block.Header.BlockHash, blkHash)
		})
	}
}

func TestGetBlockByTxId(t *testing.T) {
	rawBlkInfo, err := proto.Marshal(&common.BlockInfo{
		Block: &common.Block{
			Header: &common.BlockHeader{
				BlockHash: []byte{},
			},
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawBlkInfo,
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

			_, err = cc.GetBlockByTxId("txid", false)
			require.Nil(t, err)
		})
	}
}

func TestGetLastConfigBlock(t *testing.T) {
	rawBlkInfo, err := proto.Marshal(&common.BlockInfo{
		Block: &common.Block{
			Header: &common.BlockHeader{
				BlockType: common.BlockType_CONFIG_BLOCK,
			},
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawBlkInfo,
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

			blk, err := cc.GetLastConfigBlock(false)
			require.Nil(t, err)
			require.Equal(t, blk.Block.Header.BlockType, common.BlockType_CONFIG_BLOCK)
		})
	}
}

func TestGetChainInfo(t *testing.T) {
	rawChainInfo, err := proto.Marshal(&discovery.ChainInfo{})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawChainInfo,
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

			chainInfo, err := cc.GetChainInfo()
			require.Nil(t, err)
			bz, err := chainInfo.Marshal()
			require.Nil(t, err)
			require.Equal(t, bz, rawChainInfo)
		})
	}
}

func TestGetNodeChainList(t *testing.T) {
	rawChainList, err := proto.Marshal(&discovery.ChainList{})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawChainList,
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

			chainList, err := cc.GetNodeChainList()
			require.Nil(t, err)
			bz, err := chainList.Marshal()
			require.Nil(t, err)
			require.Equal(t, bz, rawChainList)
		})
	}
}

func TestGetFullBlockByHeight(t *testing.T) {
	var height uint64 = 1
	rawBlk, err := proto.Marshal(&store.BlockWithRWSet{
		Block: &common.Block{
			Header: &common.BlockHeader{
				BlockHeight: height,
			},
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawBlk,
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

			blk, err := cc.GetFullBlockByHeight(height)
			require.Nil(t, err)
			bz, err := blk.Marshal()
			require.Nil(t, err)
			require.Equal(t, bz, rawBlk)
		})
	}
}

func TestGetArchivedBlockHeight(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte("1"),
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

			height, err := cc.GetArchivedBlockHeight()
			require.Nil(t, err)
			height2, err := strconv.ParseUint(string(tt.serverTxResp.ContractResult.Result), 10, 64)
			require.Nil(t, err)
			require.Equal(t, height, height2)
		})
	}
}

func TestGetBlockHeightByTxId(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte("1"),
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

			height, err := cc.GetBlockHeightByTxId("txid")
			require.Nil(t, err)
			height2, err := strconv.ParseUint(string(tt.serverTxResp.ContractResult.Result), 10, 64)
			require.Nil(t, err)
			require.Equal(t, height, height2)
		})
	}
}

func TestGetBlockHeightByHash(t *testing.T) {
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  []byte("1"),
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

			height, err := cc.GetBlockHeightByHash("block hash")
			require.Nil(t, err)
			height2, err := strconv.ParseUint(string(tt.serverTxResp.ContractResult.Result), 10, 64)
			require.Nil(t, err)
			require.Equal(t, height, height2)
		})
	}
}

func TestGetLastBlock(t *testing.T) {
	var height uint64 = 1
	rawBlk, err := proto.Marshal(&common.BlockInfo{
		Block: &common.Block{
			Header: &common.BlockHeader{
				BlockHeight: height,
			},
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawBlk,
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

			blk, err := cc.GetLastBlock(false)
			require.Nil(t, err)
			bz, err := blk.Marshal()
			require.Nil(t, err)
			require.Equal(t, bz, rawBlk)
		})
	}
}

func TestGetCurrentBlockHeight(t *testing.T) {
	var height uint64 = 1
	rawBlk, err := proto.Marshal(&common.BlockInfo{
		Block: &common.Block{
			Header: &common.BlockHeader{
				BlockHeight: height,
			},
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  rawBlk,
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

			height, err := cc.GetCurrentBlockHeight()
			require.Nil(t, err)
			var blk = &common.BlockInfo{}
			err = proto.Unmarshal(tt.serverTxResp.ContractResult.Result, blk)
			require.Nil(t, err)
			require.Equal(t, height, blk.Block.Header.BlockHeight)
		})
	}
}

func TestGetBlockHeaderByHeight(t *testing.T) {
	var height uint64 = 1
	raw, err := proto.Marshal(&common.BlockHeader{
		BlockHeight: height,
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  raw,
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

			header, err := cc.GetBlockHeaderByHeight(height)
			require.Nil(t, err)
			bz, err := header.Marshal()
			require.Nil(t, err)
			require.Equal(t, bz, raw)
		})
	}
}

func TestInvokeSystemContract(t *testing.T) {
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

			_, err = cc.InvokeSystemContract("claim", "save", "",
				[]*common.KeyValuePair{{}}, -1, false)
			require.Nil(t, err)
		})
	}
}

func TestQuerySystemContract(t *testing.T) {
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

			_, err = cc.QuerySystemContract("claim", "save",
				[]*common.KeyValuePair{{}}, -1)
			require.Nil(t, err)
		})
	}
}

func TestGetMerklePathByTxId(t *testing.T) {
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

			_, err = cc.GetMerklePathByTxId("tx id")
			require.Nil(t, err)
		})
	}
}

func TestCreateNativeContractAccessGrantPayload(t *testing.T) {
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

			_, err = cc.CreateNativeContractAccessGrantPayload([]string{"a", "b"})
			require.Nil(t, err)
		})
	}
}

func TestCreateNativeContractAccessRevokePayload(t *testing.T) {
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

			_, err = cc.CreateNativeContractAccessRevokePayload([]string{"a", "b"})
			require.Nil(t, err)
		})
	}
}

func TestGetContractInfo(t *testing.T) {
	contractName := "contract123"
	raw, err := json.Marshal(&common.Contract{
		Name: contractName,
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  raw,
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

			contract, err := cc.GetContractInfo(contractName)
			require.Nil(t, err)
			bz, err := json.Marshal(contract)
			require.Nil(t, err)
			require.Equal(t, bz, raw)
		})
	}
}

func TestGetContractList(t *testing.T) {
	contractName := "contract111"
	raw, err := json.Marshal([]*common.Contract{
		{
			Name: contractName,
		},
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  raw,
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

			contracts, err := cc.GetContractList()
			require.Nil(t, err)
			bz, err := json.Marshal(contracts)
			require.Nil(t, err)
			require.Equal(t, bz, raw)
		})
	}
}

func TestGetDisabledNativeContractList(t *testing.T) {
	contractName := "contract1"
	raw, err := json.Marshal([]string{
		contractName,
	})
	require.Nil(t, err)
	tests := []struct {
		name         string
		serverTxResp *common.TxResponse
		serverTxErr  error
	}{
		{
			"good",
			&common.TxResponse{ContractResult: &common.ContractResult{
				Result:  raw,
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

			contracts, err := cc.GetDisabledNativeContractList()
			require.Nil(t, err)
			bz, err := json.Marshal(contracts)
			require.Nil(t, err)
			require.Equal(t, bz, raw)
		})
	}
}

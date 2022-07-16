/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/proto"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/store"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func (cc *ChainClient) CreateArchiveBlockPayload(targetBlockHeight uint64) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [Archive] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.ArchiveBlock_BLOCK_HEIGHT.String(),
			Value: utils.U64ToBytes(targetBlockHeight),
		},
	}

	payload := cc.CreatePayload("", common.TxType_ARCHIVE, syscontract.SystemContract_ARCHIVE_MANAGE.String(),
		syscontract.ArchiveFunction_ARCHIVE_BLOCK.String(), pairs, defaultSeq, nil)

	return payload, nil
}

func (cc *ChainClient) CreateRestoreBlockPayload(fullBlock []byte) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [restore] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.RestoreBlock_FULL_BLOCK.String(),
			Value: fullBlock,
		},
	}

	payload := cc.CreatePayload("", common.TxType_ARCHIVE, syscontract.SystemContract_ARCHIVE_MANAGE.String(),
		syscontract.ArchiveFunction_RESTORE_BLOCK.String(), pairs, defaultSeq, nil)

	return payload, nil
}

func (cc *ChainClient) SignArchivePayload(payload *common.Payload) (*common.Payload, error) {
	return payload, nil
}

func (cc *ChainClient) SendArchiveBlockRequest(payload *common.Payload, timeout int64) (*common.TxResponse, error) {
	return cc.sendContractRequest(payload, nil, timeout, false)
}

func (cc *ChainClient) SendRestoreBlockRequest(payload *common.Payload, timeout int64) (*common.TxResponse, error) {
	return cc.sendContractRequest(payload, nil, timeout, false)
}

func (cc *ChainClient) GetArchivedFullBlockByHeight(blockHeight uint64) (*store.BlockWithRWSet, error) {
	fullBlock, err := cc.GetFromArchiveStore(blockHeight)
	if err != nil {
		return nil, err
	}

	return fullBlock, nil
}

func (cc *ChainClient) GetArchivedBlockByHeight(blockHeight uint64, withRWSet bool) (*common.BlockInfo, error) {
	fullBlock, err := cc.GetFromArchiveStore(blockHeight)
	if err != nil {
		return nil, err
	}

	blockInfo := &common.BlockInfo{
		Block: fullBlock.Block,
	}

	if withRWSet {
		blockInfo.RwsetList = fullBlock.TxRWSets
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetArchivedBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error) {
	blockHeight, err := cc.GetBlockHeightByTxId(txId)
	if err != nil {
		return nil, err
	}

	return cc.GetArchivedBlockByHeight(blockHeight, withRWSet)
}

func (cc *ChainClient) GetArchivedBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error) {
	blockHeight, err := cc.GetBlockHeightByHash(blockHash)
	if err != nil {
		return nil, err
	}

	return cc.GetArchivedBlockByHeight(blockHeight, withRWSet)
}

func (cc *ChainClient) GetArchivedTxByTxId(txId string) (*common.TransactionInfo, error) {
	blockHeight, err := cc.GetBlockHeightByTxId(txId)
	if err != nil {
		return nil, err
	}

	blockInfo, err := cc.GetArchivedBlockByHeight(blockHeight, false)
	if err != nil {
		return nil, err
	}

	for idx, tx := range blockInfo.Block.Txs {
		if tx.Payload.TxId == txId {
			return &common.TransactionInfo{
				Transaction: tx,
				BlockHeight: blockInfo.Block.Header.BlockHeight,
				BlockHash:   blockInfo.Block.Header.BlockHash,
				TxIndex:     uint32(idx),
			}, nil
		}
	}

	return nil, fmt.Errorf("CANNOT BE HERE! unknown tx [%s] in archive block [%d]", txId, blockHeight)
}

func (cc *ChainClient) GetFromArchiveStore(blockHeight uint64) (*store.BlockWithRWSet, error) {
	archiveType := utils.Config.ChainClientConfig.ArchiveConfig.Type
	if archiveType == "mysql" {
		return cc.GetArchivedBlockFromMySQL(blockHeight)
	}

	return nil, fmt.Errorf("unsupport archive type [%s]", archiveType)
}

func (cc *ChainClient) GetArchivedBlockFromMySQL(blockHeight uint64) (*store.BlockWithRWSet, error) {

	var (
		blockWithRWSetBytes []byte
		hmac                string
		blockWithRWSet      store.BlockWithRWSet
	)

	dest := utils.Config.ChainClientConfig.ArchiveConfig.Dest
	destList := strings.Split(dest, ":")
	if len(destList) != 4 {
		return nil, fmt.Errorf("invalid archive dest")
	}

	user, pwd, host, port := destList[0], destList[1], destList[2], destList[3]
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s_%s?charset=utf8mb4",
		user, pwd, host, port, utils.MysqlDBNamePrefix, cc.chainId))
	if err != nil {
		return nil, fmt.Errorf("mysql init failed, %s", err.Error())
	}
	defer db.Close()

	err = db.QueryRow(fmt.Sprintf("SELECT Fblock_with_rwset, Fhmac from %s_%d WHERE Fblock_height=?",
		utils.MysqlTableNamePrefix, blockHeight/utils.RowsPerBlockInfoTable+1), blockHeight).Scan(&blockWithRWSetBytes, &hmac)
	if err != nil {
		return nil, fmt.Errorf("select from mysql failed, %s", err.Error())
	}

	err = proto.Unmarshal(blockWithRWSetBytes, &blockWithRWSet)
	if err != nil {
		return nil, fmt.Errorf("unmarshal store.BlockWithRWSet failed, %s", err.Error())
	}

	return &blockWithRWSet, nil
}

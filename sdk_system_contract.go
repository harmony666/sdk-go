/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/discovery"
	"chainmaker.org/chainmaker/pb-go/v2/store"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/gogo/protobuf/proto"
)

func (cc *ChainClient) GetTxByTxId(txId string) (*common.TransactionInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]",
		syscontract.ChainQueryFunction_GET_TX_BY_TX_ID, txId)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_TX_BY_TX_ID.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractTxId,
				Value: []byte(txId),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	transactionInfo := &common.TransactionInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, transactionInfo); err != nil {
		return nil, fmt.Errorf("GetTxByTxId unmarshal transaction info payload failed, %s", err)
	}

	return transactionInfo, nil
}

func (cc *ChainClient) GetTxWithRWSetByTxId(txId string) (*common.TransactionInfoWithRWSet, error) {
	cc.logger.Debugf("[SDK] begin GetTxWithRWSetByTxId, [method:%s]/[txId:%s]",
		syscontract.ChainQueryFunction_GET_TX_BY_TX_ID, txId)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_TX_BY_TX_ID.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractTxId,
				Value: []byte(txId),
			},
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte("true"),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	tx := &common.TransactionInfoWithRWSet{}
	if err = proto.Unmarshal(resp.ContractResult.Result, tx); err != nil {
		return nil, fmt.Errorf("GetTxWithRWSetByTxId unmarshal failed, %s", err)
	}

	return tx, nil
}

func (cc *ChainClient) GetBlockByHeight(blockHeight uint64, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HEIGHT, blockHeight, withRWSet)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HEIGHT.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHeight,
				Value: []byte(strconv.FormatUint(blockHeight, 10)),
			},
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByHeight unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHash:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HASH, blockHash, withRWSet)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HASH.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHash,
				Value: []byte(blockHash),
			},
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByHash unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_BLOCK_BY_TX_ID, txId, withRWSet)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_BY_TX_ID.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractTxId,
				Value: []byte(txId),
			},
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByTxId unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetLastConfigBlock(withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_LAST_CONFIG_BLOCK, withRWSet)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_LAST_CONFIG_BLOCK.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByTxId unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetChainInfo() (*discovery.ChainInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]", syscontract.ChainQueryFunction_GET_CHAIN_INFO)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_CHAIN_INFO.String(), []*common.KeyValuePair{}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	chainInfo := &discovery.ChainInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainInfo); err != nil {
		return nil, fmt.Errorf("GetChainInfo unmarshal chain info payload failed, %s", err)
	}

	return chainInfo, nil
}

func (cc *ChainClient) GetNodeChainList() (*discovery.ChainList, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.ChainQueryFunction_GET_NODE_CHAIN_LIST)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_NODE_CHAIN_LIST.String(), []*common.KeyValuePair{}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	chainList := &discovery.ChainList{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainList); err != nil {
		return nil, fmt.Errorf("GetNodeChainList unmarshal chain list payload failed, %s", err)
	}

	return chainList, nil
}

func (cc *ChainClient) GetFullBlockByHeight(blockHeight uint64) (*store.BlockWithRWSet, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]",
		syscontract.ChainQueryFunction_GET_FULL_BLOCK_BY_HEIGHT, blockHeight)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_FULL_BLOCK_BY_HEIGHT.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHeight,
				Value: []byte(strconv.FormatUint(blockHeight, 10)),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	fullBlockInfo := &store.BlockWithRWSet{}
	if err = proto.Unmarshal(resp.ContractResult.Result, fullBlockInfo); err != nil {
		return nil, fmt.Errorf("GetFullBlockByHeight unmarshal block info payload failed, %s", err)
	}

	return fullBlockInfo, nil
}

func (cc *ChainClient) GetArchivedBlockHeight() (uint64, error) {
	return cc.getBlockHeight("", "")
}

func (cc *ChainClient) GetBlockHeightByTxId(txId string) (uint64, error) {
	return cc.getBlockHeight(txId, "")
}

func (cc *ChainClient) GetBlockHeightByHash(blockHash string) (uint64, error) {
	return cc.getBlockHeight("", blockHash)
}

func (cc *ChainClient) getBlockHeight(txId, blockHash string) (uint64, error) {
	var (
		txType       = common.TxType_QUERY_CONTRACT
		contractName = syscontract.SystemContract_CHAIN_QUERY.String()
		method       string
		pairs        []*common.KeyValuePair
	)

	if txId != "" {
		method = syscontract.ChainQueryFunction_GET_BLOCK_HEIGHT_BY_TX_ID.String()
		pairs = []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractTxId,
				Value: []byte(txId),
			},
		}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]", method, txId)
	} else if blockHash != "" {
		method = syscontract.ChainQueryFunction_GET_BLOCK_HEIGHT_BY_HASH.String()
		pairs = []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHash,
				Value: []byte(blockHash),
			},
		}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHash:%s]", method, blockHash)
	} else {
		method = syscontract.ChainQueryFunction_GET_ARCHIVED_BLOCK_HEIGHT.String()
		pairs = []*common.KeyValuePair{}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]", method)
	}

	payload := cc.CreatePayload("", txType, contractName, method, pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockHeight, err := strconv.ParseUint(string(resp.ContractResult.Result), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s, parse block height failed, %s", payload.TxType, err)
	}

	return blockHeight, nil
}

func (cc *ChainClient) GetLastBlock(withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_LAST_BLOCK, withRWSet)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_LAST_BLOCK.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetLastBlock unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetCurrentBlockHeight() (uint64, error) {
	block, err := cc.GetLastBlock(false)
	if err != nil {
		return 0, err
	}

	return block.Block.Header.BlockHeight, nil
}

func (cc *ChainClient) GetBlockHeaderByHeight(blockHeight uint64) (*common.BlockHeader, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]",
		syscontract.ChainQueryFunction_GET_BLOCK_HEADER_BY_HEIGHT, blockHeight)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_HEADER_BY_HEIGHT.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHeight,
				Value: []byte(strconv.FormatUint(blockHeight, 10)),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockHeader := &common.BlockHeader{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockHeader); err != nil {
		return nil, fmt.Errorf("GetBlockHeaderByHeight unmarshal block header payload failed, %s", err)
	}

	return blockHeader, nil
}

func (cc *ChainClient) InvokeSystemContract(contractName, method, txId string, kvs []*common.KeyValuePair,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	cc.logger.Debugf("[SDK] begin to INVOKE system contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, kvs)

	payload := cc.CreatePayload(txId, common.TxType_INVOKE_CONTRACT, contractName, method, kvs, defaultSeq, nil)

	return cc.sendContractRequest(payload, nil, timeout, withSyncResult)
}

func (cc *ChainClient) QuerySystemContract(contractName, method string, kvs []*common.KeyValuePair,
	timeout int64) (*common.TxResponse, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [contractName:%s]/[method:%s]/[params:%+v]",
		contractName, method, kvs)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, contractName, method, kvs, defaultSeq, nil)

	resp, err := cc.proposalRequestWithTimeout(payload, nil, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	return resp, nil
}

func (cc *ChainClient) GetMerklePathByTxId(txId string) ([]byte, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]",
		syscontract.ChainQueryFunction_GET_MERKLE_PATH_BY_TX_ID, txId)

	kvs := []*common.KeyValuePair{
		{
			Key:   utils.KeyBlockContractTxId,
			Value: []byte(txId),
		},
	}

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_MERKLE_PATH_BY_TX_ID.String(), kvs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}
	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) createNativeContractAccessPayload(method string,
	accessContractList []string) (*common.Payload, error) {
	val, err := json.Marshal(accessContractList)
	if err != nil {
		return nil, err
	}
	kvs := []*common.KeyValuePair{
		{
			Key:   syscontract.ContractAccess_NATIVE_CONTRACT_NAME.String(),
			Value: val,
		},
	}
	return cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CONTRACT_MANAGE.String(),
		method, kvs, defaultSeq, nil), nil
}

func (cc *ChainClient) CreateNativeContractAccessGrantPayload(grantContractList []string) (*common.Payload, error) {
	return cc.createNativeContractAccessPayload(syscontract.ContractManageFunction_GRANT_CONTRACT_ACCESS.String(),
		grantContractList)
}

func (cc *ChainClient) CreateNativeContractAccessRevokePayload(revokeContractList []string) (*common.Payload, error) {
	return cc.createNativeContractAccessPayload(syscontract.ContractManageFunction_REVOKE_CONTRACT_ACCESS.String(),
		revokeContractList)
}

func (cc *ChainClient) GetContractInfo(contractName string) (*common.Contract, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.ContractQueryFunction_GET_CONTRACT_INFO)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CONTRACT_MANAGE.String(),
		syscontract.ContractQueryFunction_GET_CONTRACT_INFO.String(), []*common.KeyValuePair{
			{
				Key:   syscontract.GetContractInfo_CONTRACT_NAME.String(),
				Value: []byte(contractName),
			},
		}, defaultSeq, nil,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	contract := &common.Contract{}
	if err = json.Unmarshal(resp.ContractResult.Result, contract); err != nil {
		return nil, fmt.Errorf("GetContractInfo unmarshal failed, %s", err)
	}

	return contract, nil
}

func (cc *ChainClient) GetContractList() ([]*common.Contract, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.ContractQueryFunction_GET_CONTRACT_LIST)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CONTRACT_MANAGE.String(),
		syscontract.ContractQueryFunction_GET_CONTRACT_LIST.String(), nil, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	var contracts []*common.Contract
	if err = json.Unmarshal(resp.ContractResult.Result, &contracts); err != nil {
		return nil, fmt.Errorf("GetContractList unmarshal failed, %s", err)
	}

	return contracts, nil
}

func (cc *ChainClient) GetDisabledNativeContractList() ([]string, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.ContractQueryFunction_GET_DISABLED_CONTRACT_LIST)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CONTRACT_MANAGE.String(),
		syscontract.ContractQueryFunction_GET_DISABLED_CONTRACT_LIST.String(), nil, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	var contractNames []string
	if err = json.Unmarshal(resp.ContractResult.Result, &contractNames); err != nil {
		return nil, fmt.Errorf("GetDisabledNativeContractList unmarshal failed, %s", err)
	}

	return contractNames, nil
}

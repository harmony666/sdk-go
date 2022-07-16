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

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/hibe"
	"chainmaker.org/chainmaker/common/v2/serialize"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

// hibe msg's Keys
const (
	// HibeMsgKey as a payload parameter
	HibeMsgKey = "hibe_msg"

	// HibeMsgIdKey Key as a hibeMsgMap parameter
	HibeMsgIdKey = "tx_id"

	// HibeMsgCipherTextKey Key as a hibeMsgMap parameter
	// The value of the key (CT) is the hibe_msg's message (ciphertext)
	HibeMsgCipherTextKey = "CT"

	// HibeParamsKey The value of the key (org_id) is the unique identifier of a HIBE params
	HibeParamsKey = "org_id"

	// HibeParamsValueKey The value of the key (params) is the Hibe's params
	HibeParamsValueKey = "params"
)

func (cc *ChainClient) CreateHibeInitParamsTxPayloadParams(orgId string,
	hibeParams []byte) ([]*common.KeyValuePair, error) {
	if err := hibe.ValidateId(orgId); err != nil {
		return nil, err
	}

	if len(hibeParams) == 0 {
		return nil, errors.New("invalid parameters, hibe params is nil")
	}

	payloadParams := []*common.KeyValuePair{
		{
			Key:   HibeParamsKey,
			Value: []byte(orgId),
		},
		{
			Key:   HibeParamsValueKey,
			Value: hibeParams,
		},
	}

	return payloadParams, nil
}

func (cc *ChainClient) CreateHibeTxPayloadParamsWithHibeParams(plaintext []byte, receiverIds []string,
	paramsBytesList [][]byte, txId string, keyType crypto.KeyType) ([]*common.KeyValuePair, error) {
	if len(paramsBytesList) == 0 {
		return nil, errors.New("invalid parameters, paramsBytesList is nil")
	}

	if len(receiverIds) != len(paramsBytesList) {
		return nil, errors.New("invalid parameters, receiverIds and paramsList do not match, place check them")
	}

	for _, paramsBytes := range paramsBytesList {
		if len(paramsBytes) == 0 {
			return nil, errors.New("invalid parameters, there are empty paramsBytes in the ParamsBytesList")
		}
	}

	paramsList := make([]*hibe.Params, len(paramsBytesList))
	for i, bytes := range paramsBytesList {
		params, ok := new(hibe.Params).Unmarshal(bytes)
		if !ok {
			return nil, errors.New("paramsBytesList unmarshal failed, please check it")
		}

		paramsList[i] = params
	}

	hibeMsg, err := hibe.EncryptHibeMsg(plaintext, receiverIds, paramsList, keyType)
	if err != nil {
		return nil, err
	}

	hibeMsgBytes, err := json.Marshal(hibeMsg)
	if err != nil {
		return nil, err
	}

	payloadParams := []*common.KeyValuePair{
		{
			Key:   HibeMsgIdKey,
			Value: []byte(txId),
		},
		{
			Key:   HibeMsgKey,
			Value: hibeMsgBytes,
		},
	}

	return payloadParams, nil
}

func (cc *ChainClient) CreateHibeTxPayloadParamsWithoutHibeParams(contractName, queryParamsMethod string,
	plaintext []byte, receiverIds []string, receiverOrgIds []string, txId string, keyType crypto.KeyType,
	timeout int64) ([]*common.KeyValuePair, error) {
	hibeParamsBytesList := make([][]byte, len(receiverOrgIds))
	for i, id := range receiverOrgIds {
		hibeParamsBytes, err := cc.QueryHibeParamsWithOrgId(contractName, queryParamsMethod, id, timeout)
		if err != nil {
			return nil, err
		}

		if len(hibeParamsBytes) == 0 {
			return nil, fmt.Errorf("no souch params of %s's org, please check it", id)
		}

		hibeParamsBytesList[i] = hibeParamsBytes
	}

	return cc.CreateHibeTxPayloadParamsWithHibeParams(plaintext, receiverIds, hibeParamsBytesList, txId, keyType)
}

func (cc *ChainClient) QueryHibeParamsWithOrgId(contractName, method, orgId string, timeout int64) ([]byte, error) {
	if err := hibe.ValidateId(orgId); err != nil {
		return nil, err
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   HibeParamsKey,
			Value: []byte(orgId),
		},
	}

	resp, err := cc.QueryContract(contractName, method, pairs, timeout)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", common.TxType_QUERY_CONTRACT.String(), err.Error())
	}
	if resp.ContractResult.Result == nil {
		return nil, errors.New("no such params, please check orgId")
	}

	// resp -> hibe.params
	result := serialize.EasyUnmarshal(resp.ContractResult.Result)
	resultMap := serialize.EasyCodecItemToParamsMap(result)

	hibeParamsBytes := resultMap[HibeParamsValueKey]
	hibeParams := new(hibe.Params)
	hibeParams, ok := hibeParams.Unmarshal(hibeParamsBytes)
	if !ok {
		return nil, errors.New("hibe.Params.Unmarshal failed")
	}
	return hibeParams.Marshal(), nil
}

func (cc *ChainClient) DecryptHibeTxByTxId(localId string, hibeParams []byte, hibePrvKey []byte, txId string,
	keyType crypto.KeyType) ([]byte, error) {
	if txId == "" {
		return nil, errors.New("invalid parameters, txId is empty")
	}

	transactionInfo, err := cc.GetTxByTxId(txId)
	if err != nil {
		return nil, err
	}

	return utils.DecryptHibeTx(localId, hibeParams, hibePrvKey, transactionInfo.Transaction, keyType)
}

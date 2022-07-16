/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/hibe"
	"chainmaker.org/chainmaker/pb-go/v2/common"
)

// hibe msg's Keys
const (
	HibeMsgKey           = "hibe_msg"
	HibeMsgIdKey         = "tx_id"
	HibeMsgCipherTextKey = "CT"
	HibeParamsKey        = "org_id"
	HibeParamsValueKey   = "params"
)

// ReadHibeParamsWithFilePath Returns the serialized byte array of hibeParams
func ReadHibeParamsWithFilePath(hibeParamsFilePath string) ([]byte, error) {
	paramsBytes, err := ioutil.ReadFile(hibeParamsFilePath)
	if err != nil {
		return nil, fmt.Errorf("open hibe params file failed, [err:%s]", err)
	}

	return paramsBytes, nil
}

// ReadHibePrvKeysWithFilePath Returns the serialized byte array of hibePrvKey
func ReadHibePrvKeysWithFilePath(hibePrvKeyFilePath string) ([]byte, error) {
	prvKeyBytes, err := ioutil.ReadFile(hibePrvKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("open hibe privateKey file failed, [err:%s]", err)
	}

	return prvKeyBytes, nil
}

func DecryptHibeTx(localId string, hibeParams []byte, hibePrvKey []byte, tx *common.Transaction,
	keyType crypto.KeyType) ([]byte, error) {
	localParams, ok := new(hibe.Params).Unmarshal(hibeParams)
	if !ok {
		return nil, errors.New("hibe.Params.Unmarshal failed, please check your file")
	}

	prvKey, ok := new(hibe.PrivateKey).Unmarshal(hibePrvKey)
	if !ok {
		return nil, errors.New("hibe.PrivateKey.Unmarshal failed, please check your file")
	}

	// get hibe_msg from tx
	// requestPayload := &common.Payload{}
	// err := proto.Unmarshal(tx.Payload, requestPayload)

	//if err != nil {
	//	return nil, err
	//}

	hibeMsgMap := make(map[string]string)
	// for _, item := range requestPayload.Parameters {
	for _, item := range tx.Payload.Parameters {
		if item.Key == HibeMsgKey {
			err := json.Unmarshal([]byte(item.Value), &hibeMsgMap)
			if err != nil {
				return nil, err
			}
		}
	}

	if hibeMsgMap == nil {
		return nil, errors.New("no such message, please check transaction")
	}

	return hibe.DecryptHibeMsg(localId, localParams, prvKey, hibeMsgMap, keyType)
}

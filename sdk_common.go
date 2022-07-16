/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package chainmaker_sdk_go

import (
	"fmt"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"

	"chainmaker.org/chainmaker/common/v2/crypto"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

const (
	// DefaultRetryLimit 默认轮训交易结果最大次数
	DefaultRetryLimit = 10
	// DefaultRetryInterval 默认每次轮训交易结果时的等待时间，单位ms
	DefaultRetryInterval = 500
	// defaultSeq default sequence
	defaultSeq = 0
)

func (cc *ChainClient) GetSyncResult(txId string) (*common.Result, error) {
	var (
		txInfo *common.TransactionInfo
		err    error
	)

	err = retry.Retry(func(uint) error {
		txInfo, err = cc.GetTxByTxId(txId)
		if err != nil {
			return err
		}

		return nil
	},
		strategy.Wait(time.Duration(cc.retryInterval)*time.Millisecond),
		strategy.Limit(uint(cc.retryLimit)),
	)

	if err != nil {
		return nil, fmt.Errorf("get tx by txId [%s] failed, %s", txId, err.Error())
	}

	if txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.Result == nil {
		return nil, fmt.Errorf("get result by txId [%s] failed, %+v", txId, txInfo)
	}

	return txInfo.Transaction.Result, nil
}

func (cc *ChainClient) sendContractRequest(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {

	resp, err := cc.proposalRequestWithTimeout(payload, endorsers, timeout)
	if err != nil {
		return resp, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if withSyncResult {
			result, err := cc.GetSyncResult(payload.TxId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}
			resp.Code = result.Code
			resp.Message = result.Message
			resp.ContractResult = result.ContractResult
			resp.TxId = payload.TxId
		}
	}

	return resp, nil
}

func (cc *ChainClient) CreatePayload(txId string, txType common.TxType, contractName, method string,
	kvs []*common.KeyValuePair, seq uint64, limit *common.Limit) *common.Payload {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	payload := utils.NewPayload(
		utils.WithChainId(cc.chainId),
		utils.WithTxType(txType),
		utils.WithTxId(txId),
		utils.WithTimestamp(time.Now().Unix()),
		utils.WithContractName(contractName),
		utils.WithMethod(method),
		utils.WithParameters(kvs),
		utils.WithSequence(seq),
		utils.WithLimit(limit),
	)

	return payload
}

func (cc *ChainClient) SignPayload(payload *common.Payload) (*common.EndorsementEntry, error) {
	var (
		sender    *accesscontrol.Member
		signBytes []byte
		err       error
	)
	if cc.authType == PermissionedWithCert {

		hashalgo, err := bcx509.GetHashFromSignatureAlgorithm(cc.userCrt.SignatureAlgorithm)
		if err != nil {
			return nil, fmt.Errorf("invalid algorithm: %s", err.Error())
		}

		signBytes, err = utils.SignPayloadWithHashType(cc.privateKey, hashalgo, payload)
		if err != nil {
			return nil, fmt.Errorf("SignPayload failed, %s", err)
		}

		sender = &accesscontrol.Member{
			OrgId:      cc.orgId,
			MemberInfo: cc.userCrtBytes,
			MemberType: accesscontrol.MemberType_CERT,
		}

	} else {
		signBytes, err = utils.SignPayloadWithHashType(cc.privateKey, crypto.HashAlgoMap[cc.hashType], payload)
		if err != nil {
			return nil, fmt.Errorf("SignPayload failed, %s", err.Error())
		}
		sender = &accesscontrol.Member{
			OrgId:      cc.orgId,
			MemberInfo: cc.pkBytes,
			MemberType: accesscontrol.MemberType_PUBLIC_KEY,
		}
	}

	entry := &common.EndorsementEntry{
		Signer:    sender,
		Signature: signBytes,
	}
	return entry, nil
}

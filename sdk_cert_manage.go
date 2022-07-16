/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func (cc *ChainClient) GetCertHash() ([]byte, error) {
	chainConfig, err := cc.GetChainConfig()

	if err != nil {
		return nil, fmt.Errorf("get cert hash failed, %s", err.Error())
	}

	certHash, err := utils.GetCertificateId(cc.userCrtBytes, chainConfig.Crypto.Hash)
	if err != nil {
		return nil, fmt.Errorf("calc cert hash failed, %s", err.Error())
	}

	return certHash, nil
}

func (cc *ChainClient) QueryCert(certHashes []string) (*common.CertInfos, error) {
	cc.logger.Infof("[SDK] begin to query cert, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(), syscontract.CertManageFunction_CERTS_QUERY.String())

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyCertHashes,
			Value: []byte(strings.Join(certHashes, ",")),
		},
	}

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CERT_MANAGE.String(),
		syscontract.CertManageFunction_CERTS_QUERY.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	certInfos := &common.CertInfos{}
	if err = proto.Unmarshal(resp.ContractResult.Result, certInfos); err != nil {
		return nil, fmt.Errorf("unmarshal cert infos payload failed, %s", err.Error())
	}

	return certInfos, nil
}

func (cc *ChainClient) AddCert() (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to add cert, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(), syscontract.CertManageFunction_CERT_ADD.String())

	certHash, err := cc.GetCertHash()
	if err != nil {
		return nil, fmt.Errorf("get cert hash in hex failed, %s", err.Error())
	}

	payload := cc.CreateCertManagePayload(syscontract.CertManageFunction_CERT_ADD.String(), nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	resp.ContractResult = &common.ContractResult{
		Code:   utils.SUCCESS,
		Result: certHash,
	}

	return resp, nil
}

func (cc *ChainClient) DeleteCert(certHashes []string) (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to delete cert, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(), syscontract.CertManageFunction_CERTS_DELETE.String())

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyCertHashes,
			Value: []byte(strings.Join(certHashes, ",")),
		},
	}

	payload := cc.CreateCertManagePayload(syscontract.CertManageFunction_CERTS_DELETE.String(), pairs)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) CreateCertManagePayload(method string, kvs []*common.KeyValuePair) *common.Payload {
	cc.logger.Debugf("[SDK] create CertManagePayload, method: %s", method)
	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CERT_MANAGE.String(),
		method, kvs, defaultSeq, nil)
	return payload
}

func (cc *ChainClient) CreateCertManageFrozenPayload(certs []string) *common.Payload {
	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyCerts,
			Value: []byte(strings.Join(certs, ",")),
		},
	}

	return cc.CreateCertManagePayload(syscontract.CertManageFunction_CERTS_FREEZE.String(), pairs)
}

func (cc *ChainClient) CreateCertManageUnfrozenPayload(certs []string) *common.Payload {
	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyCerts,
			Value: []byte(strings.Join(certs, ",")),
		},
	}

	return cc.CreateCertManagePayload(syscontract.CertManageFunction_CERTS_UNFREEZE.String(), pairs)
}

func (cc *ChainClient) CreateCertManageRevocationPayload(certCrl string) *common.Payload {
	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyCertCrl,
			Value: []byte(certCrl),
		},
	}

	return cc.CreateCertManagePayload(syscontract.CertManageFunction_CERTS_REVOKE.String(), pairs)
}

func (cc *ChainClient) SignCertManagePayload(payload *common.Payload) (*common.EndorsementEntry, error) {
	return cc.SignPayload(payload)
}

func (cc *ChainClient) SendCertManageRequest(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.sendContractRequest(payload, endorsers, timeout, withSyncResult)
}

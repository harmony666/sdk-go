/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package chainmaker_sdk_go

import (
	"fmt"
	"strings"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/gogo/protobuf/proto"
)

const (
	KEY_ALIAS   = "alias"
	KEY_ALIASES = "aliases"
	KEY_CERT    = "cert"
)

func (cc *ChainClient) AddAlias() (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to add alias, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(), syscontract.CertManageFunction_CERT_ALIAS_ADD.String())

	kvs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIAS,
			Value: []byte(cc.alias),
		},
	}

	payload := cc.CreateCertManagePayload(syscontract.CertManageFunction_CERT_ALIAS_ADD.String(), kvs)

	resp, err := cc.sendContractRequest(payload, nil, -1, true)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	resp.ContractResult = &common.ContractResult{
		Code:   utils.SUCCESS,
		Result: []byte(cc.alias),
	}

	return resp, nil
}

func (cc *ChainClient) QueryCertsAlias(aliases []string) (*common.AliasInfos, error) {
	cc.logger.Infof("[SDK] begin to query cert by aliases, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(), syscontract.CertManageFunction_CERTS_ALIAS_QUERY.String())

	kvs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIASES,
			Value: []byte(strings.Join(aliases, ",")),
		},
	}

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CERT_MANAGE.String(),
		syscontract.CertManageFunction_CERTS_ALIAS_QUERY.String(), kvs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	aliasInfos := &common.AliasInfos{}
	if err = proto.Unmarshal(resp.ContractResult.Result, aliasInfos); err != nil {
		return nil, fmt.Errorf("unmarshal cert alias infos payload failed, %s", err.Error())
	}

	return aliasInfos, nil
}

func (cc *ChainClient) CreateUpdateCertByAliasPayload(alias, newCertPEM string) *common.Payload {
	cc.logger.Debugf("[SDK] create [UpdateCertByAlias] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIAS,
			Value: []byte(alias),
		},
		{
			Key:   KEY_CERT,
			Value: []byte(newCertPEM),
		},
	}

	return cc.CreateCertManagePayload(syscontract.CertManageFunction_CERT_ALIAS_UPDATE.String(), pairs)
}

func (cc *ChainClient) SignUpdateCertByAliasPayload(payload *common.Payload) (*common.EndorsementEntry, error) {
	return cc.SignCertManagePayload(payload)
}

func (cc *ChainClient) UpdateCertByAlias(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.SendCertManageRequest(payload, endorsers, timeout, withSyncResult)
}

func (cc *ChainClient) CreateDeleteCertsAliasPayload(aliases []string) *common.Payload {
	cc.logger.Debugf("[SDK] create [DeleteAlias] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIASES,
			Value: []byte(strings.Join(aliases, ",")),
		},
	}

	return cc.CreateCertManagePayload(syscontract.CertManageFunction_CERTS_ALIAS_DELETE.String(), pairs)
}

func (cc *ChainClient) SignDeleteAliasPayload(payload *common.Payload) (*common.EndorsementEntry, error) {
	return cc.SignCertManagePayload(payload)
}

func (cc *ChainClient) DeleteCertsAlias(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.SendCertManageRequest(payload, endorsers, timeout, withSyncResult)
}

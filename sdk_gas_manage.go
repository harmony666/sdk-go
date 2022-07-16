/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"errors"
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func (cc *ChainClient) CreateSetGasAdminPayload(address string) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateSetGasAdminPayload] payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasAddressKey,
			Value: []byte(address),
		},
	}

	return cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_SET_ADMIN.String(), pairs, defaultSeq, nil), nil
}

func (cc *ChainClient) GetGasAdmin() (string, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.GasAccountFunction_GET_ADMIN)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_GET_ADMIN.String(), nil, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return "", fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return "", fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	return string(resp.ContractResult.Result), nil
}

func (cc *ChainClient) CreateRechargeGasPayload(rechargeGasList []*syscontract.RechargeGas) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateRechargeGasPayload] payload")

	rechargeGasReq := syscontract.RechargeGasReq{BatchRechargeGas: rechargeGasList}
	bz, err := rechargeGasReq.Marshal()
	if err != nil {
		return nil, err
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasBatchRecharge,
			Value: bz,
		},
	}

	return cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_RECHARGE_GAS.String(), pairs, defaultSeq, nil), nil
}

func (cc *ChainClient) GetGasBalance(address string) (int64, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.GasAccountFunction_GET_BALANCE)

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_GET_BALANCE.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyGasAddressKey,
				Value: []byte(address),
			},
		}, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err := utils.CheckProposalRequestResp(resp, true); err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	balance, err := strconv.ParseInt(string(resp.ContractResult.Result), 10, 64)
	if err != nil {
		return 0, fmt.Errorf(errStringFormat, "strconv.ParseInt", err)
	}

	return balance, nil
}

func (cc *ChainClient) CreateRefundGasPayload(address string, amount int64) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateRefundGasPayload] payload")

	if amount <= 0 {
		return nil, errors.New("amount must > 0")
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasAddressKey,
			Value: []byte(address),
		},
		{
			Key:   utils.KeyGasChargeGasAmount,
			Value: []byte(strconv.FormatInt(amount, 10)),
		},
	}

	return cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_REFUND_GAS.String(), pairs, defaultSeq, nil), nil
}

func (cc *ChainClient) CreateFrozenGasAccountPayload(address string) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateFrozenGasAccountPayload] payload")

	return cc.createFrozenUnfrozenGasAccountPayload(syscontract.GasAccountFunction_FROZEN_ACCOUNT.String(), address)
}

func (cc *ChainClient) CreateUnfrozenGasAccountPayload(address string) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateFrozenGasAccountPayload] payload")

	return cc.createFrozenUnfrozenGasAccountPayload(syscontract.GasAccountFunction_UNFROZEN_ACCOUNT.String(), address)
}

func (cc *ChainClient) GetGasAccountStatus(address string) (bool, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.GasAccountFunction_ACCOUNT_STATUS)

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasAddressKey,
			Value: []byte(address),
		},
	}

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_ACCOUNT_STATUS.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return false, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err := utils.CheckProposalRequestResp(resp, true); err != nil {
		return false, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	return string(resp.ContractResult.Result) == "0", nil
}

func (cc *ChainClient) SendGasManageRequest(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	cc.logger.Debug("[SDK] begin to get chain config")
	return cc.sendContractRequest(payload, endorsers, timeout, withSyncResult)
}

func (cc *ChainClient) createFrozenUnfrozenGasAccountPayload(method string,
	address string) (*common.Payload, error) {

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasAddressKey,
			Value: []byte(address),
		},
	}

	return cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		method, pairs, defaultSeq, nil), nil
}

func (cc *ChainClient) AttachGasLimit(payload *common.Payload, limit *common.Limit) *common.Payload {
	payload.Limit = limit
	return payload
}

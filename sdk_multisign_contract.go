/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func (cc *ChainClient) MultiSignContractReq(payload *common.Payload) (*common.TxResponse, error) {

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) MultiSignContractVote(multiSignReqPayload *common.Payload,
	endorser *common.EndorsementEntry, isAgree bool) (*common.TxResponse, error) {
	agree := syscontract.VoteStatus_AGREE
	if !isAgree {
		agree = syscontract.VoteStatus_REJECT
	}
	msvi := &syscontract.MultiSignVoteInfo{
		Vote:        agree,
		Endorsement: endorser,
	}
	msviByte, _ := msvi.Marshal()
	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiVote_VOTE_INFO.String(),
			Value: msviByte,
		},
		{
			Key:   syscontract.MultiVote_TX_ID.String(),
			Value: []byte(multiSignReqPayload.TxId),
		},
	}
	payload := cc.createMultiSignVotePayload(pairs)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) MultiSignContractQuery(txId string) (*common.TxResponse, error) {

	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiVote_TX_ID.String(),
			Value: []byte(txId),
		},
	}
	payload := cc.createMultiSignQueryPayload(pairs)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) CreateMultiSignReqPayload(pairs []*common.KeyValuePair) *common.Payload {

	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_REQ.String(), pairs, defaultSeq, nil)
	return payload
}

func (cc *ChainClient) createMultiSignVotePayload(pairs []*common.KeyValuePair) *common.Payload {
	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_VOTE.String(), pairs, defaultSeq, nil)

	return payload
}

func (cc *ChainClient) createMultiSignQueryPayload(pairs []*common.KeyValuePair) *common.Payload {
	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_QUERY.String(), pairs, defaultSeq, nil)

	return payload
}

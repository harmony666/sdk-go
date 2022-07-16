/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"
	"io"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
)

func (cc *ChainClient) SubscribeBlock(ctx context.Context, startBlock, endBlock int64, withRWSet,
	onlyHeader bool) (<-chan interface{}, error) {

	payload := cc.CreateSubscribeBlockPayload(startBlock, endBlock, withRWSet, onlyHeader)

	return cc.Subscribe(ctx, payload)
}

func (cc *ChainClient) SubscribeTx(ctx context.Context, startBlock, endBlock int64, contractName string,
	txIds []string) (<-chan interface{}, error) {

	payload := cc.CreateSubscribeTxPayload(startBlock, endBlock, contractName, txIds)

	return cc.Subscribe(ctx, payload)
}

func (cc *ChainClient) SubscribeContractEvent(ctx context.Context, startBlock, endBlock int64,
	contractName, topic string) (<-chan interface{}, error) {

	payload := cc.CreateSubscribeContractEventPayload(startBlock, endBlock, contractName, topic)

	return cc.Subscribe(ctx, payload)
}

func (cc *ChainClient) Subscribe(ctx context.Context, payload *common.Payload) (<-chan interface{}, error) {

	req, err := cc.GenerateTxRequest(payload, nil)
	if err != nil {
		return nil, err
	}

	client, err := cc.pool.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.rpcNode.Subscribe(ctx, req, grpc.MaxCallSendMsgSize(client.rpcMaxSendMsgSize),
		grpc.MaxCallRecvMsgSize(client.rpcMaxRecvMsgSize))
	if err != nil {
		return nil, err
	}

	c := make(chan interface{})
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := resp.Recv()
				if err == io.EOF {
					cc.logger.Debugf("[SDK] Subscriber got EOF and stop recv msg")
					return
				}

				if err != nil {
					cc.logger.Errorf("[SDK] Subscriber receive failed, %s", err)
					return
				}

				var ret interface{}
				switch payload.Method {
				case syscontract.SubscribeFunction_SUBSCRIBE_BLOCK.String():
					blockInfo := &common.BlockInfo{}
					if err = proto.Unmarshal(result.Data, blockInfo); err == nil {
						ret = blockInfo
						break
					}

					blockHeader := &common.BlockHeader{}
					if err = proto.Unmarshal(result.Data, blockHeader); err == nil {
						ret = blockHeader
						break
					}

					cc.logger.Error("[SDK] Subscriber receive block failed, %s", err)
					close(c)
					return
				case syscontract.SubscribeFunction_SUBSCRIBE_TX.String():
					tx := &common.Transaction{}
					if err = proto.Unmarshal(result.Data, tx); err != nil {
						cc.logger.Error("[SDK] Subscriber receive tx failed, %s", err)
						close(c)
						return
					}
					ret = tx
				case syscontract.SubscribeFunction_SUBSCRIBE_CONTRACT_EVENT.String():
					events := &common.ContractEventInfoList{}
					if err = proto.Unmarshal(result.Data, events); err != nil {
						cc.logger.Error("[SDK] Subscriber receive contract event failed, %s", err)
						close(c)
						return
					}
					for _, event := range events.ContractEvents {
						c <- event
					}
					continue

				default:
					ret = result.Data
				}

				c <- ret
			}
		}
	}()

	return c, nil
}

func (cc *ChainClient) CreateSubscribeBlockPayload(startBlock, endBlock int64,
	withRWSet, onlyHeader bool) *common.Payload {

	return cc.CreatePayload("", common.TxType_SUBSCRIBE, syscontract.SystemContract_SUBSCRIBE_MANAGE.String(),
		syscontract.SubscribeFunction_SUBSCRIBE_BLOCK.String(), []*common.KeyValuePair{
			{
				Key:   syscontract.SubscribeBlock_START_BLOCK.String(),
				Value: utils.I64ToBytes(startBlock),
			},
			{
				Key:   syscontract.SubscribeBlock_END_BLOCK.String(),
				Value: utils.I64ToBytes(endBlock),
			},
			{
				Key:   syscontract.SubscribeBlock_WITH_RWSET.String(),
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
			{
				Key:   syscontract.SubscribeBlock_ONLY_HEADER.String(),
				Value: []byte(strconv.FormatBool(onlyHeader)),
			},
		}, defaultSeq, nil,
	)
}

func (cc *ChainClient) CreateSubscribeTxPayload(startBlock, endBlock int64,
	contractName string, txIds []string) *common.Payload {

	return cc.CreatePayload("", common.TxType_SUBSCRIBE, syscontract.SystemContract_SUBSCRIBE_MANAGE.String(),
		syscontract.SubscribeFunction_SUBSCRIBE_TX.String(), []*common.KeyValuePair{
			{
				Key:   syscontract.SubscribeTx_START_BLOCK.String(),
				Value: utils.I64ToBytes(startBlock),
			},
			{
				Key:   syscontract.SubscribeTx_END_BLOCK.String(),
				Value: utils.I64ToBytes(endBlock),
			},
			{
				Key:   syscontract.SubscribeTx_CONTRACT_NAME.String(),
				Value: []byte(contractName),
			},
			{
				Key:   syscontract.SubscribeTx_TX_IDS.String(),
				Value: []byte(strings.Join(txIds, ",")),
			},
		}, defaultSeq, nil,
	)
}

func (cc *ChainClient) CreateSubscribeContractEventPayload(startBlock, endBlock int64,
	contractName, topic string) *common.Payload {

	return cc.CreatePayload("", common.TxType_SUBSCRIBE, syscontract.SystemContract_SUBSCRIBE_MANAGE.String(),
		syscontract.SubscribeFunction_SUBSCRIBE_CONTRACT_EVENT.String(), []*common.KeyValuePair{
			{
				Key:   syscontract.SubscribeContractEvent_START_BLOCK.String(),
				Value: utils.I64ToBytes(startBlock),
			},
			{
				Key:   syscontract.SubscribeContractEvent_END_BLOCK.String(),
				Value: utils.I64ToBytes(endBlock),
			},
			{
				Key:   syscontract.SubscribeContractEvent_CONTRACT_NAME.String(),
				Value: []byte(contractName),
			},
			{
				Key:   syscontract.SubscribeContractEvent_TOPIC.String(),
				Value: []byte(topic),
			},
		}, defaultSeq, nil,
	)
}

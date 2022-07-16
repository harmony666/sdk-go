/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	commonPb "chainmaker.org/chainmaker/pb-go/v2/common"
)

type PayloadOption func(*commonPb.Payload)

func NewPayload(opts ...PayloadOption) *commonPb.Payload {
	config := &commonPb.Payload{}
	for _, opt := range opts {
		opt(config)
	}

	return config
}

// WithChainId set chainId of payload
func WithChainId(chainId string) PayloadOption {
	return func(config *commonPb.Payload) {
		config.ChainId = chainId
	}
}

// WithTxType set TxType of payload
func WithTxType(txType commonPb.TxType) PayloadOption {
	return func(config *commonPb.Payload) {
		config.TxType = txType
	}
}

// WithTxId set TxId of payload
func WithTxId(txId string) PayloadOption {
	return func(config *commonPb.Payload) {
		config.TxId = txId
	}
}

// WithTimestamp set Timestamp of payload
func WithTimestamp(timestamp int64) PayloadOption {
	return func(config *commonPb.Payload) {
		config.Timestamp = timestamp
	}
}

// WithExpirationTime set ExpirationTime of payload
func WithExpirationTime(expirationTime int64) PayloadOption {
	return func(config *commonPb.Payload) {
		config.ExpirationTime = expirationTime
	}
}

// WithContractName set ContractName of payload
func WithContractName(contractName string) PayloadOption {
	return func(config *commonPb.Payload) {
		config.ContractName = contractName
	}
}

// WithMethod set Method of payload
func WithMethod(method string) PayloadOption {
	return func(config *commonPb.Payload) {
		config.Method = method
	}
}

// WithParameters set Parameters of payload
func WithParameters(parameters []*commonPb.KeyValuePair) PayloadOption {
	return func(config *commonPb.Payload) {
		config.Parameters = parameters
	}
}

// AddParameter add one Parameter of payload
func AddParameter(parameter *commonPb.KeyValuePair) PayloadOption {
	return func(config *commonPb.Payload) {
		config.Parameters = append(config.Parameters, parameter)
	}
}

// WithSequence set Sequence of payload
func WithSequence(sequence uint64) PayloadOption {
	return func(config *commonPb.Payload) {
		config.Sequence = sequence
	}
}

// WithLimit set Limit of payload
func WithLimit(limit *commonPb.Limit) PayloadOption {
	return func(config *commonPb.Payload) {
		config.Limit = limit
	}
}

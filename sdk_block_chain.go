/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/config"
	"google.golang.org/grpc"
)

func (cc *ChainClient) CheckNewBlockChainConfig() error {
	cc.logger.Debug("[SDK] begin to send check new block chain config command")
	req := &config.CheckNewBlockChainConfigRequest{}
	client, err := cc.pool.getClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := client.rpcNode.CheckNewBlockChainConfig(ctx, req, grpc.MaxCallSendMsgSize(client.rpcMaxSendMsgSize))
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("check new block chain config failed, %s", res.Message)
	}
	return nil
}

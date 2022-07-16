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

func (cc *ChainClient) GetChainMakerServerVersion() (string, error) {
	cc.logger.Debug("[SDK] begin to get chainmaker server version")
	req := &config.ChainMakerVersionRequest{}
	client, err := cc.pool.getClient()
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	res, err := client.rpcNode.GetChainMakerVersion(ctx, req, grpc.MaxCallSendMsgSize(client.rpcMaxSendMsgSize))
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf("get chainmaker server version failed, %s", res.Message)
	}
	return res.Version, nil
}

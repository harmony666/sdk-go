/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithNodeAddr(t *testing.T) {
	addr := "127.0.0.1:12301"
	opt := WithNodeAddr(addr)
	var config = &NodeConfig{}
	opt(config)
	require.Equal(t, config.addr, addr)
}

func TestWithNodeConnCnt(t *testing.T) {
	connCnt := 10
	opt := WithNodeConnCnt(connCnt)
	var config = &NodeConfig{}
	opt(config)
	require.Equal(t, config.connCnt, connCnt)
}

func TestWithNodeUseTLS(t *testing.T) {
	opt := WithNodeUseTLS(true)
	var config = &NodeConfig{}
	opt(config)
	require.Equal(t, config.useTLS, true)
}

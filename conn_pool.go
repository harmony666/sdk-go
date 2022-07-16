/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"chainmaker.org/chainmaker/common/v2/ca"
	"chainmaker.org/chainmaker/pb-go/v2/api"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	networkClientRetryInterval = 500 // 获取可用客户端连接对象重试时间间隔，单位：ms
	networkClientRetryLimit    = 5   // 获取可用客户端连接对象最大重试次数
)

var _ ConnectionPool = (*ClientConnectionPool)(nil)

type ConnectionPool interface {
	initGRPCConnect(nodeAddr string, useTLS bool, caPaths, caCerts []string, tlsHostName string) (*grpc.ClientConn, error)
	getClient() (*networkClient, error)
	getClientWithIgnoreAddrs(ignoreAddrs map[string]struct{}) (*networkClient, error)
	getLogger() utils.Logger
	Close() error
}

// 客户端连接结构定义
type networkClient struct {
	rpcNode           api.RpcNodeClient
	conn              *grpc.ClientConn
	nodeAddr          string
	useTLS            bool
	caPaths           []string
	caCerts           []string
	tlsHostName       string
	ID                string
	rpcMaxRecvMsgSize int
	rpcMaxSendMsgSize int
}

func (cli *networkClient) sendRequestWithTimeout(txReq *common.TxRequest, timeout int64) (*common.TxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel() // releases resources if SendRequest completes before timeout elapses
	return cli.rpcNode.SendRequest(ctx, txReq, grpc.MaxCallSendMsgSize(cli.rpcMaxSendMsgSize))
}

// ClientConnectionPool 客户端连接池结构定义
type ClientConnectionPool struct {
	connections       []*networkClient
	logger            utils.Logger
	userKeyBytes      []byte
	userCrtBytes      []byte
	rpcMaxRecvMsgSize int
	rpcMaxSendMsgSize int
}

// NewConnPool 创建连接池
func NewConnPool(config *ChainClientConfig) (*ClientConnectionPool, error) {
	pool := &ClientConnectionPool{
		logger:            config.logger,
		userKeyBytes:      config.userKeyBytes,
		userCrtBytes:      config.userCrtBytes,
		rpcMaxRecvMsgSize: config.rpcClientConfig.rpcClientMaxReceiveMessageSize * 1024 * 1024,
		rpcMaxSendMsgSize: config.rpcClientConfig.rpcClientMaxSendMessageSize * 1024 * 1024,
	}

	for idx, node := range config.nodeList {
		for i := 0; i < node.connCnt; i++ {
			cli := &networkClient{
				nodeAddr:          node.addr,
				useTLS:            node.useTLS,
				caPaths:           node.caPaths,
				caCerts:           node.caCerts,
				tlsHostName:       node.tlsHostName,
				ID:                fmt.Sprintf("%v-%v-%v", idx, node.addr, node.tlsHostName),
				rpcMaxRecvMsgSize: pool.rpcMaxRecvMsgSize,
				rpcMaxSendMsgSize: pool.rpcMaxSendMsgSize,
			}
			pool.connections = append(pool.connections, cli)
		}
	}

	// 打散，用作负载均衡
	pool.connections = shuffle(pool.connections)

	return pool, nil
}

// 初始化GPRC客户端连接
func (pool *ClientConnectionPool) initGRPCConnect(nodeAddr string, useTLS bool, caPaths, caCerts []string,
	tlsHostName string) (*grpc.ClientConn, error) {
	var tlsClient ca.CAClient
	if useTLS {
		if len(caCerts) != 0 {
			tlsClient = ca.CAClient{
				ServerName: tlsHostName,
				CaCerts:    caCerts,
				CertBytes:  pool.userCrtBytes,
				KeyBytes:   pool.userKeyBytes,
				Logger:     pool.logger,
			}
		} else {
			tlsClient = ca.CAClient{
				ServerName: tlsHostName,
				CaPaths:    caPaths,
				CertBytes:  pool.userCrtBytes,
				KeyBytes:   pool.userKeyBytes,
				Logger:     pool.logger,
			}
		}

		c, err := tlsClient.GetCredentialsByCA()
		if err != nil {
			return nil, err
		}
		return grpc.Dial(nodeAddr, grpc.WithTransportCredentials(*c),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(pool.rpcMaxRecvMsgSize)))
	}
	return grpc.Dial(nodeAddr, grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(pool.rpcMaxRecvMsgSize)))
}

// 获取空闲的可用客户端连接对象
func (pool *ClientConnectionPool) getClient() (*networkClient, error) {
	return pool.getClientWithIgnoreAddrs(nil)
}

func (pool *ClientConnectionPool) getClientWithIgnoreAddrs(ignoreAddrs map[string]struct{}) (*networkClient, error) {
	var nc *networkClient

	err := retry.Retry(func(uint) error {
		for _, cli := range pool.connections {

			if ignoreAddrs != nil {
				if _, ok := ignoreAddrs[cli.ID]; ok {
					continue
				}
			}

			if cli.conn == nil || cli.conn.GetState() == connectivity.Shutdown {

				conn, err := pool.initGRPCConnect(cli.nodeAddr, cli.useTLS, cli.caPaths, cli.caCerts, cli.tlsHostName)
				if err != nil {
					pool.logger.Errorf("init grpc connection [nodeAddr:%s] failed, %s", cli.ID, err.Error())
					continue
				}

				cli.conn = conn
				cli.rpcNode = api.NewRpcNodeClient(conn)
				nc = cli
				return nil
			}

			s := cli.conn.GetState()
			if s == connectivity.Idle || s == connectivity.Ready {
				nc = cli
				return nil
			}
		}

		return fmt.Errorf("all client connections are busy")

	}, strategy.Wait(networkClientRetryInterval*time.Millisecond), strategy.Limit(networkClientRetryLimit))

	if err != nil {
		return nil, err
	}

	return nc, nil
}

func (pool *ClientConnectionPool) getLogger() utils.Logger {
	return pool.logger
}

// Close 关闭连接池
func (pool *ClientConnectionPool) Close() error {
	for _, c := range pool.connections {
		if c.conn == nil {
			continue
		}

		if err := c.conn.Close(); err != nil {
			pool.logger.Errorf("stop %s connection failed, %s",
				c.nodeAddr, err.Error())

			continue
		}
	}

	return nil
}

//nolint
// 数组打散
func shuffle(vals []*networkClient) []*networkClient {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]*networkClient, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}

	return ret
}

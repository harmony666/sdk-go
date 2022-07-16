/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"errors"
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker/common/v2/cert"
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	"chainmaker.org/chainmaker/common/v2/crypto/pkcs11"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/common/v2/log"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"go.uber.org/zap"
)

const (
	// MaxConnCnt 单ChainMaker节点最大连接数
	MaxConnCnt = 1024
	// DefaultGetTxTimeout 查询交易超时时间
	DefaultGetTxTimeout = 10
	// DefaultSendTxTimeout 发送交易超时时间
	DefaultSendTxTimeout = 10
	// DefaultRpcClientMaxReceiveMessageSize 默认grpc客户端接收message最大值 4M
	DefaultRpcClientMaxReceiveMessageSize = 4
	// DefaultRpcClientMaxSendMessageSize 默认grpc客户端发送message最大值 4M
	DefaultRpcClientMaxSendMessageSize = 4
)

var (
	// global thread-safe pkcs11 handler
	p11Handle *pkcs11.P11Handle
)

func GetP11Handle() *pkcs11.P11Handle {
	return p11Handle
}

// NodeConfig 节点配置
type NodeConfig struct {
	// 必填项
	// 节点地址
	addr string
	// 节点连接数
	connCnt int
	// 选填项
	// 是否启用TLS认证
	useTLS bool
	// CA ROOT证书路径
	caPaths []string
	// CA ROOT证书内容（同时配置caPaths和caCerts以caCerts为准）
	caCerts []string
	// TLS hostname
	tlsHostName string
}

type NodeOption func(config *NodeConfig)

// WithNodeAddr 设置节点地址
func WithNodeAddr(addr string) NodeOption {
	return func(config *NodeConfig) {
		config.addr = addr
	}
}

// WithNodeConnCnt 设置节点连接数
func WithNodeConnCnt(connCnt int) NodeOption {
	return func(config *NodeConfig) {
		config.connCnt = connCnt
	}
}

// WithNodeUseTLS 设置是否启动TLS开关
func WithNodeUseTLS(useTLS bool) NodeOption {
	return func(config *NodeConfig) {
		config.useTLS = useTLS
	}
}

// WithNodeCAPaths 添加CA证书路径
func WithNodeCAPaths(caPaths []string) NodeOption {
	return func(config *NodeConfig) {
		config.caPaths = caPaths
	}
}

// WithNodeCACerts 添加CA证书内容
func WithNodeCACerts(caCerts []string) NodeOption {
	return func(config *NodeConfig) {
		config.caCerts = caCerts
	}
}

func WithNodeTLSHostName(tlsHostName string) NodeOption {
	return func(config *NodeConfig) {
		config.tlsHostName = tlsHostName
	}
}

// ArchiveConfig Archive配置
type ArchiveConfig struct {
	// 非必填
	// secret key
	secretKey string
}

type ArchiveOption func(config *ArchiveConfig)

// WithSecretKey 设置Archive的secret key
func WithSecretKey(key string) ArchiveOption {
	return func(config *ArchiveConfig) {
		config.secretKey = key
	}
}

// RPCClientConfig RPC Client 链接配置
type RPCClientConfig struct {
	// grpc客户端接收和发送消息时，允许单条message大小的最大值(MB)
	rpcClientMaxReceiveMessageSize, rpcClientMaxSendMessageSize int
}

type RPCClientOption func(config *RPCClientConfig)

// WithRPCClientMaxReceiveMessageSize 设置RPC Client的Max Receive Message Size
func WithRPCClientMaxReceiveMessageSize(size int) RPCClientOption {
	return func(config *RPCClientConfig) {
		config.rpcClientMaxReceiveMessageSize = size
	}
}

// WithRPCClientMaxSendMessageSize 设置RPC Client的Max Send Message Size
func WithRPCClientMaxSendMessageSize(size int) RPCClientOption {
	return func(config *RPCClientConfig) {
		config.rpcClientMaxSendMessageSize = size
	}
}

// Pkcs11Config pkcs11配置
type Pkcs11Config struct {
	// 是否开启pkcs11, 如果为 ture 则下面所有的字段都是必填
	Enabled bool
	// path to the .so file of pkcs11 interface
	Library string
	// label for the slot to be used
	Label string
	// password to logon the HSM(Hardware security module)
	Password string
	// size of HSM session cache
	SessionCacheSize int
	// hash algorithm used to compute SKI, eg, SHA256
	Hash string
}

type AuthType uint32

const (
	// permissioned with certificate
	PermissionedWithCert AuthType = iota + 1

	// permissioned with public key
	PermissionedWithKey

	// public key
	Public
)

const (
	// DefaultAuthType is default cert auth type
	DefaultAuthType = ""
)

var AuthTypeToStringMap = map[AuthType]string{
	PermissionedWithCert: "permissionedwithcert",
	PermissionedWithKey:  "permissionedwithkey",
	Public:               "public",
}

var StringToAuthTypeMap = map[string]AuthType{
	"permissionedwithcert": PermissionedWithCert,
	"permissionedwithkey":  PermissionedWithKey,
	"public":               Public,
}

type ChainClientConfig struct {
	// logger若不设置，将采用默认日志文件输出日志，建议设置，以便采用集成系统的统一日志输出
	logger utils.Logger

	// 链客户端相关配置
	// 方式1：配置文件指定（方式1与方式2可以同时使用，参数指定的值会覆盖配置文件中的配置）
	confPath string

	// 方式2：参数指定（方式1与方式2可以同时使用，参数指定的值会覆盖配置文件中的配置）
	orgId    string
	chainId  string
	nodeList []*NodeConfig

	// 以下xxxPath和xxxBytes同时指定的话，优先使用Bytes
	userKeyFilePath     string
	userCrtFilePath     string
	userSignKeyFilePath string // 公钥模式下使用该字段
	userSignCrtFilePath string

	userKeyBytes     []byte
	userCrtBytes     []byte
	userSignKeyBytes []byte // 公钥模式下使用该字段
	userSignCrtBytes []byte

	// 以下字段为经过处理后的参数
	privateKey crypto.PrivateKey // 证书和公钥身份模式都使用该字段存储私钥

	// 证书模式下
	userCrt *bcx509.Certificate

	// 公钥模式下
	userPk crypto.PublicKey
	crypto *CryptoConfig

	// 归档特性的配置
	archiveConfig *ArchiveConfig

	// rpc客户端设置
	rpcClientConfig *RPCClientConfig

	// pkcs11的配置
	pkcs11Config *Pkcs11Config

	// AuthType
	authType AuthType

	// retry config
	retryLimit    int // if <=0 then use DefaultRetryLimit
	retryInterval int // if <=0 then use DefaultRetryInterval

	// alias
	alias string

	enableNormalKey bool
}

type CryptoConfig struct {
	hash string
}

type CryptoOption func(config *CryptoConfig)

// WithHashType 公钥模式下：添加用户哈希算法配置
func WithHashAlgo(hashType string) CryptoOption {
	return func(config *CryptoConfig) {
		config.hash = hashType
	}
}

type ChainClientOption func(*ChainClientConfig)

func WithAuthType(authType string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.authType = StringToAuthTypeMap[authType]
	}
}

func WithEnableNormalKey(enableNormalKey bool) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.enableNormalKey = enableNormalKey
	}
}

// WithConfPath 设置配置文件路径
func WithConfPath(confPath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.confPath = confPath
	}
}

// AddChainClientNodeConfig 添加ChainMaker节点地址及连接数配置
func AddChainClientNodeConfig(nodeConfig *NodeConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.nodeList = append(config.nodeList, nodeConfig)
	}
}

// WithUserKeyFilePath 添加用户私钥文件路径配置
func WithUserKeyFilePath(userKeyFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userKeyFilePath = userKeyFilePath
	}
}

// WithUserCrtFilePath 添加用户证书文件路径配置
func WithUserCrtFilePath(userCrtFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userCrtFilePath = userCrtFilePath
	}
}

// WithUserSignKeyFilePath 添加用户签名私钥文件路径配置
func WithUserSignKeyFilePath(userSignKeyFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignKeyFilePath = userSignKeyFilePath
	}
}

// WithUserSignCrtFilePath 添加用户签名证书文件路径配置
func WithUserSignCrtFilePath(userSignCrtFilePath string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignCrtFilePath = userSignCrtFilePath
	}
}

// WithUserKeyBytes 添加用户私钥文件内容配置
func WithUserKeyBytes(userKeyBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userKeyBytes = userKeyBytes
	}
}

// WithUserCrtBytes 添加用户证书文件内容配置
func WithUserCrtBytes(userCrtBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userCrtBytes = userCrtBytes
	}
}

// WithUserSignKeyBytes 添加用户签名私钥文件内容配置
func WithUserSignKeyBytes(userSignKeyBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignKeyBytes = userSignKeyBytes
	}
}

// WithUserSignCrtBytes 添加用户签名证书文件内容配置
func WithUserSignCrtBytes(userSignCrtBytes []byte) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.userSignCrtBytes = userSignCrtBytes
	}
}

// WithChainClientOrgId 添加OrgId
func WithChainClientOrgId(orgId string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.orgId = orgId
	}
}

// WithChainClientChainId 添加ChainId
func WithChainClientChainId(chainId string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.chainId = chainId
	}
}

// WithRetryLimit 设置 chain client 同步模式下，轮训获取交易结果时的最大轮训次数
func WithRetryLimit(limit int) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.retryLimit = limit
	}
}

// WithRetryInterval 设置 chain client 同步模式下，每次轮训交易结果时的等待时间，单位：ms
func WithRetryInterval(interval int) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.retryInterval = interval
	}
}

func WithChainClientAlias(alias string) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.alias = alias
	}
}

// WithChainClientLogger 设置Logger对象，便于日志打印
func WithChainClientLogger(logger utils.Logger) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.logger = logger
	}
}

// WithArchiveConfig 设置Archive配置
func WithArchiveConfig(conf *ArchiveConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.archiveConfig = conf
	}
}

// WithRPCClientConfig 设置grpc客户端配置
func WithRPCClientConfig(conf *RPCClientConfig) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.rpcClientConfig = conf
	}
}

// WithPkcs11Config 设置pkcs11配置
func WithPkcs11Config(conf *Pkcs11Config) ChainClientOption {
	return func(config *ChainClientConfig) {
		config.pkcs11Config = conf
	}
}

// 生成SDK配置并校验合法性
func generateConfig(opts ...ChainClientOption) (*ChainClientConfig, error) {
	config := &ChainClientConfig{}
	for _, opt := range opts {
		opt(config)
	}

	// 校验config参数合法性
	if err := checkConfig(config); err != nil {
		return nil, err
	}

	// 进一步处理config参数
	if err := dealConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func setAuthType(config *ChainClientConfig) {
	if config.authType == 0 {
		if utils.Config.ChainClientConfig.AuthType == "" {
			config.authType = PermissionedWithCert
		} else {
			config.authType = StringToAuthTypeMap[utils.Config.ChainClientConfig.AuthType]
		}
	}
}

func setCrypto(config *ChainClientConfig) {
	if config.authType == PermissionedWithCert {
		config.crypto = &CryptoConfig{}
		return
	}

	if utils.Config.ChainClientConfig.Crypto != nil && config.crypto == nil {
		config.crypto = &CryptoConfig{
			hash: utils.Config.ChainClientConfig.Crypto.Hash,
		}
	}
}

func setChainConfig(config *ChainClientConfig) {
	if utils.Config.ChainClientConfig.ChainId != "" && config.chainId == "" {
		config.chainId = utils.Config.ChainClientConfig.ChainId
	}

	if utils.Config.ChainClientConfig.OrgId != "" && config.orgId == "" {
		config.orgId = utils.Config.ChainClientConfig.OrgId
	}

	if utils.Config.ChainClientConfig.Alias != "" && config.alias == "" {
		config.alias = utils.Config.ChainClientConfig.Alias
	}

	config.enableNormalKey = utils.Config.ChainClientConfig.EnableNormalKey
}

// 如果参数没有设置，便使用配置文件的配置
func setUserConfig(config *ChainClientConfig) {
	if config.authType == PermissionedWithKey || config.authType == Public { // 公钥身份或公链模式
		if utils.Config.ChainClientConfig.UserSignKeyFilePath != "" && config.userSignKeyFilePath == "" &&
			config.userSignKeyBytes == nil {
			config.userSignKeyFilePath = utils.Config.ChainClientConfig.UserSignKeyFilePath
		}
		return
	}

	// 默认证书模式
	if utils.Config.ChainClientConfig.UserKeyFilePath != "" && config.userKeyFilePath == "" &&
		config.userKeyBytes == nil {
		config.userKeyFilePath = utils.Config.ChainClientConfig.UserKeyFilePath
	}

	if utils.Config.ChainClientConfig.UserCrtFilePath != "" && config.userCrtFilePath == "" &&
		config.userCrtBytes == nil {
		config.userCrtFilePath = utils.Config.ChainClientConfig.UserCrtFilePath
	}

	if utils.Config.ChainClientConfig.UserSignKeyFilePath != "" && config.userSignKeyFilePath == "" &&
		config.userSignKeyBytes == nil {
		config.userSignKeyFilePath = utils.Config.ChainClientConfig.UserSignKeyFilePath
	}

	if utils.Config.ChainClientConfig.UserSignCrtFilePath != "" && config.userSignCrtFilePath == "" &&
		config.userSignCrtBytes == nil {
		config.userSignCrtFilePath = utils.Config.ChainClientConfig.UserSignCrtFilePath
	}
}

func setNodeList(config *ChainClientConfig) {
	if len(utils.Config.ChainClientConfig.NodesConfig) > 0 && len(config.nodeList) == 0 {
		for _, conf := range utils.Config.ChainClientConfig.NodesConfig {
			// 只允许证书模式下启用TLS
			if config.authType == PermissionedWithKey || config.authType == Public {
				conf.EnableTLS = false
			}

			node := NewNodeConfig(
				// 节点地址，格式：127.0.0.1:12301
				WithNodeAddr(conf.NodeAddr),
				// 节点连接数
				WithNodeConnCnt(conf.ConnCnt),
				// 节点是否启用TLS认证
				WithNodeUseTLS(conf.EnableTLS),
				// 根证书路径，支持多个
				WithNodeCAPaths(conf.TrustRootPaths),
				// TLS Hostname
				WithNodeTLSHostName(conf.TLSHostName),
			)

			config.nodeList = append(config.nodeList, node)
		}
	}
}

func setArchiveConfig(config *ChainClientConfig) {
	if utils.Config.ChainClientConfig.ArchiveConfig != nil && config.archiveConfig == nil {
		archive := NewArchiveConfig(
			// secret key
			WithSecretKey(utils.Config.ChainClientConfig.ArchiveConfig.SecretKey),
		)

		config.archiveConfig = archive
	}
}

func setRPCClientConfig(config *ChainClientConfig) {
	if utils.Config.ChainClientConfig.RPCClientConfig != nil && config.rpcClientConfig == nil {
		rpcClient := NewRPCClientConfig(
			WithRPCClientMaxReceiveMessageSize(utils.Config.ChainClientConfig.RPCClientConfig.MaxRecvMsgSize),
			WithRPCClientMaxSendMessageSize(utils.Config.ChainClientConfig.RPCClientConfig.MaxSendMsgSize),
		)
		config.rpcClientConfig = rpcClient
	}
}

func setPkcs11Config(config *ChainClientConfig) {
	if config.authType == PermissionedWithCert {
		if utils.Config.ChainClientConfig.Pkcs11Config != nil && config.pkcs11Config == nil {
			config.pkcs11Config = NewPkcs11Config(
				utils.Config.ChainClientConfig.Pkcs11Config.Enabled,
				utils.Config.ChainClientConfig.Pkcs11Config.Library,
				utils.Config.ChainClientConfig.Pkcs11Config.Label,
				utils.Config.ChainClientConfig.Pkcs11Config.Password,
				utils.Config.ChainClientConfig.Pkcs11Config.SessionCacheSize,
				utils.Config.ChainClientConfig.Pkcs11Config.Hash,
			)
		}
	} else {
		config.pkcs11Config = &Pkcs11Config{
			Enabled: false,
		}
		//if utils.Config.ChainClientConfig.Pkcs11Config != nil && config.pkcs11Config == nil {
		//	config.pkcs11Config = &Pkcs11Config{
		//	}
		//}
	}
}

func setRetryConfig(config *ChainClientConfig) {
	config.retryLimit = utils.Config.ChainClientConfig.RetryLimit
	config.retryInterval = utils.Config.ChainClientConfig.RetryInterval
}

func readConfigFile(config *ChainClientConfig) error {
	// 若没有配置配置文件
	if config.confPath == "" {
		return nil
	}

	if err := utils.InitConfig(config.confPath); err != nil {
		return fmt.Errorf("init config failed, %s", err.Error())
	}

	setAuthType(config)

	setCrypto(config)

	setChainConfig(config)

	setUserConfig(config)

	setNodeList(config)

	setArchiveConfig(config)

	setRPCClientConfig(config)

	setPkcs11Config(config)

	setRetryConfig(config)

	return nil
}

// SDK配置校验
func checkConfig(config *ChainClientConfig) error {

	var (
		err error
	)

	if err = readConfigFile(config); err != nil {
		return fmt.Errorf("read sdk config file failed, %s", err.Error())
	}

	// 如果logger未指定，使用默认zap logger
	if config.logger == nil {
		config.logger = getDefaultLogger()
	}

	if err = checkNodeListConfig(config); err != nil {
		return err
	}

	if err = checkUserConfig(config); err != nil {
		return err
	}

	if err = checkChainConfig(config); err != nil {
		return err
	}

	if err = checkArchiveConfig(config); err != nil {
		return err
	}

	if err = checkPkcs11Config(config); err != nil {
		return err
	}

	return checkRPCClientConfig(config)
}

func checkNodeListConfig(config *ChainClientConfig) error {
	// 连接的节点地址不可为空
	if len(config.nodeList) == 0 {
		return fmt.Errorf("connect chainmaker node address is empty")
	}

	// 已配置的节点地址连接数，需要在合理区间
	for _, node := range config.nodeList {
		if node.connCnt <= 0 || node.connCnt > MaxConnCnt {
			return fmt.Errorf("node connection count should >0 && <=%d",
				MaxConnCnt)
		}

		if node.useTLS {
			// 如果开启了TLS认证，CA路径必填
			if len(node.caPaths) == 0 && len(node.caCerts) == 0 {
				return fmt.Errorf("if node useTLS is open, should set caPaths or caCerts")
			}

			// 如果开启了TLS认证，需配置TLS HostName
			if node.tlsHostName == "" {
				return fmt.Errorf("if node useTLS is open, should set tls hostname")
			}
		}
	}

	return nil
}

func checkUserConfig(config *ChainClientConfig) error {
	if config.authType == PermissionedWithCert {
		// 用户私钥不可为空
		if config.userKeyFilePath == "" && config.userKeyBytes == nil {
			return fmt.Errorf("user key cannot be empty")
		}

		// 用户证书不可为空
		if config.userCrtFilePath == "" && config.userCrtBytes == nil {
			return fmt.Errorf("user crt cannot be empty")
		}
	} else {
		if config.userSignKeyFilePath == "" && config.userSignKeyBytes == nil {
			return fmt.Errorf("user key cannot be empty")
		}
	}

	return nil
}

func checkChainConfig(config *ChainClientConfig) error {
	if config.authType == PermissionedWithCert || config.authType == PermissionedWithKey {
		// OrgId不可为空
		if config.orgId == "" {
			return fmt.Errorf("orgId cannot be empty")
		}
	}

	// ChainId不可为空
	if config.chainId == "" {
		return fmt.Errorf("chainId cannot be empty")
	}

	return nil
}

func checkArchiveConfig(config *ChainClientConfig) error {
	return nil
}

func checkPkcs11Config(config *ChainClientConfig) error {
	if config.pkcs11Config == nil || !config.pkcs11Config.Enabled {
		return nil
	}
	// 如果config.pkcs11Config.Enabled == true 则其他参数不能为空
	if config.pkcs11Config.Library == "" {
		return errors.New("config.pkcs11Config.Library must not empty")
	}
	if config.pkcs11Config.Label == "" {
		return errors.New("config.pkcs11Config.Label must not empty")
	}
	if config.pkcs11Config.Password == "" {
		return errors.New("config.pkcs11Config.Password must not empty")
	}
	if config.pkcs11Config.SessionCacheSize == 0 {
		return errors.New("config.pkcs11Config.SessionCacheSize must > 0")
	}
	if config.pkcs11Config.Hash == "" {
		return errors.New("config.pkcs11Config.Hash must not empty")
	}
	return nil
}

func checkRPCClientConfig(config *ChainClientConfig) error {
	if config.rpcClientConfig == nil {
		rpcClient := NewRPCClientConfig(
			WithRPCClientMaxReceiveMessageSize(DefaultRpcClientMaxReceiveMessageSize),
			WithRPCClientMaxSendMessageSize(DefaultRpcClientMaxSendMessageSize),
		)
		config.rpcClientConfig = rpcClient
	} else {
		if config.rpcClientConfig.rpcClientMaxReceiveMessageSize <= 0 {
			config.rpcClientConfig.rpcClientMaxReceiveMessageSize = DefaultRpcClientMaxReceiveMessageSize
		}
		if config.rpcClientConfig.rpcClientMaxSendMessageSize <= 0 {
			config.rpcClientConfig.rpcClientMaxSendMessageSize = DefaultRpcClientMaxSendMessageSize
		}
	}
	return nil
}

func dealConfig(config *ChainClientConfig) error {
	var err error
	if err = dealRetryConfig(config); err != nil {
		return err
	}

	// PermissionedWithKey & Public
	if config.authType == PermissionedWithKey || config.authType == Public {
		return dealUserSignKeyConfig(config)
	}

	// PermissionedWithCert
	if err = dealUserCrtConfig(config); err != nil {
		return err
	}

	if err = dealUserKeyConfig(config); err != nil {
		return err
	}

	if err = dealUserSignCrtConfig(config); err != nil {
		return err
	}

	return dealUserSignKeyConfig(config)
}

func dealUserCrtConfig(config *ChainClientConfig) (err error) {

	if config.userCrtBytes == nil {
		// 读取用户证书
		config.userCrtBytes, err = ioutil.ReadFile(config.userCrtFilePath)
		if err != nil {
			return fmt.Errorf("read user crt file failed, %s", err.Error())
		}
	}

	// 将证书转换为证书对象
	if config.userCrt, err = utils.ParseCert(config.userCrtBytes); err != nil {
		return fmt.Errorf("utils.ParseCert failed, %s", err.Error())
	}

	return nil
}

func dealUserKeyConfig(config *ChainClientConfig) (err error) {

	if config.userKeyBytes == nil {
		// 从私钥文件读取用户私钥，转换为privateKey对象
		config.userKeyBytes, err = ioutil.ReadFile(config.userKeyFilePath)
		if err != nil {
			return fmt.Errorf("read user key file failed, %s", err)
		}
	}

	config.privateKey, err = asym.PrivateKeyFromPEM(config.userKeyBytes, nil)
	if err != nil {
		return fmt.Errorf("parse user key file to privateKey obj failed, %s", err)
	}

	return nil
}

func dealUserSignCrtConfig(config *ChainClientConfig) (err error) {

	if config.userSignCrtBytes == nil {
		if config.userSignCrtFilePath == "" {
			config.userSignCrtBytes = config.userCrtBytes
			return nil
		}

		config.userSignCrtBytes, err = ioutil.ReadFile(config.userSignCrtFilePath)
		if err != nil {
			return fmt.Errorf("read user sign crt file failed, %s", err.Error())
		}

	}

	if config.userCrt, err = utils.ParseCert(config.userSignCrtBytes); err != nil {
		return fmt.Errorf("utils.ParseCert failed, %s", err.Error())
	}

	return nil
}

func dealUserSignKeyConfig(config *ChainClientConfig) (err error) {
	// PermissionedWithKey, Public
	if config.authType == PermissionedWithKey || config.authType == Public {
		if config.userSignKeyBytes == nil {
			config.userSignKeyBytes, err = ioutil.ReadFile(config.userSignKeyFilePath)
			if err != nil {
				return fmt.Errorf("read user private Key file failed, %s", err.Error())
			}

			config.privateKey, err = asym.PrivateKeyFromPEM(config.userSignKeyBytes, nil)
			if err != nil {
				return fmt.Errorf("PrivateKeyFromPEM failed, %s", err.Error())
			}

			config.userPk = config.privateKey.PublicKey()
		}
		return nil
	}

	// PermissionedWithCert
	if config.userSignKeyBytes == nil {
		if config.userSignKeyFilePath == "" {
			config.userSignKeyBytes = config.userKeyBytes
			return nil
		}

		config.userSignKeyBytes, err = ioutil.ReadFile(config.userSignKeyFilePath)
		if err != nil {
			return fmt.Errorf("read user sign key file failed, %s", err.Error())
		}
	}

	if config.pkcs11Config.Enabled {
		p11Handle, err = pkcs11.New(config.pkcs11Config.Library, config.pkcs11Config.Label,
			config.pkcs11Config.Password, config.pkcs11Config.SessionCacheSize, config.pkcs11Config.Hash)
		if err != nil {
			return fmt.Errorf("failed to initialize pkcs11 handle, %s", err)
		}
		config.privateKey, err = cert.ParseP11PrivKey(p11Handle, config.userSignKeyBytes)
		if err != nil {
			return fmt.Errorf("cert.ParseP11PrivKey failed, %s", err)
		}
	} else {
		config.privateKey, err = asym.PrivateKeyFromPEM(config.userSignKeyBytes, nil)
		if err != nil {
			return fmt.Errorf("parse user key file to privateKey obj failed, %s", err)
		}
	}

	return nil
}

func dealRetryConfig(config *ChainClientConfig) (err error) {

	if config.retryLimit <= 0 {
		config.retryLimit = DefaultRetryLimit
	}

	if config.retryInterval <= 0 {
		config.retryInterval = DefaultRetryInterval
	}

	return nil
}

func getDefaultLogger() *zap.SugaredLogger {
	config := log.LogConfig{
		Module:       "[SDK]",
		LogPath:      "./sdk.log",
		LogLevel:     log.LEVEL_DEBUG,
		MaxAge:       30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: false,
	}

	logger, _ := log.InitSugarLogger(&config)
	return logger
}

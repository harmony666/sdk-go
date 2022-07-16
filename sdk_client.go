/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errStringFormat    = "%s failed, %s"
	sdkErrStringFormat = "[SDK] %s"
)

var _ SDKInterface = (*ChainClient)(nil)

type ChainClient struct {
	// common config
	logger  utils.Logger
	pool    ConnectionPool
	chainId string
	orgId   string

	userCrtBytes []byte
	userCrt      *bcx509.Certificate
	privateKey   crypto.PrivateKey
	publicKey    crypto.PublicKey
	pkBytes      []byte

	// cert hash config
	enabledCrtHash bool
	userCrtHash    []byte

	// archive config
	archiveConfig *ArchiveConfig

	// grpc client config
	rpcClientConfig *RPCClientConfig

	// pkcs11 config
	pkcs11Config *Pkcs11Config

	hashType string
	authType AuthType
	// retry config
	retryLimit    int // if <=0 then use DefaultRetryLimit
	retryInterval int // if <=0 then use DefaultRetryInterval

	// alias support
	enabledAlias bool
	alias        string

	// default TimestampKey , true NormalKey support
	enableNormalKey bool
}

func NewNodeConfig(opts ...NodeOption) *NodeConfig {
	config := &NodeConfig{}
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func NewConnPoolWithOptions(opts ...ChainClientOption) (*ClientConnectionPool, error) {
	config, err := generateConfig(opts...)
	if err != nil {
		return nil, err
	}

	return NewConnPool(config)
}

func NewArchiveConfig(opts ...ArchiveOption) *ArchiveConfig {
	config := &ArchiveConfig{}
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func NewRPCClientConfig(opts ...RPCClientOption) *RPCClientConfig {
	config := &RPCClientConfig{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func NewPkcs11Config(enabled bool, libPath, label, password string,
	sessionCacheSize int, hashAlgo string) *Pkcs11Config {
	return &Pkcs11Config{
		Enabled:          enabled,
		Library:          libPath,
		Label:            label,
		Password:         password,
		SessionCacheSize: sessionCacheSize,
		Hash:             hashAlgo,
	}
}

func NewChainClient(opts ...ChainClientOption) (*ChainClient, error) {
	config, err := generateConfig(opts...)
	if err != nil {
		return nil, err
	}

	pool, err := NewConnPool(config)
	if err != nil {
		return nil, err
	}

	var hashType = ""
	var publicKey crypto.PublicKey
	var pkBytes []byte
	var pkPem string
	if config.authType == PermissionedWithKey || config.authType == Public {
		hashType = config.crypto.hash
		publicKey = config.userPk
		pkPem, err = publicKey.String()
		if err != nil {
			return nil, err
		}

		pkBytes = []byte(pkPem)
	}

	cc := &ChainClient{
		pool:            pool,
		logger:          config.logger,
		chainId:         config.chainId,
		orgId:           config.orgId,
		alias:           config.alias,
		userCrtBytes:    config.userSignCrtBytes,
		userCrt:         config.userCrt,
		privateKey:      config.privateKey,
		archiveConfig:   config.archiveConfig,
		rpcClientConfig: config.rpcClientConfig,
		pkcs11Config:    config.pkcs11Config,

		publicKey: publicKey,
		hashType:  hashType,
		authType:  config.authType,
		pkBytes:   pkBytes,

		retryLimit:    config.retryLimit,
		retryInterval: config.retryInterval,

		enableNormalKey: config.enableNormalKey,
	}

	// 若设置了别名，便启用
	if config.authType == PermissionedWithCert && len(cc.alias) > 0 {
		if err := cc.EnableAlias(); err != nil {
			return nil, err
		}
	}

	return cc, nil
}

func (cc *ChainClient) IsEnableNormalKey() bool {
	return cc.enableNormalKey
}

func (cc *ChainClient) Stop() error {
	return cc.pool.Close()
}

func (cc *ChainClient) proposalRequest(payload *common.Payload,
	endorsers []*common.EndorsementEntry) (*common.TxResponse, error) {
	return cc.proposalRequestWithTimeout(payload, endorsers, -1)
}

func (cc *ChainClient) proposalRequestWithTimeout(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64) (*common.TxResponse, error) {

	req, err := cc.GenerateTxRequest(payload, endorsers)
	if err != nil {
		return nil, err
	}

	return cc.sendTxRequest(req, timeout)
}

func (cc *ChainClient) GenerateTxRequest(payload *common.Payload,
	endorsers []*common.EndorsementEntry) (*common.TxRequest, error) {
	var (
		signer    *accesscontrol.Member
		signBytes []byte
		err       error
	)

	// 构造Sender
	if cc.authType == PermissionedWithCert {

		if cc.enabledAlias && len(cc.alias) > 0 {
			signer = &accesscontrol.Member{
				OrgId:      cc.orgId,
				MemberInfo: []byte(cc.alias),
				MemberType: accesscontrol.MemberType_ALIAS,
			}
		} else if cc.enabledCrtHash && len(cc.userCrtHash) > 0 {
			signer = &accesscontrol.Member{
				OrgId:      cc.orgId,
				MemberInfo: cc.userCrtHash,
				MemberType: accesscontrol.MemberType_CERT_HASH,
			}
		} else {
			signer = &accesscontrol.Member{
				OrgId:      cc.orgId,
				MemberInfo: cc.userCrtBytes,
				MemberType: accesscontrol.MemberType_CERT,
			}
		}
	} else {
		signer = &accesscontrol.Member{
			OrgId:      cc.orgId,
			MemberInfo: cc.pkBytes,
			MemberType: accesscontrol.MemberType_PUBLIC_KEY,
		}
	}

	req := &common.TxRequest{
		Payload: payload,
		Sender: &common.EndorsementEntry{
			Signer:    signer,
			Signature: nil,
		},
		Endorsers: endorsers,
	}

	if cc.authType == PermissionedWithCert {
		hashalgo, err := bcx509.GetHashFromSignatureAlgorithm(cc.userCrt.SignatureAlgorithm)
		if err != nil {
			return nil, fmt.Errorf("invalid algorithm: %v", err.Error())
		}

		signBytes, err = utils.SignPayloadWithHashType(cc.privateKey, hashalgo, payload)
		if err != nil {
			return nil, fmt.Errorf("SignPayload failed, %s", err.Error())
		}
	} else {
		signBytes, err = utils.SignPayloadWithHashType(cc.privateKey, crypto.HashAlgoMap[cc.hashType], payload)
		if err != nil {
			return nil, fmt.Errorf("SignPayload failed, %s", err.Error())
		}
	}

	req.Sender.Signature = signBytes

	return req, nil
}

func (cc *ChainClient) sendTxRequest(txRequest *common.TxRequest, timeout int64) (*common.TxResponse, error) {

	var (
		errMsg string
	)

	if timeout < 0 {
		timeout = DefaultSendTxTimeout
		if strings.HasPrefix(txRequest.Payload.TxType.String(), "QUERY") {
			timeout = DefaultGetTxTimeout
		}
	}

	ignoreAddrs := make(map[string]struct{})
	for {
		client, err := cc.pool.getClientWithIgnoreAddrs(ignoreAddrs)
		if err != nil {
			return nil, err
		}

		if len(ignoreAddrs) > 0 {
			cc.logger.Debugf("[SDK] begin try to connect node [%s]", client.ID)
		}

		resp, err := client.sendRequestWithTimeout(txRequest, timeout)
		if err != nil {
			resp := &common.TxResponse{
				Message: err.Error(),
				TxId:    txRequest.Payload.TxId,
			}

			statusErr, ok := status.FromError(err)
			if ok && (statusErr.Code() == codes.DeadlineExceeded ||
				// desc = "transport: Error while dialing dial tcp 127.0.0.1:12301: connect: connection refused"
				statusErr.Code() == codes.Unavailable) {

				resp.Code = common.TxStatusCode_TIMEOUT
				errMsg = fmt.Sprintf("call [%s] meet network error, try to connect another node if has, %s",
					client.ID, err.Error())

				cc.logger.Errorf(sdkErrStringFormat, errMsg)
				ignoreAddrs[client.ID] = struct{}{}
				continue
			}

			cc.logger.Errorf("statusErr.Code() : %s", statusErr.Code())

			resp.Code = common.TxStatusCode_INTERNAL_ERROR
			errMsg = fmt.Sprintf("client.call failed, %+v", err)
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return resp, fmt.Errorf(errMsg)
		}

		resp.TxId = txRequest.Payload.TxId
		cc.logger.Debugf("[SDK] proposalRequest resp: %+v", resp)
		return resp, nil
	}
}

// EnableCertHash Cert Hash logic
func (cc *ChainClient) EnableCertHash() error {
	var (
		err error
	)

	// 优先使用别名，如果开启了别名，直接忽略压缩证书
	if cc.enabledAlias {
		return nil
	}

	if cc.GetAuthType() != PermissionedWithCert {
		return errors.New("cert hash is not supported")
	}

	// 0.已经启用压缩证书
	if cc.enabledCrtHash {
		return nil
	}

	// 1.如尚未获取证书Hash，便进行获取
	if len(cc.userCrtHash) == 0 {
		// 获取证书Hash
		cc.userCrtHash, err = cc.GetCertHash()
		if err != nil {
			errMsg := fmt.Sprintf("get cert hash failed, %s", err.Error())
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}
	}

	// 2.链上查询证书是否存在
	ok, err := cc.getCheckCertHash()
	if err != nil {
		errMsg := fmt.Sprintf("enable cert hash, get and check cert hash failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	// 3.1 若证书已经上链，直接返回
	if ok {
		cc.enabledCrtHash = true
		return nil
	}

	// 3.2 若证书未上链，添加证书
	resp, err := cc.AddCert()
	if err != nil {
		errMsg := fmt.Sprintf("enable cert hash AddCert failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		errMsg := fmt.Sprintf("enable cert hash AddCert got invalid resp, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	// 循环检查证书是否成功上链
	err = cc.checkUserCertOnChain()
	if err != nil {
		errMsg := fmt.Sprintf("check user cert on chain failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	cc.enabledCrtHash = true

	return nil
}

func (cc *ChainClient) DisableCertHash() error {
	cc.enabledCrtHash = false
	return nil
}

func (cc *ChainClient) GetEnabledCrtHash() bool {
	return cc.enabledCrtHash
}

func (cc *ChainClient) GetUserCrtHash() []byte {
	return cc.userCrtHash
}

func (cc *ChainClient) GetHashType() string {
	return cc.hashType
}

func (cc *ChainClient) GetAuthType() AuthType {
	return cc.authType
}

func (cc *ChainClient) GetPublicKey() crypto.PublicKey {
	return cc.publicKey
}

func (cc *ChainClient) GetPrivateKey() crypto.PrivateKey {
	return cc.privateKey
}

func (cc *ChainClient) GetCertPEM() []byte {
	return cc.userCrtBytes
}

func (cc *ChainClient) GetLocalCertAlias() string {
	return cc.alias
}

// ChangeSigner change ChainClient siger. signerCrt passes nil in Public or PermissionedWithKey mode
func (cc *ChainClient) ChangeSigner(signerPrivKey crypto.PrivateKey, signerCrt *bcx509.Certificate) error {
	signerPubKey := signerPrivKey.PublicKey()
	pkPem, err := signerPubKey.String()
	if err != nil {
		return err
	}

	cc.pkBytes = []byte(pkPem)
	cc.publicKey = signerPubKey
	cc.privateKey = signerPrivKey

	if signerCrt != nil {
		crtPem := pem.EncodeToMemory(&pem.Block{Bytes: signerCrt.Raw, Type: "CERTIFICATE"})
		cc.userCrtBytes = crtPem
		cc.userCrt = signerCrt
	}
	return nil
}

// 检查证书是否成功上链
func (cc *ChainClient) checkUserCertOnChain() error {
	err := retry.Retry(func(uint) error {
		ok, err := cc.getCheckCertHash()
		if err != nil {
			errMsg := fmt.Sprintf("check user cert on chain, get and check cert hash failed, %s", err.Error())
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}

		if !ok {
			errMsg := "user cert havenot on chain yet, and try again"
			cc.logger.Debugf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}

		return nil
	}, strategy.Limit(10), strategy.Wait(time.Second))

	if err != nil {
		errMsg := fmt.Sprintf("check user upload cert on chain failed, try again later, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func (cc *ChainClient) getCheckCertHash() (bool, error) {
	// 根据已缓存证书Hash，查链上是否存在
	certInfo, err := cc.QueryCert([]string{hex.EncodeToString(cc.userCrtHash)})
	if err != nil {
		errMsg := fmt.Sprintf("QueryCert failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return false, errors.New(errMsg)
	}

	if len(certInfo.CertInfos) == 0 {
		return false, nil
	}

	// 返回链上证书列表长度不为1，即报错
	if len(certInfo.CertInfos) > 1 {
		errMsg := "CertInfos != 1"
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return false, errors.New(errMsg)
	}

	// 如果链上证书Hash不为空
	if len(certInfo.CertInfos[0].Cert) > 0 {
		// 如果和缓存的证书Hash不一致则报错
		if hex.EncodeToString(cc.userCrtHash) != certInfo.CertInfos[0].Hash {
			errMsg := fmt.Sprintf("not equal certHash, [expected:%s]/[actual:%s]",
				cc.userCrtHash, certInfo.CertInfos[0].Hash)
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return false, errors.New(errMsg)
		}

		// 如果和缓存的证书Hash一致，则说明已经上传好了证书，具备提交压缩证书交易的能力
		return true, nil
	}

	return false, nil
}

func (cc *ChainClient) Pkcs11Config() *Pkcs11Config {
	return cc.pkcs11Config
}

func CreateChainClient(pool ConnectionPool, userCrtBytes, privKey, userCrtHash []byte, orgId, chainId string,
	enabledCrtHash int) (*ChainClient, error) {
	cert, err := utils.ParseCert(userCrtBytes)
	if err != nil {
		return nil, err
	}

	priv, err := asym.PrivateKeyFromPEM(privKey, nil)
	if err != nil {
		return nil, err
	}

	chain := &ChainClient{
		pool:         pool,
		logger:       pool.getLogger(),
		chainId:      chainId,
		orgId:        orgId,
		userCrtBytes: userCrtBytes,
		userCrt:      cert,
		privateKey:   priv,
	}

	return chain, nil
}

func (cc *ChainClient) EnableAlias() error {
	var (
		err error
	)

	// 已经启用别名，直接返回
	if cc.enabledAlias {
		return nil
	}

	// 查询别名是否上链
	ok, err := cc.getCheckAlias()
	if err != nil {
		errMsg := fmt.Sprintf("enable alias, get and check alias failed, %s", err.Error())
		cc.logger.Debugf(sdkErrStringFormat, errMsg)
		//return errors.New(errMsg)
	}

	// 别名已上链
	if ok {
		cc.enabledAlias = true
		return nil
	}

	// 添加别名
	resp, err := cc.AddAlias()
	if err != nil {
		errMsg := fmt.Sprintf("enable alias AddAlias failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		errMsg := fmt.Sprintf("enable alias AddAlias got invalid resp, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	// 循环检查别名是否成功上链
	err = cc.checkAliasOnChain()
	if err != nil {
		errMsg := fmt.Sprintf("check alias on chain failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	cc.enabledAlias = true

	return nil
}

func (cc *ChainClient) getCheckAlias() (bool, error) {
	aliasInfos, err := cc.QueryCertsAlias([]string{cc.alias})
	if err != nil {
		errMsg := fmt.Sprintf("QueryCertsAlias failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return false, errors.New(errMsg)
	}

	if len(aliasInfos.AliasInfos) != 1 {
		return false, errors.New("alias not found")
	}

	if aliasInfos.AliasInfos[0].Alias != cc.alias {
		return false, errors.New("alias not equal")
	}

	if aliasInfos.AliasInfos[0].NowCert.Cert == nil {
		return false, errors.New("alias has been deleted")
	}

	return true, nil
}

func (cc *ChainClient) checkAliasOnChain() error {
	err := retry.Retry(func(uint) error {
		ok, err := cc.getCheckAlias()
		if err != nil {
			errMsg := fmt.Sprintf("check alias on chain, get and check alias failed, %s", err.Error())
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}

		if !ok {
			errMsg := "alias havenot on chain yet, and try again"
			cc.logger.Debugf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}

		return nil
	}, strategy.Limit(10), strategy.Wait(time.Second))

	if err != nil {
		errMsg := fmt.Sprintf("check upload alias on chain failed, try again later, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	return nil
}

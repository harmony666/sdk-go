/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

//nolint
const (
	// System Block Contract keys
	KeyBlockContractWithRWSet   = "withRWSet"
	KeyBlockContractBlockHash   = "blockHash"
	KeyBlockContractBlockHeight = "blockHeight"
	KeyBlockContractTxId        = "txId"

	// System Chain Config Contract keys
	KeyChainConfigContractRoot              = "root"
	KeyChainConfigContractOrgId             = "org_id"
	KeyChainConfigAddrType                  = "addr_type"
	KeyChainConfigContractNodeId            = "node_id"
	KeyChainConfigContractNewNodeId         = "new_node_id"
	KeyChainConfigContractNodeIds           = "node_ids"
	KeyChainConfigContractBlockHeight       = "block_height"
	KeyChainConfigContractTrustMemberOrgId  = "org_id"
	KeyChainConfigContractTrustMemberInfo   = "member_info"
	KeyChainConfigContractTrustMemberNodeId = "node_id"
	KeyChainConfigContractTrustMemberRole   = "role"

	// CoreConfig keys
	KeyTxSchedulerTimeout         = "tx_scheduler_timeout"
	KeyTxSchedulerValidateTimeout = "tx_scheduler_validate_timeout"

	// BlockConfig keys
	KeyTxTimeOut       = "tx_timeout"
	KeyBlockTxCapacity = "block_tx_capacity"
	KeyBlockSize       = "block_size"
	KeyBlockInterval   = "block_interval"
	KeyTxParamterSize  = "tx_parameter_size"

	// CertManage keys
	KeyCertHashes = "cert_hashes"
	KeyCerts      = "certs"
	KeyCertCrl    = "cert_crl"

	// PrivateCompute keys
	KeyOrderId      = "order_id"
	KeyPrivateDir   = "private_dir"
	KeyContractName = "contract_name"
	KeyCodeHash     = "code_hash"
	KeyResult       = "result"
	KeyCodeHeader   = "code_header"
	KeyVersion      = "version"
	KeyIsDeploy     = "is_deploy"
	KeyRWSet        = "rw_set"
	KeyEvents       = "events"
	KeyReportHash   = "report_hash"
	KeySign         = "sign"
	KeyKey          = "key"
	KeyPayload      = "payload"
	KeyOrgIds       = "org_ids"
	KeySignPairs    = "sign_pairs"
	KeyCaCert       = "ca_cert"
	KeyEnclaveId    = "enclave_id"
	KeyReport       = "report"
	KeyProof        = "proof"
	KeyDeployReq    = "deploy_req"
	KeyPrivateReq   = "private_req"

	// PubkeyManage keys
	KeyPubkey      = "pubkey"
	KeyPubkeyRole  = "role"
	KeyPubkeyOrgId = "org_id"

	// Gas management
	KeyGasAddressKey       = "address_key"
	KeyGasPublicKey        = "public_key"
	KeyGasBatchRecharge    = "batch_recharge"
	KeyGasBalancePublicKey = "balance_public_key"
	KeyGasChargePublicKey  = "charge_public_key"
	KeyGasChargeGasAmount  = "charge_gas_amount"
	KeyGasFrozenPublicKey  = "frozen_public_key"
)

//nolint
const (
	// ArchiveConfig consts
	MysqlDBNamePrefix     = "cm_archived_chain"
	MysqlTableNamePrefix  = "t_block_info"
	RowsPerBlockInfoTable = 100000
)

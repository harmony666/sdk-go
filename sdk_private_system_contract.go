/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/proto"

	"chainmaker.org/chainmaker/common/v2/json"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

const ContractResultCode_OK uint32 = 0 //todo pb create const

func (cc *ChainClient) SaveDir(orderId, txId string,
	privateDir *common.StrSlice, withSyncResult bool, timeout int64) (*common.TxResponse, error) {

	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	cc.logger.Infof("[SDK] begin to save dir , [contract:%s]/[method:%s]/[txId:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_DIR.String(),
		txId,
	)

	// 构造Payload
	priDirBytes, err := privateDir.Marshal()
	if err != nil {
		return nil, fmt.Errorf("serielized private dir failed, %s", err.Error())
	}

	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyOrderId:    []byte(orderId),
		utils.KeyPrivateDir: priDirBytes,
	})

	payload := cc.CreatePayload(txId, common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_DIR.String(), pairs, defaultSeq, nil)

	resp, err := cc.sendContractRequest(payload, nil, timeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (cc *ChainClient) GetContract(contractName, codeHash string) (*common.PrivateGetContract, error) {

	cc.logger.Infof("[SDK] begin to get contract , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_CONTRACT.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyContractName: []byte(contractName),
		utils.KeyCodeHash:     []byte(codeHash),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_CONTRACT.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	contractInfo := &common.PrivateGetContract{}
	if err = proto.Unmarshal(resp.ContractResult.Result, contractInfo); err != nil {
		return nil, fmt.Errorf("GetContract unmarshal contract info payload failed, %s", err.Error())
	}

	return contractInfo, nil
}

func (cc *ChainClient) SaveData(contractName string, contractVersion string, isDeployment bool, codeHash []byte,
	reportHash []byte, result *common.ContractResult, codeHeader []byte, txId string, rwSet *common.TxRWSet,
	sign []byte, events *common.StrSlice, privateReq []byte, withSyncResult bool,
	timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	cc.logger.Infof("[SDK] begin to save data , [contract:%s]/[method:%s]/[txId:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_DATA.String(),
		txId,
	)

	// 构造Payload
	var rwSetStr string
	if rwSet != nil {
		rwb, err := rwSet.Marshal()
		if err != nil {
			return nil, fmt.Errorf("construct save data payload failed, %s", err.Error())
		}
		rwSetStr = string(rwb)
	}

	var eventsStr string
	if events != nil {
		eb, err := events.Marshal()
		if err != nil {
			return nil, fmt.Errorf("construct save data payload failed, %s", err.Error())
		}
		eventsStr = string(eb)
	}

	var resultStr string
	if result != nil {
		result, err := result.Marshal()
		if err != nil {
			return nil, fmt.Errorf("construct save data payload failed, %s", err.Error())
		}
		resultStr = string(result)
	}

	deployStr := strconv.FormatBool(isDeployment)
	pairsMap := map[string][]byte{
		utils.KeyResult:       []byte(resultStr),
		utils.KeyCodeHeader:   codeHeader,
		utils.KeyContractName: []byte(contractName),
		utils.KeyVersion:      []byte(contractVersion),
		utils.KeyIsDeploy:     []byte(deployStr),
		utils.KeyCodeHash:     codeHash,
		utils.KeyRWSet:        []byte(rwSetStr),
		utils.KeyEvents:       []byte(eventsStr),
		utils.KeyReportHash:   reportHash,
		utils.KeySign:         sign,
	}

	if isDeployment {
		pairsMap[utils.KeyDeployReq] = privateReq
	} else {
		pairsMap[utils.KeyPrivateReq] = privateReq
	}

	pairs := paramsMap2KVPairs(pairsMap)

	payload := cc.CreatePayload(txId, common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_DATA.String(), pairs, defaultSeq, nil)

	resp, err := cc.sendContractRequest(payload, nil, timeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (cc *ChainClient) GetData(contractName, key string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_DATA.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyContractName: []byte(contractName),
		utils.KeyKey:          []byte(key),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_DATA.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetDir(orderId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_DATA.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyOrderId: []byte(orderId),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_DIR.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) CheckCallerCertAuth(payload string, orgIds []string, signPairs []*syscontract.SignInfo) (
	*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to check caller cert auth  , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_CHECK_CALLER_CERT_AUTH.String(),
	)

	orgIdsJson, err := json.Marshal(orgIds)
	if err != nil {
		return nil, fmt.Errorf("json marshal orgIds failed, err: %v", err)
	}
	signPairsJson, err := json.Marshal(signPairs)
	if err != nil {
		return nil, fmt.Errorf("json marshal signPairs failed, err: %v", err)
	}
	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyPayload:   []byte(payload),
		utils.KeyOrgIds:    orgIdsJson,
		utils.KeySignPairs: signPairsJson,
	})

	payloadBytes := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_CHECK_CALLER_CERT_AUTH.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payloadBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payloadBytes.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) SaveEnclaveCACert(
	enclaveCACert, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	cc.logger.Infof("[SDK] begin to save ca cert , [contract:%s]/[method:%s]/[txId:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_CA_CERT.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyCaCert: []byte(enclaveCACert),
	})

	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_CA_CERT.String(), pairs, defaultSeq, nil)

	resp, err := cc.sendContractRequest(payload, nil, timeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) SaveEnclaveReport(
	enclaveId, report, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	cc.logger.Infof("[SDK] begin to save enclave report , [contract:%s]/[method:%s]/[txId:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_ENCLAVE_REPORT.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
		utils.KeyReport:    []byte(report),
	})

	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_ENCLAVE_REPORT.String(), pairs, defaultSeq, nil)

	resp, err := cc.sendContractRequest(payload, nil, timeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) CreateSaveEnclaveCACertPayload(enclaveCACert string, txId string) (*common.Payload, error) {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	cc.logger.Infof("[SDK] begin to save ca cert , [contract:%s]/[method:%s]/[txId:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_CA_CERT.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyCaCert: []byte(enclaveCACert),
	})

	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_CA_CERT.String(), pairs, defaultSeq, nil)

	return payload, nil
}

func (cc *ChainClient) GetEnclaveCACert() ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get ca cert , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_CA_CERT.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_CA_CERT.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) CreateSaveEnclaveReportPayload(enclaveId, report, txId string) (*common.Payload, error) {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	cc.logger.Infof("[SDK] begin to save enclave report , [contract:%s]/[method:%s]/[txId:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_ENCLAVE_REPORT.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
		utils.KeyReport:    []byte(report),
	})

	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_ENCLAVE_REPORT.String(), pairs, defaultSeq, nil)

	return payload, nil
}

func (cc *ChainClient) SaveRemoteAttestationProof(proof, txId string, withSyncResult bool,
	timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	cc.logger.Infof("[SDK] begin to save_remote_attestation_proof , [contract:%s]/[method:%s]/[txId:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_REMOTE_ATTESTATION.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyProof: []byte(proof),
	})

	payload := cc.CreatePayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_SAVE_REMOTE_ATTESTATION.String(), pairs, defaultSeq, nil)

	resp, err := cc.sendContractRequest(payload, nil, timeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) GetEnclaveEncryptPubKey(enclaveId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin get_enclave_encrypt_pub_key() , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_ENCRYPT_PUB_KEY.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_ENCRYPT_PUB_KEY.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetEnclaveVerificationPubKey(enclaveId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_VERIFICATION_PUB_KEY.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_VERIFICATION_PUB_KEY.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetEnclaveReport(enclaveId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_REPORT.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_REPORT.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetEnclaveChallenge(enclaveId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_CHALLENGE.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_CHALLENGE.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetEnclaveSignature(enclaveId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_SIGNATURE.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_SIGNATURE.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetEnclaveProof(enclaveId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_VERIFICATION_PUB_KEY.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string][]byte{
		utils.KeyEnclaveId: []byte(enclaveId),
	})

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_PRIVATE_COMPUTE.String(),
		syscontract.PrivateComputeFunction_GET_ENCLAVE_PROOF.String(), pairs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func paramsMap2KVPairs(params map[string][]byte) (kvPairs []*common.KeyValuePair) {
	for key, val := range params {
		kvPair := &common.KeyValuePair{
			Key:   key,
			Value: val,
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return
}

func checkProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != ContractResultCode_OK {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

func (cc *ChainClient) SendMultiSigningRequest(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {

	return cc.sendContractRequest(payload, endorsers, timeout, withSyncResult)
}

/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/config"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	testKey     = "key001"
	nodePeerId1 = "QmQVkTSF6aWzRSddT3rro6Ve33jhKpsHFaQoVxHKMWzhuN"
	nodePeerId2 = "QmQVkTSF6aWzRSddT3rro6Ve33jhKpsHFaQoVxHKMWzhuN"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkConfigPKUser1Path     = "../sdk_configs/sdk_config_pk_user1.yml"
)

func main() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	//client, err := examples.CreateChainClientWithSDKConf(sdkConfigPKUser1Path)
	if err != nil {
		log.Fatalln(err)
	}

	//testChainConfig(client)
	testChainConfigGasEnable(client)
	testChainConfigAlterAddrType(client)
}

func testChainConfig(client *sdk.ChainClient) {
	var (
		chainConfig *config.ChainConfig
		ok          bool
	)

	fmt.Println("====================== 根据区块高度获取链配置 ======================")
	testGetChainConfigByBlockHeight(client, 1)

	fmt.Println("====================== 获取链Sequence ======================")
	testGetChainConfigSeq(client)

	fmt.Println("====================== 更新CoreConfig ======================")
	rand.Seed(time.Now().UnixNano())
	txSchedulerTimeout := uint64(rand.Intn(61))
	txSchedulerValidateTimeout := uint64(rand.Intn(61))
	testChainConfigCoreUpdate(client, txSchedulerTimeout, txSchedulerValidateTimeout, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(5 * time.Second)
	chainConfig = testGetChainConfig(client)
	fmt.Printf("txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n", txSchedulerTimeout, txSchedulerValidateTimeout)
	fmt.Printf("chainConfig txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n",
		chainConfig.Core.TxSchedulerTimeout, chainConfig.Core.TxSchedulerValidateTimeout)
	if txSchedulerTimeout != chainConfig.Core.TxSchedulerTimeout {
		log.Fatalln("require txSchedulerTimeout == int(chainConfig.Core.TxSchedulerTimeout)")
	}
	if txSchedulerValidateTimeout != chainConfig.Core.TxSchedulerValidateTimeout {
		log.Fatalln("require txSchedulerValidateTimeout == int(chainConfig.Core.TxSchedulerValidateTimeout)")
	}

	fmt.Println("====================== 更新BlockConfig ======================")
	txTimestampVerify := rand.Intn(2) == 0
	txTimeout := uint32(rand.Intn(1000)) + 600
	blockTxCapacity := uint32(rand.Intn(1000)) + 1
	blockSize := uint32(rand.Intn(10)) + 1
	blockInterval := uint32(rand.Intn(10000)) + 10
	txParameterSize := uint32(rand.Intn(100))
	testChainConfigBlockUpdate(client, txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval, txParameterSize, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	fmt.Printf("tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n", strconv.FormatBool(txTimestampVerify), txTimeout, blockTxCapacity, blockSize, blockInterval)
	fmt.Printf("chainConfig txSchedulerTimeout: tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n",
		strconv.FormatBool(chainConfig.Block.TxTimestampVerify), chainConfig.Block.TxTimeout, chainConfig.Block.BlockTxCapacity, chainConfig.Block.BlockSize, chainConfig.Block.BlockInterval)
	if chainConfig.Block.TxTimestampVerify != txTimestampVerify {
		log.Fatalln("require chainConfig.Block.TxTimestampVerify == txTimestampVerify")
	}
	if txTimeout != chainConfig.Block.TxTimeout {
		log.Fatalln("require txTimeout == int(chainConfig.Block.TxTimeout)")
	}
	if blockTxCapacity != chainConfig.Block.BlockTxCapacity {
		log.Fatalln("require equal")
	}
	if blockSize != chainConfig.Block.BlockSize {
		log.Fatalln("require equal")
	}
	if blockInterval != chainConfig.Block.BlockInterval {
		log.Fatalln("require equal")
	}

	fmt.Println("====================== 新增trust root ca ======================")
	trustCount := len(testGetChainConfig(client).TrustRoots)
	raw, err := ioutil.ReadFile("../../testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}
	trustRootOrgId := examples.OrgId5
	trustRootCrt := string(raw)
	testChainConfigTrustRootAdd(client, trustRootOrgId, []string{trustRootCrt}, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount+1 != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[trustCount].OrgId {
		log.Fatalln("require equal")
	}

	//for _, root := range chainConfig.TrustRoots[trustCount].Root {
	//	if trustRootCrt == root {
	//		ok = true
	//		break
	//	}
	//}
	if !ok {
		log.Fatalln("require equal")
	}

	fmt.Println("====================== 更新trust root ca ======================")
	raw, err = ioutil.ReadFile("../../testdata/crypto-config/wx-org6.chainmaker.org/ca/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootUpdate(client, trustRootOrgId, []string{trustRootCrt}, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg5Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount+1 != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[trustCount].OrgId {
		log.Fatalln("require equal")
	}
	//for _, root := range chainConfig.TrustRoots[trustCount].Root {
	//	if trustRootCrt == root {
	//		ok = true
	//		break
	//	}
	//}
	if !ok {
		log.Fatalln("require equal")
	}

	fmt.Println("====================== 删除trust root ca ======================")
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootDelete(client, trustRootOrgId, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg5Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
	}

	// 6) [PermissionAdd]
	permissionCount := len(testGetChainConfig(client).ResourcePolicies)
	permissionResourceName := "TEST_PREMISSION"
	policy := &accesscontrol.Policy{
		Rule: "ANY",
	}
	testChainConfigPermissionAdd(client, permissionResourceName, policy, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if permissionCount+1 != len(chainConfig.ResourcePolicies) {
		log.Fatalln("require equal")
	}
	if !proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy) {
		log.Fatalln("require true")
	}

	// 7) [PermissionUpdate]
	permissionResourceName = "TEST_PREMISSION"
	policy = &accesscontrol.Policy{
		Rule: "ANY",
	}
	testChainConfigPermissionUpdate(client, permissionResourceName, policy, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if permissionCount+1 != len(chainConfig.ResourcePolicies) {
		log.Fatalln("require equal")
	}
	if !proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy) {
		log.Fatalln("require true")
	}

	// 8) [PermissionDelete]
	testChainConfigPermissionDelete(client, permissionResourceName, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if permissionCount != len(chainConfig.ResourcePolicies) {
		log.Fatalln("require equal")
	}

	// 9) [ConsensusNodeAddrAdd]
	nodeOrgId := examples.OrgId4
	nodeIds := []string{nodePeerId1}
	testChainConfigConsensusNodeIdAdd(client, nodeOrgId, nodeIds, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		log.Fatalln("require equal")
	}
	if 2 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[3].NodeId[1] {
		log.Fatalln("require equal")
	}

	// 10) [ConsensusNodeAddrUpdate]
	nodeOrgId = examples.OrgId4
	nodeOldId := nodePeerId1
	nodeNewId := nodePeerId2
	testChainConfigConsensusNodeIdUpdate(client, nodeOrgId, nodeOldId, nodeNewId, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		log.Fatalln("require equal")
	}
	if 2 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeNewId != chainConfig.Consensus.Nodes[3].NodeId[1] {
		log.Fatalln("require equal")
	}

	// 11) [ConsensusNodeAddrDelete]
	nodeOrgId = examples.OrgId4
	testChainConfigConsensusNodeIdDelete(client, nodeOrgId, nodeNewId, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		log.Fatalln("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		log.Fatalln("require equal")
	}

	// 12) [ConsensusNodeOrgAdd]
	raw, err = ioutil.ReadFile("../../testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootAdd(client, trustRootOrgId, []string{trustRootCrt}, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[4].OrgId {
		log.Fatalln("require equal")
	}
	//for _, root := range chainConfig.TrustRoots[trustCount].Root {
	//	if trustRootCrt == root {
	//		ok = true
	//		break
	//	}
	//}
	if !ok {
		log.Fatalln("require equal")
	}
	nodeOrgId = examples.OrgId5
	nodeIds = []string{nodePeerId1}
	testChainConfigConsensusNodeOrgAdd(client, nodeOrgId, nodeIds, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.Consensus.Nodes) {
		log.Fatalln("require equal")
	}
	if nodeOrgId != chainConfig.Consensus.Nodes[4].OrgId {
		log.Fatalln("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[4].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[4].NodeId[0] {
		log.Fatalln("require equal")
	}

	// 13) [ConsensusNodeOrgUpdate]
	nodeOrgId = examples.OrgId5
	nodeIds = []string{nodePeerId2}
	testChainConfigConsensusNodeOrgUpdate(client, nodeOrgId, nodeIds, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.Consensus.Nodes) {
		log.Fatalln("require equal")
	}
	if nodeOrgId != chainConfig.Consensus.Nodes[4].OrgId {
		log.Fatalln("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[4].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[4].NodeId[0] {
		log.Fatalln("require equal")
	}

	// 14) [ConsensusNodeOrgDelete]
	nodeOrgId = examples.OrgId5
	testChainConfigConsensusNodeOrgDelete(client, nodeOrgId, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 4 != len(chainConfig.Consensus.Nodes) {
		log.Fatalln("require equal")
	}

	// 15) [ConsensusExtAdd]
	kvs := []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: []byte("test_value"),
		},
	}
	testChainConfigConsensusExtAdd(client, kvs, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 2 != len(chainConfig.Consensus.ExtConfig) {
		log.Fatalln("require equal")
	}
	if !proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]) {
		log.Fatalln("require equal")
	}

	// 16) [ConsensusExtUpdate]
	kvs = []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: []byte("updated_value"),
		},
	}
	testChainConfigConsensusExtUpdate(client, kvs, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 2 != len(chainConfig.Consensus.ExtConfig) {
		log.Fatalln("require equal")
	}
	if !proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]) {
		log.Fatalln("require equal")
	}

	// 16) [ConsensusExtDelete]
	keys := []string{testKey}
	testChainConfigConsensusExtDelete(client, keys, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 1 != len(chainConfig.Consensus.ExtConfig) {
		log.Fatalln("require equal")
	}
}

func testGetChainConfig(client *sdk.ChainClient) *config.ChainConfig {
	resp, err := client.GetChainConfig()
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("GetChainConfig resp: %+v\n", resp)
	return resp
}

func testGetChainConfigByBlockHeight(client *sdk.ChainClient, blockHeight uint64) {
	resp, err := client.GetChainConfigByBlockHeight(blockHeight)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
}

func testGetChainConfigSeq(client *sdk.ChainClient) {
	seq, err := client.GetChainConfigSequence()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("chainconfig seq: %d\n", seq)
}

func testChainConfigCoreUpdate(client *sdk.ChainClient, txSchedulerTimeout, txSchedulerValidateTimeout uint64, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigCoreUpdatePayload(
		txSchedulerTimeout, txSchedulerValidateTimeout)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigBlockUpdate(client *sdk.ChainClient, txTimestampVerify bool,
	txTimeout, blockTxCapacity, blockSize, blockInterval, txParameterSize uint32, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigBlockUpdatePayload(
		txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval, txParameterSize)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigTrustRootAdd(client *sdk.ChainClient, trustRootOrgId string, trustRootCrt []string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigTrustRootUpdate(client *sdk.ChainClient, trustRootOrgId string, trustRootCrt []string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigTrustRootDelete(client *sdk.ChainClient, trustRootOrgId string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigTrustRootDeletePayload(trustRootOrgId)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigPermissionAdd(client *sdk.ChainClient, permissionResourceName string, policy *accesscontrol.Policy, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigPermissionAddPayload(permissionResourceName, policy)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigPermissionUpdate(client *sdk.ChainClient, permissionResourceName string, policy *accesscontrol.Policy, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigPermissionUpdatePayload(permissionResourceName, policy)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigPermissionDelete(client *sdk.ChainClient, permissionResourceName string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigPermissionDeletePayload(permissionResourceName)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusNodeIdAdd(client *sdk.ChainClient, nodeAddrOrgId string, nodeIds []string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusNodeIdAddPayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusNodeIdUpdate(client *sdk.ChainClient, nodeAddrOrgId, nodeOldIds, nodeNewIds string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusNodeIdUpdatePayload(nodeAddrOrgId, nodeOldIds, nodeNewIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusNodeIdDelete(client *sdk.ChainClient, nodeAddrOrgId, nodeId string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusNodeIdDeletePayload(nodeAddrOrgId, nodeId)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusNodeOrgAdd(client *sdk.ChainClient, nodeAddrOrgId string, nodeIds []string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusNodeOrgAddPayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusNodeOrgUpdate(client *sdk.ChainClient, nodeAddrOrgId string, nodeIds []string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusNodeOrgUpdatePayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusNodeOrgDelete(client *sdk.ChainClient, nodeAddrOrgId string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusNodeOrgDeletePayload(nodeAddrOrgId)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusExtAdd(client *sdk.ChainClient, kvs []*common.KeyValuePair, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusExtAddPayload(kvs)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusExtUpdate(client *sdk.ChainClient, kvs []*common.KeyValuePair, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusExtUpdatePayload(kvs)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigConsensusExtDelete(client *sdk.ChainClient, keys []string, usernames ...string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigConsensusExtDeletePayload(keys)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, payload, usernames...)
}

func testChainConfigGasEnable(client *sdk.ChainClient) {
	fmt.Println("====================== 启用/停用Gas计费 ======================")
	chainConfig := testGetChainConfig(client)
	if chainConfig.AccountConfig == nil {
		fmt.Println("unsupport gas")
		return
	}

	isGasEnabled := chainConfig.AccountConfig.EnableGas
	fmt.Printf("is gas enable: %+v\n", isGasEnabled)
	payload, err := client.CreateChainConfigEnableOrDisableGasPayload()
	if err != nil {
		log.Fatalln(err)
	}
	signAndSendRequest(client, payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)
	chainConfig = testGetChainConfig(client)
	isGasEnabled2 := chainConfig.AccountConfig.EnableGas
	fmt.Printf("is gas enable: %+v\n", isGasEnabled2)

	if isGasEnabled != !isGasEnabled2 {
		log.Fatalf("gas enable failed")
	}
}

func testChainConfigAlterAddrType(client *sdk.ChainClient) {
	fmt.Println("====================== 修改地址类型 ======================")
	chainConfig := testGetChainConfig(client)
	addrType := chainConfig.Vm.AddrType
	fmt.Printf("current address type is: %s\n", addrType.String())

	newAddrType := strconv.FormatInt(int64((addrType+1)%2), 10)

	payload, err := client.CreateChainConfigAlterAddrTypePayload(newAddrType)
	if err != nil {
		log.Fatalln(err)
	}
	signAndSendRequest(client, payload, examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1)

	chainConfig = testGetChainConfig(client)
	addrType = chainConfig.Vm.AddrType
	fmt.Printf("after alter address type is: %s\n", addrType.String())
}

func signAndSendRequest(client *sdk.ChainClient, payload *common.Payload, usernames ...string) {
	// 各组织Admin权限用户签名
	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(payload, endorsers, -1, true)
	if err != nil {
		log.Fatalln(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Printf("ChainConfigUpdate resp: %+v\n", resp)
}

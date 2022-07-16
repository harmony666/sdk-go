/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/config"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestGetCertHash(t *testing.T) {
	goodChainConfigBz, err := proto.Marshal(&config.ChainConfig{Crypto: &config.CryptoConfig{Hash: "SHA256"}})
	require.Nil(t, err)

	tests := []struct {
		name             string
		cliUserCertBytes []byte
		serverTxResp     *common.TxResponse
		serverTxErr      error

		expHash string
		wantErr bool
	}{
		{
			"good",
			[]byte("-----BEGIN CERTIFICATE-----\nMIICijCCAi+gAwIBAgIDBS9vMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTIwMTIwODA2NTM0M1oXDTI1\nMTIwNzA2NTM0M1owgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjbGllbnQxLDAqBgNVBAMTI2NsaWVudDEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE56xayRx0\n/a8KEXPxRfiSzYgJ/sE4tVeI/ZbjpiUX9m0TCJX7W/VHdm6WeJLOdCDuLLNvjGTy\nt8LLyqyubJI5AKN7MHkwDgYDVR0PAQH/BAQDAgGmMA8GA1UdJQQIMAYGBFUdJQAw\nKQYDVR0OBCIEIMjAiM2eMzlQ9HzV9ePW69rfUiRZVT2pDBOMqM4WVJSAMCsGA1Ud\nIwQkMCKAIDUkP3EcubfENS6TH3DFczH5dAnC2eD73+wcUF/bEIlnMAoGCCqGSM49\nBAMCA0kAMEYCIQCWUHL0xisjQoW+o6VV12pBXIRJgdeUeAu2EIjptSg2GAIhAIxK\nLXpHIBFxIkmWlxUaanCojPSZhzEbd+8LRrmhEO8n\n-----END CERTIFICATE-----"),
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: goodChainConfigBz,
				},
			},
			nil,
			"e77c9238c51e3446d942f94bd8803cc4f351254f8771f972146d7bfc6e0be7f4",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			cli.userCrtBytes = tt.cliUserCertBytes

			hash, err := cli.GetCertHash()
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.expHash, fmt.Sprintf("%x", hash))
		})
	}
}

func TestQueryCert(t *testing.T) {
	goodCertInfos := &common.CertInfos{
		CertInfos: []*common.CertInfo{
			{
				Hash: "e77c9238c51e3446d942f94bd8803cc4f351254f8771f972146d7bfc6e0be7f4",
				Cert: []byte("-----BEGIN CERTIFICATE-----\nMIICijCCAi+gAwIBAgIDBS9vMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTIwMTIwODA2NTM0M1oXDTI1\nMTIwNzA2NTM0M1owgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjbGllbnQxLDAqBgNVBAMTI2NsaWVudDEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE56xayRx0\n/a8KEXPxRfiSzYgJ/sE4tVeI/ZbjpiUX9m0TCJX7W/VHdm6WeJLOdCDuLLNvjGTy\nt8LLyqyubJI5AKN7MHkwDgYDVR0PAQH/BAQDAgGmMA8GA1UdJQQIMAYGBFUdJQAw\nKQYDVR0OBCIEIMjAiM2eMzlQ9HzV9ePW69rfUiRZVT2pDBOMqM4WVJSAMCsGA1Ud\nIwQkMCKAIDUkP3EcubfENS6TH3DFczH5dAnC2eD73+wcUF/bEIlnMAoGCCqGSM49\nBAMCA0kAMEYCIQCWUHL0xisjQoW+o6VV12pBXIRJgdeUeAu2EIjptSg2GAIhAIxK\nLXpHIBFxIkmWlxUaanCojPSZhzEbd+8LRrmhEO8n\n-----END CERTIFICATE-----"),
			},
		},
	}
	goodCertInfosBz, err := proto.Marshal(goodCertInfos)
	require.Nil(t, err)

	tests := []struct {
		name         string
		certHashes   []string
		serverTxResp *common.TxResponse
		serverTxErr  error

		wantErr bool
	}{
		{
			"good",
			[]string{"e77c9238c51e3446d942f94bd8803cc4f351254f8771f972146d7bfc6e0be7f4"},
			&common.TxResponse{
				Code: common.TxStatusCode_SUCCESS,
				ContractResult: &common.ContractResult{
					Result: goodCertInfosBz,
				},
			},
			nil,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverTxErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			certInfos, err := cli.QueryCert(tt.certHashes)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, goodCertInfos, certInfos)
		})
	}
}

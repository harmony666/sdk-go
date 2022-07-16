/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func TestSignPayload(t *testing.T) {
	tests := []struct {
		name             string
		userSignKeyBytes []byte
		userSignCrtBytes []byte
		unsignedPayload  *common.Payload
		wantErr          bool
	}{
		{
			"good",
			[]byte("-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIByphjR4auvodMAWeaWsDXuADlGVi0ODAZtOh7tcIr2hoAoGCCqGSM49\nAwEHoUQDQgAE56xayRx0/a8KEXPxRfiSzYgJ/sE4tVeI/ZbjpiUX9m0TCJX7W/VH\ndm6WeJLOdCDuLLNvjGTyt8LLyqyubJI5AA==\n-----END EC PRIVATE KEY-----"),
			[]byte("-----BEGIN CERTIFICATE-----\nMIICijCCAi+gAwIBAgIDBS9vMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTIwMTIwODA2NTM0M1oXDTI1\nMTIwNzA2NTM0M1owgZExCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjbGllbnQxLDAqBgNVBAMTI2NsaWVudDEuc2lnbi53eC1vcmcx\nLmNoYWlubWFrZXIub3JnMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE56xayRx0\n/a8KEXPxRfiSzYgJ/sE4tVeI/ZbjpiUX9m0TCJX7W/VHdm6WeJLOdCDuLLNvjGTy\nt8LLyqyubJI5AKN7MHkwDgYDVR0PAQH/BAQDAgGmMA8GA1UdJQQIMAYGBFUdJQAw\nKQYDVR0OBCIEIMjAiM2eMzlQ9HzV9ePW69rfUiRZVT2pDBOMqM4WVJSAMCsGA1Ud\nIwQkMCKAIDUkP3EcubfENS6TH3DFczH5dAnC2eD73+wcUF/bEIlnMAoGCCqGSM49\nBAMCA0kAMEYCIQCWUHL0xisjQoW+o6VV12pBXIRJgdeUeAu2EIjptSg2GAIhAIxK\nLXpHIBFxIkmWlxUaanCojPSZhzEbd+8LRrmhEO8n\n-----END CERTIFICATE-----"),
			&common.Payload{
				ChainId: "chain1",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(nil, nil, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			cli.privateKey, err = asym.PrivateKeyFromPEM(tt.userSignKeyBytes, nil)
			require.Nil(t, err)

			cli.userCrt, err = utils.ParseCert(tt.userSignCrtBytes)
			require.Nil(t, err)

			e, err := cli.SignPayload(tt.unsignedPayload)
			require.Equal(t, tt.wantErr, err != nil)

			payloadBz, err := proto.Marshal(tt.unsignedPayload)
			require.Nil(t, err)

			var opts crypto.SignOpts
			hashalgo, err := bcx509.GetHashFromSignatureAlgorithm(cli.userCrt.SignatureAlgorithm)
			require.Nil(t, err)

			opts.Hash = hashalgo
			opts.UID = crypto.CRYPTO_DEFAULT_UID

			verified, err := cli.userCrt.PublicKey.VerifyWithOpts(payloadBz, e.Signature, &opts)
			require.Nil(t, err)
			require.True(t, verified)
		})
	}
}

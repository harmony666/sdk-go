/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker/common/v2/crypto"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func TestSignPayloadWithPath(t *testing.T) {
	tests := []struct {
		name            string
		unsignedPayload *common.Payload
		wantErr         bool
	}{
		{
			"good",
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

			e, err := SignPayloadWithPath(utils.Config.ChainClientConfig.UserSignKeyFilePath,
				utils.Config.ChainClientConfig.UserSignCrtFilePath, tt.unsignedPayload)
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

func TestZXAddress(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			"pkHex",
			"3059301306072a8648ce3d020106082a811ccf5501822d034200044a4c24cf037b0c7a027e634b994a5fdbcd0faa718ce9053e3f75fcb9a865523a605aff92b5f99e728f51a924d4f18d5819c42f9b626bdf6eea911946efe7442d",
			"ZXaaa6f45415493ffb832ca28faa14bef5c357f5f0",
		},
		{
			"pkPEM",
			`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAESkwkzwN7DHoCfmNLmUpf280PqnGM
6QU+P3X8uahlUjpgWv+Stfmeco9RqSTU8Y1YGcQvm2Jr327qkRlG7+dELQ==
-----END PUBLIC KEY-----`,
			"ZXaaa6f45415493ffb832ca28faa14bef5c357f5f0",
		},
		{
			"certPEM",
			`-----BEGIN CERTIFICATE-----
MIICzjCCAi+gAwIBAgIDCzLUMAoGCCqGSM49BAMCMGoxCzAJBgNVBAYTAkNOMRAw
DgYDVQQIEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMRAwDgYDVQQKEwd3eC1v
cmcxMRAwDgYDVQQLEwdyb290LWNhMRMwEQYDVQQDEwp3eC1vcmcxLWNhMB4XDTIw
MTAyOTEzMzgxMFoXDTMwMTAyNzEzMzgxMFowcDELMAkGA1UEBhMCQ04xEDAOBgNV
BAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWppbmcxEDAOBgNVBAoTB3d4LW9yZzEx
EzARBgNVBAsTCkNoYWluTWFrZXIxFjAUBgNVBAMTDXVzZXIxLnd4LW9yZzEwgZsw
EAYHKoZIzj0CAQYFK4EEACMDgYYABAGLEJZriYzK9Se/vMGfkwjhU55eEZsM2iKM
emSZICh/HY37uR0BFAVUjMYEj84tJBzEEzlpD+AUAe44/b11b+GCMwDXPKcsjHK0
jsAPrN5LH7uptXsjMFpN2bbOqvj6sAIDfTV9chuF91LxCjYnh+Lya0ikextGkpbp
HOvi5eQ/yUHSQaN7MHkwDgYDVR0PAQH/BAQDAgGmMA8GA1UdJQQIMAYGBFUdJQAw
KQYDVR0OBCIEIAp+6tWmoiE0KmdtpLFBZpBj1Ni7JH8g2XPgoQwhQS8qMCsGA1Ud
IwQkMCKAIMsnP+UWEyGuyEHBn7JkJzb+tfBqsRCBUIPyMZH4h1HPMAoGCCqGSM49
BAMCA4GMADCBiAJCAIENc8ip2BP4yJpj9SdR9pvZc4/qbBzKucZQaD/GT2sj0FxH
hp8YLjSflgw1+uWlMb/WCY60MyxZr/RRsTYpHu7FAkIBSMAVxw5RYySsf4J3bpM0
CpIO2ZrxkJ1Nm/FKZzMLQjp7Dm//xEMkpCbqqC6koOkRP2MKGSnEGXGfRr1QgBvr
8H8=
-----END CERTIFICATE-----`,
			"ZX0787b8affa4cbdb9994548010c80d9741113ae78",
		},
		{
			"certPath",
			`-----BEGIN CERTIFICATE-----
MIICzjCCAi+gAwIBAgIDCzLUMAoGCCqGSM49BAMCMGoxCzAJBgNVBAYTAkNOMRAw
DgYDVQQIEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMRAwDgYDVQQKEwd3eC1v
cmcxMRAwDgYDVQQLEwdyb290LWNhMRMwEQYDVQQDEwp3eC1vcmcxLWNhMB4XDTIw
MTAyOTEzMzgxMFoXDTMwMTAyNzEzMzgxMFowcDELMAkGA1UEBhMCQ04xEDAOBgNV
BAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWppbmcxEDAOBgNVBAoTB3d4LW9yZzEx
EzARBgNVBAsTCkNoYWluTWFrZXIxFjAUBgNVBAMTDXVzZXIxLnd4LW9yZzEwgZsw
EAYHKoZIzj0CAQYFK4EEACMDgYYABAGLEJZriYzK9Se/vMGfkwjhU55eEZsM2iKM
emSZICh/HY37uR0BFAVUjMYEj84tJBzEEzlpD+AUAe44/b11b+GCMwDXPKcsjHK0
jsAPrN5LH7uptXsjMFpN2bbOqvj6sAIDfTV9chuF91LxCjYnh+Lya0ikextGkpbp
HOvi5eQ/yUHSQaN7MHkwDgYDVR0PAQH/BAQDAgGmMA8GA1UdJQQIMAYGBFUdJQAw
KQYDVR0OBCIEIAp+6tWmoiE0KmdtpLFBZpBj1Ni7JH8g2XPgoQwhQS8qMCsGA1Ud
IwQkMCKAIMsnP+UWEyGuyEHBn7JkJzb+tfBqsRCBUIPyMZH4h1HPMAoGCCqGSM49
BAMCA4GMADCBiAJCAIENc8ip2BP4yJpj9SdR9pvZc4/qbBzKucZQaD/GT2sj0FxH
hp8YLjSflgw1+uWlMb/WCY60MyxZr/RRsTYpHu7FAkIBSMAVxw5RYySsf4J3bpM0
CpIO2ZrxkJ1Nm/FKZzMLQjp7Dm//xEMkpCbqqC6koOkRP2MKGSnEGXGfRr1QgBvr
8H8=
-----END CERTIFICATE-----`,
			"ZX0787b8affa4cbdb9994548010c80d9741113ae78",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err  error
				addr string
			)

			if tt.name == "pkHex" {
				addr, err = GetZXAddressFromPKHex(tt.input)
			} else if tt.name == "pkPEM" {
				addr, err = GetZXAddressFromPKPEM(tt.input)
			} else if tt.name == "certPEM" {
				addr, err = GetZXAddressFromCertPEM(tt.input)
			} else if tt.name == "certPath" {
				var tmpFile *os.File
				tmpFile, err = ioutil.TempFile(os.TempDir(), "zx-*.crt")
				require.Nil(t, err)

				defer os.Remove(tmpFile.Name())

				_, err = tmpFile.Write([]byte(tt.input))
				require.Nil(t, err)

				addr, err = GetZXAddressFromCertPath(tmpFile.Name())
			}

			fmt.Printf("ZX address: %s\n", addr)
			require.Nil(t, err)
			require.Equal(t, tt.expect, addr)
		})
	}
}

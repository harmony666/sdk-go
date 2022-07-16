/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"io/ioutil"
	"testing"

	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/stretchr/testify/require"
)

func TestChainClient_ChangeSigner(t *testing.T) {
	tests := []struct {
		name                string
		userSignCrtFilePath string
		userSignKeyFilePath string
	}{
		{
			"PermissionedWithCert mode",
			"./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt",
			"./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key",
		},
		{
			"Public or PermissionedWithKey mode",
			"",
			"./testdata/crypto-config-pk/public/admin/admin1/admin1.key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSignKeyPem, err := ioutil.ReadFile(tt.userSignKeyFilePath)
			require.Nil(t, err)
			privKey, err := asym.PrivateKeyFromPEM(userSignKeyPem, nil)
			require.Nil(t, err)
			pubKey := privKey.PublicKey()

			var userSignCrtPem []byte
			var userSignCrt *bcx509.Certificate
			if tt.userSignCrtFilePath != "" {
				userSignCrtPem, err = ioutil.ReadFile(tt.userSignCrtFilePath)
				require.Nil(t, err)
				userSignCrt, err = utils.ParseCert(userSignCrtPem)
				require.Nil(t, err)
			}

			cc := &ChainClient{}
			err = cc.ChangeSigner(privKey, userSignCrt)
			require.Nil(t, err)

			pkPem, err := pubKey.String()
			require.Nil(t, err)

			require.Equal(t, cc.publicKey, pubKey)
			require.Equal(t, cc.pkBytes, []byte(pkPem))

			require.Equal(t, cc.privateKey, privKey)

			require.Equal(t, cc.userCrtBytes, userSignCrtPem)
			require.Equal(t, cc.userCrt, userSignCrt)
		})
	}
}

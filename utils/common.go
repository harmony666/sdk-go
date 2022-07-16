/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto/hash"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	guuid "github.com/google/uuid"
)

const (
	// SUCCESS ContractResult success code
	SUCCESS uint32 = 0
	// Separator chainmker hex ca
	Separator = byte(202)
)

func GetRandTxId() string {
	return uuid.GetUUID() + uuid.GetUUID()
}

// GetTimestampTxId by current time, see: GetTimestampTxIdByNano
// eg: 687dca1d9c4fdf1652fdfc072182654f53622c496aa94c05b47d34263cd99ec9
func GetTimestampTxId() string {
	return GetTimestampTxIdByNano(time.Now().UnixNano())
}

// GetTimestampTxIdByNano nanosecond
func GetTimestampTxIdByNano(nano int64) string {
	b := make([]byte, 16, 32)
	binary.BigEndian.PutUint64(b, uint64(nano))
	/*
		Read generates len(p) random bytes from the default Source and
		writes them into p. It always returns len(p) and a nil error.
		Read, unlike the Rand.Read method, is safe for concurrent use.
	*/
	b[8] = Separator
	_, _ = rand.Read(b[9:16])
	u := guuid.New()
	b = append(b, u[:]...)
	return hex.EncodeToString(b)
}

// GetNanoByTimestampTxId validate and parse 160 ... 223 22
func GetNanoByTimestampTxId(timestampTxId string) (nano int64, err error) {
	b, err := hex.DecodeString(timestampTxId)
	if err != nil {
		return
	}
	if b[8] != Separator {
		err = errors.New("not timestamp tx id")
		return
	}
	nano = int64(binary.BigEndian.Uint64(b[:8]))
	return
}

func CheckProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		if resp.Message == "" {
			if resp.ContractResult != nil && resp.ContractResult.Code != SUCCESS {
				return errors.New(resp.ContractResult.Message)
			}
			return errors.New(resp.Code.String())
		}
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != SUCCESS {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

func GetCertificateId(certPEM []byte, hashType string) ([]byte, error) {
	if certPEM == nil {
		return nil, fmt.Errorf("get cert certPEM == nil")
	}
	certDer, _ := pem.Decode(certPEM)
	if certDer == nil {
		return nil, fmt.Errorf("invalid certificate")
	}
	return GetCertificateIdFromDER(certDer.Bytes, hashType)
}

func GetCertificateIdFromDER(certDER []byte, hashType string) ([]byte, error) {
	if certDER == nil {
		return nil, fmt.Errorf("get cert from der certDER == nil")
	}
	id, err := hash.GetByStrType(hashType, certDER)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func ParseCert(crtPEM []byte) (*bcx509.Certificate, error) {
	certBlock, _ := pem.Decode(crtPEM)
	if certBlock == nil {
		return nil, fmt.Errorf("decode pem failed, invalid certificate")
	}

	cert, err := bcx509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 parse cert failed, %s", err)
	}

	return cert, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsArchived(txStatusCode common.TxStatusCode) bool {
	return txStatusCode == common.TxStatusCode_ARCHIVED_BLOCK || txStatusCode == common.TxStatusCode_ARCHIVED_TX
}

func IsArchivedString(txStatusCode string) bool {
	return txStatusCode == common.TxStatusCode_ARCHIVED_BLOCK.String() ||
		txStatusCode == common.TxStatusCode_ARCHIVED_TX.String()
}

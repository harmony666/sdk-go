package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	cert, err := ioutil.ReadFile("admin1.sign.crt")
	if err != nil {
		log.Fatal(err)
	}
	input := []byte(cert)

	// 演示base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	// fmt.Println(encodeString)

	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(decodeBytes))

	// 对pem证书进行解码在解析
	certBlock, _ := pem.Decode(decodeBytes)
	if certBlock == nil {
		fmt.Println(err)
		return
	}

	certBody, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	org := certBody.Subject.Organization[0]
	fmt.Println(org)

}

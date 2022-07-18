//https://www.cnblogs.com/unqiang/p/6677208.html
package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	input := []byte("helloworld")

	// 演示base64编码、EncodeToString返回src的base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	fmt.Println("编码后：" + encodeString)

	// 对上面的编码结果进行base64解码、DecodeString返回由base64字符串s表示的字节。
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("解码后：" + string(decodeBytes))

	// fmt.Println("///////////////////////////////")

	// // 如果要用在url中，需要使用URLEncoding
	// uEnc := base64.URLEncoding.EncodeToString([]byte(input))
	// fmt.Println(uEnc)

	// uDec, err := base64.URLEncoding.DecodeString(uEnc)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(string(uDec))
}

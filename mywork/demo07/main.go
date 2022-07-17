package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
	"github.com/yalp/jsonpath"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Config Config
}

type Config struct {
	Models []Model
	Acls   []Acl
}

type Model struct {
	Name   string
	Schema string
}

type Acl struct {
	Model     string
	Role      string
	Operation string
	Action    string
}

func main() {
	engine := gin.Default()
	engine.POST("/contractver", func(context *gin.Context) {
		cert := context.GetHeader("Authorization")
		bodyAsByteArray, _ := ioutil.ReadAll(context.Request.Body)
		jsonBody := string(bodyAsByteArray)
		raw := []byte(jsonBody)
		var data interface{}
		json.Unmarshal(raw, &data)
		// 拿到post请求发送的json中的datatype的值
		datatype, err := jsonpath.Read(data, "$.datatype")
		if err != nil {
			panic(err)
		}

		// 获取证书参数
		decodeBytes, err := base64.StdEncoding.DecodeString(cert)
		if err != nil {
			log.Fatalln(err)
		}
		certBlock, _ := pem.Decode([]byte(string(decodeBytes)))
		if certBlock == nil {
			fmt.Println("1")
			return
		}
		certBody, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			fmt.Println(err)
			return
		}

		org := certBody.Subject.Organization[0]

		//读取yaml文件
		conf := GetConf()
		num := len(conf.Config.Models) //打印结果num为2
		//找到schema文件路径
		i := 0
		schemapath := ""
	Loop:
		for i < num {
			if conf.Config.Models[i].Name == datatype {
				schemapath = conf.Config.Models[i].Schema
				break Loop
			}
			i = i + 1
		}

		if schemapath == "" {
			// 返回401
			context.JSON(401, gin.H{"message": "no that item"})
		}

		// 判断是否通过schema校验，通过校验后检查组织名称是否在acl中，接着判断这个组织是否有访问接口的权限
		j := 0
		flag := 0
		num2 := len(conf.Config.Acls)
		if TestSchema(schemapath, bodyAsByteArray) {
		Loop2:
			for j < num2 {
				if org == conf.Config.Acls[j].Role {
					if conf.Config.Acls[j].Operation == "create" {
						context.JSON(200, gin.H{"message": "You can access the interface!"})
						flag = 1
						break Loop2
					} else {
						// 返回401
						context.JSON(401, gin.H{"message": "No permission!"})
					}

				}
				j = j + 1
			}
			if flag == 0 {
				context.JSON(401, gin.H{"message": "The organization cannot be found!"})
			}
		} else {
			context.JSON(401, gin.H{"message": "JSON verification failed"})
		}
	})
	engine.Run()

}

func GetConf() Conf {
	var conf Conf
	// 加载文件
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	// 将读取的yaml文件解析为响应的 struct
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return conf
}

func TestSchema(schemapath string, json []byte) bool {
	schemaContent, err := ioutil.ReadFile(schemapath)

	if err != nil {
		panic(err.Error())
	}
	jsonContent := json
	// fmt.Println(jsonContent)
	if err != nil {
		panic(err.Error())
	}

	loader1 := gojsonschema.NewStringLoader(string(schemaContent))
	schema, err := gojsonschema.NewSchema(loader1)
	if err != nil {
		panic(err.Error())
	}

	documentLoader := gojsonschema.NewStringLoader(string(jsonContent))
	result, err := schema.Validate(documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
		return true
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return false
	}
}

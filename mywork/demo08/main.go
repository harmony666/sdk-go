//package conf
package main

import (
	"fmt"
	"io/ioutil"

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

func GetConf() Conf {
	var conf Conf // 加载文件
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	} // 将读取的yaml文件解析为响应的 struct
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return conf
}

func main() {
	conf := GetConf()
	fmt.Println(conf.Config.Models[0].Name)
	fmt.Println(conf.Config.Models[0].Schema)

	fmt.Println(conf.Config.Models[1].Name)
	fmt.Println(conf.Config.Models[1].Schema)

	fmt.Println(conf.Config.Acls[0].Model)
	fmt.Println(conf.Config.Acls[0].Role)
	fmt.Println(conf.Config.Acls[0].Operation)
	fmt.Println(conf.Config.Acls[0].Action)

	fmt.Println(conf.Config.Acls[1].Model)
	fmt.Println(conf.Config.Acls[1].Role)
	fmt.Println(conf.Config.Acls[1].Operation)
	fmt.Println(conf.Config.Acls[1].Action)

}

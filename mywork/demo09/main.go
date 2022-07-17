package main

import (
	"fmt"
	"io/ioutil"

	"github.com/xeipuuv/gojsonschema"
)

func main() {
	schemapath := "schema/contract_schema.json"
	json, _ := ioutil.ReadFile("document.json")
	flag := TestSchema(schemapath, json)
	fmt.Println(flag)
}
func TestSchema(schemapath string, json []byte) bool {
	schemaContent, err := ioutil.ReadFile(schemapath)
	if err != nil {
		panic(err.Error())
	}
	jsonContent := json
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

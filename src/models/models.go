package models

import "io/ioutil"
import "gopkg.in/yaml.v3"

type ApiKey struct {
	XApiKey string `yaml:"x-api-key"`
}

func (ApiKey) GetXApiKey() *ApiKey {

	yamlFile, _ := ioutil.ReadFile("resources/ApiKey.yaml")
	data := &ApiKey{}
	_ = yaml.Unmarshal(yamlFile, data)

	return data
}

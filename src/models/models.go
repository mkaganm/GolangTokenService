package models

import (
	"GolangJWTService/src/errors"
	"io/ioutil"
)
import "gopkg.in/yaml.v3"

type ApiKey struct {
	XApiKey string `yaml:"x-api-key"`
}

func (ApiKey) GetXApiKey() *ApiKey {

	yamlFile, _ := ioutil.ReadFile("resources/ApiKey.yaml")
	data := &ApiKey{}
	err := yaml.Unmarshal(yamlFile, data)
	errors.CheckErr(err)

	return data
}

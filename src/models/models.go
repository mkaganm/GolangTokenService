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

	yamlFile, err := ioutil.ReadFile("resources/ApiKey.yaml")
	errors.CheckErr(err)

	data := &ApiKey{}

	err = yaml.Unmarshal(yamlFile, data)
	errors.CheckErr(err)

	return data
}

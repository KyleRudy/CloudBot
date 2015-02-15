package main

import (
	"encoding/json"
	"os"
)

type CloudbotConfiguration struct {
	Server          string `json=server`
	Nick            string `json=nick`
	User            string `json=user`
	OperPwd         string `json=operpwd`
	LookupLocations string `json=lookuplocations`
}

func RetrieveConfig() (CloudbotConfiguration, error) {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := CloudbotConfiguration{}
	err := decoder.Decode(&configuration)
	return configuration, err
}

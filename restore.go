package main

import (
	"encoding/base64"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"strings"
)

func Restore(ipAddress string, token string, infile string) {

	config := api.DefaultConfig()
	config.Address = ipAddress
	config.Token = token

	data, err := ioutil.ReadFile(infile)
	if err != nil {
		panic(err)
	}

	client, _ := api.NewClient(config)
	kv := client.KV()

	pairs := DeleteEmpty(strings.Split(string(data), "\n"))

	for _, element := range pairs {
		pair := strings.Split(element, ":")

		if len(pair) > 1 {
			val, decodeErr := base64.StdEncoding.DecodeString(pair[1])
			if decodeErr != nil {
				panic(decodeErr)
			}

			p := &api.KVPair{Key: pair[0], Value: val}
			_, err := kv.Put(p, nil)
			if err != nil {
				panic(err)
			}
		}
	}
	logger.Infof("Wrote %d KV pairs to Consul at %s from file %s", len(pairs), ipAddress, infile)
}

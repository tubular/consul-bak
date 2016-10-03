package main

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
	"sort"
)

// ByCreateIndex sorts consul KV pairs
type ByCreateIndex api.KVPairs

func (a ByCreateIndex) Len() int           { return len(a) }
func (a ByCreateIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreateIndex) Less(i, j int) bool { return a[i].CreateIndex < a[j].CreateIndex }

// Backup talks to Consul and starts reading KV pais and writes them to a file
func Backup(ipAddress string, token string, outfile string, exclusion []string, inclusion []string) {

	config := api.DefaultConfig()
	config.Address = ipAddress
	config.Token = token

	client, _ := api.NewClient(config)
	kv := client.KV()

	pairs, _, err := kv.List("/", nil)
	if err != nil {
		panic(err)
	}

	sort.Sort(ByCreateIndex(pairs))

	if len(pairs) == 0 {
		logger.Info("No KV pairs found. Nothing to do. Exiting.")
		os.Exit(0)
	} else {
		logger.Infof("Found %d KV pairs", len(pairs))
	}

	outstring := ""
	if len(exclusion) > 0 {
		for _, element := range pairs {
			if !StartsWith(exclusion, element.Key) {
				encodedValue := base64.StdEncoding.EncodeToString(element.Value)
				outstring += fmt.Sprintf("%s:%s\n", element.Key, encodedValue)
			}
		}
	} else if len(inclusion) > 0 {
		for _, element := range pairs {
			if StartsWith(inclusion, element.Key) {
				encodedValue := base64.StdEncoding.EncodeToString(element.Value)
				outstring += fmt.Sprintf("%s:%s\n", element.Key, encodedValue)
			}
		}
	} else {
		for _, element := range pairs {
			encodedValue := base64.StdEncoding.EncodeToString(element.Value)
			outstring += fmt.Sprintf("%s:%s\n", element.Key, encodedValue)
		}
	}

	file, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}

	if _, err := file.Write([]byte(outstring)[:]); err != nil {
		panic(err)
	}

	logger.Infof("Wrote %d KV pairs to %s", len(pairs), outfile)
}

// BackupACLs talks to Consul and starts printing ACL rules currently in place
func BackupACLs(ipAddress string, token string, outfile string) {

	config := api.DefaultConfig()
	config.Address = ipAddress
	config.Token = token

	client, _ := api.NewClient(config)
	acl := client.ACL()

	tokens, _, err := acl.List(nil)
	if err != nil {
		panic(err)
	}

	outstring := ""
	for _, element := range tokens {
		outstring += fmt.Sprintf("====\nID: %s\nName: %s\nType: %s\nRules:\n%s\n", element.ID, element.Name, element.Type, element.Rules)
	}

	file, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}

	if _, err := file.Write([]byte(outstring)[:]); err != nil {
		panic(err)
	}
}

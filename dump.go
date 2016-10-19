package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Dump(ipAddress string, token string, rootPath string) {

	var standardDirPerms os.FileMode = 0775
	logger.Infof("Starting dump from %s to %s.", ipAddress, rootPath)

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

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)

	absPath, pathErr := filepath.Abs(rootPath)
	logger.Debugf("Path to dump kv in is: %s", absPath)
	Check(pathErr)
	ensureErr := EnsureDir(absPath, standardDirPerms)
	Check(ensureErr)
	chdirErr := os.Chdir(absPath)
	Check(chdirErr)
	baseDir, _ := os.Getwd()
	logger.Debugf("Changed directory to: %s", baseDir)

	for _, element := range pairs {
		var path string
		path = fmt.Sprintf("%s/%s", baseDir, element.Key)
		logger.Infof("Writing %s", path)
		if strings.HasSuffix(path, "/") {
			err := EnsureDir(path, standardDirPerms)
			Check(err)
		} else {
			// ensure directory is present before creating file
			err := EnsureDir(filepath.Dir(path), standardDirPerms)
			Check(err)

			f, createErr := os.Create(path)
			Check(createErr)
			f.Write(element.Value)
		}
	}
}

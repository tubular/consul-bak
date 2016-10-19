package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Sync reads a repo specified by `gitString` and writes it to Consul at `ipAddress`
func Sync(ipAddress string, token string, gitString string) {
	logger.Infof("Got git string: %s", gitString)
	const cacheDir = "/var/cache/consul_backup/"

	parts := strings.Split(gitString, "|")
	gitURL := parts[0]
	kvRootPath := parts[1]
	logger.Debugf("URL: %s, KV path: %s", gitURL, kvRootPath)

	if !strings.HasSuffix(kvRootPath, "/") {
		kvRootPath = kvRootPath + "/"
	}

	var rootDir string = fmt.Sprintf("%s%s", cacheDir, kvRootPath)

	if !Which("git") {
		logger.Error("You must have git installed to use syncgit mode.")
		os.Exit(1)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	repoExists, existsErr := PathExists(fmt.Sprintf("%s.git", cacheDir))
	Check(existsErr)

	if !repoExists {
		logger.Infof("Repo at %s does not exist yet, creating.", cacheDir)
		var perm os.FileMode = 0775
		mkdirErr := os.MkdirAll(cacheDir, perm)
		Check(mkdirErr)
		GitBinaryCall(fmt.Sprintf("clone %s %s", gitURL, cacheDir))
		chdirErr := os.Chdir(cacheDir)
		Check(chdirErr)
		GitBinaryCall("config core.sparseCheckout true")
		f, createErr := os.Create(".git/info/sparse-checkout")
		Check(createErr)
		f.WriteString(kvRootPath)
	}
	changeDirErr := os.Chdir(rootDir)
	Check(changeDirErr)
	currentDir, _ := os.Getwd()

	logger.Info("Pulling latest info from git.")
	GitBinaryCall("checkout master")
	GitBinaryCall("remote update")
	GitBinaryCall("pull origin master")

	config := api.DefaultConfig()
	config.Address = ipAddress
	config.Token = token

	client, _ := api.NewClient(config)
	kv := client.KV()

	logger.Info("Creating a list of keys")
	fileList := []string{}
	walkErr := filepath.Walk(currentDir, func(path string, f os.FileInfo, err error) error {
		logger.Info(path)
		fileList = append(fileList, path)
		return nil
	})
	Check(walkErr)

	for _, fileName := range fileList {
		fi, statErr := os.Stat(fileName)
		Check(statErr)

		var key string
		var val []byte
		var readErr error

		switch mode := fi.Mode(); {
		case mode.IsDir():
			key = fmt.Sprintf("%s/", fileName)
		case mode.IsRegular():
			key = fileName
			val, readErr = ioutil.ReadFile(fileName)
			Check(readErr)
		}
		key = strings.Replace(key, rootDir, "", 1)
		if key != "" {
			p := &api.KVPair{Key: key, Value: val}
			_, putErr := kv.Put(p, nil)
			Check(putErr)
		}
	}
	logger.Infof("Wrote %d KV pairs to Consul at %s from Git repo %s", len(fileList), ipAddress, gitURL)
}

package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"os"
	"strings"
)

func main() {

	usage := fmt.Sprintf(`Consul KV and ACL Backup with KV Restore tool.

Version: %s (Commit: %s)

Usage:
  consul-backup [-i IP] [--http-port HTTPPORT] [--rpc-port RPCPORT]
                [-l] [-t TOKEN] [-a] [-b ACLBACKUPFILE] [-n INPREFIX]...
                [-x EXPREFIX]... [--restore] [--no-prompt] <filename>
  consul-backup -h | --help
  consul-backup --version

Options:
  -h --help                          Show this screen.
  --version                          Show version.
  -l, --leader-only                  Create backup only on consul leader.
  --rpc-port=RPCPORT                 RPC port [default: 8400].
  --http-port=HTTPPORT               HTTP endpoint port [default: 8500].
  -i, --address=IP                   The HTTP endpoint of Consul [default: 127.0.0.1].
  -t, --token=TOKEN                  An ACL Token with proper permissions in Consul [default: ].
  -a, --aclbackup                    Backup ACLs, does nothing in restore mode. ACL restore not available at this time.
  -b, --aclbackupfile=ACLBACKUPFILE  ACL Backup Filename [default: acl.bkp].
  -x, --exclude-prefix=[EXPREFIX]    Repeatable option for keys starting with prefix to exclude from the backup.
  -n, --include-prefix=[INPREFIX]    Repeatable option for keys starting with prefix to include in the backup.
  -r, --restore                      Activate restore mode.
  --no-prompt                        Don't prompt, force overwrite in restore mode.`, Version, GitCommit)

	arguments, _ := docopt.Parse(usage, nil, true, fmt.Sprintf("consul-backup %s (%s)", Version, GitCommit), false)
	SetupLogging()

	var (
		httpEndpoint = fmt.Sprintf("%s:%s", arguments["--address"], arguments["--http-port"])
		rpcEndpoint  = fmt.Sprintf("%s:%s", arguments["--address"], arguments["--rpc-port"])
		rpcOptString = fmt.Sprintf("-rpc-addr=%s", rpcEndpoint)
	)

	logger.Infof("Verifying HTTP endpoint: %s", httpEndpoint)
	CheckSocket(httpEndpoint)
	logger.Infof("Verifying RPC endpoint: %s", rpcEndpoint)
	CheckSocket(rpcEndpoint)

	if arguments["--leader-only"] == true {
		logger.Info("Running in leader only mode, only running backup/restore on Consul leader.")
		// if consul client is not available we keep running
		if Which("consul") {
			var out = ConsulBinaryCall("info", rpcOptString)
			if strings.Contains(out, "leader = true") {
				logger.Info("Consul leader, continuing.")
			} else {
				logger.Error("Not a Consul leader, stopping.")
				os.Exit(1)
			}
		} else {
			logger.Error("Could not find `consul` utility. Is your $PATH setup properly?")
			os.Exit(1)
		}
	}

	if arguments["--restore"] == true {
		logger.Info("Running in restore mode.")
		if (len(arguments["--exclude-prefix"].([]string)) > 0) || (len(arguments["--include-prefix"].([]string)) > 0) {
			logger.Error("--exclude-prefix, -x and --include-prefix, -n can be used only for backups")
			os.Exit(1)
		}
		if arguments["--no-prompt"] == false {
			fmt.Printf("\nWarning! This will overwrite existing kv. Press [enter] to continue; CTL-C to exit")
			fmt.Scanln()
		}
		logger.Infof("Restoring KV from file: %s", arguments["<filename>"].(string))
		Restore(httpEndpoint, arguments["--token"].(string), arguments["<filename>"].(string))
	} else {
		logger.Info("Running in backup mode.")
		if (len(arguments["--exclude-prefix"].([]string)) > 0) && (len(arguments["--include-prefix"].([]string)) > 0) {
			logger.Error("--exclude-prefix and --include-prefix cannot be used together")
			os.Exit(1)
		}
		if len(arguments["--exclude-prefix"].([]string)) > 0 {
			logger.Infof("excluding keys with prefix(es): %s", arguments["--exclude-prefix"].([]string))
		}
		if len(arguments["--include-prefix"].([]string)) > 0 {
			logger.Infof("including only keys with prefix(es): %s", arguments["--include-prefix"].([]string))
		}
		logger.Infof("KV store will be backed up to file: %s", arguments["<filename>"].(string))
		Backup(httpEndpoint, arguments["--token"].(string), arguments["<filename>"].(string), arguments["--exclude-prefix"].([]string), arguments["--include-prefix"].([]string))
		if arguments["--aclbackup"] == true {

			logger.Infof("ACL Tokens will be backed up to file: %s", arguments["--aclbackupfile"].(string))
			BackupACLs(httpEndpoint, arguments["--token"].(string), arguments["--aclbackupfile"].(string))
		}
	}
}

# Consul Bak [![Build Status](https://travis-ci.org/Tubular/consul-bak.png)](https://travis-ci.org/Tubular/consul-bak)

Backs up, restores, syncs and dumps KV pairs in a Consul cluster using the consul-api Go library.

## Changelog / Releases

[See the changelog here](CHANGELOG.md)


## Modes


### Backup 

This backups your Consul master's KV store into a file. You can run this via a crontab and sync 
it to S3 for example:

```
filename="$(date +%Y%m%d).txt"
/usr/local/bin/consul-bak backup --leader-only $filename
/usr/local/bin/aws s3 mv $filename s3://<my_bucket>/
```


### Restore

This restores a file created by consul-bak into the Consul master's KV store. Example:

```
/usr/local/bin/consul-bak restore $filename
```


### Syncgit

This synchronises a filesystem representation of the KV tree into a Consul master's KV store. The
filesystem representation is expected to be a git repository. This is especially useful if you want
a code-first approach to creating configuration in Consul. Example:

```
/usr/local/bin/consul-bak sync git@github.com:Tubular/consul-bak.git|path/to/root/of/kv/tree
```


### Dumptree

This dumps the values of the Consul master's KV as a tree onto the filesystem. If you want to start
using sync, you could use this to kickstart your project. Example:

```
/usr/local/bin/consul-bak dumptree /tmp/kv_tree
```


## Usage examples

```sh
Usage:
  consul-bak (backup|restore|aclbackup)
                [--leader-only]
                [--rpc-port RPCPORT]
                [--http-port HTTPPORT]
                [--address IP]
                [--include-prefix INPREFIX]...
                [--exclude-prefix EXPREFIX]...
                [--token TOKEN]
                [--no-prompt]
                <filename>
  consul-bak dumptree
                [--leader-only]
                [--rpc-port RPCPORT]
                [--http-port HTTPPORT]
                [--address IP]
                <pathname>
  consul-bak syncgit
                [--leader-only]
                [--rpc-port RPCPORT]
                [--http-port HTTPPORT]
                [--address IP]
                <git-url>
  consul-bak -h | --help
  consul-bak -v | --version

Options:
  -h, --help                         Show this screen.
  -v, --version                      Show version.
	--mode=MODE                        Set mode, can be one of backup,restore,syncgit,dumptree,aclbackup [default: backup]
  --leader-only                      Only run on consul leader.
  --rpc-port=RPCPORT                 RPC port [default: 8400].
  --http-port=HTTPPORT               HTTP endpoint port [default: 8500].
  --address=IP                       The HTTP endpoint of Consul [default: 127.0.0.1].
  --include-prefix=[INPREFIX]        Repeatable option for keys starting with prefix to include in the backup.
  --exclude-prefix=[EXPREFIX]        Repeatable option for keys starting with prefix to exclude from the backup.
  --token=TOKEN                      An ACL Token with proper permissions in Consul [default: ].
  --force                            Don't prompt, force overwrite.
```



## Development

```
# Somewhere in your GOPATH
git clone git@github.com:Tubular/consul-bak.git

# Get requirements
glide install

# Run
make run
```

## Release

To bump a version, change the `Version` constant in version.go and use git tag:

```
vi version.go && git commit -m "Bump version to #.#.#" && git tag -a "Released #.#.#"
```


*Code inspired by work done here: https://github.com/kailunshi/consul-backup*

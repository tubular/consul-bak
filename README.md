# Consul Backup [![Build Status](https://travis-ci.org/Tubular/consul-backup.png)](https://travis-ci.org/Tubular/consul-backup)

Backs up and restores KV pairs in a Consul cluster using the consul-api Go library.

## Changelog / Releases

[See the changelog here](CHANGELOG.md)

## Usage examples

```sh
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
  --no-prompt                        Don't prompt, force overwrite in restore mode.
```



## Development

```
# Somewhere in your GOPATH
git clone git@github.com:Tubular/consul-backug.git

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


*Code forked from work done here: https://github.com/kailunshi/consul-backup*

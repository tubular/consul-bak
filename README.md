# `consul-backup` - Consul Backup and Restore tool

This will use consul-api (Go library) to recursively backup and restore all your key/value pairs.
You need to set up your Go environment, then run `make`, which will generate executable named `consul-backup`.

## Usage examples

```sh
Usage:
  consul-backup [-i IP:PORT] [-t TOKEN] [--aclbackup] [--aclbackupfile ACLBACKUPFILE] [--restore] <filename>
  consul-backup -h | --help
  consul-backup --version

Options:
  -h --help                          Show this screen.
  --version                          Show version.
  -i, --address=IP:PORT              The HTTP endpoint of Consul [default: 127.0.0.1:8500].
  -t, --token=TOKEN                  An ACL Token with proper permissions in Consul [default: ].
  -a, --aclbackup                    Backup ACLs, does nothing in restore mode. ACL restore not available at this time.
  -b, --aclbackupfile=ACLBACKUPFILE  ACL Backup Filename [default: acl.bkp].
  -r, --restore                      Activate restore mode
```


## Development / Release



To bump a version, change the `Version` constant in version.go and use git tag:

vi version.go && git commit -m "Bump version to #.#.#" && git tag -a "Released #.#.#"


*Code forked from work done here: https://github.com/kailunshi/consul-backup*

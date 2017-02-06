## Building the Custom telegraf agent with support for sync_gateway and gateload

Note: There is an issue building telegraf from the couchbaselabs/telegraf fork directly. The following steps get the influxdb/telegraf repo source and copy over the new/modified files from couchbaselabs/telegraf.

### From Source:
Telegraf manages dependencies via [gdm](https://github.com/sparrc/gdm),
which gets installed via the Makefile
if you don't have it already. You also must build with golang version 1.5+.

1. [Install Go](https://golang.org/doc/install)
1. [Setup your GOPATH](https://golang.org/doc/code.html#GOPATH)
1. Run `mkdir $HOME/gotelegraf`
1. Run `export GOPATH=$HOME/gotelegraf`
1. Run `export PATH=$PATH:$GOPATH/bin`
1. Run `mkdir -p $GOPATH/src/github.com/influxdata`
1. Run `cd $GOPATH/src/github.com/influxdata`
1. Run `git clone git@github.com:couchbaselabs/telegraf`
1. Run `cd telegraf`
1. Run `make`

### Deploying binary to mobile-testkit repo

1. Repeat "From Source" instructions on Linux
1. `cp ./cmd/telegraf /path/to/mobile-testkit/libraries/provision/ansible/playbooks/files`
1. Commit and push


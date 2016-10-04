## Building the Custom telegraf agent with support for sync_gateway and gateload

Note: There is an issue building telegraf from the couchbaselabs/telegraf fork directly. The following steps get the influxdb/telegraf repo source and copy over the new/modified files from couchbaselabs/telegraf.

### From Source:
Telegraf manages dependencies via [gdm](https://github.com/sparrc/gdm),
which gets installed via the Makefile
if you don't have it already. You also must build with golang version 1.5+.

1. [Install Go](https://golang.org/doc/install)
2. [Setup your GOPATH](https://golang.org/doc/code.html#GOPATH)
3. Run `go get github.com/influxdata/telegraf`
4. Run `go get github.com/couchbaselabs/telegraf`
5. Run `cd $GOPATH/src/github.com/influxdata/telegraf`
6. Run `cp $GOPATH/src/github.com/couchbaselabs/telegraf/plugins/inputs/all/all.go plugins/inputs/all/all.go`
7. Run `mkdir -p plugins/inputs/syncgateway`
8. Run `cp $GOPATH/src/github.com/couchbaselabs/telegraf/plugins/inputs/syncgateway/syncgateway.go plugins/inputs/syncgateway/syncgateway.go`
8. Run `mkdir -p plugins/inputs/gateload`
9. Run `cp $GOPATH/src/github.com/couchbaselabs/telegraf/plugins/inputs/gateload/gateload.go plugins/inputs/gateload/gateload.go`
10. 9. Run `cp $GOPATH/src/github.com/couchbaselabs/telegraf/plugins/inputs/gateload/flatten.go plugins/inputs/gateload/flatten.go`
7. Run `make`
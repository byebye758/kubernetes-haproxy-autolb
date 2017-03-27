package main

import (
	"kubernetes-haproxy-autolb/backends/con"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	//"encoding/json"
	"fmt"
	//"net"
	//"strings"
	//"unsafe"
	"flag"
	"strings"
)

func main() {
	enps := flag.String("Endpoints", "10.1.10.201:2379", "etcdserverip eg:--endpoints=10.1.10.201,10.1.10.202:2379 ")
	devname := flag.String("Networkdevname", "ens160", "eg:  --Serviceip=ens160")
	flag.Parse()
	endpoints := strings.SplitN(*enps, ",", -1)
	dev := *devname

	a := new(etcd3client.Register)
	a.Endpoints = endpoints
	a.Key = "/autohaproxy/haproxyip/" + *con.Devip(dev)
	a.Value = *con.Devip(dev)
	go a.PutLease()

	fmt.Println(<-a.Ch)
}

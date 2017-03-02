package main

import (
	"fmt"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	"kubernetes-haproxy-autolb/backends/node"
	"kubernetes-haproxy-autolb/backends/watch"
)

var (
	endpoints = []string{"10.1.10.201:2379"}
	serviceip = "192.168.110.0/24"
)

func main() {
	ch := make(chan string)

	f := etcd3client.Node{endpoints, "/autohaproxy/node/nodeip/"}

	go watch.Nodeiproutewatch("/autohaproxy/node/nodeip/", endpoints, f, ch)

	node.Iproute(f, endpoints)
	for {

		fmt.Println(<-ch)

	}
}

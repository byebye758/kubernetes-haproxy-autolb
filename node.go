package main

import (
	"fmt"
	//"kubernetes-haproxy-autolb/backends/con"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	"kubernetes-haproxy-autolb/backends/node"
	"kubernetes-haproxy-autolb/backends/watch"
	//"time"
)

var (
	endpoints = []string{"10.1.10.201:2379"}
	serviceip = "192.168.110.0/24"
)

func main() {
	ch := make(chan string)
	//nodeip := con.HostIP()
	//dockerip, _ := con.Getdockerip()
	a := etcd3client.Autotable{endpoints, "/autohaproxy/autotable/"}

	f := etcd3client.Node{endpoints, "/autohaproxy/node/nodeip/"}

	// g := etcd3client.NodeRegister{
	// 	endpoints,
	// 	"/autohaproxy/node/nodeip/" + nodeip,
	// 	nodeip,
	// 	dockerip,
	// 	ch,
	// }
	//g.NodePutLease() /*注册node到etcd中*/
	//go node.Noderegister(g, ch)
	//node.Iproute(f, endpoints)
	// fmt.Println("2")
	//node.Noderoute(a, endpoints)
	// fmt.Println("3")
	node.Serviceiproute(serviceip)

	go watch.Nodeiproutewatch("/autohaproxy/node/nodeip/", endpoints, f, ch)
	go watch.Nodenoderoutewatch("/autohaproxy/autotable/", endpoints, a, ch)
	// go node.Iproute(f, endpoints)
	// go node.Noderoute(a, endpoints)
	for {

		fmt.Println(<-ch)

	}
}

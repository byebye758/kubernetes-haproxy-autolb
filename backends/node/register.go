package node

import (
	//"autohalb/backends/con"
	"autohalb/backends/etcd3client"
	//"fmt"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// )

func Noderegister(g etcd3client.NodePutLeaser, ch chan string) {

	//nodeip := con.HostIP()
	//dockerip, _ := con.Getdockerip()
	//ch := make(chan string)

	// g := etcd3client.NodeRegister{
	// 	endpoints,
	// 	"/autohaproxy/node/nodeip/" + nodeip,
	// 	nodeip,
	// 	dockerip,
	// 	ch,
	// }

	go g.NodePutLease()
	// b := etcd3client.Node{
	// 	endpoints,
	// 	"/autohaproxy/node/nodeip/",
	// }
	// fmt.Println(b.AGet())
	//fmt.Println(<-a.Ch)

}

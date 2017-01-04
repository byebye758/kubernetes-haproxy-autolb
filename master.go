package main

import (
	"./backends/autotable"
	//"autohalb/backends/con"
	"./backends/etcd3client"
	"./backends/watch"
	"flag"
	"fmt"
	"strings"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// )

func main() {

	var enps = flag.String(endpoints, "10.1.10.201", "etcdserverip eg:--endpoints=10.1.10.201,10.1.10.202 ")
	endpoints := strings.SplitN(enps, ",", -1)
	a := etcd3client.Autotable{endpoints, "/autohaproxy/autotable/"}

	b := etcd3client.Projecttable{endpoints, "/autohaproxy/project/"}

	c := etcd3client.Routetable{endpoints, "/autohaproxy/haproxytable/"}
	d := etcd3client.Pod{endpoints, "/registry/pods/default/"}
	e := etcd3client.Haip{endpoints, "/autohaproxy/haproxyip"}
	autotable.Autotable(a, b, c, e, d, endpoints)
	//autotable.Autotabledelete(d, a, endpoints)
	ch := make(chan string)

	go watch.Watch("/registry/services/endpoints/default/", endpoints, a, b, c, e, d, ch)
	go watch.Watch("/autohaproxy/haproxyip/", endpoints, a, b, c, e, d, ch)
	go watch.Watch("/registry/pods/default/", endpoints, a, b, c, e, d, ch)

	for {

		fmt.Println(<-ch)

	}
	//fmt.Println(<-ch)

}

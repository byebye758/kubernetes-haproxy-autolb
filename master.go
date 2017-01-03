package main

import (
	"autohalb/backends/autotable"
	//"autohalb/backends/con"
	"autohalb/backends/etcd3client"
	"autohalb/backends/watch"
	"fmt"
)

var (
	endpoints = []string{"10.1.10.201:2379"}
)

func main() {
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

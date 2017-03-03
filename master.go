package main

import (
	"kubernetes-haproxy-autolb/backends/autotable"
	//"autohalb/backends/con"
	"flag"
	"fmt"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	"kubernetes-haproxy-autolb/backends/watch"
	"strings"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// )

func main() {

	var enps = flag.String("Endpoints", "10.1.10.201:2379", "etcdserverip eg:--endpoints=10.1.10.201,10.1.10.202:2379 ")
	flag.Parse()

	endpoints := strings.SplitN(*enps, ",", -1)
	//fmt.Println(endpoints)

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
	go watch.Watch("/autohaproxy/project/", endpoints, a, b, c, e, d, ch)

	for {

		fmt.Println(<-ch)

	}
	//fmt.Println(<-ch)

}

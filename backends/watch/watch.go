package watch

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"kubernetes-haproxy-autolb/backends/autotable"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	"kubernetes-haproxy-autolb/backends/node"
	"log"
	"time"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 2 * time.Second
)

//func Watch(key string, endpoints []string, a, b, c, e etcd3client.AGetr, d etcd3client.Getr, add func(etcd3client.AGetr, etcd3client.AGetr, etcd3client.AGetr, etcd3client.AGetr, etcd3client.Getr, []string) /*, ch chan string*/) {

func Watch(key string, endpoints []string, a, b, c, e etcd3client.AGetr, d etcd3client.Getr, ch chan string) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(key, endpoints)
	rch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	//fmt.Println(<-rch)
	fmt.Println("abc")

	for wresp := range rch {
		fmt.Println(456)
		for _, ev := range wresp.Events {
			fmt.Println(ev.Type)

			go autotable.Autotable(a, b, c, e, d, endpoints)
			ch <- key
		}
	}

}

func Nodeiproutewatch(key string, endpoints []string, f etcd3client.AGetr, ch chan string) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(key, endpoints)
	rch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	//fmt.Println(<-rch)
	//fmt.Println("abc")

	for wresp := range rch {
		fmt.Println(key)
		for _, ev := range wresp.Events {
			fmt.Println(ev.Type)

			go node.Iproute(f, endpoints)
			ch <- key
		}
	}
}

func Nodenoderoutewatch(key string, endpoints []string, a etcd3client.AGetr, ch chan string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(key, endpoints)
	rch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	//fmt.Println(<-rch)
	//fmt.Println("abc")

	for wresp := range rch {
		fmt.Println(key)
		for _, ev := range wresp.Events {
			fmt.Println(ev.Type)

			go node.Noderoute(a, endpoints)
			ch <- key
		}
	}

}

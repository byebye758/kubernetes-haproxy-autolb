package main

import (
	"encoding/json"
	"fmt"
	//"os"
	"context"
	"flag"
	"github.com/coreos/etcd/clientv3"
	"log"
	"strings"
	"time"
	"unsafe"
)

func main() {
	enps := flag.String("Endpoints", "10.1.10.201:2379", "etcdserverip eg:--endpoints=10.1.10.201,10.1.10.202:2379 ")
	hip := flag.String("Haproxyip", "10.1.10.208", "ep: 10.1.10.208")
	htable := flag.String("Haproxytable", "10", "ep: 10")
	flag.Parse()
	haproxyip := *hip
	haproxytable := *htable
	endpoints := strings.SplitN(*enps, ",", -1)
	type Autotable struct {
		Haproxyip    string
		Haproxytable string
	}

	group := Autotable{
		Haproxyip:    haproxyip,
		Haproxytable: haproxytable,
	}

	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	Eput("/autohaproxy/haproxytable/"+group.Haproxyip, B2S(b), endpoints)
	fmt.Println(B2S(b))
}
func B2S(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}
func Eput(k, v string, endpoints []string) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = cli.Put(ctx, k, v)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}

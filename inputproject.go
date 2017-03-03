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

	projectport := make(map[string]string)
	enps := flag.String("Endpoints", "10.1.10.201:2379", "etcdserverip eg:--endpoints=10.1.10.201,10.1.10.202:2379 ")

	p := flag.String("Projectport", "80:8080,8080:8080", "ep:   80:8080,8080:8080")
	pn := flag.String("Projectname", "k8s-apiv2", "ep: k8s-apiv2")
	haip := flag.String("Haproxyip", "10.1.10.208,10.1.10.206", "ep: 10.1.10.208,10.1.10.206")

	flag.Parse()

	prot := strings.SplitN(*p, ",", -1)
	haproxyip := strings.SplitN(*haip, ",", -1)
	projectname := *pn
	endpoints := strings.SplitN(*enps, ",", -1)

	for _, v := range prot {
		v := (strings.SplitN(v, ":", -1))
		//fmt.Println(v[0], v[1])
		projectport[v[0]] = v[1]

	}

	type Project struct {
		Projectname string
		Projectport map[string]string
		Haproxyip   []string
		//Podname      string
		//Haproxytable string
		//Nodeip       string
		//Podip        string
	}

	group := Project{
		Projectname: projectname,
		Projectport: projectport,
		Haproxyip:   haproxyip,
		//Haproxytable: "20",
		//Podname:      "k8s-apiv2-1434876282-odeff",
		//Nodeip:       "10.1.10.202",
		//Podip:        "192.168.192.3",
	}

	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	Eput("/autohaproxy/project/"+group.Projectname, B2S(b), endpoints)
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

package main

import (
	"kubernetes-haproxy-autolb/backends/con"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	//"encoding/json"
	"fmt"
	//"net"
	//"strings"
	//"unsafe"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// )

func main() {
	enps := flag.String("Endpoints", "10.1.10.201:2379", "etcdserverip eg:--endpoints=10.1.10.201,10.1.10.202:2379 ")
	// sip := flag.String("Serviceip", "192.168.110.0/24", "eg:  --Serviceip=192.168.110.0/24")
	flag.Parse()
	endpoints := strings.SplitN(*enps, ",", -1)

	a := new(etcd3client.Register)
	a.Endpoints = endpoints
	a.Key = "/autohaproxy/haproxyip/" + con.HostIP()

	a.Value = con.HostIP()
	go a.PutLease()

	fmt.Println(<-a.Ch)
}

// func HostIP() (a string) {
// 	conn, err := net.Dial("udp", "1.1.1.1:80")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	defer conn.Close()

// 	a = strings.Split(conn.LocalAddr().String(), ":")[0]
// 	return
// }

// func JsonCreate(v string) string {
// 	type Ha struct {
// 		Haproxyip string
// 	}
// 	group := Ha{
// 		Haproxyip: v,
// 	}
// 	b, err := json.Marshal(group)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	return B2S(b)
// 	//os.Stdout.Write(b)
// }
// func B2S(buf []byte) string {
// 	return *(*string)(unsafe.Pointer(&buf))
// }

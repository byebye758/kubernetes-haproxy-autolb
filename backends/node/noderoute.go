package node

import (
	"fmt"
	"kubernetes-haproxy-autolb/backends/con"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	"kubernetes-haproxy-autolb/backends/node/cmd"
	"strings"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// )
/*添加pod返回到haproxy的回程路由*/

func Noderoute(a etcd3client.AGetr, endpoints []string) {
	nodemap := make(map[string]map[string]string) //etcd get  autotable   format map
	//a := etcd3client.Autotable{endpoints, "/autohaproxy/autotable/"}
	b, _ := cmd.NoderuleGet()
	fmt.Println(b, "----------NoderuleGet")

	autotable := a.AGet()
	nodeip := con.HostIP()

	fmt.Println(autotable, b, nodeip, "Noderoute")
	for _, v := range autotable {
		v := v.(map[string]interface{})
		etcdnodeip := v["Nodeip"].(string)
		fmt.Println(etcdnodeip)
		etcdhaproxyip := v["Haproxyip"].(string)
		etcdhaproxytable := v["Haproxytable"].(string)
		etcdpodip := v["Podip"].(string)
		if strings.EqualFold(nodeip, etcdnodeip) {
			test := map[string]string{
				"Nodeip":       etcdnodeip,
				"Haproxyip":    etcdhaproxyip,
				"Haproxytable": etcdhaproxytable,
				"Podip":        etcdpodip,
			}
			nodemap[etcdpodip] = test

		}

	}
	fmt.Println(nodemap, "Nodemap---------------")

	for _, v := range b {
		fmt.Println(v, "-----------bbbbb-------------")

		if _, ok := nodemap[v["Podip"]]; ok {
			fmt.Println(v["Podip"], "NodeIPOK------------------")

		} else {
			fmt.Println(v["Podip"], "PodIP------------------")
			cmd.Routetablecmd("ip rule del from "+v["Podip"], "")

		}

	}
	b1, _ := cmd.NoderuleGet()
	noderulemap := make(map[string]map[string]string)
	fmt.Println(b1, "print BBBB")
	for _, v := range b1 {
		fmt.Println(v)

		noderulemap[v["Podip"]] = v

	}
	fmt.Println(nodemap, "NODE RULE MAP ||", noderulemap)

	for _, v := range nodemap {
		podip := v["Podip"]
		haproxyip := v["Haproxyip"]
		haproxytable := v["Haproxytable"]

		if _, ok := noderulemap[podip]; ok {

		} else {
			cmd.Routetablecmd("ip rule add from "+podip+"/32 pref 30000 table ", haproxytable)
			cmd.Routetablecmd("ip route replace default via "+haproxyip+" table ", haproxytable)
		}
	}

}

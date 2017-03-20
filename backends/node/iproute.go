package node

import (
	"errors"
	"fmt"
	"kubernetes-haproxy-autolb/backends/con"
	"kubernetes-haproxy-autolb/backends/etcd3client"
	"kubernetes-haproxy-autolb/backends/log"
	"os/exec"
	"strings"
)

// var (
// 	endpoints = []string{"10.1.10.201:2379"}
// )
/* 规则7 为 本地 lo接口   路由规则   k8s service ip地址   指向 lo接口  */

/*增加到 直连路由增加到table7 中  */
func Serviceiproute(serviceip string) {

	net := con.Hosipnetwork()
	for _, v := range net {

		Routetablecmd("ip route replace "+v["ip"]+" dev "+v["devname"]+" proto kernel  scope link table ", "7")

	}

	Routetablecmd("ip route replace "+serviceip+" dev lo  scope link table ", "7")

	RuleAdd("7")

}

/*根据   etcd get 到的 路由表 刷新  node 本地 路由策略*/
func Iproute(f etcd3client.AGetr, endpoints []string) {
	//f := etcd3client.Node{endpoints, "/autohaproxy/node/nodeip/"}
	etcdnodeGet := f.AGet()
	nodeospfipGet, _ := NodeOspfIpGet()
	localnodeip := con.HostIP()
	//localdockerip, _ := con.Getdockerip()
	fmt.Println(localnodeip)

	fmt.Println(etcdnodeGet, "ETCDnodeget-----------nodeospfipGet", nodeospfipGet)
	//fmt.Println(nodeospfipGet)

	for k, v := range nodeospfipGet {

		if _, ok := etcdnodeGet[k]; ok {

		} else {
			fmt.Println(v["Dockerip"], "--------------nodeospfipGet[Dockerip]")
			if k == "docker0" {
				continue

			} else {
				ip := v["Dockerip"]
				Routetablecmd("ip route del "+ip+" table ", "8")

				fmt.Println("Delete  ---------------", ip)
			} //需要注意空格

		}

	}

	//fmt.Println("route ok   ", etcdnodeGet)
	//Routetablecmd("ip route replace "+localdockerip+" dev docker0  scope link table ", "8")
	for k, v := range etcdnodeGet {
		v := v.(map[string]interface{})
		dockerip := v["Dockerip"].(string)
		nodeip := v["Nodeip"].(string)
		if strings.EqualFold(k, localnodeip) {
			fmt.Println(dockerip, "DOCKERIP")
			fmt.Println("addtableok")
			Routetablecmd("ip route replace "+dockerip+" dev docker0  scope link table ", "8")

		} else {
			fmt.Println(dockerip, nodeip)
			Routetablecmd("ip route replace "+dockerip+" via "+nodeip+" table ", "8")
			fmt.Println(dockerip, "OKOKKKKK----------")

		}

	}
	RuleAdd("8")
}

func NodeOspfIpGet() (iproutemap map[string]map[string]string, e error) {
	iproutemap = make(map[string]map[string]string)

	cmd := exec.Command("/bin/sh", "-c", "ip route show table 8")

	out, err := cmd.CombinedOutput()
	if err != nil {
		e := errors.New("cmd  Error")
		log.Log("/bin/sh -c ip route show table 8", "NodeOspfIpGet")
		panic("/bin/sh -c ip route show table 8 NodeOspfIpGet")
		return iproutemap, e
	}

	abc := string(out)
	a1 := strings.Replace(abc, " ", ",", -1)
	a2 := strings.Split(a1, "\n")
	//fmt.Println(a2)
	for _, v := range a2 {
		b1 := strings.Split(v, ",")
		if len(b1) > 1 {
			test := map[string]string{
				"Nodeip":   b1[2],
				"Dockerip": b1[0],
			}
			iproutemap[b1[2]] = test
			//fmt.Println(b1, len(b1))

		}

	}

	e = nil
	return iproutemap, e

}

func Routetablecmd(routecmd, tableid string) /*error */ {
	cmd := exec.Command("/bin/sh", "-c", routecmd+tableid)
	_, err := cmd.CombinedOutput()
	if err != nil {
		//return errors.New("cmd  Error1")
		//fmt.Println(err)
		log.Log("cmd  exe  error Routetablecmd /bin/sh -c", routecmd+tableid)
		panic("cmd  exe  error Routetablecmd /bin/sh -c" + routecmd + tableid)
	}

}

func RuleAdd(id string) {
	cmd := exec.Command("/bin/sh", "-c", "ip rule show |awk -F ':' '{print $1}'")
	out, err := cmd.CombinedOutput()
	if err != nil {

		fmt.Println(err)

	}
	//fmt.Println(string(out))
	abc := string(out)
	a1 := strings.Fields(abc)
	//fmt.Println(a1)
	cunzai := "0"
	for _, v := range a1 {
		if v == id {
			cunzai = "1"
		}
	}
	if cunzai == "0" {
		Routetablecmd("ip rule add from all pref "+id+" table ", id)
	}
}

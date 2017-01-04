package autotable

import (
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/con"
	"github.com/byebye758/kubernetes-haproxy-autolb/backends/etcd3client"
	//"encoding/json"
	//"fmt"
	//"github.com/coreos/etcd/clientv3"
	"strings"
	"unsafe"
)

// a := etcd3client.Autotable{endpoints, "/autohaproxy/autotable/"}

// 	b := etcd3client.Projecttable{endpoints, "/autohaproxy/project/"}

// 	c := etcd3client.Routetable{endpoints, "/autohaproxy/haproxytable/"}
// 	d := etcd3client.Pod{endpoints, "/registry/pods/default/"}
// 	e := etcd3client.Haip{endpoints, "/autohaproxy/haproxyip"}
func Autotable(a, b, c, e etcd3client.AGetr, d etcd3client.Getr, endpoints []string) {

	amap := make(map[string]interface{})
	for _, v := range b.AGet() { //循环项目列表
		v := v.(map[string]interface{})
		//fmt.Println(v)
		for _, v1 := range d.Get() { //循环  k8s pod
			//判断   pod中项目是否在 项目列表中
			if strings.EqualFold(v["Projectname"].(string), v1["Projectname"]) {
				//fmt.Println(v1["Podname"])
				//生成匹配   项目列表 的 pod列表  生成如下   map结构
				newmap := map[string]interface{}{
					"Podip":       v1["Podip"],
					"Nodeip":      v1["Nodeip"],
					"Podname":     v1["Podname"],
					"Projectname": v1["Projectname"],
					"Haproxyip":   v["Haproxyip"],
					"Projectport": v["Projectport"],
				}
				amap[v1["Podname"]] = newmap //匹配项目名称的  根据pod 生成map

			} else { //没有匹配到项目暂时不做处理

			}

		}
	}
	//fmt.Println(d.Get())
	//fmt.Println(amap)

	for k2, v2 := range amap { //循环 pod匹配项目列表后的 map     amap

		v2 := v2.(map[string]interface{})
		if _, ok := a.AGet()[k2]; ok { //判断是否在  autotable 列表中   在的话不做处理
			//fmt.Println(amap[k2])

		} else {
			//

			z1 := con.Haipnum(a, b, e, v2["Projectname"].(string))
			// w1 := v2["Haproxyip"].([]interface{})
			// u := make([]string, 0)

			// for _, v3 := range w1 {
			// 	//fmt.Println(v3.(string))
			// 	//判断  projecttable 是否存在于  haproxyip
			// 	if _, ok := e.AGet()[v3.(string)]; ok {
			// 		u = append(u, v3.(string))
			// 	}
			// }

			// //fmt.Println(v3.(string))

			// //fmt.Println(u)
			y1, err := con.Mix(z1)
			var y2 string
			if err == nil {
				y2 = y1
			} else {
				break
			}
			//fmt.Println(y1)

			v2["Haproxyip"] = y2
			p := c.AGet()
			for k4, v4 := range p { // 添加匹配  haproxyip 的haproxytable

				v4 := v4.(map[string]interface{})
				if strings.EqualFold(k4, v2["Haproxyip"].(string)) {
					v2["Haproxytable"] = v4["Haproxytable"]
				}

			}

			//fmt.Println(v2["Haproxyip"].(string))
			etcd3client.Autotableput(endpoints, v2["Haproxyip"].(string), v2)

		}

	}
	//fmt.Println(e.AGet())
	Autotabledelete(d, a, endpoints)

}
func B2S(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}
func Autotabledelete(a etcd3client.Getr, b etcd3client.AGetr, endpoints []string) { //a  为 k8s pod列表    b  为 autotable
	m := make(map[string]string)
	for _, v := range a.Get() {
		m[v["Podname"]] = v["Podip"]
	}
	//fmt.Println(m)

	for k, _ := range b.AGet() {
		if _, ok := m[k]; ok {
		} else {
			//删除   autotable中无用pod

			d := etcd3client.Pod{
				Endpoints: endpoints,
				Key:       "/autohaproxy/autotable/" + k,
			}
			d.Delete()
		}

	}

}

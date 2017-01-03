package con

import (
	"autohalb/backends/etcd3client"
	"strings"
)

//a 为  Autotable   x为Projecttable   haproip  为   haproxy注册进去的 ip地址  b 为 Projectname
func Haipnum(a, x, haproip etcd3client.AGetr, b string) map[string]int {
	back := make(map[string]int)   // back  记录    haproxyip  与 pod数量的 map
	f := make(map[string][]string) //生成匹配项目的 autotable
	AB := a.AGet()                 //获取 autotable  列表
	for _, v := range AB {
		c := v.(map[string]interface{})
		if strings.EqualFold(c["Projectname"].(string), b) { // 判断autotable  匹配主程序项目名称项目
			d := c["Haproxyip"].(string)
			e := c["Podname"].(string)
			f[d] = append(f[d], e)
		}

	}
	//fmt.Println(f)

	for _, v1 := range x.AGet() {
		c1 := v1.(map[string]interface{})
		if strings.EqualFold(c1["Projectname"].(string), b) {
			//fmt.Println(v1)
			d1 := c1["Haproxyip"].([]interface{})
			for _, v2 := range d1 {
				//fmt.Println(v2.(string))
				if _, ok := f[v2.(string)]; ok {
					//fmt.Println(v2.(string))

				} else {
					back[v2.(string)] = 0
				}
			}

		}
	}

	for k3, v3 := range f {
		back[k3] = len(v3)
	}
	p := make(map[string]int)
	for k4, v4 := range back {
		if _, ok := haproip.AGet()[k4]; ok {
			p[k4] = v4
		}
	}

	return p //匹配  haproxyip 后  返回    key为haproxyip vale 为 上面代理pod数量
}

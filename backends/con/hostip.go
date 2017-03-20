package con

import (
	"fmt"
	"kubernetes-haproxy-autolb/backends/log"
	"net"
	"strings"
)

/*获取主机IP地址*/
func HostIP() (a string) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Log("Hostip error", "HostIP")
		panic("Hostip error HostIP")
		return
	}
	defer conn.Close()

	a = strings.Split(conn.LocalAddr().String(), ":")[0]
	return
}

func Hosipnetwork() (network map[string]map[string]string) {
	network = make(map[string]map[string]string)
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {

		add, _ := inter.Addrs()
		if len(add) == 2 {
			ip := add[1].String()
			devname := inter.Name
			mac := inter.HardwareAddr.String() //获取本机MAC地址
			fmt.Println(mac, ip, devname)

			network[mac] = map[string]string{"ip": ip, "devname": devname}

		}
	}
	return network
}

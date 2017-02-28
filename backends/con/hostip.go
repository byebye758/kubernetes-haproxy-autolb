package con

import (
	//"fmt"
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

package con

import (
	"fmt"
	"net"
	"strings"
)

func HostIP() (a string) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	a = strings.Split(conn.LocalAddr().String(), ":")[0]
	return
}

package con

import (
	"errors"
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

func Hosipnetwork() (network map[string]map[string]string) {
	interfaces, err := net.Interfaces()
	network = make(map[string]map[string]string)
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		add, _ := inter.Addrs()
		//fmt.Println(len(add), add)

		if len(add) == 2 {
			//fmt.Println(add, inter.Name)

			if !(strings.EqualFold(inter.Name, "lo") || strings.EqualFold(inter.Name, "docker0")) {

				ip := add[0].String()
				ip, err = formatip(ip)
				if err != nil {
					continue
				}
				devname := inter.Name
				mac := inter.HardwareAddr.String() //获取本机MAC地址
				//fmt.Println(mac, ip, devname)

				network[mac] = map[string]string{"ip": ip, "devname": devname}

			}
		}
	}
	return network

}

func formatip(ip string) (fip string, err error) {
	a := strings.Split(ip, ".")
	if len(a) == 4 {
		b := strings.Split(a[3], "/")
		//fmt.Println(b[1])
		fip = a[0] + "." + a[1] + "." + a[2] + "." + "0/" + b[1]
		//fmt.Println(fip)
		err = nil
		return fip, err
	}

	err = errors.New("this is a new error")
	return fip, err
}

//Linux  中根据网卡名称取ip
func Devip(dev string) (ip *string) {
	interfaces, err := net.Interfaces()
	//network = make(map[string]map[string]string)
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		add, _ := inter.Addrs()

		//fmt.Println(add, inter.Name)

		if strings.EqualFold(inter.Name, dev) {
			ad := add[0].String()
			a := strings.Split(ad, "/")
			ip = &a[0]

		} else {
			log.Log("get  dev ip   error,no matching equipment", "Devip")
			panic("get  dev ip   error,no matching equipment")
		}

	}
	return ip
}

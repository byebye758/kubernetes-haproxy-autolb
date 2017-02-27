package cmd

import (
	//"autohalb/backends/con"
	//"autohalb/backends/etcd3client"
	"errors"
	"fmt"
	//"kubernetes-haproxy-autolb/backends/log"
	"os/exec"
	"strings"
)

func NodeOspfIpGet() (iproutemap map[string]map[string]string, e error) {
	iproutemap = make(map[string]map[string]string)

	cmd := exec.Command("/bin/sh", "-c", "ip route show table 8")

	out, err := cmd.CombinedOutput()
	if err != nil {
		e := errors.New("cmd  Error")
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
		fmt.Println(err)

	}
	//return nil
}
func NoderuleGet() (ipruleslice []map[string]string, e error) {
	ipruleslice = make([]map[string]string, 0)

	excludeid := map[string]string{"0": "0", "1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7", "8": "8", "9": "9", "32766": "32766", "32767": "32767"}
	cmd := exec.Command("/bin/sh", "-c", "ip rule show")

	out, err := cmd.CombinedOutput()
	if err != nil {
		e = errors.New("cmd  Error")
		return ipruleslice, e
	}

	abc := string(out)
	a1 := strings.Replace(abc, ":", " ", -1)
	//a1 = strings.Replace(a1, " ", ",", -1)
	a2 := strings.Split(a1, "\n")
	//fmt.Println(a2)
	for _, v := range a2 {
		//fmt.Println(v)

		b1 := strings.Split(v, " ")
		if len(b1) > 1 {

			ruleid := b1[0]
			if _, ok := excludeid[ruleid]; ok {

			} else {
				test := map[string]string{
					"Nodeip": b1[2],
					"Ruleid": b1[0],
				}

				ipruleslice = append(ipruleslice, test)
			}
		}
	}

	e = nil
	return ipruleslice, e

}

package con

import (
	"errors"
	"os/exec"
	"strings"
)

/*获取node端 docker ip地址*/
func Getdockerip() (ip string, e error) {
	cmd := exec.Command("/bin/sh", "-c", "ip addr | grep  docker0$ | awk '{print $2}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		e := errors.New("cmd  Error")
		return ip, e
	} else {
		panic("get docker ip  not ok", "Getdockerip")
	}
	abc := string(out)

	a1 := strings.Fields(abc)

	e = nil

	a2 := strings.SplitN(a1[0], "/", 2)
	b := a2[0]
	c := strings.SplitAfterN(b, ".", 4)
	ip = c[0] + c[1] + c[2] + "0" + "/" + a2[1]
	return ip, e

}

package con

import (
	"errors"
	"sort"
)

// func Mix(newslice []string, newmap map[string]int) (ip string) {

// 	a := newmap[newslice[0]]
// 	b := newslice[0]
// 	for k, v := range newmap {
// 		if a >= v {
// 			a = v
// 			b = k
// 		}

// 	}
// 	//fmt.Println(a, b, a)
// 	return b
// }
func Mix(newmap map[string]int) (ip string, err error) { //按照  map的 value 进行排序 返回  最小k的值
	var p *string
	if len(newmap) == 0 {
		err = errors.New("error")
		ip = "ip"
		return
	}
	values := make([]int, len(newmap))

	i := 0
	for _, v := range newmap {
		values[i] = v
		i++
	}

	sort.Ints(values)
	for k, v := range newmap {
		if values[0] == v {
			p = &k
			break
		}

	}
	ip = *p
	err = nil
	return
}

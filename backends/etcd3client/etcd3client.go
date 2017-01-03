package etcd3client

import (
	//"autohalb/backends/autotable"
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
	"unsafe"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 2 * time.Second
)

type NodePutLeaser interface {
	NodePutLease()
}

type AutotableDleter interface {
	Delete()
}
type Getr interface {
	Get() []map[string]string
}
type HaPut interface {
	Put() (ch <-chan string)
}
type AGetr interface { //提供autotable  提供projecttable 、Routetable  get数据
	AGet() map[string]interface{}
}
type Pod struct {
	Endpoints []string
	Key       string
}
type Node struct {
	Endpoints []string
	Key       string
}
type Haip struct {
	Endpoints []string
	Key       string
}

type Autotable struct {
	Endpoints []string
	Key       string
}
type Projecttable struct {
	Endpoints []string
	Key       string
}
type Routetable struct {
	Endpoints []string
	Key       string
}
type Register struct {
	Endpoints []string
	Key       string
	Value     string
	Ch        chan string
}
type NodeRegister struct {
	Endpoints []string
	Key       string
	Nodeip    string
	Dockerip  string
	Ch        chan string
}

func (e Pod) Get() (AA []map[string]string) {

	// var (
	// 	dialTimeout    = 5 * time.Second
	// 	requestTimeout = 2 * time.Second
	// )

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, e.Key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	//AA := make([]map[string]string, 0)

	for _, ev := range resp.Kvs {

		var I map[string]interface{}
		if err := json.Unmarshal([]byte(ev.Value), &I); err == nil {
			x := I["status"]
			x1 := I["metadata"]

			//fmt.Println(x)
			y := x.(map[string]interface{})
			//fmt.Println(y)
			y1 := x1.(map[string]interface{})
			//fmt.Println(y1["labels"])
			//break

			a := y["podIP"]
			var a1 string
			if v, ok := a.(string); ok {
				a1 = v
			} else {
				continue
			}

			//fmt.Println(a1)
			b := y["hostIP"]

			//fmt.Println(b)
			c := y["containerStatuses"]
			d := y1["name"]
			//fmt.Println(c)

			if v1, ok := c.([]interface{}); ok {
				//fmt.Println(v1[0])
				if v2, ok := v1[0].(map[string]interface{}); ok {

					if f, ok := v2["name"].(string); ok {
						//fmt.Println(f)
						test := map[string]string{
							"Podip":       a1,
							"Nodeip":      b.(string),
							"Podname":     d.(string),
							"Projectname": f,
						}
						AA = append(AA, test)
						//fmt.Println(test)

					}
				}
			}
		}
		//break
	}
	//fmt.Println(AA)
	return AA

}

func (p Pod) Delete() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   p.Endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// count keys about to be deleted
	gres, err := cli.Delete(ctx, p.Key)
	if err != nil {
		log.Fatal(err)
		fmt.Println(gres)
	}

}

// func (t Routetable) Get() (AA []map[string]string) {
// 	// var (
// 	// 	dialTimeout    = 5 * time.Second
// 	// 	requestTimeout = 2 * time.Second
// 	// )

// 	cli, err := clientv3.New(clientv3.Config{
// 		Endpoints:   t.Endpoints,
// 		DialTimeout: dialTimeout,
// 	})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cli.Close() // make sure to close the client
// 	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
// 	resp, err := cli.Get(ctx, t.Key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
// 	cancel()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, ev := range resp.Kvs {

// 		test := map[string]string{
// 			B2S(ev.Key): B2S(ev.Value),
// 		}
// 		AA = append(AA, test)

// 	}
// 	return AA
// }
func (h Register) PutLease() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   h.Endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	//fmt.Println("aaa")
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Grant create")
	type Ha struct {
		Haproxyip string
		Lease     clientv3.LeaseID
	}
	group := Ha{
		Haproxyip: h.Value,
		Lease:     resp.ID,
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}

	_, err = cli.Put(context.TODO(), h.Key, B2S(b), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// the key 'foo' will be kept forever
	ch1, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	for {

		ka := <-ch1
		fmt.Println("ttl:", ka.TTL)

	}

	h.Ch <- "error"

}

func (n NodeRegister) NodePutLease() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   n.Endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	//fmt.Println("aaa")
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Grant create")
	type Node struct {
		Dockerip string
		Nodeip   string
		Lease    clientv3.LeaseID
	}
	group := Node{
		Dockerip: n.Dockerip,
		Nodeip:   n.Nodeip,
		Lease:    resp.ID,
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}

	_, err = cli.Put(context.TODO(), n.Key, B2S(b), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// the key 'foo' will be kept forever
	ch1, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	for {

		ka := <-ch1
		fmt.Println("ttl:", ka.TTL)

	}

	n.Ch <- "error"

}
func (h Haip) AGet() map[string]interface{} {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   h.Endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, h.Key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	M := make(map[string]interface{})

	for _, ev := range resp.Kvs {
		// fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		var I map[string]interface{}
		if err := json.Unmarshal([]byte(ev.Value), &I); err == nil {
			b := I["Haproxyip"].(string)
			M[b] = I
			//fmt.Println(I["Podname"].(string))
		}
	}
	return M
}
func (n Node) AGet() map[string]interface{} {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   n.Endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, n.Key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	M := make(map[string]interface{})

	for _, ev := range resp.Kvs {
		// fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		var I map[string]interface{}
		if err := json.Unmarshal([]byte(ev.Value), &I); err == nil {
			b := I["Nodeip"].(string)
			M[b] = I
			//fmt.Println(I["Podname"].(string))
		}
	}
	return M
}

func (a Autotable) AGet() map[string]interface{} {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   a.Endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, a.Key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	M := make(map[string]interface{})

	for _, ev := range resp.Kvs {
		// fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		var I map[string]interface{}
		if err := json.Unmarshal([]byte(ev.Value), &I); err == nil {
			b := I["Podname"].(string)
			M[b] = I
			//fmt.Println(I["Podname"].(string))
		}
	}
	return M
}

func (p Projecttable) AGet() map[string]interface{} {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   p.Endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, p.Key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	M := make(map[string]interface{})
	for _, ev := range resp.Kvs {
		// fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		var I map[string]interface{}

		if err := json.Unmarshal([]byte(ev.Value), &I); err == nil {
			b := I["Projectname"].(string)
			M[b] = I
			//fmt.Println(I["Podname"].(string))
		}
	}
	return M
}

func (r Routetable) AGet() map[string]interface{} {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   r.Endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, r.Key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	M := make(map[string]interface{})
	for _, ev := range resp.Kvs {
		// fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		var I map[string]interface{}

		if err := json.Unmarshal([]byte(ev.Value), &I); err == nil {
			b := I["Haproxyip"].(string)
			M[b] = I
			//fmt.Println(I["Podname"].(string))
		}
	}
	return M
}
func B2S(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

//value  为  生成需要的autotable
func Autotableput(endpoints []string, haip string, value map[string]interface{}) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "/autohaproxy/haproxyip/"+haip)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	type Ha struct {
		Haproxyip string
		Lease     clientv3.LeaseID
		//Lease int64
	}
	var st Ha
	for _, ev := range resp.Kvs {

		err = json.Unmarshal([]byte(ev.Value), &st)
		if err != nil {
			fmt.Println("Unmarshal faild")
		}

	}
	//fmt.Println(st.Lease, value["Projectname"].(string), value["Projectport"].(map[string]interface{}), value["Haproxyip"].(string), value["Haproxytable"].(string), value["Podname"].(string), value["Podip"].(string), value["Nodeip"].(string))

	type autotable struct {
		Projectname  string
		Projectport  map[string]interface{}
		Haproxyip    string
		Haproxytable string
		Podname      string
		Podip        string
		Nodeip       string
		Lease        clientv3.LeaseID
	}
	group := autotable{
		Projectname:  value["Projectname"].(string),
		Projectport:  value["Projectport"].(map[string]interface{}),
		Haproxyip:    value["Haproxyip"].(string),
		Haproxytable: value["Haproxytable"].(string),
		Podname:      value["Podname"].(string),
		Podip:        value["Podip"].(string),
		Nodeip:       value["Nodeip"].(string),
		Lease:        st.Lease,
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	_, err = cli.Put(context.TODO(), "/autohaproxy/autotable/"+value["Podname"].(string), B2S(b), clientv3.WithLease(st.Lease))
	if err != nil {
		log.Fatal(err)
	}

}

// func Watch(key string, endpoints []string, a, b, c, e AGetr, d Getr, add func(AGetr, AGetr, AGetr, AGetr, Getr, []string), ch chan string) {
// 	cli, err := clientv3.New(clientv3.Config{
// 		Endpoints:   endpoints,
// 		DialTimeout: dialTimeout,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(key, endpoints)
// 	rch := cli.Watch(context.Background(), "key", clientv3.WithPrefix())
// 	//fmt.Println(<-rch)

// 	for wresp := range rch {
// 		fmt.Println(456)
// 		for _, ev := range wresp.Events {
// 			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
// 			add(a, b, c, e, d, endpoints)
// 			ch <- key
// 		}
// 	}

// }

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type pid = int

type dispatch = map[string]interface{}

type nodeSpec struct {
	nodeID       int
	nodeLocators map[string]string
}

type pSpec struct {
	pid  pid
	node nodeSpec
}

type msg struct {
	pid  pid
	name string
	args []interface{}
}

// https://jhall.io/posts/go-json-tricks-array-as-structs/
func (r *msg) UnmarshalJSON(p []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	r.pid = tmp[0].(pid)
	r.name = tmp[1].(string)
	r.args = tmp[2:]
	return nil
}

type receiver = func() msg

var pidCounter int = 1

var ps = make(map[pid]chan msg)

var myNode = nodeSpec{nodeID: 4, nodeLocators: make(map[string]string)}

func request_handler(w http.ResponseWriter, r *http.Request) {
	var m msg

	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	send(m.pid, m.name, m.args...)
}

func createNode(f func()) {
	myNode.nodeLocators["http"] = "http://0.0.0.0:8090"
	go f()
	http.HandleFunc("/", request_handler)
	http.ListenAndServe(":8090", nil)
}

func spawnM(m map[string]interface{}) pid {
	fmt.Println("map ty", reflect.TypeOf(m))
	return spawn(func(receive receiver) {
		for true {
			ms := receive()
			f := m[ms.name].(func(...interface{}))
			fmt.Println("spawnm msg got type", reflect.TypeOf(f))
			f(ms.args...)
		}
	})
}

func spawn(f func(receiver)) pid {

	pid := pidCounter
	pidCounter++

	ch := make(chan msg)

	ps[pid] = ch

	r := func() msg {
		return <-ch
	}

	go f(r)

	return pid
}

func send(p pid, m string, as ...interface{}) {
	ch := ps[p]
	ch <- msg{pid: p, name: m, args: as}
}

func af(receive receiver) {
	for true {
		m := receive()
		switch m.name {
		case "foo":
			fmt.Println("got foo")
		case "hello":
			fmt.Println("got hello", m.args)
		}
	}
}

func prType(x interface{}) {

	fmt.Println("got type", reflect.TypeOf(x))

}

func main() {

	prType(func(x int) int { return x + 1 })

	createNode(func() {

		mpid:=spawnM(dispatch{
			"foo": func(x int) {
				fmt.Println("FOO GOT", x)
				return
			},
		})
		send(mpid, "foo", 51)

		spid := spawn(af)
		fmt.Println("x type", spid)
		send(spid, "foo")
		send(spid, "hello", 42)
	})
}

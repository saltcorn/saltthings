package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type pid = string

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

//https://stackoverflow.com/a/31832326/19839414
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
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

type receiver = func() (string,  []interface{})

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
	rand.Seed(time.Now().UnixNano())
	myNode.nodeLocators["http"] = "http://0.0.0.0:8090"
	go f()
	http.HandleFunc("/", request_handler)
	http.ListenAndServe(":8090", nil)
}


func spawn(f func(receiver)) pid {

	pid := RandStringBytes(16)

	ch := make(chan msg)

	ps[pid] = ch

	r := func() (string, []interface{}) {
		m:=<-ch
		return m.name, m.args
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
		name, args := receive()
		switch name {
		case "foo":
			fmt.Println("got foo")
		case "hello":
			fmt.Println("got hello", args)
		}
	}
}


func main() {

	createNode(func() {

		/*mpid:=spawnM(dispatch{
			"foo": func(x int) {
				fmt.Println("FOO GOT", x)
				return
			},
		})
		send(mpid, "foo", 51)*/

		spid := spawn(af)
		fmt.Println("x type", spid)
		send(spid, "foo")
		send(spid, "hello", 42)
	})
}

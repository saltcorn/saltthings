package main

import (
	"bytes"
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
	nodeID       string
	nodeLocators map[string]string
}

type pSpec struct {
	pid  pid
	node nodeSpec
}

type msg struct {
	proc pSpec
	name string
	args []interface{}
}

// https://stackoverflow.com/a/31832326/19839414
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
	r.proc = tmp[0].(pSpec)
	r.name = tmp[1].(string)
	r.args = tmp[2:]
	return nil
}

type receiver = func() (string, []interface{})

var ps = make(map[pid]chan msg)

var myNode = nodeSpec{nodeID: "", nodeLocators: make(map[string]string)}

func request_handler(w http.ResponseWriter, r *http.Request) {
	var m msg

	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	send(m.proc, m.name, m.args...)
}

func createNode(f func()) {
	rand.Seed(time.Now().UnixNano())

	myNode.nodeID = RandStringBytes(16)
	myNode.nodeLocators["http"] = "http://0.0.0.0:8090"

	go f()
	http.HandleFunc("/", request_handler)
	http.ListenAndServe(":8090", nil)
}

func spawn(f func(receiver, pSpec)) pSpec {

	pid := RandStringBytes(16)

	ch := make(chan msg)

	ps[pid] = ch

	p := pSpec{pid: pid, node: myNode}

	r := func() (string, []interface{}) {
		m := <-ch
		return m.name, m.args
	}

	go f(r, p)

	return p
}

func send(p pSpec, m string, as ...interface{}) {
	if p.node.nodeID != myNode.nodeID {
		msgArray := [](interface{}){p, m}
		msgArray = append(msgArray, as...)
		postBody, _ := json.Marshal(msgArray)
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(
			p.node.nodeLocators["http"],
			"application/json",
			responseBody)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	} else {
		ch := ps[p.pid]
		ch <- msg{proc: p, name: m, args: as}
	}
}

func af(receive receiver, self pSpec) {
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

		spid := spawn(af)
		fmt.Println("x type", spid)
		send(spid, "foo")
		send(spid, "hello", 42)
	})
}

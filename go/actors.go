package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type pid = int

type nodeSpec struct {
	nodeID int
	nodeLocators map[string]string
}


type pSpec struct {
	pid pid
	node nodeSpec
}

type msg struct {
	pid pid
	name string	
	args []interface{}
}

//https://jhall.io/posts/go-json-tricks-array-as-structs/
func (r *msg) UnmarshalJSON(p []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(p, &tmp); err != nil {
			return err
	}
	r.pid = tmp[0].(pid)
	r.name = tmp[1].(string)
	r.args = tmp[2].([]interface{})
	return nil
}

type receiver = func() msg

var pidCounter int =1

var ps= make(map[pid]chan msg)

func request_handler(w http.ResponseWriter, r *http.Request) {
	var m msg

	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(b, &m); err != nil {
			panic(err)
	}
	send(m.pid,m.name, m.args...)
}

func createNode(f func()) {
	go f()
	http.HandleFunc("/", request_handler)
	http.ListenAndServe(":8090", nil)
}

func spawn(f func(receiver)) pid {

	pid:=pidCounter
	pidCounter++

	ch := make(chan msg)

	ps[pid] = ch

	r := func () msg {
		return <- ch
	}

	go f(r)

	return pid
}

func send(p pid, m string, as ...interface{}) {
	ch := ps[p]
	ch <- msg{pid:p, name:m, args: as}
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

func main() {
	createNode(func() {
		spid := spawn(af)
		fmt.Println("got PID", spid)
		send(spid, "foo")
		send(spid, "hello", 42)
	})
}
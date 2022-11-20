package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type pid = int

type msg struct {
	pid pid
	name string	
	args []interface{}
}

type receiver = func() msg

var pidCounter int =1

var ps= make(map[pid]chan msg)

func request_handler(w http.ResponseWriter, r *http.Request) {
	var req msg
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
			panic(err)
	}
	send(req.pid,req.name, req.args...)
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
	spid := spawn(af)
	fmt.Println("got PID", spid)
	send(spid, "foo")
	send(spid, "hello", 42)
	http.HandleFunc("/", request_handler)
	http.ListenAndServe(":8090", nil)
}
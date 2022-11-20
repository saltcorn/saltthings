package main

import "fmt"

type pid = int

type msg struct {
	pid pid
	name string	
}

type procStore struct {
	pid pid
	ch chan msg
}
type receiver = func() msg

var pidCounter int =1
var procs []procStore

func spawn(f func(receiver)) pid {

	pid:=pidCounter
	pidCounter++

	ch := make(chan msg)

	procs = append(procs, procStore{pid: pid, ch: ch})

	r := func () msg {
		return <- ch
	}

	go f(r)

	return pid
}

func af(receive receiver) {
	m := receive()
	switch m.name {
	case "foo":
			fmt.Println("got foo")
	}

}


func main() {
	spid := spawn(af)
	fmt.Println("got PID", spid)
	select {} // block forever
}
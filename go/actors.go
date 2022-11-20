package main

import "fmt"

type pid = int

type msg struct {
	pid pid
	name string	
	args []interface{}
}

type receiver = func() msg

var pidCounter int =1

var ps= make(map[pid]chan msg)


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
	m := receive()
	switch m.name {
	case "foo":
			fmt.Println("got foo")
	case "hello":
			fmt.Println("got hello", m.args)
	}
}

func main() {
	spid := spawn(af)
	fmt.Println("got PID", spid)
	//send(spid, "foo")
	send(spid, "hello", 42)
	select {} // block forever
}
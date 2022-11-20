package main

import "fmt"

type pid = int

type msg struct {
	pid pid
	name string	
}

type receiver = func() msg

var pidCounter int =1

func spawn(f func(receiver)) pid {

	pid:=pidCounter
	pidCounter++
	
	r := func () msg {
		return msg{pid:1, name:"foo"}
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
}
package main

import (
	"fmt"
	"reflect"
)


func spawnM(m map[string]interface{}) pSpec {
	fmt.Println("map ty", reflect.TypeOf(m))
	return spawn(func(receive receiver, self pSpec) {
		for true {
			name, args := receive()
			f := m[name].(func(...interface{}))
			fmt.Println("spawnm msg got type", reflect.TypeOf(f))
			f(args...)
		}
	})
}


func prType(x interface{}) {

	fmt.Println("got type", reflect.TypeOf(x))

}

		/*mpid:=spawnM(dispatch{
			"foo": func(x int) {
				fmt.Println("FOO GOT", x)
				return
			},
		})
		send(mpid, "foo", 51)*/
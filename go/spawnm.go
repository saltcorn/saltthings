package main

import (
	"fmt"
	"reflect"
)


func spawnM(m map[string]interface{}) pid {
	fmt.Println("map ty", reflect.TypeOf(m))
	return spawn(func(receive receiver) {
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
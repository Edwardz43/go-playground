package main

import (
	"fmt"
)

type v int

func main() {
	var a = []interface{}{1, 2, 3}
	fmt.Println(a)
	fmt.Println(a...)

	x := [3]int{1, 2, 3}

	func(arr [3]int) {
		arr[0] = 7
		fmt.Println(arr)
	}(x)

	fmt.Println(x)

	m := map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
	}

	for k, v := range m {
		println(k, v)
	}
}

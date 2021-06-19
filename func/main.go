package main

// 具名函数
func Add(a, b int) int {
	return a + b
}

// 匿名函数
var Add2 = func(a, b int) int {
	return a + b
}

func Inc() (v int) {
	defer func() {
		v++
		println(v)
	}()
	return 42
}

func main() {
	// print(Add(1, 2))
	// print(Add2(1, 2))
	println(Inc())
}

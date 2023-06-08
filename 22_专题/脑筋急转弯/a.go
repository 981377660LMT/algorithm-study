package main

func main() {
	foo([]string{"a", "b", "c"})
	foo([]int{1, 2, 3})
	foo(1)
}

func foo(arg interface{}) {

}

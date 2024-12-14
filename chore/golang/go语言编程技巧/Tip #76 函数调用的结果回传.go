// !result、err 作为函数的参数.(回调函数风格的处理函数)

package main

func main() {
	var doSomething func() (int, error)
	var processResult func(int, error)

	result, err := doSomething()
	processResult(result, err)
}

// 如果没有error的时候它会返回结果，否则就会停止运行.
func Must[T any](result T, err error) T {
	panic("todo")
}

// 通过简单的'defer'关键字，你可以借助一个小技巧实现在另一个函数的开头和结尾处执行一个函数.
// 返回一个cleanup有点useEffect的感觉了.

package main

import "fmt"

func main() {
	defer MultistageDefer()()
	fmt.Println("Main")
}

func MultistageDefer() func() {
	fmt.Println("Init")
	return func() {
		fmt.Println("Clean up")
	}
}

// Output:
// Init
// Main
// Clean up

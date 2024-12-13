// !Go语言的设计哲学是专注于最常见的需求，而非极端特例。
// !go 的 case 自带一个 break.
// 除非使用 fallthrough 关键字，否则不会继续执行下一个 case 语句。
// 如果我们在循环内使用 switch 并且想要完全跳出循环，我们可以使用 break 标签(label)。

package main

import "math/rand"

func main() {
	testSwitch()
}

func testSwitch() {
	n := rand.Intn(10)
	switch n {
	case 0:
		println("0")
	case 1, 2:
		println("1")
	case 3:
		println("2")
		fallthrough
	case 4:
		println("3")
	default:
		println("4")
	}
}

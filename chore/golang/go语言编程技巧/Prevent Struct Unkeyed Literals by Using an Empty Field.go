// https://colobu.com/gotips/055.html
// 限制使用命名参数

package main

func main() {

	a := Point{X: 1, Y: 2}
	// b := Point{3, 4}
	c := Point2{X: 3, Y: 4}
}

type Point struct {
	_ struct{}
	X float64
	Y float64
}

type noUnkeyed = struct{}
type Point2 struct {
	noUnkeyed
	X float64
	Y float64
}

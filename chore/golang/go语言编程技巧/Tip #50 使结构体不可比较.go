package main

func main() {
	a, b := Point{X: 1, Y: 2}, Point{X: 3, Y: 4}
	_ = a == b
}

type Point struct {
	_    [0]func()
	X, Y float64
}

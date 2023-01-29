package main

import "fmt"

func main() {
	type Point struct {
		x, y int
	}

	points := []Point{
		{1, 2},
		{2, 3},
		{3, 4},
	}

	copu := make([]Point, len(points))
	fmt.Println(copu)

	copy(copu, points)
	fmt.Println(copu)
}

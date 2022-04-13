package main

import "fmt"

// golang暴力 1176 ms
func circleGame(toys [][]int, circles [][]int, r int) (res int) {
	for _, toy := range toys {
		if toy[2] > r {
			continue
		}

		for _, circle := range circles {
			if (circle[0]-toy[0])*(circle[0]-toy[0])+(circle[1]-toy[1])*(circle[1]-toy[1]) <= (r-toy[2])*(r-toy[2]) {
				res++
				break
			}
		}
	}

	return res
}

func main() {
	fmt.Println(circleGame([][]int{{0, 0, 1}, {1, 0, 2}, {2, 0, 3}, {0, 1, 1}, {1, 1, 2}, {2, 1, 3}, {0, 2, 1}, {1, 2, 2}, {2, 2, 3}}, [][]int{{-1, -1, 0}, {10, 10, 11}}, 1))
}

package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	n := rand.Intn(10000) + 2
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		points[i] = [2]int{rand.Intn(10000) - 50, rand.Intn(10000) - 50}
	}

	dist2 := func(p, q [2]int) int { return (p[0]-q[0])*(p[0]-q[0]) + (p[1]-q[1])*(p[1]-q[1]) }
	bestI, bestJ := ClosestPair(points)
	best := dist2(points[bestI], points[bestJ])

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if dist2(points[i], points[j]) < best {
				panic("wrong answer")
			}
		}
	}
	fmt.Println("ok")
}

// 平面最近点对，返回两个点的下标.
//  len(points) >= 2
func ClosestPair(points [][2]int) (int, int) {
	mp := map[[2]int]int{}
	n := len(points)
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	rand.Shuffle(n, func(i, j int) { order[i], order[j] = order[j], order[i] })
	newPoints := make([][2]int, n)
	for i, idx := range order {
		newPoints[i] = points[idx]
	}
	points = newPoints

	calc := func(i, j int) int {
		xi, yi := points[i][0], points[i][1]
		xj, yj := points[j][0], points[j][1]
		return (xj-xi)*(xj-xi) + (yj-yi)*(yj-yi)
	}

	best := calc(0, 1)
	res := [2]int{0, 1}
	w := int(math.Sqrt(float64(best)))
	nxt := make([]int, n)
	for i := range nxt {
		nxt[i] = -1
	}

	insert := func(i int) {
		key := [2]int{int(points[i][0] / w), int(points[i][1] / w)}
		if j, ok := mp[key]; ok {
			nxt[i] = j
		}
		mp[key] = i
	}

	query := func(i int) bool {
		a := (points[i][0] / w)
		b := (points[i][1] / w)
		upd := false

		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				key := [2]int{a + dx, b + dy}
				j := -1
				if res, ok := mp[key]; ok {
					j = res
				}
				for j != -1 {
					cand := calc(i, j)
					if cand < best {
						best = cand
						res = [2]int{i, j}
						w = int(math.Sqrt(float64(best)))
						upd = true
					}
					j = nxt[j]
				}
			}
		}
		return upd
	}

	insert(0)
	insert(1)

	for i := 2; i < n; i++ {
		if query(i) {
			if best == 0 {
				break
			}
			mp = map[[2]int]int{}
			for j := 0; j < i; j++ {
				insert(j)
			}
		}
		insert(i)
	}

	return order[res[0]], order[res[1]]
}

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(RandomTree(5))
}

// RandomGraph 生成随机图的边.
func RandomGraph(n int, directed bool, simple bool) [][2]int {
	edges := make([][2]int, 0)
	cand := make([][2]int, 0)
	for a := 0; a < n; a++ {
		for b := 0; b < n; b++ {
			if simple && a == b {
				continue
			}
			if !directed && a > b {
				continue
			}
			cand = append(cand, [2]int{a, b})
		}
	}
	m := rand.Intn(len(cand) + 1)
	s := make(map[int]bool)
	for i := 0; i < m; i++ {
		for {
			i := rand.Intn(len(cand))
			if simple && s[i] {
				continue
			}
			s[i] = true
			a, b := cand[i][0], cand[i][1]
			edges = append(edges, [2]int{a, b})
			break
		}
	}
	randomRelabel(n, edges)
	return edges
}

// RandomTree 生成随机树的边.
func RandomTree(n int) [][2]int {
	edges := make([][2]int, 0)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{rand.Intn(i), i})
	}
	randomRelabel(n, edges)
	return edges
}

// randomRelabel 随机重标号.
func randomRelabel(n int, edges [][2]int) {
	rand.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	for i := range edges {
		edges[i][0], edges[i][1] = a[edges[i][0]], a[edges[i][1]]
	}
}

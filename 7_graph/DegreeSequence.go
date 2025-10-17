// 从度数构建图
// 判断度数序列是否合法的 Erdős–Gallai定理 (度数合法)
// 什么叫合法的度数序列？就是存在一个无重边无自环（简单）的无向图，使得图中每个点的度数构成的序列为给定的序列。
// http://kanari.logdown.com/posts/2014/03/09/erdos-gallai-theorem-conditions-for-a-sequence-to-be-graphical
// https://maspypy.github.io/library/graph/degree_sequence.hpp
// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/graph/is_graphic.cc#L6

package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(ConstructFromDegreeSequence([]int{2, 2, 2, 2}))
}

func simpleGraphExists(degrees []int) bool {
	return CheckDegreeSequence(degrees)
}

// 判断度数序列是否属于一个无重边无自环（简单）的无向图.
func CheckDegreeSequence(deg []int) bool {
	n := len(deg)
	if n == 0 {
		return true
	}
	if maxs(deg) >= n {
		return false
	}
	if sum(deg)%2 != 0 {
		return false
	}

	d := append([]int(nil), deg...)
	sort.Slice(d, func(i, j int) bool { return d[i] > d[j] })
	pref := make([]int, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + d[i]
	}

	for k := 1; k <= n; k++ {
		idx := sort.Search(n, func(i int) bool { return d[i] <= k })
		right := k * (k - 1)
		if idx > k {
			right += (idx-k)*k + (pref[n] - pref[idx])
		} else {
			right += pref[n] - pref[k]
		}
		if pref[k] > right {
			return false
		}
	}
	return true
}

// 从度数序列构建一个无重边无自环（简单）的无向图.
func ConstructFromDegreeSequence(deg []int) (edges [][2]int) {
	if !CheckDegreeSequence(deg) {
		return nil
	}
	n := len(deg)
	data := make([][]int, n)
	for v := 0; v < n; v++ {
		data[deg[v]] = append(data[deg[v]], v)
	}
	mx := n - 1
	for i := 0; i < n; i++ {
		for mx >= 0 && len(data[mx]) == 0 {
			mx--
		}
		v := data[mx][len(data[mx])-1]
		data[mx] = data[mx][:len(data[mx])-1]
		nbd := make([]int, 0)
		k := mx
		for len(nbd) < deg[v] {
			if k == 0 {
				return nil
			}
			if len(data[k]) == 0 {
				k--
				continue
			}
			x := data[k][len(data[k])-1]
			data[k] = data[k][:len(data[k])-1]
			nbd = append(nbd, x)
		}
		for _, x := range nbd {
			edges = append(edges, [2]int{v, x})
			deg[x]--
			data[deg[x]] = append(data[deg[x]], x)
		}
		deg[v] = 0
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(a []int) int {
	m := a[0]
	for _, v := range a {
		if v > m {
			m = v
		}
	}
	return m
}

func sum(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}

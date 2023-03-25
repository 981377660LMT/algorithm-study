// 从度数构建图
// 判断度数序列是否合法的 Erdős–Gallai定理
// 什么叫合法的度数序列？就是存在一个无重边无自环（简单）的无向图，使得图中每个点的度数构成的序列为给定的序列。
// http://kanari.logdown.com/posts/2014/03/09/erdos-gallai-theorem-conditions-for-a-sequence-to-be-graphical
// https://maspypy.github.io/library/graph/degree_sequence.hpp

package main

import "fmt"

func main() {
	fmt.Println(ConstructFromDegreeSequence([]int{1, 2, 1}))
}

// 判断度数序列是否属于一个无重边无自环（简单）的无向图.
func CheckDegreeSequence(deg []int) bool {
	n := len(deg)
	if maxs(deg) >= n {
		return false
	}
	if sum(deg)%2 != 0 {
		return false
	}
	counter := make([]int, n)
	for _, v := range deg {
		counter[v]++
	}
	p := 0
	for i := 0; i < n; i++ {
		for j := 0; j < counter[i]; j++ {
			deg[p] = i
			p++
		}
	}

	A := make([]int, n+1)
	B := make([]int, n+1)
	for i, d := range deg {
		A[i+1] += 2*i - d
		if d < i {
			B[0]++
			B[d]--
			A[d] += d
			A[i+1] -= d
		}
		if d >= i {
			B[0]++
			B[i+1]--
		}
	}

	for i := 1; i < n+1; i++ {
		A[i] += A[i-1]
		B[i] += B[i-1]
	}

	for k := 0; k < n+1; k++ {
		x := A[k] + B[k]*k
		if x < 0 {
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
		for len(data[mx]) == 0 {
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

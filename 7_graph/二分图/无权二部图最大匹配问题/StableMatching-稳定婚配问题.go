// StableMatching-稳定婚配问题(安定マッチング )
// O(n+m)logm

package main

import (
	"fmt"
	"sort"
)

func main() {
	boy, girl := int32(3), int32(3)
	S := NewStableMatching(boy, girl)

	// 索引代表男性编号（0，1，2），数组代表喜好女性的列表
	boys := [][]int32{{0, 1, 2}, {1, 0, 2}, {0, 1, 2}}
	// 索引代表女性编号（0，1，2），数组代表喜好男性的列表
	girls := [][]int32{{1, 2, 0}, {0, 1, 2}, {0, 1, 2}}

	for i := int32(0); i < boy; i++ {
		for j := int32(0); j < girl; j++ {
			S.Add(i, j, boys[i][j], girls[j][i])
		}
	}

	fmt.Println(S.Calc()) // [[0 1] [1 0] [2 2]]
}

const INF32 int32 = 1e9 + 10

type StableMatching struct {
	n1, n2                     int32
	data                       [][][3]int32
	match1, match2, val1, val2 []int32
}

func NewStableMatching(n1, n2 int32) *StableMatching {
	return &StableMatching{n1: n1, n2: n2, data: make([][][3]int32, n1)}
}

// 男生v1和女生v2的偏爱度分别为x1和x2.
// x：价值，大的优先
func (s *StableMatching) Add(v1, v2, x1, x2 int32) {
	s.data[v1] = append(s.data[v1], [3]int32{v2, x1, x2})
}

// 返回稳定婚配的结果(男生、女生).
func (s *StableMatching) Calc() [][2]int32 {
	for v1 := int32(0); v1 < s.n1; v1++ {
		vs := s.data[v1]
		sort.Slice(vs, func(i, j int) bool { return vs[i][1] < vs[j][1] })
	}
	match1 := make([]int32, s.n1)
	for i := range match1 {
		match1[i] = -1
	}
	match2 := make([]int32, s.n2)
	for i := range match2 {
		match2[i] = -1
	}
	val1 := make([]int32, s.n1)
	for i := range val1 {
		val1[i] = -INF32
	}
	val2 := make([]int32, s.n2)
	for i := range val2 {
		val2[i] = -INF32
	}
	stack := make([]int32, s.n1)
	for i := range stack {
		stack[i] = int32(i)
	}
	for len(stack) > 0 {
		v1 := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		match1[v1] = -1
		val1[v1] = -INF32
		if len(s.data[v1]) == 0 {
			continue
		}
		item := s.data[v1][len(s.data[v1])-1]
		s.data[v1] = s.data[v1][:len(s.data[v1])-1]
		v2, x1, x2 := item[0], item[1], item[2]
		if !chmax32(&val2[v2], x2) {
			stack = append(stack, v1)
			continue
		}
		if match2[v2] != -1 {
			stack = append(stack, match2[v2])
		}
		match1[v1] = v2
		match2[v2] = v1
		val1[v1] = x1
	}
	res := make([][2]int32, 0)
	for v1 := int32(0); v1 < s.n1; v1++ {
		v2 := match1[v1]
		if v2 != -1 {
			res = append(res, [2]int32{v1, v2})
		}
	}
	return res
}

func chmax32(a *int32, b int32) bool {
	if *a < b {
		*a = b
		return true
	}
	return false
}

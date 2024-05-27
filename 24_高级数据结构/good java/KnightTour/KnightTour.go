// 骑士路线/马的路线
// 找到一个骑士的移动顺序，使得 board 中每个单元格都 恰好 被访问一次（起始单元格已被访问，不应 再次访问）
// 输入的数据保证在给定条件下至少存在一种访问所有单元格的移动顺序。

package main

import "math/rand"

// 2664. 巡逻的骑士
// https://leetcode.cn/problems/the-knights-tour/description/
func tourOfKnight(m int, n int, r int, c int) [][]int {
	K := NewKnightTour(int32(m), int32(n), int32(r), int32(c))
	order := K.Order
	res := make([][]int, m)
	for i := range res {
		res[i] = make([]int, n)
		for j := range res[i] {
			res[i][j] = int(order[i][j])
		}
	}
	return res
}

const INF32 int32 = 1e9 + 10

var DIR8 = [][]int32{{1, 2}, {-1, 2}, {1, -2}, {-1, -2}, {2, 1}, {-2, 1}, {-2, -1}, {2, -1}}

type KnightTour struct {
	Order [][]int32
	ids   []int32
	n, m  int32
	step  int32
}

func NewKnightTour(n, m, r, c int32) *KnightTour {
	K := &KnightTour{}
	K.Order = make([][]int32, n)
	for i := range K.Order {
		K.Order[i] = make([]int32, m)
	}
	K.ids = make([]int32, n*m)
	for i := range K.ids {
		K.ids[i] = int32(i)
	}
	K.n, K.m = n, m
	for {
		for _, row := range K.Order {
			for j := range row {
				row[j] = -1
			}
		}
		rand.Shuffle(len(K.ids), func(i, j int) { K.ids[i], K.ids[j] = K.ids[j], K.ids[i] })
		K.step = 0
		if K.dfs(r, c) {
			break
		}
	}
	return K
}

func (K *KnightTour) possible(i, j int32) bool {
	return 0 <= i && i < K.n && 0 <= j && j < K.m && K.Order[i][j] == -1
}

func (K *KnightTour) degree(i, j int32) int32 {
	var d int32
	for _, dir := range DIR8 {
		if K.possible(i+dir[0], j+dir[1]) {
			d++
		}
	}
	return d
}

func (K *KnightTour) dfs(i, j int32) bool {
	K.Order[i][j] = K.step
	K.step++
	bestX, bestY := int32(-1), int32(-1)
	bestDeg := INF32
	bestId := int32(-1)
	for _, dir := range DIR8 {
		x, y := i+dir[0], j+dir[1]
		if !K.possible(x, y) {
			continue
		}
		d := K.degree(x, y)
		if d < bestDeg {
			bestDeg = d
			bestId = -1
		}
		if d == bestDeg && bestId < K.ids[x*K.m+y] {
			bestId = K.ids[x*K.m+y]
			bestX, bestY = x, y
		}
	}
	if bestDeg == INF32 {
		return K.step == K.n*K.m
	}
	return K.dfs(bestX, bestY)
}

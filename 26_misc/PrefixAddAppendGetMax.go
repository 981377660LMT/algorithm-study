// 一般用于数据结构优化dp

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	demo()
	// Solve()
}

func demo() {
	S := NewPrefixAddAppendGetMax(10)
	S.Append(1)
	S.PrefixAdd(1, 1)
	fmt.Println(S.GetMax())
	S.Append(2)
	S.PrefixAdd(2, 10)
	fmt.Println(S.GetMax())
	S.Append(3)
	S.PrefixAdd(3, 5)
	fmt.Println(S.GetMax())
}

// L. Ticket to Ride
// https://qoj.ac/contest/1472/problem/7905
//
// 有 (n + 1) 个点排成一条线，编号从 0 到 n。
// 还有 n 条线段，第 i 条线段连接点 (i − 1) 和 i。
// 给定 q 个区间 [li,ri]，每个区间还有一个分数 vi。
// 可以选择将一些线段涂红，如果点 li 到 ri 之间的线段都是红的，那么得 vi 分。
// 求恰好涂红 1, 2, · · · , n 条线段的最大得分。
//
// 表示前 i 条线段涂了 j 条，而且第 i 条线段是红色的最大得分
func Solve() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	run := func() {
		var N, M int32
		fmt.Fscan(in, &N, &M)
		RtoLX := make([][][2]int, N+1)
		for i := int32(0); i < M; i++ {
			var a, b, c int
			fmt.Fscan(in, &a, &b, &c)
			RtoLX[b] = append(RtoLX[b], [2]int{a, c})
		}

		var res []int
		dp := make([]int, N+2)
		for i := range dp {
			dp[i] = -INF
		}
		dp[0] = 0
		for i := int32(0); i < N; i++ {
			ndp := make([]int, N+2)
			for j := range ndp {
				ndp[j] = -INF
			}
			A := NewPrefixAddAppendGetMax(N + 1)
			for j := int32(0); j <= N; j++ {
				A.Append(dp[j])
				for _, pair := range RtoLX[j] {
					A.PrefixAdd(int32(pair[0]+1), pair[1])
				}
				ndp[j+1] = max(ndp[j+1], A.GetMax())
			}
			dp = ndp
			res = append(res, dp[N+1])
		}

		for i := len(res) - 1; i >= 0; i-- {
			fmt.Fprint(out, res[i], " ")
		}
		fmt.Fprintln(out)
	}

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		run()
	}
}

const INF int = 1e18

// 支持三种操作：
// 1. 在末尾追加一个数
// 2. 对前缀进行加法操作
// 3. 查询全体数的最大值
type PrefixAddAppendGetMax struct {
	n      int32
	uf     *UnionFindArraySimple32
	ranges [][2]int32
	data   []int
	max    int
}

func NewPrefixAddAppendGetMax(maxAppendOp int32) *PrefixAddAppendGetMax {
	return &PrefixAddAppendGetMax{
		n:      maxAppendOp,
		uf:     NewUnionFindArraySimple32(maxAppendOp),
		ranges: make([][2]int32, maxAppendOp),
		data:   make([]int, 0, maxAppendOp),
		max:    -INF,
	}
}

func (pa *PrefixAddAppendGetMax) GetMax() int {
	return pa.max
}

func (pa *PrefixAddAppendGetMax) Append(x int) {
	i := int32(len(pa.data))
	pa.ranges[i] = [2]int32{i, i + 1}
	if i == 0 {
		pa.max = x
		pa.data = append(pa.data, x)
		return
	}
	if x > pa.max {
		pa.data = append(pa.data, x-pa.max)
		pa.max = x
		return
	}
	pa.data = append(pa.data, 0)
	pa.merge(i)
}

// [0,end), +x
func (pa *PrefixAddAppendGetMax) PrefixAdd(end int32, x int) {
	if end == 0 {
		return
	}
	if !(end <= int32(len(pa.data)) && x >= 0) {
		panic("invalid input")
	}
	pa.data[0] += x
	pa.max += x
	for x > 0 {
		tmp := pa.uf.Find(end)
		l, r := pa.ranges[tmp][0], pa.ranges[tmp][1]
		var p int32
		if l == end {
			p = l
		} else {
			p = r
		}
		if p == int32(len(pa.data)) {
			return
		}
		y := pa.data[p]
		z := min(x, y)
		pa.data[p] -= z
		x -= z
		pa.max -= z
		if pa.data[p] == 0 {
			pa.merge(p)
			continue
		}
	}
}

func (pa *PrefixAddAppendGetMax) merge(i int32) {
	a, b := pa.uf.Find(i-1), pa.uf.Find(i)
	pa.uf.Union(a, b, nil)
	c := pa.uf.Find(a)
	pa.ranges[c][0] = pa.ranges[a][0]
	pa.ranges[c][1] = pa.ranges[b][1]
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

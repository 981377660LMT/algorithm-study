// // 带有时间序列的莫队,时间复杂度O(n^5/3)
// // https://maspypy.github.io/library/ds/offline_query/mo_3d.hpp
// // https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/mo.go

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"sort"
// )

// var _pool = make(map[interface{}]int)

// func id(o interface{}) int {
// 	if v, ok := _pool[o]; ok {
// 		return v
// 	}
// 	v := len(_pool)
// 	_pool[o] = v
// 	return v
// }

// func main() {
// 	// https://www.luogu.com.cn/problem/P1903
// 	// Q L R 查询第L支画笔到第R支画笔中共有几种不同颜色的画笔。
// 	// R P Col 把第P支画笔替换为颜色 Col

// 	// n,q<=1e5
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n, q int
// 	fmt.Fscan(in, &n, &q)
// 	nums := make([]int, n)
// 	for i := 0; i < n; i++ {
// 		fmt.Fscan(in, &nums[i])
// 	}
// 	for i, v := range nums { // 离散化
// 		nums[i] = id(v)
// 	}

// 	modify := [][3]int{} // [pos, old, new]

// }

// // 支持单点修改的莫队算法.时间复杂度O(n^5/3).
// type MoModify struct {
// 	query [][3]int // [time, left, right]
// }

// func NewMoModify() *MoModify {
// 	return &MoModify{}
// }

// // 添加一个查询，查询范围为`左闭右开区间` [start, end).
// //  0 <= left <= right <= n.
// //  time: 查询的时间点, 即当前修改的次数`len(modify)`.
// func (mm *MoModify) AddQuery(time int, start, end int) {
// 	mm.query = append(mm.query, [3]int{time, start, end})
// }

// // 返回每个查询的结果.
// //  add: 将数据添加到窗口. delta: 1 表示向右移动，-1 表示向左移动.
// //  remove: 将数据从窗口移除. delta: 1 表示向右移动，-1 表示向左移动.
// //  addUpdate: 添加修改. time: 修改的时间点. start, end: 当前窗口的范围.
// //  removeUpdate: 移除修改. time: 修改的时间点. start, end: 当前窗口的范围.
// //  query: 查询窗口内的数据.
// func (mm *MoModify) Run(
// 	add func(index int, delta int),
// 	remove func(index int, delta int),
// 	addUpdate func(time int, start, end int),
// 	removeUpdate func(time int, start, end int),
// 	query func(qid int),
// ) {
// 	q := max(1, len(mm.query))
// 	blockSize := 1
// 	for blockSize*blockSize*blockSize < q*q {
// 		blockSize++
// 	}
// 	order := mm._getOrder(blockSize)
// 	t, l, r := 0, 0, 0
// 	for _, qid := range order {
// 		item := mm.query[qid]
// 		nt, nl, nr := item[0], item[1], item[2]
// 		for l > nl {
// 			l--
// 			remove(l, -1)
// 		}
// 		for r < nr {
// 			add(r, 1)
// 			r++
// 		}
// 		for l < nl {
// 			remove(l, 1)
// 			l++
// 		}
// 		for r > nr {
// 			r--
// 			add(r, -1)
// 		}
// 		for t < nt {
// 			addUpdate(t, l, r)
// 			t++
// 		}
// 		for t > nt {
// 			t--
// 			removeUpdate(t, l, r)
// 		}
// 		query(qid)
// 	}
// }

// func (mm *MoModify) _getOrder(blockSize int) []int {
// 	k := 1 << 20
// 	q := len(mm.query)
// 	key := make([]int, q)
// 	for i, item := range mm.query {
// 		t := item[0] / blockSize
// 		l := item[1] / blockSize
// 		x := item[2]
// 		if l&1 == 1 {
// 			x = -x
// 		}
// 		x += l * k
// 		if t&1 == 1 {
// 			x = -x
// 		}
// 		x += t * k * k
// 		key[i] = x
// 	}

// 	order := make([]int, q)
// 	for i := range order {
// 		order[i] = i
// 	}
// 	sort.Slice(order, func(i, j int) bool { return key[order[i]] < key[order[j]] })

// 	cost := func(a, b int) int {
// 		q1 := mm.query[order[a]]
// 		q2 := mm.query[order[b]]
// 		return abs(q1[0]-q2[0]) + abs(q1[1]-q2[1]) + abs(q1[2]-q2[2])
// 	}
// 	for k := 0; k < q-5; k++ {
// 		if cost(k, k+2)+cost(k+1, k+3) < cost(k, k+1)+cost(k+2, k+3) {
// 			order[k+1], order[k+2] = order[k+2], order[k+1]
// 		}
// 		if cost(k, k+3)+cost(k+1, k+4) < cost(k, k+1)+cost(k+3, k+4) {
// 			order[k+1], order[k+3] = order[k+3], order[k+1]
// 		}
// 	}

// 	return order
// }

// func abs(x int) int {
// 	if x < 0 {
// 		return -x
// 	}
// 	return x
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

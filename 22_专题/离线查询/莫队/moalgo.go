package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	// G - Range Pairing Query
	// https://atcoder.jp/contests/abc242/tasks/abc242_g
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const N int = 1e5 + 10
	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	mo := NewMoAlgo(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		mo.AddQuery(l, r)
	}

	pair := 0
	counter := [N + 1]int{}
	res := make([]int, q)
	add := func(i, _ int) {
		v := nums[i]
		pair -= counter[v] / 2
		counter[v]++
		pair += counter[v] / 2
	}
	remove := func(i, _ int) {
		v := nums[i]
		pair -= counter[v] / 2
		counter[v]--
		pair += counter[v] / 2
	}
	query := func(qid int) { res[qid] = pair }

	mo.Run(add, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoAlgo struct {
	queryOrder int
	chunkSize  int
	buckets    [][]query
}

type query struct{ qi, left, right int }

func NewMoAlgo(n, q int) *MoAlgo {
	chunkSize := max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	buckets := make([][]query, n/chunkSize+1)
	return &MoAlgo{chunkSize: chunkSize, buckets: buckets}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//  0 <= left <= right <= n
func (mo *MoAlgo) AddQuery(left, right int) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right})
	mo.queryOrder++
}

// 返回每个查询的结果.
//  add: 将数据添加到窗口. delta: 1 表示向右移动，-1 表示向左移动.
//  remove: 将数据从窗口移除. delta: 1 表示向右移动，-1 表示向左移动.
//  query: 查询窗口内的数据.
func (mo *MoAlgo) Run(
	add func(index, delta int),
	remove func(index, delta int),
	query func(qid int),
) {
	left, right := 0, 0

	for i, bucket := range mo.buckets {
		if i&1 == 1 {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right < bucket[j].right })
		} else {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right > bucket[j].right })
		}

		for _, q := range bucket {
			// !窗口扩张
			for left > q.left {
				left--
				add(left, -1)
			}
			for right < q.right {
				add(right, 1)
				right++
			}

			// !窗口收缩
			for left < q.left {
				remove(left, 1)
				left++
			}
			for right > q.right {
				right--
				remove(right, -1)
			}

			query(q.qi)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

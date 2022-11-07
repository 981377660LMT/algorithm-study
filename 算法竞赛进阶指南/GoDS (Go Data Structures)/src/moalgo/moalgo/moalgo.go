package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// G - Range Pairing Query
// https://atcoder.jp/contests/abc242/tasks/abc242_g
func main() {
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

	pair := 0
	counter := [N + 1]int{}
	mo := NewMoAlgo(nums, &op{
		add: func(v, _, _, _ int) {
			pair -= counter[v] / 2
			counter[v]++
			pair += counter[v] / 2
		},
		remove: func(v, _, _, _ int) {
			pair -= counter[v] / 2
			counter[v]--
			pair += counter[v] / 2
		},
		query: func() int { return pair },
	})

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		r--
		mo.AddQuery(l, r)
	}

	res := mo.Work()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoAlgo struct {
	n          int
	queryOrder int
	chunkSize  int
	data       []int
	buckets    [][]query
	op         *op
}

type query struct {
	qi, left, right int
}

type op struct {
	// 将数据添加到窗口
	add func(value, index, qLeft, qRight int)
	// 将数据从窗口中移除
	remove func(value, index, qLeft, qRight int)
	// 更新当前窗口的查询结果
	query func() int
}

func NewMoAlgo(data []int, op *op) *MoAlgo {
	n := len(data)
	chunkSize := int(math.Ceil(math.Sqrt(float64(n))))
	buckets := make([][]query, n/chunkSize+1)
	return &MoAlgo{n: n, data: data, chunkSize: chunkSize, buckets: buckets, op: op}
}

// 0 <= left <= right < n
func (mo *MoAlgo) AddQuery(left, right int) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right + 1})
	mo.queryOrder++
}

// 返回每个查询的结果
func (mo *MoAlgo) Work() []int {
	data, buckets, q := mo.data, mo.buckets, mo.queryOrder
	res := make([]int, q)
	left, right := 0, 0

	for i, bucket := range buckets {
		if i&1 == 1 {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right > bucket[j].right })
		} else {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right < bucket[j].right })
		}

		for _, q := range bucket {
			// !窗口收缩
			for right > q.right {
				right--
				mo.op.remove(data[right], right, q.left, q.right-1)
			}
			for left < q.left {
				mo.op.remove(data[left], left, q.left, q.right-1)
				left++
			}

			// !窗口扩张
			for right < q.right {
				mo.op.add(data[right], right, q.left, q.right-1)
				right++
			}
			for left > q.left {
				left--
				mo.op.add(data[left], left, q.left, q.right-1)
			}

			res[q.qi] = mo.op.query()
		}
	}

	return res
}

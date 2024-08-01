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

	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	// 读一个整数，支持负数
	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	const N int32 = 1e5 + 10
	n := int32(NextInt())
	nums := make([]int32, n)
	for i := range nums {
		nums[i] = int32(NextInt())
	}

	q := int32(NextInt())
	mo := NewMoAlgo(n, q)
	for i := int32(0); i < q; i++ {
		l, r := int32(NextInt()), int32(NextInt())
		l--
		mo.AddQuery(l, r)
	}

	pair := 0
	counter := [N + 1]int{}
	res := make([]int, q)
	add := func(i int32) {
		v := nums[i]
		pair -= counter[v] >> 1
		counter[v]++
		pair += counter[v] >> 1
	}
	remove := func(i int32) {
		v := nums[i]
		pair -= counter[v] >> 1
		counter[v]--
		pair += counter[v] >> 1
	}
	query := func(qid int32) { res[qid] = pair }

	mo.Run(add, add, remove, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoAlgo struct {
	queryOrder int32
	chunkSize  int32
	buckets    [][]query
}

type query struct{ qi, left, right int32 }

func NewMoAlgo(n, q int32) *MoAlgo {
	chunkSize := max32(1, n/max32(1, int32(math.Sqrt(float64(q*2/3)))))
	buckets := make([][]query, n/chunkSize+1)
	return &MoAlgo{chunkSize: chunkSize, buckets: buckets}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//
//	0 <= left <= right <= n
func (mo *MoAlgo) AddQuery(left, right int32) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right})
	mo.queryOrder++
}

func (mo *MoAlgo) Run(
	addLeft func(i int32),
	addRight func(i int32),
	removeLeft func(i int32),
	removeRight func(i int32),
	query func(qi int32),
) {
	left, right := int32(0), int32(0)

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
				addLeft(left)
			}
			for right < q.right {
				addRight(right)
				right++
			}
			// !窗口收缩
			for left < q.left {
				removeLeft(left)
				left++
			}
			for right > q.right {
				right--
				removeRight(right)
			}
			query(q.qi)
		}
	}
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

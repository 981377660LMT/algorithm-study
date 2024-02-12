// https://www.luogu.com.cn/problem/U47402
// 有n个元素，每个元素有x和y属性.
// 求 aj<=ai、bj<=bi、i!=j 的数对 (i,j) 的个数.
// n<=1e5.1<=ai<=bi<=n

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	xs, ys := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
	}

	smaller := Solve(xs, ys)

	res := 0
	for _, v := range smaller {
		res += v
	}
	fmt.Fprintln(out, res)
}

// 解法: 排序+树状数组
func Solve(xs, ys []int) []int {
	n := len(xs)
	res := make([]int, n)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		a, b := order[i], order[j]
		if xs[a] != xs[b] {
			return xs[a] < xs[b]
		}
		return ys[a] < ys[b]
	})

	bit := NewBit(n + 10)
	for i := 0; i < n; i++ {
		group := []int{order[i]}
		x, y := xs[order[i]], ys[order[i]]
		for i+1 < n && xs[order[i+1]] == x && ys[order[i+1]] == y { // 处理重叠的点
			i++
			group = append(group, order[i])
		}
		bit.Add(y, len(group))
		smaller := bit.QueryPrefix(y + 1)
		for _, qi := range group {
			res[qi] = smaller - 1 // !减去自己
		}
	}

	return res
}

type BIT struct {
	n    int
	data []int
}

func NewBit(n int) *BIT {
	res := &BIT{n: n, data: make([]int, n)}
	return res
}

func (b *BIT) Add(index int, v int) {
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BIT) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BIT) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

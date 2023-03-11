// RangeXorBasisQuery
// 查询区间线性基

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func atcoder233H() {
	// https://atcoder.jp/contests/abc223/tasks/abc223_h
	// 给定一个正整数数组和q个查询,每个查询给定一个区间[lefti, righti],
	// !问是否存在一个数xi,使得从原数组区间[lefti, righti]中选取至少一个数,能凑出数异或和为xi
	// !n<=4e5 q<=2e5 nums[i]<=1<<60
	// https://maspypy.com/%e5%88%86%e5%89%b2%e7%b5%b1%e6%b2%bb%e3%81%ab%e3%82%88%e3%82%8b%e9%9d%99%e7%9a%84%e5%88%97%e3%81%ae%e5%8c%ba%e9%96%93%e7%a9%8d%e3%82%af%e3%82%a8%e3%83%aa
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	R := NewRangeXorBasisQuery(nums)
	needs := make([]int, q)
	for i := 0; i < q; i++ {
		var l, r, x int
		fmt.Fscan(in, &l, &r, &x)
		R.AddQuery(l-1, r)
		needs[i] = x
	}

	res := make([]bool, q)
	R.Run(func(qi int, basis []int) {
		need := needs[qi]
		for _, b := range basis {
			need = min(need, need^b)
		}
		res[qi] = need == 0
	})

	for _, b := range res {
		if b {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

func main() {
	Q := NewRangeXorBasisQuery([]int{3, 1, 4, 1, 5})
	Q.AddQuery(0, 3)
	Q.AddQuery(1, 5)
	Q.Run(func(qi int, basis []int) {
		fmt.Println(qi, basis)
	})
}

type RangeXorBasisQuery struct {
	n, q, log int
	nums      []int
	query     [][][2]int
}

func NewRangeXorBasisQuery(nums []int) *RangeXorBasisQuery {
	max_ := 0
	if len(nums) > 0 {
		max_ = maxs(nums...)
	}
	log := bits.Len(uint(max_))
	return &RangeXorBasisQuery{
		n:     len(nums),
		log:   log,
		nums:  nums,
		query: make([][][2]int, len(nums)+1),
	}
}

// [start, end)
//  0<=start<=end<=n
func (r *RangeXorBasisQuery) AddQuery(start, end int) {
	r.query[end] = append(r.query[end], [2]int{r.q, start})
	r.q++
}

// f: qi, basis
//  qi: 查询编号
//  basis: basis[-1]为区间线性基,basis[-2]为大于等于1<<1的线性基,basis[-3]为大于等于1<<2的线性基...
//         basis[0]为大于等于1<<log的线性基
//         也可以理解为高斯消元后从上到下的线性基
func (r *RangeXorBasisQuery) Run(f func(qi int, basis []int)) {
	d := make([][2]int, r.log)
	for k := 0; k < r.log; k++ {
		d[k] = [2]int{1 << k, -1}
	}

	basis := []int{}
	for i := 0; i < r.n+1; i++ {
		for _, q := range r.query[i] {
			qi, start := q[0], q[1]
			basis = basis[:0]
			for k := r.log - 1; k > -1; k-- {
				if d[k][1] >= start {
					basis = append(basis, d[k][0])
				}
				f(qi, basis)
			}
		}
		if i == r.n {
			break
		}
		p := [2]int{r.nums[i], i}
		for k := r.log - 1; k > -1; k-- {
			if p[0]>>k&1 == 0 {
				continue
			}
			if p[1] > d[k][1] {
				p, d[k] = d[k], p
			}
			p[0] ^= d[k][0]
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

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

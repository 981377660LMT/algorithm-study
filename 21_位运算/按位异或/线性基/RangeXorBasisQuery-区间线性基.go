// RangeXorBasisQuery
// 查询区间线性基/区间异或最大值/前缀线性基

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// demo()
	// cf1100f()
	atcoder233H()
}

func demo() {
	Q := NewRangeXorBasisQuery([]int{3, 1, 4, 1, 5})
	Q.AddQuery(0, 3)
	Q.AddQuery(1, 5)
	Q.Run(func(qi int, basis []int) {
		fmt.Println(qi, basis)
	})
}

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

// Ivan and Burgers
// https://codeforces.com/contest/1100/problem/F
// 给定一个数组和q个查询,每个查询给定一个区间[lefti, righti],
// 求在原数组区间[lefti, righti]中选取任意个数,能凑出的最大异或和
func cf1100f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	X := NewRangeXorBasisQuery(nums)
	for i := 0; i < q; i++ {
		var left, right int
		fmt.Fscan(in, &left, &right)
		left--
		X.AddQuery(left, right)
	}

	res := make([]int, q)
	X.Run(func(qi int, basis []int) {
		for _, b := range basis {
			res[qi] = max(res[qi], res[qi]^b)
		}
	})

	for _, r := range res {
		fmt.Fprintln(out, r)
	}
}

type RangeXorBasisQuery struct {
	n, q, log int
	nums      []int
	query     [][][2]int
}

func NewRangeXorBasisQuery(nums []int) *RangeXorBasisQuery {
	max_ := 0
	if len(nums) > 0 {
		max_ = maxs(nums)
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
//
//	0<=start<=end<=n
func (r *RangeXorBasisQuery) AddQuery(start, end int) {
	r.query[end] = append(r.query[end], [2]int{r.q, start})
	r.q++
}

// f: qi, basis
//
//	qi: 查询编号
//	basis: 高斯消元后从上到下的线性基，例如 [4 3 1].
func (r *RangeXorBasisQuery) Run(f func(qi int, basis []int)) {
	data := make([][2]int, r.log) // 第i位的值，对第i位造成影响的数的编号
	for k := 0; k < r.log; k++ {
		data[k] = [2]int{1 << k, -1}
	}

	// sweep line
	curBasis := []int{}
	for end := 0; end < r.n+1; end++ {
		// query
		for _, q := range r.query[end] {
			qi, start := q[0], q[1]
			curBasis = curBasis[:0]
			for k := r.log - 1; k > -1; k-- {
				if data[k][1] >= start {
					curBasis = append(curBasis, data[k][0])
				}
			}
			f(qi, curBasis)
		}
		if end == r.n {
			break
		}

		// add
		p := [2]int{r.nums[end], end}
		for k := r.log - 1; k > -1; k-- {
			if (p[0]>>k)&1 == 1 {
				if p[1] > data[k][1] {
					p, data[k] = data[k], p
				}
				p[0] ^= data[k][0]
			}
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

func maxs(nums []int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

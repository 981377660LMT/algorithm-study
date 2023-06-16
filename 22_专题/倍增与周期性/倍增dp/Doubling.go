// https://maspypy.github.io/library/ds/doubling.hpp
// 倍增doubling

// !状態 a から 1 回操作すると、状態 b に遷移し、モノイドの元 x を加える。
//  行き先がない場合：-1 （add 不要）

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1097
	// 给定一个数组和q次查询
	// 初始时res为0，每次查询会执行k次操作:
	//  将res加上nums[res%n]的值
	// 求每次查询后res的值
	// (n,q<=2e5,k<=1e12)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	db := NewDoubling(n, 1e12+10)
	for i := 0; i < n; i++ {
		db.Add(i, (i+nums[i])%n, nums[i]) // res的模从i变为(i+nums[i])%n，res加上nums[i]
	}
	db.Build()

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		_, res := db.Jump(0, k)
		fmt.Fprintln(out, res)
	}
}

type E = int

// monoidAdd
func (*Doubling) e() E          { return 0 }
func (*Doubling) op(e1, e2 E) E { return e1 + e2 }

type Doubling struct {
	n          int
	log        int
	isPrepared bool
	to         [][]int
	dp         [][]E
}

func NewDoubling(n, maxStep int) *Doubling {
	res := &Doubling{}
	res.n = n
	res.log = bits.Len(uint(maxStep))
	res.to = make([][]int, res.log)
	res.dp = make([][]E, res.log)
	for i := 0; i < res.log; i++ {
		res.to[i] = make([]int, n)
		res.dp[i] = make([]E, n)
		for j := 0; j < n; j++ {
			res.to[i][j] = -1
			res.dp[i][j] = res.e()
		}
	}
	return res
}

// 初始状态(leaves):从 `from` 状态到 `to` 状态，边权为 `weight`.
//  0 <= from, to < n
func (d *Doubling) Add(from, to int, weight E) {
	if d.isPrepared {
		panic("Doubling is prepared")
	}
	if to < -1 || to >= d.n {
		panic("to is out of range")
	}

	d.to[0][from] = to
	d.dp[0][from] = weight
}

func (d *Doubling) Build() {
	if d.isPrepared {
		panic("Doubling is prepared")
	}

	d.isPrepared = true
	for k := 0; k < d.log-1; k++ {
		for v := 0; v < d.n; v++ {
			w := d.to[k][v]
			if w == -1 {
				d.to[k+1][v] = -1
				d.dp[k+1][v] = d.dp[k][v]
				continue
			}
			d.to[k+1][v] = d.to[k][w]
			d.dp[k+1][v] = d.op(d.dp[k][v], d.dp[k][w])
		}
	}
}

// 从 `from` 状态开始，执行 `step` 次操作，返回最终状态的编号和操作的结果。
//  0 <= from < n
//  如果最终状态不存在，返回 -1, e()
func (d *Doubling) Jump(from, step int) (to int, res E) {
	if !d.isPrepared {
		panic("Doubling is not prepared")
	}
	if step >= 1<<d.log {
		panic("step is over max step")
	}

	res = d.e()
	to = from
	for k := 0; k < d.log; k++ {
		if to == -1 {
			break
		}

		if step&(1<<k) != 0 {
			res = d.op(res, d.dp[k][to])
			to = d.to[k][to]
		}
	}
	return
}

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step`.
func (d *Doubling) MaxStep(from int, check func(e E) bool) int {
	if !d.isPrepared {
		panic("Doubling is not prepared")
	}

	res := d.e()
	step := 0

	for k := d.log - 1; k >= 0; k-- {
		to := d.to[k][from]
		next := d.op(res, d.dp[k][from])
		if check(next) {
			step |= 1 << k
			from = to
			res = next
		}
	}

	return step
}

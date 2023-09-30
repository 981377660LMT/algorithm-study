package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func PointSetRangeType(nums []int, operations [][3]int) []int {}

func main() {
	// https://www.luogu.com.cn/problem/P1903
	// Q L R 查询第L支画笔到第R支画笔中共有几种不同颜色的画笔。
	// R P Col 把第P支画笔替换为颜色 Col

	// n,q<=1e5
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	operations := [][3]int{}

	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "Q" {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			operations = append(operations, [3]int{0, l, r})
		} else {
			var p, col int
			fmt.Fscan(in, &p, &col)
			p--
			operations = append(operations, [3]int{1, p, col})
		}
	}

	res := PointSetRangeType(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type Bitset []uint

func NewBitset(n int) Bitset { return make(Bitset, n>>6+1) } // (n+_w-1)>>6

func (b Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

func (b Bitset) Copy() Bitset {
	res := make(Bitset, len(b))
	copy(res, b)
	return res
}

// 借用 bits 库中的一些方法的名字
func (b Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}

// 将 c 的元素合并进 b
func (b Bitset) IOr(c Bitset) Bitset {
	for i, v := range c {
		b[i] |= v
	}
	return b
}

func (b Bitset) Or(c Bitset) Bitset {
	res := NewBitset(len(b))
	for i, v := range b {
		res[i] = v | c[i]
	}
	return res
}

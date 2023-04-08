
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/2265
	// !给定一个长为2的幂次的数组
	// 1 x y 将s[x]变为y
	// 2 left right indexXor 字符串s[i^indexXor](left<=i<right)的所有子串的数值和
	// !eg:f(123) = 123+12+23+13+1+2+3=177
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	var s string
	fmt.Fscan(in, &n, &s, &q)
	n = 1 << n

	leaves := make([]S, n)
	for i := range leaves {
		leaves[i] = createS(int(s[i] - '0'))
	}
	seg := NewXorSegTree(leaves)

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var x, y int
			fmt.Fscan(in, &x, &y)
			seg.Set(x, createS(y))
		} else {
			var left, right, indexXor int
			fmt.Fscan(in, &left, &right, &indexXor)
			right++
			res := seg.Query(left, right, indexXor)
			fmt.Fprintln(out, res.allSum)
		}
	}

}

const MOD int = 998244353

var POW2 [1 << 18]int
var POW11 [1 << 18]int

func init() {
	POW2[0] = 1
	POW11[0] = 1
	for i := 1; i < len(POW2); i++ {
		POW2[i] = POW2[i-1] * 2 % MOD
		POW11[i] = POW11[i-1] * 11 % MOD
	}
}

func createS(x int) S { return S{1, x} }

type S = struct{ size, allSum int }

func (*XorSegTree) e() S { return S{} }
func (*XorSegTree) op(a, b S) S { // !f(a+b)=f(a)*11^len(b)+f(b)*2^len(a)
	return S{a.size + b.size, (a.allSum*POW11[b.size] + b.allSum*POW2[a.size]) % MOD}
}

type XorSegTree struct {
	n, log, size int
	data         [][]S
	h            int
	unit         S
}

// XorSegTree 支持半群的单点修改和区间查询.
//  op:只需要满足结合律 op(op(a,b),c) = op(a,op(b,c)).
//  !区间查询时下标可以异或上 indexXor.
//  !nums 的长度必须要是2的幂.
func NewXorSegTree(nums []S) *XorSegTree {
	res := &XorSegTree{}
	n := len(nums)
	log := 1
	for (1 << log) < n {
		log++
	}
	size := 1 << log
	if n != size {
		panic("len(nums) must be power of 2")
	}
	H := log >> 1
	data := make([][]S, H+1)
	for i := range data {
		data[i] = make([]S, size)
	}
	for i := range nums {
		data[0][i] = nums[i]
	}

	res.n, res.log, res.size = n, log, size
	res.data = data
	res.h = H
	res.unit = res.e()

	for h := 1; h <= H; h++ {
		for i := 0; i < n>>h; i++ {
			res._update(h, i)
		}
	}

	return res
}

func (st *XorSegTree) Get(i int) S { return st.data[0][i] }

func (st *XorSegTree) Set(i int, x S) {
	st.data[0][i] = x
	for h := 1; h <= st.h; h++ {
		i >>= 1
		st._update(h, i)
	}
}

func (st *XorSegTree) Update(i int, x S) {
	st.data[0][i] = st.op(st.data[0][i], x)
	for h := 1; h <= st.h; h++ {
		i >>= 1
		st._update(h, i)
	}
}

// Calculate prod_{l<=i<r} A[x xor i], in O(log N) time.
func (st *XorSegTree) Query(start, end, indexXor int) S {
	H := st.h
	x1, x2 := st.unit, st.unit
	for h := 0; h < H; h++ {
		if start >= end {
			break
		}
		if start&(1<<h) != 0 {
			x1 = st.op(x1, st.data[h][start^indexXor])
			start += 1 << h
		}
		if end&(1<<h) != 0 {
			end -= 1 << h
			x2 = st.op(st.data[h][end^indexXor], x2)
		}
	}
	for start < end {
		x1 = st.op(x1, st.data[H][start^indexXor])
		start += 1 << H
	}
	return st.op(x1, x2)
}

func (st *XorSegTree) GetAll(i int) []S { return st.data[0] }

func (st *XorSegTree) _update(h, i int) {
	count := 1 << (h - 1)
	a := 1 << h
	b := a + count
	for k := 0; k < count; k++ {
		st.data[h][a+k] = st.op(st.data[h-1][a+k], st.data[h-1][b+k])
		st.data[h][b+k] = st.op(st.data[h-1][b+k], st.data[h-1][a+k])
	}
}

func pow(base, exp, mod int) int {
	base %= mod
	res := 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

// dynamicHashString/hashStringDynamic
// 动态哈希.
//
// api:
//  1. Insert(i, c): 在第i个位置插入字符c.
//  2. Pop(i): 删除第i个位置的字符.
//  3. Set(i, c): 将第i个位置的字符设置为c.
//  4. Reverse(start, end): 将[start, end)范围内的字符翻转.
//  5. Get(start, end): 获取[start, end)的哈希值.

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	abc331_f()
}

func demo() {
	s := "asezfvgbadpihoamgkcmco"
	base := NewHashStringBase(len(s), 37)
	hs := NewHashString(len(s), func(i int) uint { return uint(s[i]) }, base, true)
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(1, 2))
	fmt.Println(hs.Get(2, 3))
	hs.Set(0, 1)
	fmt.Println(hs.Get(0, 1))
}

// F - Palindrome Query
// https://atcoder.jp/contests/abc331/tasks/abc331_f
func abc331_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)

	base := NewHashStringBase(len(s), 0)
	hasher1 := NewHashString(len(s), func(i int) uint { return uint(s[i]) }, base, true)
	hasher2 := NewHashString(len(s), func(i int) uint { return uint(s[n-i-1]) }, base, true)

	isPalindrome := func(start, end int) bool {
		return hasher1.Get(start, end) == hasher2.Get(n-end, n-start)
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var pos int
			var c string
			fmt.Fscan(in, &pos, &c)
			pos--
			hasher1.Set(pos, uint(c[0]))
			hasher2.Set(n-pos-1, uint(c[0]))
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			if isPalindrome(l, r) {
				fmt.Fprintln(out, "Yes")
			} else {
				fmt.Fprintln(out, "No")
			}
		}
	}
}

// https://leetcode.cn/problems/sum-of-scores-of-built-strings/description/
func sumScores(s string) int64 {
	n := len(s)
	base := NewHashStringBase(n, 0)
	hasher := NewHashString(n, func(i int) uint { return uint(s[i]) }, base, false)
	countPre := func(curLen, start int) int {
		left, right := 1, curLen
		for left <= right {
			mid := (left + right) >> 1
			hash1 := hasher.Get(start, start+mid)
			hash2 := hasher.Get(0, mid)
			if hash1 == hash2 {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		return right
	}

	res := 0
	for i := 1; i < n+1; i++ {
		if s[0] != s[n-i] {
			continue
		}
		count := countPre(i, n-i)
		res += count
	}

	return int64(res)

}

const (
	hashStringMod    uint = (1 << 61) - 1
	hashStringMask30 uint = (1 << 30) - 1
	hashStringMask31 uint = (1 << 31) - 1
	hashStringMASK61 uint = hashStringMod
)

type HashStringBase struct {
	n    int
	powb []uint
	invb []uint
}

// base: 0 表示随机生成
func NewHashStringBase(n int, base uint) *HashStringBase {
	res := &HashStringBase{}
	if base == 0 {
		base = uint(37 + rand.Intn(1e9))
	}
	powb := make([]uint, n+1)
	invb := make([]uint, n+1)
	powb[0] = 1
	invb[0] = 1

	var exgcd func(a, b int) (gcd, x, y int)
	exgcd = func(a, b int) (gcd, x, y int) {
		if b == 0 {
			return a, 1, 0
		}
		gcd, y, x = exgcd(b, a%b)
		y -= a / b * x
		return
	}
	modInv := func(a, m int) int {
		g, x, _ := exgcd(a, m)
		if g != 1 && g != -1 {
			return -1
		}
		res := x % m
		if res < 0 {
			res += m
		}
		return res
	}
	invbpow := uint(modInv(int(base), int(hashStringMod)))
	for i := 1; i <= n; i++ {
		powb[i] = res.Mul(powb[i-1], base)
		invb[i] = res.Mul(invb[i-1], invbpow)
	}

	res.n = n
	res.powb = powb
	res.invb = invb
	return res
}

// h1 <- h2. len(h2) == k.
func (hsb *HashStringBase) Concat(h1, h2, h2Len uint) uint {
	return hsb.Mod(hsb.Mul(h1, hsb.powb[h2Len]) + h2)
}

// a*b % (2^61-1)
func (hsb *HashStringBase) Mul(a, b uint) uint {
	au := a >> 31
	ad := a & hashStringMask31
	bu := b >> 31
	bd := b & hashStringMask31
	mid := ad*bu + au*bd
	midu := mid >> 30
	midd := mid & hashStringMask30
	return hsb.Mod(au*bu<<1 + midu + (midd << 31) + ad*bd)
}

// x % (2^61-1)
func (hsb *HashStringBase) Mod(x uint) uint {
	xu := x >> 61
	xd := x & hashStringMASK61
	res := xu + xd
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
}

type HashString struct {
	n       int
	base    *HashStringBase
	presum  []uint
	updated bool
	seg     *segmentTree
}

func NewHashString(n int, f func(i int) uint, base *HashStringBase, update bool) *HashString {
	data := make([]uint, n)
	presum := make([]uint, n+1)
	powb := base.powb
	for i := 0; i < n; i++ {
		c := f(i)
		data[i] = base.Mul(powb[n-i-1], c)
		presum[i+1] = base.Mod(presum[i] + data[i])
	}
	res := &HashString{n: n, base: base, presum: presum, updated: false}
	if update {
		res.seg = newSegmentTreeFrom(data)
	}
	return res
}

func (hs *HashString) Get(start, end int) uint {
	if start < 0 {
		start = 0
	}
	if end > hs.n {
		end = hs.n
	}
	if start >= end {
		return 0
	}
	if hs.updated {
		return hs.base.Mul(hs.seg.Query(start, end), hs.base.invb[hs.n-end])
	} else {
		diff := uint(0)
		if v1, v2 := hs.presum[end], hs.presum[start]; v1 >= v2 {
			diff = v1 - v2
		} else {
			diff = v1 + hashStringMod - v2
		}
		return hs.base.Mul(diff, hs.base.invb[hs.n-end])
	}
}

func (hs *HashString) Set(index int, c uint) {
	hs.updated = true
	hs.seg.Set(index, hs.base.Mul(hs.base.powb[hs.n-index-1], c))
}

func (hs *HashString) Len() int { return hs.n }

type E = uint

func (*segmentTree) e() E { return 0 }
func (*segmentTree) op(a, b E) E {
	res := a + b
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
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

type segmentTree struct {
	n, size int
	seg     []E
}

func newSegmentTreeFrom(leaves []E) *segmentTree {
	res := &segmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *segmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *segmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *segmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *segmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *segmentTree) QueryAll() E { return st.seg[1] }
func (st *segmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

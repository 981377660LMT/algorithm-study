/*

References

[1]   Tiskin, A. (2008).Semi-local string comparison: Algorithmic techniques
      and applications. Mathematics in Computer Science, 1(4), 571-603.

[2]   Claude, F., Navarro, G., & Ordónez, A. (2015). The wavelet matrix: An
      efficient wavelet tree for large alphabets. Information Systems, 47,
      15-32.

*/

// TODO 有问题

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	nums := []int{3, 1, 0, 2, 4}
	lis := NewStaticRangeLISQuery(nums)
	fmt.Println(lis.Query(0, 5))
}

const None int = -1

type StaticRangeLISQuery struct {
	n  int
	wm *Wm
}

// 给定一个0-n-1的排列perm，求区间[start,end)的最长上升子序列的长度.
func NewStaticRangeLISQuery(perm []int) *StaticRangeLISQuery {
	return &StaticRangeLISQuery{
		n:  len(perm),
		wm: convert(perm),
	}
}

func (lis *StaticRangeLISQuery) Query(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > lis.n {
		end = lis.n
	}
	if start >= end {
		return 0
	}
	return (end - start) - lis.wm.CountLessThan(start, lis.n, end)
}

type _Bv struct {
	v []bvPair
}

type bvPair struct {
	bit uint64
	sum int
}

func NewBv(n int) *_Bv {
	return &_Bv{
		v: make([]bvPair, n>>6+1),
	}
}

func (b *_Bv) Set(i int) {
	b.v[i>>6].bit |= 1 << i
	b.v[i>>6].sum += 1
}

func (b *_Bv) Build() {
	for i := 1; i < len(b.v); i++ {
		b.v[i].sum += b.v[i-1].sum
	}
}

func (b *_Bv) Rank(i int) int {
	return b.v[i>>6].sum - bits.OnesCount64(uint64(b.v[i>>6].bit&(^uint64(0)<<(i&63))))
}

func (b *_Bv) One() int {
	return b.v[len(b.v)-1].sum
}

type Wm struct {
	mat []*_Bv
}

func NewWm(bitLength int, a []int) *Wm {
	mat := make([]*_Bv, bitLength)
	for i := range mat {
		mat[i] = NewBv(len(a))
	}
	a0 := make([]int, 0, len(a))
	for p := bitLength - 1; p >= 0; p-- {
		v := mat[p]
		itr := 0
		for i := 0; i < len(a); i++ {
			if test(a[i], p) {
				v.Set(i)
				a[itr] = a[i]
				itr++
			} else {
				a0 = append(a0, a[i])
			}
		}
		v.Build()
		copy(a[itr:], a0)
		a0 = a0[:0]
	}
	return &Wm{
		mat: mat,
	}
}

func (w *Wm) CountLessThan(l, r, key int) int {
	len_ := r - l
	for p := len(w.mat) - 1; p >= 0; p-- {
		v := w.mat[p]
		rankL := v.Rank(l)
		rankR := v.Rank(r)
		if test(key, p) {
			l = rankL
			r = rankR
		} else {
			len_ -= rankR - rankL
			o := v.One()
			l += o - rankL
			r += o - rankR
		}
	}
	return len_ - (r - l)
}

func test(x, k int) bool {
	return (x & (1 << k)) != 0
}

func inverse(p []int) []int {
	n := len(p)
	q := make([]int, n)
	for i := 0; i < n; i++ {
		q[i] = None
	}
	for i := 0; i < n; i++ {
		if p[i] != None {
			q[p[i]] = i
		}
	}
	return q
}

type Iter struct {
	data []int
	pos  int
}

func NewIter(data []int, pos int) *Iter {
	return &Iter{
		data: data,
		pos:  pos,
	}
}
func (it *Iter) Next(n int)      { it.pos += n }
func (it *Iter) Prev(n int)      { it.pos -= n }
func (it *Iter) Value() int      { return it.data[it.pos] }
func (it *Iter) Set(v int)       { it.data[it.pos] = v }
func (it *Iter) GetAt(i int) int { return it.data[i] }
func (it *Iter) SetAt(i, v int)  { it.data[i] = v }
func (it *Iter) Copy() *Iter {
	return &Iter{
		pos:  it.pos,
		data: it.data,
	}
}

type DIter struct{ delta, col int }

func NewDIter() *DIter { return &DIter{} }
func (it *DIter) Copy() *DIter {
	return &DIter{
		delta: it.delta,
		col:   it.col,
	}
}

func unitMongeMul(n int, stack *Iter, a, b *Iter) {
	if n == 1 {
		stack.SetAt(0, 0)
		return
	}

	cRow := stack.Copy()
	stack.Next(n)
	cCol := stack.Copy()
	stack.Next(n)
	mp := func(len int, f func(int) bool, g func(int) int) {
		ah := stack.Copy()
		am := stack.Copy()
		am.Next(len)
		bh := stack.Copy()
		bh.Next(2 * len)
		bm := stack.Copy()
		bm.Next(3 * len)

		split := func(v, vh, vm *Iter) {
			for i := 0; i < n; i++ {
				if f(v.GetAt(i)) {
					vh.Set(g(v.GetAt(i)))
					vh.Next(1)
					vm.Set(i)
					vm.Next(1)
				}
			}
		}

		split(a, ah.Copy(), am.Copy())
		split(b, bh.Copy(), bm.Copy())
		c := stack.Copy()
		c.Next(4 * len)
		unitMongeMul(len, c.Copy(), ah.Copy(), bh.Copy())
		for i := 0; i < len; i++ {
			row := am.GetAt(i)
			col := bm.GetAt(c.GetAt(i))
			cRow.SetAt(row, col)
			cCol.SetAt(col, row)
		}
	}

	mid := n >> 1
	mp(mid, func(x int) bool { return x < mid }, func(x int) int { return x })
	mp(n-mid, func(x int) bool { return x >= mid }, func(x int) int { return x - mid })

	row := n
	right := func(it *DIter) {
		if b.GetAt(it.col) < mid {
			if cCol.GetAt(it.col) >= row {
				it.delta += 1
			}
		} else {
			if cCol.GetAt(it.col) < row {
				it.delta += 1
			}
		}
		it.col += 1
	}
	up := func(it *DIter) {
		if a.GetAt(row) < mid {
			if cRow.GetAt(row) >= it.col {
				it.delta -= 1
			}
		} else {
			if cRow.GetAt(row) < it.col {
				it.delta -= 1
			}
		}
	}

	neg, pos := NewDIter(), NewDIter()
	for row != 0 {
		for pos.col != n {
			temp := pos.Copy()
			right(temp)
			if temp.delta == 0 {
				pos = temp
			} else {
				break
			}
		}
		row--
		up(neg)
		up(pos)
		for neg.delta != 0 {
			right(neg)
		}
		if neg.col > pos.col {
			cRow.SetAt(row, pos.col)
		}
	}
}

func subunitMongeMul(a, b []int) []int {
	n := len(a)
	aInv := inverse(a)
	bInv := inverse(b)
	b, bInv = bInv, b
	var aMap, bMap []int
	for i := n - 1; i >= 0; i-- {
		if a[i] != None {
			aMap = append(aMap, i)
			a[n-len(aMap)] = a[i]
		}
	}

	for i, j := 0, len(aMap)-1; i < j; i, j = i+1, j-1 {
		aMap[i], aMap[j] = aMap[j], aMap[i]
	}

	{
		cnt := 0
		for i := 0; i < n; i++ {
			if aInv[i] == None {
				a[cnt] = i
				cnt += 1
			}
		}
	}

	for i := 0; i < n; i++ {
		if b[i] != None {
			b[len(bMap)] = b[i]
			bMap = append(bMap, i)
		}
	}

	{
		cnt := len(bMap)
		for i := 0; i < n; i++ {
			if bInv[i] == None {
				b[cnt] = i
				cnt += 1
			}
		}
	}

	cSize := func(n int) int {
		res := 0
		for n > 1 {
			res += 2 * n
			n = (n + 1) / 2
			res += 4 * n
		}
		return res + 1
	}(n)
	c := make([]int, cSize)

	unitMongeMul(n, NewIter(c, 0), NewIter(a, 0), NewIter(b, 0))
	cPad := make([]int, n)
	for i := range cPad {
		cPad[i] = None
	}

	for i := 0; i < len(aMap); i++ {
		t := c[n-len(aMap)+i]
		if t < len(bMap) {
			cPad[aMap[i]] = bMap[t]
		}
	}

	return cPad
}

func seaweedDoubling(p []int) []int {
	n := len(p)
	if n == 1 {
		return []int{None}
	}
	mid := n / 2
	var lo, hi []int
	var loMap, hiMap []int
	for i := 0; i < n; i++ {
		e := p[i]
		if e < mid {
			lo = append(lo, e)
			loMap = append(loMap, i)
		} else {
			hi = append(hi, e-mid)
			hiMap = append(hiMap, i)
		}
	}
	lo = seaweedDoubling(lo)
	hi = seaweedDoubling(hi)
	loPad, hiPad := make([]int, n), make([]int, n)
	for i := range loPad {
		loPad[i] = i
		hiPad[i] = i
	}

	for i := 0; i < mid; i++ {
		if lo[i] == None {
			loPad[loMap[i]] = None
		} else {
			loPad[loMap[i]] = loMap[lo[i]]
		}
	}
	for i := 0; mid+i < n; i++ {
		if hi[i] == None {
			hiPad[hiMap[i]] = None
		} else {
			hiPad[hiMap[i]] = hiMap[hi[i]]
		}
	}

	return subunitMongeMul(loPad, hiPad)
}

func isPermutation(p []int) bool {
	n := len(p)
	used := make([]bool, n)
	for _, e := range p {
		if e < 0 || n <= e || used[e] {
			return false
		}
		used[e] = true
	}
	return true
}

func convert(p []int) *Wm {
	if !isPermutation(p) {
		panic("not permutation")
	}
	n := len(p)
	var row []int
	if n != 0 {
		row = seaweedDoubling(p)
	}
	for i := range row {
		if row[i] == None {
			row[i] = n
		}
	}
	bitLength := 0
	for n > 0 {
		bitLength += 1
		n /= 2
	}
	return NewWm(bitLength, row)
}

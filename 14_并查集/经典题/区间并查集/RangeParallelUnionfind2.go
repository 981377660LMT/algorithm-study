// nlogn区间合并
// 支持在线时时询问和修改 RangeParallelUnionfindOnline
// 复杂度浪费在了，有一些点本来已经是相同的了，但是我们又把他 merge 了一遍，
// 所以浪费了大量时间。我们考虑对序列进行哈希，每次二分找到最近的不相同的节点，
// 然后把他们 merge 起来，再把所有和他们位置相同的节点的哈希值改掉。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	P3295()
	// atc2018()
}

// https://www.luogu.com.cn/problem/P3295
// 给定一个长度为n的大数，每个大数元素为0到9之间的整数(注意不能有前导零)。
// 再给定一些约束条件，形如[start1,end1,start2,end2]，表示[start1,end1)区间内的数和[start2,end2)区间内的数相等。
// 问满足以上所有条件的数有多少个，对1e9+7取模。
func P3295() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7
	qpow := func(a, b int) int {
		res := 1
		for b > 0 {
			if b&1 == 1 {
				res = res * a % MOD
			}
			a = a * a % MOD
			b >>= 1
		}
		return res
	}

	var n, m int32
	fmt.Fscan(in, &n, &m)
	uf := NewParallelUnionFind(n)
	for i := int32(0); i < m; i++ {
		var start1, end1, start2, end2 int32
		fmt.Fscan(in, &start1, &end1, &start2, &end2)
		start1, start2 = start1-1, start2-1
		uf.UnionRange(start1, end1, start2, end2, func(v, fromRoot, toRoot int32) {})
	}

	leader := make(map[int32]struct{})
	for i := int32(0); i < n; i++ {
		leader[uf.Find(i)] = struct{}{}
	}
	part := len(leader)
	fmt.Fprintln(out, 9*qpow(10, part-1)%MOD)

}

func main2() {
	// https://atcoder.jp/contests/yahoo-procon2018-final/tasks/yahoo_procon2018_final_d
	// !前缀和后缀的LCP为lens[i]的字符串
	// !LCP => 并查集
	// 给定长为n的数组lens, 问是否存在一个长度为s的字符串,满足:
	// !s[0:i+1] 和 s[n-(i+1):n] 的最长公共前缀为 lens[i] (0<=i<n)
	// n<=3e5 0<=lens[i]<=i+1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	lens := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &lens[i])
	}

	uf := NewParallelUnionFind(n)
	for i := int32(0); i < n; i++ {
		// uf.UnionParallelly(0, n-(i+1), lens[i]) // 各个位置的字符相同
		start1, start2 := int32(0), n-(i+1)
		end1, end2 := lens[i], n-(i+1)+lens[i]
		uf.UnionRange(start1, end1, start2, end2, func(v, y, x int32) {})
	}

	for i := int32(0); i < n; i++ {
		if lens[i] == i+1 {
			continue
		}
		if uf.IsConnected(lens[i], n-(i+1)+lens[i]) { // !s[len[i]]!=s[n-(i+1)+len[i]] (因为前后缀LCP只有len[i])
			fmt.Fprintln(out, "No")
			return
		}
	}
	fmt.Fprintln(out, "Yes")
}

const BASE uint64 = 13331

type ParallelUnionFind struct {
	n, log int32
	seg    []uint64
	pow    []uint64
	dat    []int32
	nxt    []int32
}

func NewParallelUnionFind(n int32) *ParallelUnionFind {
	res := &ParallelUnionFind{}
	res.n = n

	log := int32(1)
	for 1<<log < n {
		log++
	}
	res.log = log

	pow := make([]uint64, 1<<log)
	pow[0] = 1
	for i := 0; i < 1<<log-1; i++ {
		pow[i+1] = pow[i] * BASE
	}
	res.pow = pow

	seg := make([]uint64, 2<<log)
	for i := int32(0); i < n; i++ {
		seg[i+(1<<log)] = uint64(i)
	}
	res.seg = seg

	for i := int32(1<<log - 1); i >= 1; i-- {
		res._update(i)
	}

	data, next := make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
		next[i] = -1
	}

	res.dat = data
	res.nxt = next
	return res
}

func (p *ParallelUnionFind) Find(x int32) int32 {
	if p.dat[x] < 0 {
		return x
	}
	return p.dat[x]
}

func (p *ParallelUnionFind) GetSize(x int32) int32 {
	return -p.dat[p.Find(x)]
}

// 并行合并[(a,c),(a+1,c+1),...,(b-1,d-1)].
// f(v, fromRoot, toRoot): root(v) = fromRoot -> root(v) = toRoot
func (p *ParallelUnionFind) UnionRange(start1, end1, start2, end2 int32, f func(v, fromRoot, toRoot int32)) {
	if end1-start1 != end2-start2 {
		panic("invalid range")
	}
	for {
		if p._get(start1, end1) == p._get(start2, end2) {
			break
		}
		n := p._binarySearch(func(k int32) bool {
			return p._get(start1, start1+k) == p._get(start2, start2+k)
		}, 0, end1-start1)
		x, y := p.Find(start1+n), p.Find(start2+n)
		start1, start2 = start1+n, start2+n
		if p.dat[x] > p.dat[y] {
			x, y = y, x
		}
		for p.nxt[y] != -1 {
			v := p.nxt[y]
			p.nxt[y] = p.nxt[v]
			p._set(v, uint64(x))
			f(v, y, x)
			p.dat[v] = x
			p.dat[x]--
			p.nxt[v] = p.nxt[x]
			p.nxt[x] = v
		}
		p._set(y, uint64(x))
		f(y, y, x)
		p.dat[y] = x
		p.dat[x]--
		p.nxt[y] = p.nxt[x]
		p.nxt[x] = y
	}
}

func (p *ParallelUnionFind) Union(start, end int32, f func(v, fromRoot, toRoot int32)) {
	p.UnionRange(start, start+1, end, end+1, f)
}

func (p *ParallelUnionFind) IsConnected(x, y int32) bool {
	return p.Find(x) == p.Find(y)
}

func (p *ParallelUnionFind) GetGroups() map[int32][]int32 {
	res := make(map[int32][]int32)
	for i := int32(0); i < p.n; i++ {
		root := p.Find(i)
		res[root] = append(res[root], i)
	}
	return res
}

func (p *ParallelUnionFind) String() string {
	groups := p.GetGroups()
	return fmt.Sprintf("%v", groups)
}

func (p *ParallelUnionFind) _update(i int32) {
	sz := 1 << (p.log - topbit(i) - 1)
	p.seg[i] = p.seg[i<<1]*p.pow[sz] + p.seg[i<<1|1]
}

func (p *ParallelUnionFind) _set(i int32, x uint64) {
	i += 1 << p.log
	p.seg[i] = x
	for i >>= 1; i > 0; i >>= 1 {
		p._update(i)
	}
}

func (p *ParallelUnionFind) _get(L, R int32) uint64 {
	xl, xr := uint64(0), uint64(0)
	sl, sr := 0, 0
	L += 1 << p.log
	R += 1 << p.log
	s := 1
	for L < R {
		if L&1 > 0 {
			xl = xl*p.pow[s] + p.seg[L]
			sl += s
			L++
		}
		if R&1 > 0 {
			R--
			xr = p.seg[R]*p.pow[sr] + xr
			sr += s
		}
		L >>= 1
		R >>= 1
		s <<= 1
	}
	return xl*p.pow[sr] + xr
}

func (p *ParallelUnionFind) _binarySearch(f func(int32) bool, l, r int32) int32 {
	for r-l > 1 {
		m := (l + r) >> 1
		if f(m) {
			l = m
		} else {
			r = m
		}
	}
	return l
}

func topbit(x int32) int32 {
	return int32(31 - bits.LeadingZeros32(uint32(x)))
}

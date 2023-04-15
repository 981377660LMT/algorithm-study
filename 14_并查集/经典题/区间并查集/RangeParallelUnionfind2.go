// nlogn区间合并
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main2() {
	uf := NewParallelUnionFind(10)
	fmt.Println(uf)
	uf.UnionRange(0, 5, 5, 10, func(v, from, to int) {
		fmt.Printf("v=%d from=%d to=%d\n", v, from, to)
	})
	fmt.Println(uf)
}

func main() {
	// https://atcoder.jp/contests/yahoo-procon2018-final/tasks/yahoo_procon2018_final_d
	// !前缀和后缀的LCP为lens[i]的字符串
	// !LCP => 并查集
	// 给定长为n的数组lens, 问是否存在一个长度为s的字符串,满足:
	// !s[0:i+1] 和 s[n-(i+1):n] 的最长公共前缀为 lens[i] (0<=i<n)
	// n<=3e5 0<=lens[i]<=i+1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	lens := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &lens[i])
	}

	uf := NewParallelUnionFind(n)
	for i := 0; i < n; i++ {
		// uf.UnionParallelly(0, n-(i+1), lens[i]) // 各个位置的字符相同
		start1, start2 := 0, n-(i+1)
		end1, end2 := lens[i], n-(i+1)+lens[i]
		uf.UnionRange(start1, end1, start2, end2, func(v, y, x int) {})
	}

	for i := 0; i < n; i++ {
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
	n, log int
	seg    []uint64
	pow    []uint64
	dat    []int
	nxt    []int
}

func NewParallelUnionFind(n int) *ParallelUnionFind {
	res := &ParallelUnionFind{}
	res.n = n

	log := 1
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
	for i := 0; i < n; i++ {
		seg[i+(1<<log)] = uint64(i)
	}
	res.seg = seg

	for i := 1<<log - 1; i >= 1; i-- {
		res._update(i)
	}

	data, next := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
		next[i] = -1
	}

	res.dat = data
	res.nxt = next
	return res
}

func (p *ParallelUnionFind) Find(x int) int {
	if p.dat[x] < 0 {
		return x
	}
	return p.dat[x]
}

func (p *ParallelUnionFind) GetSize(x int) int {
	return -p.dat[p.Find(x)]
}

// 并行合并[(a,c),(a+1,c+1),...,(b-1,d-1)].
func (p *ParallelUnionFind) UnionRange(start1, end1, start2, end2 int, f func(v, fromRoot, toRoot int)) {
	if end1-start1 != end2-start2 {
		panic("invalid range")
	}
	for {
		if p._get(start1, end1) == p._get(start2, end2) {
			break
		}
		n := p._binarySearch(func(k int) bool {
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

func (p *ParallelUnionFind) Union(start, end int, f func(v, fromRoot, toRoot int)) {
	p.UnionRange(start, start+1, end, end+1, f)
}

func (p *ParallelUnionFind) IsConnected(x, y int) bool {
	return p.Find(x) == p.Find(y)
}

func (p *ParallelUnionFind) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for i := 0; i < p.n; i++ {
		root := p.Find(i)
		res[root] = append(res[root], i)
	}
	return res
}

func (p *ParallelUnionFind) String() string {
	groups := p.GetGroups()
	return fmt.Sprintf("%v", groups)
}

func (p *ParallelUnionFind) _update(i int) {
	sz := 1 << (p.log - topbit(i) - 1)
	p.seg[i] = p.seg[i<<1]*p.pow[sz] + p.seg[i<<1|1]
}

func (p *ParallelUnionFind) _set(i int, x uint64) {
	i += 1 << p.log
	p.seg[i] = x
	for i >>= 1; i > 0; i >>= 1 {
		p._update(i)
	}
}

func (p *ParallelUnionFind) _get(L, R int) uint64 {
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

func (p *ParallelUnionFind) _binarySearch(f func(int) bool, l, r int) int {
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

func topbit(x int) int {
	return 31 - bits.LeadingZeros32(uint32(x))
}

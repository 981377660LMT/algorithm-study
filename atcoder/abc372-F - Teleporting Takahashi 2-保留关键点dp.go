package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

// abc372-F - Teleporting Takahashi 2-保留关键点dp
// https://atcoder.jp/contests/abc372/tasks/abc372_f
// !注意，尽量不要用map保存记忆化状态，很慢
func main() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	NextByte := func() byte {
		b := rc()
		for ; '0' > b; b = rc() {
		}
		return b
	}
	_ = NextByte

	// 读一个整数，支持负数
	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	const MOD int = 998244353

	N, M, K := int32(NextInt()), int32(NextInt()), int32(NextInt())
	if M == 0 {
		fmt.Fprintln(out, 1)
		return
	}

	adjList := make([][][2]int32, N)
	critical := NewFastSet32(N + N)
	for i := int32(0); i < M; i++ {
		x, y := int32(NextInt())-1, int32(NextInt())-1
		adjList[x] = append(adjList[x], [2]int32{y, 1})
		critical.Insert(x)
		critical.Insert(y)
		critical.Insert(x + N)
		critical.Insert(y + N)
	}
	critical.Insert(0) // 起点
	critical.Insert(N)

	for i := int32(0); i < N; i++ {
		nextStart := critical.Next(i + 1)
		adjList[i] = append(adjList[i], [2]int32{nextStart % N, nextStart - i})
	}

	// 坐标圧縮
	criticalId := make([]int32, N)
	rawId := []int32{}
	ptr := int32(0)
	critical.Enumerate(0, N, func(i int32) {
		criticalId[i] = ptr
		rawId = append(rawId, i)
		ptr++
	})

	memo := make([]int, critical.Size()*K)
	for i := range memo {
		memo[i] = -1
	}
	var dfs func(int32, int32) int
	dfs = func(v, k int32) int {
		if k <= 0 {
			return 1
		}
		hash := v*K + k
		if memo[hash] != -1 {
			return memo[hash]
		}
		res := 0
		for _, e := range adjList[rawId[v]] {
			res += dfs(criticalId[e[0]], k-e[1])
			res %= MOD
		}
		memo[hash] = res
		return res
	}

	res := dfs(criticalId[0], K)
	fmt.Fprintln(out, res)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

type FastSet32 struct {
	n, lg int32
	seg   [][]uint64
	size  int32
}

func NewFastSet32(n int32) *FastSet32 {
	res := &FastSet32{n: n}
	seg := [][]uint64{}
	n_ := n
	for {
		seg = append(seg, make([]uint64, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = int32(len(seg))
	return res
}

func NewFastSet32From(n int32, f func(i int32) bool) *FastSet32 {
	res := NewFastSet32(n)
	for i := int32(0); i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := int32(0); h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet32) Has(i int32) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet32) Insert(i int32) bool {
	if fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet32) Erase(i int32) bool {
	if !fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

func (fs *FastSet32) Next(i int32) int32 {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == int32(len(cache)) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

func (fs *FastSet32) Prev(i int32) int32 {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := int32(0); h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

func (fs *FastSet32) Enumerate(start, end int32, f func(i int32)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet32) String() string {
	res := []string{}
	for i := int32(0); i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(int(i)))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet32) Size() int32 {
	return fs.size
}

func (*FastSet32) bsr(x uint64) int32 {
	return 63 - int32(bits.LeadingZeros64(x))
}

func (*FastSet32) bsf(x uint64) int32 {
	return int32(bits.TrailingZeros64(x))
}

// 给定元素0~n-1，对数组中的某些特殊元素进行离散化.
// 返回离散化后的数组id和id对应的值.
// 特殊元素的id为0~len(idToV)-1, 非特殊元素的id为-1.
func DiscretizeSpecial(n int32, isSpecial func(i int32) bool) (vToId []int32, idToV []int32) {
	vToId = make([]int32, n)
	idToV = []int32{}
	ptr := int32(0)
	for i := int32(0); i < n; i++ {
		if isSpecial(i) {
			vToId[i] = ptr
			ptr++
			idToV = append(idToV, i)
		} else {
			vToId[i] = -1
		}
	}
	idToV = idToV[:len(idToV):len(idToV)]
	return
}

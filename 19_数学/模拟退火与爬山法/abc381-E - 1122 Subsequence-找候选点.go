// E - 11/22 Subsequence
// https://atcoder.jp/contests/abc381/tasks/abc381_e
// 给定一个字符串s，给定q个询问。
// 每个询问给定l,r，问s[l,r]的最长子序列，使得子序列中1的个数和2的个数相等，中间用/分隔。
//
// !二分求出分界点后，再向分界点左右看是否有合适的分隔符.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int32
	fmt.Fscan(in, &N, &Q)
	var S string
	fmt.Fscan(in, &S)

	presumA, presumB := make([]int32, N+1), make([]int32, N+1)
	for i := int32(0); i < N; i++ {
		if S[i] == '1' {
			presumA[i+1] = presumA[i] + 1
		} else {
			presumA[i+1] = presumA[i]
		}
		if S[i] == '2' {
			presumB[i+1] = presumB[i] + 1
		} else {
			presumB[i+1] = presumB[i]
		}
	}

	fs := NewFastSet32(N)
	for i := int32(0); i < N; i++ {
		if S[i] == '/' {
			fs.Insert(i)
		}
	}

	query := func(start, end int32) int32 {
		// 最大的upper 使得[start,upper)内的1的个数小于2的个数.
		upper := MaxRight32(
			start,
			func(right int32) bool {
				ca, cb := presumA[right]-presumA[start], presumB[end]-presumB[right]
				return ca < cb
			},
			end,
		)

		var cands []int32
		{
			i := fs.Prev(upper)
			for j := 0; j < 3; j++ {
				if i == -1 {
					break
				}
				cands = append(cands, i)
				i = fs.Prev(i - 1)
			}
			i = fs.Next(upper)
			for j := 0; j < 3; j++ {
				if i == N {
					break
				}
				cands = append(cands, i)
				i = fs.Next(i + 1)
			}
		}

		res := int32(0)
		for _, p := range cands {
			if !(start <= p && p < end) {
				continue
			}
			ca, cb := presumA[p]-presumA[start], presumB[end]-presumB[p]
			res = max32(res, min32(ca, cb)*2+1)
		}
		return res
	}

	for i := int32(0); i < Q; i++ {
		var L, R int32
		fmt.Fscan(in, &L, &R)
		L--
		fmt.Fprintln(out, query(L, R))
	}
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

// 返回大于等于i的最小元素.如果不存在,返回n.
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
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
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
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
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

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含,使用时需要right-1.
// right<=upper.
func MaxRight32(left int32, check func(right int32) bool, upper int32) int32 {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

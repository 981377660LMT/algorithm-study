// 后缀数组解决最长回文子串问题
// !回文子串即为原字符串的前缀的后缀和原字符串的后缀的前缀两部分组成
// 如果我们将字符串调转过来，在接到原字符串的后面，中间隔一个字典序最小的'板'，这样前缀就会被调转过来，‘前缀的后缀’就会变为‘后缀的前缀’。
// 问题就转化为了求某两个后缀的最长公共前缀，就可以用后缀数组求解了。

package main

import (
	"fmt"
	"index/suffixarray"
	"math/bits"
	"reflect"
	"unsafe"
)

func main() {
	fmt.Println(longestPalindrome("a"))
}

// https://leetcode.cn/problems/longest-palindromic-substring/description/
// https://www.cnblogs.com/GDOI2018/p/10296315.html
// aabaaaab
// aabaaaab$baaaabaa
// 记新字符串的长度为m.
// 长度为奇数，则查询后缀i和后缀m-i-1. -> aa|baaaab$baaaa|baa
// 长度为偶数，则查询后缀i和后缀m-i. -> aabaa|aab$baa|aabaa
func longestPalindrome(s string) string {
	n := int32(len(s))
	m := 2*n + 1
	f := func(i int32) int32 {
		if i < n {
			return int32(s[i])
		}
		if i == n {
			return -1
		}
		return int32(s[m-1-i])
	}

	S := NewSuffixArray(m, f)
	resStart, resLen := int32(0), int32(0)
	for i := int32(0); i < n; i++ {
		odd := S.Lcp(i, m, m-i-1, m)
		if tmp := 2*odd - 1; tmp > resLen {
			resLen = tmp
			resStart = i - odd + 1
		}

		even := S.Lcp(i, m, m-i, m)
		if tmp := 2 * even; tmp > resLen {
			resLen = tmp
			resStart = i - even
		}
	}

	return s[resStart : resStart+resLen]
}

type SuffixArray32 struct {
	Sa     []int32 // 排名第i的后缀是谁.
	Rank   []int32 // 后缀s[i:]的排名是多少.
	Height []int32 // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	Ords   []int32
	n      int32
	minSt  *LinearRMQ32 // 维护lcp的最小值
}

func NewSuffixArray(n int32, f func(i int32) int32) *SuffixArray32 {
	res := &SuffixArray32{n: n}
	sa, rank, lcp := SuffixArray32Simple(n, f)
	res.Sa, res.Rank, res.Height = sa, rank, lcp
	return res
}

// 求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray32) Lcp(a, b int32, c, d int32) int32 {
	if a >= b || c >= d {
		return 0
	}
	cand := suf._lcp(a, c)
	return min32(cand, min32(b-a, d-c))
}

// 比较任意两个子串s[a,b)和s[c,d)的字典序.
//
//	s[a,b) < s[c,d) 返回-1.
//	s[a,b) = s[c,d) 返回0.
//	s[a,b) > s[c,d) 返回1.
func (suf *SuffixArray32) CompareSubstr(a, b int32, c, d int32) int32 {
	len1, len2 := b-a, d-c
	lcp := suf.Lcp(a, b, c, d)
	if len1 == len2 && lcp >= len1 {
		return 0
	}
	if lcp >= len1 || lcp >= len2 { // 一个是另一个的前缀
		if len1 < len2 {
			return -1
		}
		return 1
	}
	if suf.Rank[a] < suf.Rank[c] {
		return -1
	}
	return 1
}

// 与 s[left:] 的 lcp 大于等于 k 的后缀数组(sa)上的区间.
// 如果不存在,返回(-1,-1).
func (suf *SuffixArray32) LcpRange(left int32, k int32) (start, end int32) {
	if k > suf.n-left {
		return -1, -1
	}
	if k == 0 {
		return 0, suf.n
	}
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ32(suf.Height)
	}
	i := suf.Rank[left] + 1
	start = suf.minSt.MinLeft(i, func(e int32) bool { return e >= k }) - 1 // 向左找
	end = suf.minSt.MaxRight(i, func(e int32) bool { return e >= k })      // 向右找
	return
}

// 查询s[start:end)在s中的出现次数.
func (suf *SuffixArray32) Count(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > suf.n {
		end = suf.n
	}
	if start >= end {
		return 0
	}
	a, b := suf.LcpRange(start, end-start)
	return b - a
}

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *SuffixArray32) _lcp(i, j int32) int32 {
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ32(suf.Height)
	}
	if i == j {
		return suf.n - i
	}
	r1, r2 := suf.Rank[i], suf.Rank[j]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return suf.minSt.Query(r1+1, r2+1)
}

type LinearRMQ32 struct {
	n     int32
	nums  []int32
	small []int
	large [][]int32
}

func NewLinearRMQ32(nums []int32) *LinearRMQ32 {
	n := int32(len(nums))
	res := &LinearRMQ32{n: n, nums: nums}
	stack := make([]int32, 0, 64)
	small := make([]int, 0, n)
	var large [][]int32
	large = append(large, make([]int32, 0, n>>6))
	for i := int32(0); i < n; i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] > nums[i] {
			stack = stack[:len(stack)-1]
		}
		tmp := 0
		if len(stack) > 0 {
			tmp = small[stack[len(stack)-1]]
		}
		small = append(small, tmp|(1<<(i&63)))
		stack = append(stack, i)
		if (i+1)&63 == 0 {
			large[0] = append(large[0], stack[0])
			stack = stack[:0]
		}
	}

	for i := int32(1); (i << 1) <= n>>6; i <<= 1 {
		csz := int32(n>>6 + 1 - (i << 1))
		v := make([]int32, csz)
		for k := int32(0); k < csz; k++ {
			back := large[len(large)-1]
			v[k] = res._getMin(back[k], back[k+i])
		}
		large = append(large, v)
	}

	res.small = small
	res.large = large
	return res
}

// 查询区间`[start, end)`中的最小值.
func (rmq *LinearRMQ32) Query(start, end int32) int32 {
	end--
	left := start>>6 + 1
	right := end >> 6
	if left < right {
		msb := bits.Len64(uint64(right-left)) - 1
		cache := rmq.large[msb]
		i := (left-1)<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63)))))
		cand1 := rmq._getMin(i, cache[left])
		j := right<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[end])))
		cand2 := rmq._getMin(cache[right-(1<<msb)], j)
		return rmq.nums[rmq._getMin(cand1, cand2)]
	}
	if left == right {
		i := (left-1)<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63)))))
		j := left<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[end])))
		return rmq.nums[rmq._getMin(i, j)]
	}
	return rmq.nums[right<<6+int32(bits.TrailingZeros64(uint64(rmq.small[end]&(^0<<(start&63)))))]
}

func (rmq *LinearRMQ32) _getMin(i, j int32) int32 {
	if rmq.nums[i] < rmq.nums[j] {
		return i
	}
	return j
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ32) MaxRight(left int32, check func(e int32) bool) int32 {
	if left == st.n {
		return st.n
	}
	ok, ng := left, st.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(st.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ32) MinLeft(right int32, check func(e int32) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(st.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func SuffixArray32Simple(n int32, f func(i int32) int32) (sa, rank, height []int32) {
	s := make([]byte, 0, n*4)
	for i := int32(0); i < n; i++ {
		v := f(i)
		s = append(s, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	_sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
	sa = make([]int32, 0, n)
	for _, v := range _sa {
		if v&3 == 0 {
			sa = append(sa, v>>2)
		}
	}
	rank = make([]int32, n)
	for i := int32(0); i < n; i++ {
		rank[sa[i]] = i
	}
	height = make([]int32, n)
	h := int32(0)
	for i := int32(0); i < n; i++ {
		rk := rank[i]
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && f(i+h) == f(j+h); h++ {
			}
		}
		height[rk] = h
	}
	return
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

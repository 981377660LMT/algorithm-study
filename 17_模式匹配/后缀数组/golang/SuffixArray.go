// - 支持sa/rank/lcp
// - 比较任意两个子串的字典序
// - 求出任意两个子串的最长公共前缀(lcp)

//  sa : 排第几的后缀是谁.
//  rank : 每个后缀排第几.
//  lcp : 排名相邻的两个后缀的最长公共前缀.
// 	lcp[0] = 0
// 	lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
//
//  "banana" -> sa: [5 3 1 0 4 2], rank: [3 2 5 1 4 0], lcp: [0 1 3 0 0 2]

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"
	"time"
)

// https://leetcode.cn/problems/sum-of-scores-of-built-strings/
func sumScores(s string) int64 {
	sa := NewSuffixArrayWithString(s)
	n := len(s)
	res := 0
	for i := 0; i < n; i++ {
		res += sa.Lcp(0, n, i, n)
	}
	return int64(res)
}

func main() {
	s := "banana"
	sa := NewSuffixArrayWithString(s)
	fmt.Println(sa.Sa, sa.Rank, sa.Height) // [5 3 1 0 4 2] [3 2 5 1 4 0] [0 1 3 0 0 2]

	lcpNaive := func(s string, a, b, c, d int) int {
		res := 0
		for a < b && c < d && s[a] == s[c] {
			a++
			c++
			res++
		}
		return res
	}

	compareSubtrNaive := func(s string, a, b, c, d int) int {
		for a < b && c < d && s[a] == s[c] {
			a++
			c++
		}
		if a == b {
			if c == d {
				return 0
			}
			return -1
		}
		if c == d {
			return 1
		}
		if s[a] < s[c] {
			return -1
		}
		return 1
	}

	n := 100
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteByte(byte(rand.Intn(26) + 'a'))
	}
	s = sb.String()
	sa = NewSuffixArrayWithString(s)
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			for k := j + 1; k < len(s); k++ {
				for l := k + 1; l < len(s); l++ {
					if sa.Lcp(i, j, k, l) != lcpNaive(s, i, j, k, l) {
						fmt.Println("lcp err")
					}
					if sa.CompareSubstr(i, j, k, l) != compareSubtrNaive(s, i, j, k, l) {
						fmt.Println("compare substr err!!", s[i:j], s[k:l], i, j, k, l)
					}
				}
			}
		}
	}
	fmt.Println("ok")

	time1 := time.Now()
	n = int(2e6)
	NewSuffixArray(make([]int, n))
	time2 := time.Now()
	fmt.Println((time2.Sub(time1)))

}

type SuffixArray struct {
	Sa     []int // 排名第i的后缀是谁.
	Rank   []int // 后缀s[i:]的排名是多少.
	Height []int // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	n      int
	st     *St
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray(ords []int) *SuffixArray {
	res := &SuffixArray{n: len(ords)}
	sa, rank, lcp := res._useSA(ords)
	res.Sa, res.Rank, res.Height = sa, rank, lcp
	return res
}

func NewSuffixArrayWithString(s string) *SuffixArray {
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	return NewSuffixArray(ords)
}

// 求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray) Lcp(a, b int, c, d int) int {

	cand := suf._lcp(a, c)
	return min(cand, min(b-a, d-c))
}

// 比较任意两个子串s[a,b)和s[c,d)的字典序.
//  s[a,b) < s[c,d) 返回-1.
//  s[a,b) = s[c,d) 返回0.
//  s[a,b) > s[c,d) 返回1.
func (suf *SuffixArray) CompareSubstr(a, b int, c, d int) int {
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

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *SuffixArray) _lcp(i, j int) int {
	if suf.st == nil {
		suf.st = NewStMin(suf.Height)
	}
	if i == j {
		return suf.n - i
	}
	r1, r2 := suf.Rank[i], suf.Rank[j]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return suf.st.Query(r1+1, r2+1)
}

func (suf *SuffixArray) _getSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}

	mn := mins(ords...)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords...) + 1
	isS := make([]bool, n)
	isLms := make([]bool, n)
	lms := make([]int, 0, n)
	for i := 0; i < n; i++ {
		isS[i] = true
	}
	for i := n - 2; i > -1; i-- {
		if ords[i] == ords[i+1] {
			isS[i] = isS[i+1]
		} else {
			isS[i] = ords[i] < ords[i+1]
		}
	}
	for i := 1; i < n; i++ {
		isLms[i] = !isS[i-1] && isS[i]
	}
	for i := 0; i < n; i++ {
		if isLms[i] {
			lms = append(lms, i)
		}
	}
	bin := make([]int, m)
	for _, x := range ords {
		bin[x]++
	}

	induce := func() []int {
		sa := make([]int, n)
		for i := 0; i < n; i++ {
			sa[i] = -1
		}

		saIdx := make([]int, m)
		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := len(lms) - 1; j > -1; j-- {
			i := lms[j]
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		copy(saIdx, bin)
		s := 0
		for i := 0; i < m; i++ {
			s, saIdx[i] = s+saIdx[i], s
		}
		for j := 0; j < n; j++ {
			i := sa[j] - 1
			if i < 0 || isS[i] {
				continue
			}
			x := ords[i]
			sa[saIdx[x]] = i
			saIdx[x]++
		}

		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := n - 1; j > -1; j-- {
			i := sa[j] - 1
			if i < 0 || !isS[i] {
				continue
			}
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		return sa
	}

	sa = induce()

	lmsIdx := make([]int, 0, len(sa))
	for _, i := range sa {
		if isLms[i] {
			lmsIdx = append(lmsIdx, i)
		}
	}
	l := len(lmsIdx)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = -1
	}
	ord := 0
	order[n-1] = ord
	for i := 0; i < l-1; i++ {
		j, k := lmsIdx[i], lmsIdx[i+1]
		for d := 0; d < n; d++ {
			jIsLms, kIsLms := isLms[j+d], isLms[k+d]
			if ords[j+d] != ords[k+d] || jIsLms != kIsLms {
				ord++
				break
			}
			if d > 0 && (jIsLms || kIsLms) {
				break
			}
		}
		order[k] = ord
	}
	b := make([]int, 0, l)
	for _, i := range order {
		if i >= 0 {
			b = append(b, i)
		}
	}
	var lmsOrder []int
	if ord == l-1 {
		lmsOrder = make([]int, l)
		for i, ord := range b {
			lmsOrder[ord] = i
		}
	} else {
		lmsOrder = suf._getSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

func (suf *SuffixArray) _useSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = suf._getSA(ords)

	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 lcp 也就是排名相邻的两个后缀的最长公共前缀。
	// lcp[0] = 0
	// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	lcp = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

type St struct {
	st     []int
	lookup []int
	n      int
}

func NewStMin(nums []int) *St {
	res := &St{}
	n := len(nums)
	b := bits.Len(uint(n))
	st := make([]int, b*n)
	for i := range nums {
		st[i] = nums[i]
	}
	for i := 1; i < b; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			st[i*n+j] = min(st[(i-1)*n+j], st[(i-1)*n+j+(1<<(i-1))])
		}
	}
	lookup := make([]int, n+1)
	for i := 2; i < len(lookup); i++ {
		lookup[i] = lookup[i>>1] + 1
	}
	res.st = st
	res.lookup = lookup
	res.n = n
	return res
}

func (st *St) Query(start, end int) int {
	b := st.lookup[end-start]
	return min(st.st[b*st.n+start], st.st[b*st.n+end-(1<<b)])
}

func mins(a ...int) int {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs(a ...int) int {
	mx := a[0]
	for _, x := range a {
		if x > mx {
			mx = x
		}
	}
	return mx
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

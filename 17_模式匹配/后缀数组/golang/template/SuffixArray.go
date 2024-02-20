// https://github.com/EndlessCheng/codeforces-go/blob/646deb927bbe089f60fc0e9f43d1729a97399e5f/copypasta/strings.go#L556
// https://visualgo.net/zh/suffixarray
// 常用分隔符 #(35) $(36) _(95) |(124)
// SA-IS 与 DC3 的效率对比 https://riteme.site/blog/2016-6-19/sais.html#5
// 注：Go1.13 开始使用 SA-IS 算法
//
// api:
// func GetSA(ords []int) (sa []int)
// func UseSA(ords []int) (sa, rank, lcp []int)

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"math/bits"
)

func main() {
	// NumberofSubstrings()
	testLcpRange()
}

// G3. Good Substrings
// https://codeforces.com/contest/316/submission/218670841
func CF316() {
	// itoa
}

// https://codeforces.com/contest/126/submission/227749650
func CF126() {

}

// !不同子串长度之和
// 枚举每个后缀，计算前缀总数，再减掉重复
func diffSum(s string) int {
	n := len(s)
	ords := make([]int, n)
	for i, c := range s {
		ords[i] = int(c)
	}
	_, _, height := UseSA(ords)
	res := n * (n + 1) * (n + 2) / 6 // 所有子串长度1到n的平方和
	for _, h := range height {
		res -= h * (h + 1) / 2
	}
	return res
}

// 1044. 最长重复子串
// https://leetcode.cn/problems/longest-duplicate-substring/description/
// 给你一个字符串 s ，考虑其所有 重复子串 ：即 s 的（连续）子串，在 s 中出现 2 次或更多次。这些出现之间可能存在重叠。
// 返回 任意一个 可能具有最长长度的重复子串。如果 s 不含重复子串，那么答案为 "" 。
// 子串就是后缀的前缀
// !高度数组中的最大值对应的就是可重叠最长重复子串
func longestDupSubstring(s string) string {
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	S := NewSuffixArray(ords)
	sa, height := S.Sa, S.Height
	start, max := 0, 0
	for i, h := range height {
		if h > max {
			max = h
			start = i
		}
	}
	return s[sa[start] : int(sa[start])+max]
}

// https://leetcode.cn/problems/largest-merge-of-two-strings/
// 1754. 构造字典序最大的合并字符串
func largestMerge(word1 string, word2 string) string {
	ords1, ords2 := make([]int, len(word1)), make([]int, len(word2))
	for i, c := range word1 {
		ords1[i] = int(c)
	}
	for i, c := range word2 {
		ords2[i] = int(c)
	}
	S := NewSuffixArray2(ords1, ords2)

	n1, n2 := len(word1), len(word2)
	sb := strings.Builder{}

	i, j := 0, 0
	for i < len(word1) && j < len(word2) {
		if S.CompareSubstr(i, n1, j, n2) == 1 {
			sb.WriteByte(word1[i])
			i++
		} else {
			sb.WriteByte(word2[j])
			j++
		}
	}

	sb.WriteString(word1[i:])
	sb.WriteString(word2[j:])

	return sb.String()
}

// 2261. 含最多 K 个可整除元素的子数组
// https://leetcode.cn/problems/k-divisible-elements-subarrays/
// 找出并返回满足要求的不同的子数组数，要求子数组中最多 k 个可被 p 整除的元素。
func countDistinct(nums []int, k int, p int) (res int) {
	n := len(nums)

	mods := make([]int, n)
	for i := range mods {
		mods[i] = nums[i] % p
	}

	boolToInt := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	// 1. 先用双指针O(n)的时间计算出所有满足条件的子数组的数量 注意要枚举后缀(固定left 移动right)
	right, countK := 0, 0
	suffixLen := make([]int, n) // 记录每个后缀取到的长度
	for left := 0; left < n; left++ {
		for right < n && countK+boolToInt((mods[right] == 0)) <= k {
			countK += boolToInt((mods[right] == 0))
			right++
		}
		res += right - left
		suffixLen[left] = right - left
		countK -= boolToInt(mods[left] == 0)
	}

	// 2. height数组去重
	sa, _, height := UseSA(nums)
	// 计算子串重复数量 按后缀排序的顺序枚举后缀 lcp(height)去重
	for i := 0; i < n-1; i++ {
		suffix1, suffix2 := sa[i], sa[i+1]
		subLen1, subLen2 := suffixLen[suffix1], suffixLen[suffix2]
		res -= min(height[i+1], min(subLen1, subLen2))
	}
	return
}

// https://judge.yosupo.jp/problem/number_of_substrings
// 返回 s 的不同子字符串的个数(本质不同子串数)
// 用所有子串的个数，减去相同子串的个数，就可以得到不同子串的个数。
// !子串就是后缀的前缀 按后缀排序的顺序枚举后缀，每次新增的子串就是除了与上一个后缀的 LCP 剩下的前缀
// !计算后缀数组和高度数组。根据高度数组的定义，所有高度之和就是相同子串的个数。(每一对相同子串在高度数组产生1贡献)
func NumberofSubstrings() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	ords := make([]int, n)
	for i, c := range s {
		ords[i] = int(c)
	}
	res := n * (n + 1) / 2
	_, _, height := UseSA(ords)
	for _, h := range height {
		res -= h
	}
	fmt.Fprintln(out, res)
}

func testLcpRange() {
	n := int(1e3)
	ords := make([]int, n)
	for i := 1; i < n; i++ {
		ords[i] = i * i
		ords[i] ^= ords[i-1]
	}

	S := NewSuffixArray(ords)
	S2 := NewSuffixArray(ords)
	LcpRange2 := func(left, k int) (start, end int) {
		curRank := S2.Rank[left]
		for i := curRank; i >= 0; i-- {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				start = i
			} else {
				break
			}
		}
		for i := curRank; i < n; i++ {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				end = i + 1
			} else {
				break
			}
		}
		if start == 0 && end == 0 {
			return -1, -1
		}
		return
	}

	fmt.Println(S.LcpRange(4, 0))
	fmt.Println(LcpRange2(4, 0))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			start1, end1 := S.LcpRange(i, j)
			start2, end2 := LcpRange2(i, j)
			if start1 != start2 || end1 != end2 {
				fmt.Println(i, j, start1, end1, start2, end2)
				panic("")
			}
		}
	}
	fmt.Println("pass")
}

func demo() {
	s := "abca"
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	sa, rank, height := UseSA(ords)
	fmt.Println(sa, rank, height)
}

type SuffixArray struct {
	Sa      []int // 排名第i的后缀是谁.
	Rank    []int // 后缀s[i:]的排名是多少.
	Height  []int // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	Ords    []int
	n       int
	minSt32 *St32 // 维护lcp的最小值
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray(ords []int) *SuffixArray {
	ords = append(ords[:0:0], ords...)
	res := &SuffixArray{n: len(ords), Ords: ords}
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
//
//	s[a,b) < s[c,d) 返回-1.
//	s[a,b) = s[c,d) 返回0.
//	s[a,b) > s[c,d) 返回1.
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

// 与 s[left:] 的 lcp 大于等于 k 的后缀数组(sa)上的区间.
// 如果不存在,返回(-1,-1).
func (suf *SuffixArray) LcpRange(left int, k int) (start, end int) {
	if k > suf.n-left {
		return -1, -1
	}
	if k == 0 {
		return 0, suf.n
	}
	if suf.minSt32 == nil {
		suf.minSt32 = NewStMin(suf.Height)
	}
	i := suf.Rank[left] + 1
	start = suf.minSt32.MinLeft(i, func(e int) bool { return e >= k }) - 1 // 向左找
	end = suf.minSt32.MaxRight(i, func(e int) bool { return e >= k })      // 向右找
	return
}

func (suf *SuffixArray) Print(sa, ords []int) {
	n := len(ords)
	for _, v := range sa {
		s := make([]string, 0, n-v)
		for i := v; i < n; i++ {
			s = append(s, string(ords[i]))
		}
		fmt.Println(strings.Join(s, ""))
	}
}

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *SuffixArray) _lcp(i, j int) int {
	if suf.minSt32 == nil {
		suf.minSt32 = NewStMin(suf.Height)
	}
	if i == j {
		return suf.n - i
	}
	r1, r2 := suf.Rank[i], suf.Rank[j]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return suf.minSt32.Query(r1+1, r2+1)
}

func (suf *SuffixArray) _getSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}
	mn := mins(ords)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords) + 1
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

type St32 struct {
	st     []int32
	lookup []int32
	n      int
}

func NewStMin(nums []int) *St32 {
	res := &St32{}
	n := len(nums)
	b := bits.Len(uint(n))
	st := make([]int32, b*n)
	for i := range nums {
		st[i] = int32(nums[i])
	}
	for i := 1; i < b; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			st[i*n+j] = min32(st[(i-1)*n+j], st[(i-1)*n+j+(1<<(i-1))])
		}
	}
	lookup := make([]int32, n+1)
	for i := 2; i < len(lookup); i++ {
		lookup[i] = lookup[i>>1] + 1
	}
	res.st = st
	res.lookup = lookup
	res.n = n
	return res
}

func (st *St32) Query(start, end int) int {
	b := int(st.lookup[end-start])
	return int(min32(st.st[b*st.n+start], st.st[b*st.n+end-(1<<b)]))
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (st *St32) MaxRight(left int, check func(e int) bool) int {
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
func (st *St32) MinLeft(right int, check func(e int) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
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

// 用于求解`两个字符串s和t`相关性质的后缀数组.
type SuffixArray2 struct {
	SA     *SuffixArray
	offset int
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray2(ords1, ords2 []int) *SuffixArray2 {
	newNums := append(ords1, ords2...)
	sa := NewSuffixArray(newNums)
	return &SuffixArray2{SA: sa, offset: len(ords1)}
}

func NewSuffixArray2WithString(s, t string) *SuffixArray2 {
	ords1 := make([]int, len(s))
	for i, c := range s {
		ords1[i] = int(c)
	}
	ords2 := make([]int, len(t))
	for i, c := range t {
		ords2[i] = int(c)
	}
	return NewSuffixArray2(ords1, ords2)
}

// 求任意两个子串s[a,b)和t[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray2) Lcp(a, b int, c, d int) int {
	return suf.SA.Lcp(a, b, c+suf.offset, d+suf.offset)
}

// 比较任意两个子串s[a,b)和t[c,d)的字典序.
//
//	s[a,b) < t[c,d) 返回-1.
//	s[a,b) = t[c,d) 返回0.
//	s[a,b) > t[c,d) 返回1.
func (suf *SuffixArray2) CompareSubstr(a, b int, c, d int) int {
	return suf.SA.CompareSubstr(a, b, c+suf.offset, d+suf.offset)
}

// !注意内部会修改ords.
//
//	 sa : 排第几的后缀是谁.
//	 rank : 每个后缀排第几.
//	 lcp : 排名相邻的两个后缀的最长公共前缀.
//		lcp[0] = 0
//		lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
func UseSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = GetSA(ords)

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
			for j := sa[rk-1]; i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

// 注意内部会修改ords.
func GetSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}

	mn := mins(ords)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords) + 1
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
		lmsOrder = GetSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

func mins(a []int) int {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs(a []int) int {
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

func min32(a, b int32) int32 {
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

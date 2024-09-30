package main

import (
	"fmt"
	"math/bits"
	"sort"
)

func minStartingIndex(s string, pattern string) int {
	n1, n2 := len(s), len(pattern)
	if n1 < n2 {
		return -1
	}

	S := NewSuffixArray2FromString(s, pattern)
	check := func(start, end int) bool {
		failCount := 0
		pos := start
		for pos < end {
			// 移动到下一个待匹配位置继续匹配
			if s[pos] != pattern[pos-start] {
				failCount++
				if failCount > 1 {
					return false
				}
				pos++
			} else {
				pos += S.Lcp(pos, end, n2-(end-pos), n2)
			}
		}
		return true
	}

	for i := 0; i+n2 <= n1; i++ {
		if check(i, i+n2) {
			return i
		}
	}
	return -1
}

type SuffixArray struct {
	Sa     []int // 排名第i的后缀是谁.
	Rank   []int // 后缀s[i:]的排名是多少.
	Height []int // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	Ords   []int
	n      int
	minSt  *LinearRMQ // 维护lcp的最小值
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

func NewSuffixArrayFromString(s string) *SuffixArray {
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	return NewSuffixArray(ords)
}

// 求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray) Lcp(a, b int, c, d int) int {
	if a >= b || c >= d {
		return 0
	}
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
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ(suf.Height)
	}
	i := suf.Rank[left] + 1
	start = suf.minSt.MinLeft(i, func(e int) bool { return e >= k }) - 1 // 向左找
	end = suf.minSt.MaxRight(i, func(e int) bool { return e >= k })      // 向右找
	return
}

// 查询s[start:end)在s中的出现次数.
func (suf *SuffixArray) Count(start, end int) int {
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

// 返回s在原串中所有匹配的位置(无序).
// O(len(s)*log(n))+len(result).
func (suf *SuffixArray) Lookup(target []int) (result []int) {
	sa, cur := suf.Sa, suf.Ords
	// find matching suffix index range [i:j]
	// find the first index where s would be the prefix
	i := sort.Search(len(sa), func(i int) bool {
		return suf._compareSlice(cur[sa[i]:], target) >= 0
	})
	// starting at i, find the first index at which s is not a prefix
	j := i + sort.Search(len(sa)-i, func(j int) bool {
		return !suf._hasPrefix(cur[sa[i+j]:], target)
	})
	result = make([]int, j-i)
	for k := range result {
		result[k] = sa[i+k]
	}
	return
}

func (suf *SuffixArray) Print(n int, f func(i int) int, sa []int) {
	for _, v := range sa {
		s := make([]int, 0, n-v)
		for i := v; i < n; i++ {
			s = append(s, f(i))
		}
		fmt.Println(s)
	}
}

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *SuffixArray) _lcp(i, j int) int {
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ(suf.Height)
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

// s 是否以prefix为前缀.
func (suf *SuffixArray) _hasPrefix(s []int, prefix []int) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i, v := range prefix {
		if s[i] != v {
			return false
		}
	}
	return true
}

func (suf *SuffixArray) _compareSlice(a, b []int) int {
	n1, n2 := len(a), len(b)
	ptr1, ptr2 := 0, 0
	for ptr1 < n1 && ptr2 < n2 {
		if a[ptr1] < b[ptr2] {
			return -1
		}
		if a[ptr1] > b[ptr2] {
			return 1
		}
		ptr1++
		ptr2++
	}
	if ptr1 == n1 && ptr2 == n2 {
		return 0
	}
	if ptr1 == n1 {
		return -1
	}
	return 1
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

type LinearRMQ struct {
	n     int
	nums  []int
	small []int
	large [][]int
}

// n: 序列长度.
// less: 入参为两个索引,返回值表示索引i处的值是否小于索引j处的值.
//
//	消除了泛型.
func NewLinearRMQ(nums []int) *LinearRMQ {
	n := len(nums)
	res := &LinearRMQ{n: n, nums: nums}
	stack := make([]int, 0, 64)
	small := make([]int, 0, n)
	var large [][]int
	large = append(large, make([]int, 0, n>>6))
	for i := 0; i < n; i++ {
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

	for i := 1; (i << 1) <= n>>6; i <<= 1 {
		csz := n>>6 + 1 - (i << 1)
		v := make([]int, csz)
		for k := 0; k < csz; k++ {
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
func (rmq *LinearRMQ) Query(start, end int) int {
	if start >= end {
		panic(fmt.Sprintf("start(%d) should be less than end(%d)", start, end))
	}
	end--
	left := start>>6 + 1
	right := end >> 6
	if left < right {
		msb := bits.Len64(uint64(right-left)) - 1
		cache := rmq.large[msb]
		i := (left-1)<<6 + bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63))))
		cand1 := rmq._getMin(i, cache[left])
		j := right<<6 + bits.TrailingZeros64(uint64(rmq.small[end]))
		cand2 := rmq._getMin(cache[right-(1<<msb)], j)
		return rmq.nums[rmq._getMin(cand1, cand2)]
	}
	if left == right {
		i := (left-1)<<6 + bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63))))
		j := left<<6 + bits.TrailingZeros64(uint64(rmq.small[end]))
		return rmq.nums[rmq._getMin(i, j)]
	}
	return rmq.nums[right<<6+bits.TrailingZeros64(uint64(rmq.small[end]&(^0<<(start&63))))]
}

func (rmq *LinearRMQ) _getMin(i, j int) int {
	if rmq.nums[i] < rmq.nums[j] {
		return i
	}
	return j
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ) MaxRight(left int, check func(e int) bool) int {
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
func (st *LinearRMQ) MinLeft(right int, check func(e int) bool) int {
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

func NewSuffixArray2FromString(s, t string) *SuffixArray2 {
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

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}

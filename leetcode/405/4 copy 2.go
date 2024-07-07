package main

import (
	"fmt"
	"math/bits"
	"sort"
)

const INF int = 1e18

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

func minimumCost(target string, words []string, costs []int) int {
	dp := make([]int, len(target)+1)
	for i := 1; i <= len(target); i++ {
		dp[i] = INF
	}
	pos := int32(0)

	S := NewSuffixArrayFromString(target)

	for i, char := range target {

		for _, wordIndex := range indexes[pos] {
			wordLen := len(words[wordIndex])
			if i+1 >= wordLen {
				dp[i+1] = min(dp[i+1], dp[i+1-wordLen]+costs[wordIndex])
			}
		}
	}
	if dp[len(target)] == INF {
		return -1
	}
	return dp[len(target)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个前缀.
type ACAutoMatonArray struct {
	WordPos            []int32   // WordPos[i] 表示加入的第i个模式串对应的节点编号(单词结点).
	Parent             []int32   // parent[v] 表示节点v的父节点.
	Link               []int32   // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	Children           [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	BfsOrder           []int32   // 结点的拓扑序,0表示虚拟节点.
	sigma              int32     // 字符集大小.
	offset             int32     // 字符集的偏移量.
	needUpdateChildren bool      // 是否需要更新children数组.
}

func NewACAutoMatonArray(sigma, offset int32) *ACAutoMatonArray {
	res := &ACAutoMatonArray{sigma: sigma, offset: offset}
	res.Clear()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *ACAutoMatonArray) AddString(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, s := range str {
		ord := s - trie.offset
		if trie.Children[pos][ord] == -1 {
			trie.Children[pos][ord] = trie.newNode()
			trie.Parent[len(trie.Parent)-1] = pos
		}
		pos = trie.Children[pos][ord]
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// 在pos位置添加一个字符，返回新的节点编号.
func (trie *ACAutoMatonArray) AddChar(pos, ord int32) int32 {
	ord -= trie.offset
	if trie.Children[pos][ord] != -1 {
		return trie.Children[pos][ord]
	}
	trie.Children[pos][ord] = trie.newNode()
	trie.Parent[len(trie.Parent)-1] = pos
	return trie.Children[pos][ord]
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray) Move(pos, ord int32) int32 {
	ord -= trie.offset
	if trie.needUpdateChildren {
		return trie.Children[pos][ord]
	}
	for {
		nexts := trie.Children[pos]
		if nexts[ord] != -1 {
			return nexts[ord]
		}
		if pos == 0 {
			return 0
		}
		pos = trie.Link[pos]
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray) Size() int32 {
	return int32(len(trie.Children))
}

func (trie *ACAutoMatonArray) Empty() bool {
	return len(trie.Children) == 1
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *ACAutoMatonArray) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.Link = make([]int32, len(trie.Children))
	for i := range trie.Link {
		trie.Link[i] = -1
	}
	trie.BfsOrder = make([]int32, len(trie.Children))
	head, tail := 0, 0
	trie.BfsOrder[tail] = 0
	tail++
	for head < tail {
		v := trie.BfsOrder[head]
		head++
		for i, next := range trie.Children[v] {
			if next == -1 {
				continue
			}
			trie.BfsOrder[tail] = next
			tail++
			f := trie.Link[v]
			for f != -1 && trie.Children[f][i] == -1 {
				f = trie.Link[f]
			}
			trie.Link[next] = f
			if f == -1 {
				trie.Link[next] = 0
			} else {
				trie.Link[next] = trie.Children[f][i]
			}
		}
	}
	if !needUpdateChildren {
		return
	}
	for _, v := range trie.BfsOrder {
		for i, next := range trie.Children[v] {
			if next == -1 {
				f := trie.Link[v]
				if f == -1 {
					trie.Children[v][i] = 0
				} else {
					trie.Children[v][i] = trie.Children[f][i]
				}
			}
		}
	}
}

func (trie *ACAutoMatonArray) Clear() {
	trie.WordPos = trie.WordPos[:0]
	trie.Parent = trie.Parent[:0]
	trie.Children = trie.Children[:0]
	trie.Link = trie.Link[:0]
	trie.BfsOrder = trie.BfsOrder[:0]
	trie.newNode()
}

// 获取每个状态包含的模式串的个数.
func (trie *ACAutoMatonArray) GetCounter() []int32 {
	counter := make([]int32, len(trie.Children))
	for _, pos := range trie.WordPos {
		counter[pos]++
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			counter[v] += counter[trie.Link[v]]
		}
	}
	return counter
}

// 获取每个状态包含的模式串的索引.(模式串长度和较小时使用)
// fail指针每次命中，都至少有一个比指针深度更长的单词出现，因此每个位置最坏情况下不超过O(sqrt(n))次命中
// O(n*sqrt(n))
func (trie *ACAutoMatonArray) GetIndexes() [][]int32 {
	res := make([][]int32, len(trie.Children))
	for i, pos := range trie.WordPos {
		res[pos] = append(res[pos], int32(i))
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			from, to := trie.Link[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int32, len(arr1)+len(arr2))
			i, j, k := 0, 0, 0
			for i < len(arr1) && j < len(arr2) {
				if arr1[i] < arr2[j] {
					arr3[k] = arr1[i]
					i++
				} else if arr1[i] > arr2[j] {
					arr3[k] = arr2[j]
					j++
				} else {
					arr3[k] = arr1[i]
					i++
					j++
				}
				k++
			}
			copy(arr3[k:], arr1[i:])
			k += len(arr1) - i
			copy(arr3[k:], arr2[j:])
			k += len(arr2) - j
			arr3 = arr3[:k:k]
			res[to] = arr3
		}
	}
	return res
}

// 按照拓扑序进行转移(EnumerateFail).
func (trie *ACAutoMatonArray) Dp(f func(from, to int32)) {
	for _, v := range trie.BfsOrder {
		if v != 0 {
			f(trie.Link[v], v)
		}
	}
}

func (trie *ACAutoMatonArray) BuildFailTree() [][]int32 {
	res := make([][]int32, trie.Size())
	trie.Dp(func(pre, cur int32) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (trie *ACAutoMatonArray) BuildTrieTree() [][]int32 {
	res := make([][]int32, trie.Size())
	for i := int32(1); i < trie.Size(); i++ {
		res[trie.Parent[i]] = append(res[trie.Parent[i]], i)
	}
	return res
}

// 返回str在trie树上的节点位置.如果不存在，返回0.
func (trie *ACAutoMatonArray) Search(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, char := range str {
		if pos >= int32(len(trie.Children)) || pos < 0 {
			return 0
		}
		ord := char - trie.offset
		if next := trie.Children[pos][ord]; next == -1 {
			return 0
		} else {
			pos = next
		}
	}
	return pos
}

func (trie *ACAutoMatonArray) newNode() int32 {
	trie.Parent = append(trie.Parent, -1)
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.Children = append(trie.Children, nexts)
	return int32(len(trie.Children) - 1)
}

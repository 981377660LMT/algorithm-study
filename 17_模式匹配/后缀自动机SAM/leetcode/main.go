// 100251. 数组中的最短非公共子字符串
// 给你一个数组 arr ，数组中有 n 个 非空 字符串。
// 请你求出一个长度为 n 的字符串 answer ，满足：
// answer[i] 是 arr[i] 最短 的子字符串，且它不是 arr 中其他任何字符串的子字符串。如果有多个这样的子字符串存在，answer[i] 应该是它们中字典序最小的一个。如果不存在这样的子字符串，answer[i] 为空字符串。
// 请你返回数组 answer 。
// https://leetcode.cn/problems/shortest-uncommon-substring-in-an-array/description/
//
// 1. 要处理多个子串的匹配问题，考虑将所有单词建立广义后缀自动机；
// 2. 为了处理子串包含关系，对每个单词arr[i] 的每个前缀，
//    在后缀链接树上上跳标记每个endPos，统计每个 endPos 等价类属于多少个单词，
//    不是 arr 中其他任何字符串的子字符串 等价于 endPos 恰好属于一个单词；
// 3. 对每个单词 arr[i] 的前缀，在后缀链接树上上跳更新合法的 endPos 的最小长度，
//    并记录对应的结点编号以及最小长度的子串在单词 arr[i] 中的结束下标 end，然后采用后缀数组求出字典序最小的那一个子串。

package main

import (
	"fmt"
	"index/suffixarray"
	"math/bits"
	"reflect"
	"unsafe"
)

func main() {

}

const INF int32 = 1e9 + 10

// 100251. 数组中的最短非公共子字符串
// https://leetcode.cn/problems/shortest-uncommon-substring-in-an-array/description/
func shortestSubstrings(arr []string) []string {
	sam := NewSuffixAutomatonGeneral()
	for _, v := range arr {
		sam.AddString(v)
	}

	size := sam.Size()
	nodes := sam.Nodes
	belong := make([][]int32, size)
	belongEnd := make([][]int32, size)
	visitedTime := make([]int32, size)
	for i := range visitedTime {
		visitedTime[i] = -1
	}

	// 对arr[i]的每个前缀，在后缀链接树上向上跳，标记每个endPos，表示该endPos包含了arr[i]的子串.
	// 标记次数之和不超过O(∑sqrt(∑)).
	markChain := func(sid, pos, endIndex int32) {
		for pos >= 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			if len(belong[pos]) >= 2 {
				break
			}
			belong[pos] = append(belong[pos], sid)
			belongEnd[pos] = append(belongEnd[pos], endIndex)
			pos = nodes[pos].Link
		}
	}

	// 标记所有串的子串.
	for i, w := range arr {
		pos := int32(0)
		for j, c := range w {
			pos = nodes[pos].Next[c-OFFSET]
			markChain(int32(i), pos, int32(j))
		}
	}

	res := make([]string, len(arr))
	for i := range visitedTime {
		visitedTime[i] = -1
	}
	for si, w := range arr {
		minLen, pos, end := INF, []int32{}, []int32{} // 最短的子串长度，对应的结点编号，在w中的结束下标.
		queryChain := func(si, p int32) {
			for p >= 0 && visitedTime[p] != si {
				visitedTime[p] = si
				if len(belong[p]) >= 2 {
					break
				}
				if len(belong[p]) == 1 { // 该结点只包含了当前串
					curMinLen := nodes[nodes[p].Link].MaxLen + 1
					if curMinLen < minLen {
						minLen = curMinLen
						pos = []int32{p}
						end = []int32{belongEnd[p][0]}
					} else if curMinLen == minLen {
						pos = append(pos, p)
						end = append(end, belongEnd[p][0])
					}
				}
				p = nodes[p].Link
			}
		}

		p := int32(0)
		for _, c := range w {
			p = nodes[p].Next[c-OFFSET]
			queryChain(int32(si), p)
		}

		// 后缀数组比较各个子串的字典序.
		if len(pos) > 0 {
			SA := NewSuffixArray(int32(len(w)), func(i int32) int32 { return int32(w[i]) })
			bestEnd := end[0]
			for i := 1; i < len(end); i++ {
				a, b := end[i]+1-minLen, end[i]+1
				c, d := bestEnd+1-minLen, bestEnd+1
				if SA.CompareSubstr(a, b, c, d) < 0 {
					bestEnd = end[i]
				}
			}
			res[si] = w[bestEnd+1-minLen : bestEnd+1]
		}
	}

	return res

}

const SIGMA int32 = 26   // 字符集大小
const OFFSET int32 = 'a' // 字符集的起始字符

type Node struct {
	Next   [SIGMA]int32 // SAM 转移边
	Link   int32        // 后缀链接
	MaxLen int32        // 当前节点对应的最长子串的长度
}

type SuffixAutomaton struct {
	Nodes []*Node
	n     int32 // 当前字符串长度
}

func NewSuffixAutomatonGeneral() *SuffixAutomaton {
	res := &SuffixAutomaton{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0))
	return res
}

// !需要在插入新串之前将lastPos置为0.
// eg:
//
//	sam := NewSuffixAutomatonGeneral()
//	for _,word := range words {
//	  lastPos = 0
//	  for _,c := range word {
//	    lastPos = sam.Add(lastPos,c)
//	  }
//	}
//
// 返回当前前缀对应的节点编号(lastPos).
func (sam *SuffixAutomaton) Add(lastPos int32, char int32) int32 {
	c := char - OFFSET
	sam.n++

	// 判断当前转移结点是否存在.
	if tmp := sam.Nodes[lastPos].Next[c]; tmp != -1 {
		lastNode, nextNode := sam.Nodes[lastPos], sam.Nodes[tmp]
		if lastNode.MaxLen+1 == nextNode.MaxLen {
			return tmp
		} else {
			newQ := int32(len(sam.Nodes))
			sam.Nodes = append(sam.Nodes, sam.newNode(nextNode.Link, lastNode.MaxLen+1))
			sam.Nodes[newQ].Next = nextNode.Next
			sam.Nodes[tmp].Link = newQ
			for lastPos != -1 && sam.Nodes[lastPos].Next[c] == tmp {
				sam.Nodes[lastPos].Next[c] = newQ
				lastPos = sam.Nodes[lastPos].Link
			}
			return newQ
		}
	}

	newNode := int32(len(sam.Nodes))
	// 新增一个实点以表示当前最长串
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[lastPos].MaxLen+1))
	p := lastPos
	for p != -1 && sam.Nodes[p].Next[c] == -1 {
		sam.Nodes[p].Next[c] = newNode
		p = sam.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sam.Nodes[p].Next[c]
	}
	if p == -1 || sam.Nodes[p].MaxLen+1 == sam.Nodes[q].MaxLen {
		sam.Nodes[newNode].Link = q
	} else {
		// 不够用，需要新增一个虚点
		newQ := int32(len(sam.Nodes))
		sam.Nodes = append(sam.Nodes, sam.newNode(sam.Nodes[q].Link, sam.Nodes[p].MaxLen+1))
		sam.Nodes[len(sam.Nodes)-1].Next = sam.Nodes[q].Next
		sam.Nodes[q].Link = newQ
		sam.Nodes[newNode].Link = newQ
		for p != -1 && sam.Nodes[p].Next[c] == q {
			sam.Nodes[p].Next[c] = newQ
			p = sam.Nodes[p].Link
		}
	}
	return newNode
}

func (sam *SuffixAutomaton) AddString(s string) (lastPos int32) {
	lastPos = 0
	for _, c := range s {
		lastPos = sam.Add(lastPos, c)
	}
	return
}

func (sam *SuffixAutomaton) Size() int32 {
	return int32(len(sam.Nodes))
}

func (sam *SuffixAutomaton) Print() {
	var dfs func(int32, string)
	nodeStr := make([]string, sam.Size())
	dfs = func(cur int32, s string) {
		nodeStr[cur] = s
		for i, v := range sam.Nodes[cur].Next {
			if v != -1 {
				dfs(v, s+string(OFFSET+int32(i)))
			}
		}
	}
	dfs(0, "")
	for i := int32(0); i < sam.Size(); i++ {
		fmt.Println(i, nodeStr[i])
	}
}

func (sam *SuffixAutomaton) newNode(link, maxLen int32) *Node {
	res := &Node{Link: link, MaxLen: maxLen}
	for i := int32(0); i < SIGMA; i++ {
		res.Next[i] = -1
	}
	return res
}

type SuffixArray32 struct {
	Sa     []int32 // 排名第i的后缀是谁.
	Rank   []int32 // 后缀s[i:]的排名是多少.
	Height []int32 // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
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
	if start >= end {
		panic(fmt.Sprintf("start(%d) should be less than end(%d)", start, end))
	}
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

func abs32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
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

func mins32(a []int32) int32 {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs32(a []int32) int32 {
	mx := a[0]
	for _, x := range a {
		if x > mx {
			mx = x
		}
	}
	return mx
}

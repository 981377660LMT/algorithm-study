package main

func main() {
	S := NewSubsequnceAutomatonArray(5, func(i int) byte { return "abcde"[i] })
	println(S.Match(0, 5, 3, func(i int) byte { return "ace"[i] }))
}

// 727. 最小窗口子序列
// https://leetcode.cn/problems/minimum-window-subsequence/description/
func minWindow(s1 string, s2 string) string {
	S := NewSubsequnceAutomatonArray(len(s1), func(i int) byte { return s1[i] })
	var starts []int
	for i := 0; i < len(s1); i++ {
		if s1[i] == s2[0] {
			starts = append(starts, i)
		}
	}
	res0, res1 := -1, -1
	for _, sStart := range starts {
		hit, sEnd := S.Match(sStart, len(s1), len(s2), func(i int) byte { return s2[i] })
		if hit != len(s2) {
			continue
		}
		sLen := sEnd - sStart
		if res0 == -1 || sLen < res1-res0 {
			res0, res1 = sStart, sEnd
		}
	}
	if res0 == -1 {
		return ""
	}
	return s1[res0:res1]
}

// 792. 匹配子序列的单词数
// https://leetcode.cn/problems/number-of-matching-subsequences/description/
func numMatchingSubseq(s string, words []string) int {
	n := int(len(s))
	S := NewSubsequnceAutomatonArray(n, func(i int) byte { return s[i] })
	res := 0
	for _, w := range words {
		if S.Includes(0, n, len(w), func(i int) byte { return w[i] }) {
			res++
		}
	}
	return res
}

const SIGMA byte = 26
const OFFSET byte = 97

type SubsequnceAutomatonArray struct {
	n     int
	s     func(i int) byte
	nexts [][SIGMA]int32
}

// `O(∑*n)`预处理,`∑`为字符集大小.
// `O(len(t))`查询,`len(t)`为待匹配序列的长度.
func NewSubsequnceAutomatonArray(n int, s func(i int) byte) *SubsequnceAutomatonArray {
	n32 := int32(n)
	nexts := make([][SIGMA]int32, n32)
	last := [SIGMA]int32{}
	for i := range last {
		last[i] = n32
	}
	for i := n - 1; i >= 0; i-- {
		nexts[i] = last
		last[s(i)-OFFSET] = int32(i)
	}
	return &SubsequnceAutomatonArray{n: n, s: s, nexts: nexts}
}

// 查询当前位置的下一个特定字符的位置(下标严格大于pos).
// 如果不存在，则为 n.
// 0<=pos<n.
func (s *SubsequnceAutomatonArray) Move(pos int, newValue byte) int {
	return int(s.nexts[pos][newValue-OFFSET])
}

// 查询`s[start:end)`内是否含有某序列`t`.
func (s *SubsequnceAutomatonArray) Includes(start, end int, tLen int, t func(i int) byte) bool {
	_, tPos := s.Match(start, end, tLen, t)
	return tPos == tLen
}

// 在`s[start:end)`内寻找子序列`t`.
// 返回`匹配到的t的长度`, `匹配结束时s的索引`
// 耗去的s的长度为`sEnd-start`.
func (s *SubsequnceAutomatonArray) Match(start, end int, tLen int, t func(i int) byte) (hit, sEnd int) {
	if start >= end || tLen == 0 {
		sEnd = start
		return
	}
	n := s.n
	si, ti := start, 0
	if s.s(si) == t(ti) {
		ti++
	}
	for si < end && ti < tLen {
		nextPos := s.Move(si, t(ti))
		if nextPos == n {
			hit, sEnd = ti, si+1
			return
		}
		si = nextPos
		ti++
	}
	hit, sEnd = ti, si+1
	return
}

type SubsequnceAutomatonMap struct {
	n       int
	s       func(i int) int
	indexes map[int][]int
}

// `O(n)`预处理.
// `O(len(t)logn)`查询,`len(t)`为待匹配序列的长度.
// !复杂度与字符种类数无关, 且占用空间更小.
func NewSubsequnceAutomatonMap(n int, s func(i int) int) *SubsequnceAutomatonMap {
	indexes := make(map[int][]int)
	for i := 0; i < n; i++ {
		v := s(i)
		indexes[v] = append(indexes[v], i)
	}
	return &SubsequnceAutomatonMap{n: n, s: s, indexes: indexes}
}

// 查询当前位置的下一个特定字符的位置(下标严格大于pos).
// 如果不存在，则为 n.
// 0<=pos<n.
func (s *SubsequnceAutomatonMap) Move(pos int, newValue int) int {
	indexes, ok := s.indexes[newValue]
	if !ok {
		return s.n
	}
	nextPos := s.bisectRight(indexes, pos)
	if nextPos < len(indexes) {
		return indexes[nextPos]
	} else {
		return s.n
	}
}

// 查询`s[start:end)`内是否含有某序列`t`.
func (s *SubsequnceAutomatonMap) Includes(start, end int, tLen int, t func(i int) int) bool {
	_, tPos := s.Match(start, end, tLen, t)
	return tPos == tLen
}

// 在`s[start:end)`内寻找子序列`t`.
// 返回 `匹配到的t的长度`, `匹配结束时s的索引`
// 耗去的s的长度为`sPos-start`.
func (s *SubsequnceAutomatonMap) Match(start, end int, tLen int, t func(i int) int) (hit int, sEnd int) {
	if start >= end || tLen == 0 {
		sEnd = start
		return
	}
	n := s.n
	si, ti := start, 0
	if s.s(si) == t(ti) {
		ti++
	}
	for si < end && ti < tLen {
		nextPos := s.Move(si, t(ti))
		if nextPos == n {
			hit, sEnd = ti, si+1
			return
		}
		si = nextPos
		ti++
	}
	hit, sEnd = ti, si+1
	return
}

func (s *SubsequnceAutomatonMap) bisectRight(arr []int, value int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] <= value {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

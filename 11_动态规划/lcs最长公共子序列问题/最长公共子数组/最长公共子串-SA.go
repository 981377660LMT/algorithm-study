// https://maspypy.github.io/library/string/longest_common_substring.hpp
// 两个字数串的最长公共子串
// 返回: (start1, end1, start2, end2)
// 0 <= start1 < end1 <= len(s1)
// 0 <= start2 < end2 <= len(s2)
// s1[start1:end1] == s2[start2:end2]

package main

import "fmt"

func main() {
	s1, s2 := "abcde", "cdeab"
	start1, end1, start2, end2 := longestCommonSubstring1(s1, s2)
	fmt.Println(s1[start1:end1], s2[start2:end2])
}

// 最长公共子串
func longestCommonSubstring1(s1, s2 string) (start1, end1, start2, end2 int) {
	ords1, ords2 := make([]int, len(s1)), make([]int, len(s2))
	for i, x := range s1 {
		ords1[i] = int(x)
	}
	for i, x := range s2 {
		ords2[i] = int(x)
	}
	return longestCommonSubstring2(ords1, ords2)
}

func longestCommonSubstring2(ords1, ords2 []int) (start1, end1, start2, end2 int) {
	if len(ords1) == 0 || len(ords2) == 0 {
		return
	}

	dummy := max(maxs(ords1...), maxs(ords2...)) + 1
	sb := make([]int, 0, len(ords1)+len(ords2)+1)
	sb = append(sb, ords1...)
	sb = append(sb, dummy)
	sb = append(sb, ords2...)
	sa, _, lcp := UseSA(sb)

	len_ := 0
	for i := 1; i < len(sb); i++ {
		if (sa[i-1] < len(ords1)) == (sa[i] < len(ords1)) {
			continue
		}
		if lcp[i] <= len_ {
			continue
		}
		len_ = lcp[i]

		// 来自s和t的不同子串
		// 找到了(严格)更长的公共子串,更新答案
		i1, i2 := sa[i-1], sa[i]
		if i1 > i2 {
			i1, i2 = i2, i1
		}

		start1 = i1
		end1 = start1 + len_
		start2 = i2 - len(ords1) - 1
		end2 = start2 + len_
	}

	return
}

func GetSA(ords []int) (sa []int) {
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
		lmsOrder = GetSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

//  sa : 排第几的后缀是谁.
//  rank : 每个后缀排第几.
//  lcp : 排名相邻的两个后缀的最长公共前缀.
// 	lcp[0] = 0
// 	lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
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
			for j := int(sa[rk-1]); i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
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

func sum(a ...int) int {
	s := 0
	for _, x := range a {
		s += x
	}
	return s
}

// https://blog.csdn.net/Elemmir/article/details/50988467
// 所有不相同的子串中字典序第k小的子串
// !二分出排名为K的子串是哪一个后缀的第几个未被计算过的前缀(每个后缀贡献子串数是这个后缀的长度减去其LCP)

package main

import (
	"fmt"
)

func main() {
	fmt.Println(KthSmallestSubstring("abab", 3)) // aba

}

type Sequence = string

// 字典序第k小的子串.k>=1.
func KthSmallestSubstring(s Sequence, k int) Sequence {
	n := len(s)
	ords := make([]int, n)
	for i, c := range s {
		ords[i] = int(c)
	}
	sa, _, height := UseSA(ords)
	counter := make([]int, n)
	for i, r := range sa {
		counter[i] = n - r - height[i]
	}
	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + counter[i]
	}

	left, right := 0, n
	for left <= right {
		mid := (left + right) >> 1
		if preSum[mid] < k {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	remain := k - preSum[left-1] // 排名的k的子串是第left个后缀的第remain个`未出现`过的前缀,需要加上lcp
	remain += height[left-1]
	start := sa[left-1]
	end := start + remain
	return s[start:end]
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

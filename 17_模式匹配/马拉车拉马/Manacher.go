package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/enumerate_palindromes
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s Sequence
	fmt.Fscan(in, &s)

	M := NewManacher(s)
	odd, even := M.GetOddRadius(), M.GetEvenRadius()
	for i := int32(0); i < M.Len(); i++ {
		fmt.Fprint(out, 2*odd[i]-1, " ")
		if i+1 < M.Len() {
			fmt.Fprint(out, 2*even[i+1], " ")
		}
	}
}

// https://leetcode.cn/problems/longest-palindromic-substring/description/
// 最长回文子串
func longestPalindrome(s string) string {
	M := NewManacher(s)
	start, maxLen := int32(0), int32(0)
	for i := int32(0); i < M.Len(); i++ {
		len1 := M.GetLongestOddStartsAt(i)
		if len1 > maxLen {
			maxLen = len1
			start = i
		}
		len2 := M.GetLongestEvenStartsAt(i)
		if len2 > maxLen {
			maxLen = len2
			start = i
		}
	}
	return s[start : start+maxLen]
}

// https://leetcode.cn/problems/palindromic-substrings/description/
// 回文子串个数
func countSubstrings(s string) int {
	M := NewManacher(s)
	oddRadius, evenRadius := M.GetOddRadius(), M.GetEvenRadius()
	res := 0
	for i := 0; i < len(s); i++ {
		res += int(oddRadius[i])
		res += int(evenRadius[i])
	}
	return int(res)
}

// https://leetcode.cn/problems/maximum-number-of-non-overlapping-palindrome-substrings/description/
// 2472. 不重叠回文子字符串的最大数目
func maxPalindromes(s string, k int) int {
	M := NewManacher(s)
	n := len(s)
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = dp[i-1]
		if i-k >= 0 && M.IsPalindrome(int32(i-k), int32(i)) {
			dp[i] = max(dp[i], dp[i-k]+1)
		}
		if i-k-1 >= 0 && M.IsPalindrome(int32(i-k-1), int32(i)) {
			dp[i] = max(dp[i], dp[i-k-1]+1)
		}
	}
	return dp[n]
}

type Sequence = string

type Manacher struct {
	n          int32
	seq        Sequence
	oddRadius  []int32
	evenRadius []int32
	maxOdd1    []int32
	maxOdd2    []int32
	maxEven1   []int32
	maxEven2   []int32
}

func NewManacher(seq Sequence) *Manacher {
	m := &Manacher{
		n:   int32(len(seq)),
		seq: seq,
	}
	return m
}

// 查询切片s[start:end]是否为回文串.
// 空串不为回文串.
func (ma *Manacher) IsPalindrome(start, end int32) bool {
	n := ma.n
	if start < 0 {
		start += n
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += n
	}
	if end > n {
		end = n
	}
	if start >= end {
		return false
	}

	len_ := end - start
	mid := (start + end) >> 1
	if len_&1 == 1 {
		return ma.GetOddRadius()[mid] >= len_>>1+1
	} else {
		return ma.GetEvenRadius()[mid] >= len_>>1
	}
}

// 获取每个中心点的奇回文半径`radius`.
// 回文为`[pos-radius+1:pos+radius]`.
func (ma *Manacher) GetOddRadius() []int32 {
	if ma.oddRadius != nil {
		return ma.oddRadius
	}
	n := ma.n
	ma.oddRadius = make([]int32, n)
	left, right := int32(0), int32(-1)
	for i := int32(0); i < n; i++ {
		var k int32
		if i > right {
			k = 1
		} else {
			k = minInt32(ma.oddRadius[left+right-i], right-i+1)
		}
		for i-k >= 0 && i+k < n && ma.seq[i-k] == ma.seq[i+k] {
			k++
		}
		ma.oddRadius[i] = k
		k--
		if i+k > right {
			left = i - k
			right = i + k
		}
	}
	return ma.oddRadius
}

// 获取每个中心点的偶回文半径`radius`.
// 回文为`[pos-radius:pos+radius]`.
func (ma *Manacher) GetEvenRadius() []int32 {
	if ma.evenRadius != nil {
		return ma.evenRadius
	}
	n := ma.n
	ma.evenRadius = make([]int32, n)
	left, right := int32(0), int32(-1)
	for i := int32(0); i < n; i++ {
		var k int32
		if i > right {
			k = 0
		} else {
			k = minInt32(ma.evenRadius[left+right-i+1], right-i+1)
		}
		for i-k-1 >= 0 && i+k < n && ma.seq[i-k-1] == ma.seq[i+k] {
			k++
		}
		ma.evenRadius[i] = k
		k--
		if i+k > right {
			left = i - k - 1
			right = i + k
		}
	}
	return ma.evenRadius
}

// 以s[index]开头的最长奇回文子串的长度.
func (ma *Manacher) GetLongestOddStartsAt(index int32) int32 {
	ma._initOdds()
	return ma.maxOdd1[index]
}

// 以s[index]结尾的最长奇回文子串的长度.
func (ma *Manacher) GetLongestOddEndsAt(index int32) int32 {
	ma._initOdds()
	return ma.maxOdd2[index]
}

// 以s[index]开头的最长偶回文子串的长度.
func (ma *Manacher) GetLongestEvenStartsAt(index int32) int32 {
	ma._initEvens()
	return ma.maxEven1[index]
}

// 以s[index]结尾的最长偶回文子串的长度.
func (ma *Manacher) GetLongestEvenEndsAt(index int32) int32 {
	ma._initEvens()
	return ma.maxEven2[index]
}

func (ma *Manacher) Len() int32 {
	return ma.n
}

func (ma *Manacher) _initOdds() {
	if ma.maxOdd1 != nil {
		return
	}
	n := ma.n
	ma.maxOdd1 = make([]int32, n)
	ma.maxOdd2 = make([]int32, n)
	for i := int32(0); i < n; i++ {
		ma.maxOdd1[i] = 1
		ma.maxOdd2[i] = 1
	}
	for i := int32(0); i < n; i++ {
		radius := ma.GetOddRadius()[i]
		start, end := i-radius+1, i+radius-1
		length := 2*radius - 1
		ma.maxOdd1[start] = maxInt32(ma.maxOdd1[start], length)
		ma.maxOdd2[end] = maxInt32(ma.maxOdd2[end], length)
	}
	for i := int32(0); i < n; i++ {
		if i-1 >= 0 {
			ma.maxOdd1[i] = maxInt32(ma.maxOdd1[i], ma.maxOdd1[i-1]-2)
		}
		if i+1 < n {
			ma.maxOdd2[i] = maxInt32(ma.maxOdd2[i], ma.maxOdd2[i+1]-2)
		}
	}
}

func (ma *Manacher) _initEvens() {
	if ma.maxEven1 != nil {
		return
	}
	n := ma.n
	ma.maxEven1 = make([]int32, n)
	ma.maxEven2 = make([]int32, n)
	for i := int32(0); i < n; i++ {
		radius := ma.GetEvenRadius()[i]
		if radius == 0 {
			continue
		}
		start := i - radius
		end := start + 2*radius - 1
		length := 2 * radius
		ma.maxEven1[start] = maxInt32(ma.maxEven1[start], length)
		ma.maxEven2[end] = maxInt32(ma.maxEven2[end], length)
	}
	for i := int32(0); i < n; i++ {
		if i-1 >= 0 {
			ma.maxEven1[i] = maxInt32(ma.maxEven1[i], ma.maxEven1[i-1]-2)
		}
		if i+1 < n {
			ma.maxEven2[i] = maxInt32(ma.maxEven2[i], ma.maxEven2[i+1]-2)
		}
	}
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

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

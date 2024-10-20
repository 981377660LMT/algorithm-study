// O(1)判断区间回文.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// words = ["abcd","dcba","lls","s","sssll"]
	fmt.Println(palindromePairs([]string{"abcd", "dcba", "lls", "s", "sssll"})) // [[0 1] [1 0] [3 2] [2 4]]
	// CF1326D2()
	// yosupo()
}

// 336. 回文对 (拼接回文串，拼接后是回文串的字符串对)
// https://leetcode.cn/problems/palindrome-pairs/description/
// 给定一个由唯一字符串构成的 0 索引 数组 words 。
// 回文对 是一对整数 (i, j) ，满足以下条件：
// 0 <= i, j < words.length，i != j ，并且words[i] + words[j]（两个字符串的连接）是一个回文串。返回一个数组，它包含 words 中所有满足 回文对 条件的字符串。
//
// 三种情况:
//  1. 长度相等
//     abc+cba
//  2. 长度不等
//     枚举较长的那个串，用它去和较短的串拼接
//     - len(s1)>len(s2)，abc(xxx) + cba，回文中心在s1上，说明s2的反串是A的前缀,且s1剩余后缀是回文串.
//     - len(s1)<len(s2)，cba + (xxx)abc, 回文中心在s2上，说明s1的反串是B的后缀,且s2剩余前缀是回文串.
//
// 枚举前缀用trie，判断回文用manacher.
func palindromePairs(words []string) (res [][]int) {
	manachers := make([]*Manacher, len(words))
	for i, word := range words {
		manachers[i] = NewManacher(word)
	}
	// 询问一个串的后缀s[start:)是否是回文串(这里空串不为回文串).
	isSuffixPalindrome := func(wordIndex int32, start int32) bool {
		n := int32(len(words[wordIndex]))
		if start >= n {
			return false
		}
		return manachers[wordIndex].IsPalindrome(start, n)
	}
	// 询问一个串的前缀s[:end)是否是回文串(这里空串不为回文串).
	isPrefixPalindrome := func(wordIndex int32, end int32) bool {
		if end <= 0 {
			return false
		}
		return manachers[wordIndex].IsPalindrome(0, end)
	}

	prefixTrie := NewTrie[byte]()
	suffixTrie := NewTrie[byte]()
	for wi, word := range words {
		n := int32(len(word))
		prefixTrie.Insert(n, func(i int32) byte { return word[i] }, int32(wi))
		suffixTrie.Insert(n, func(i int32) byte { return word[n-1-i] }, int32(wi))
	}

	// 枚举较长的那个串，用它去和较短的串拼接
	for wi, word := range words {
		m := int32(len(word))
		sameLen := suffixTrie.Search(m, func(i int32) byte { return word[i] })
		for _, id := range sameLen {
			if id != int32(wi) {
				res = append(res, []int{int(wi), int(id)})
			}
		}

		// AB+A，枚举前缀A
		root1 := suffixTrie.root
		for j := 0; j < len(word); j++ {
			char := word[j]
			next, ok := root1.Children[char]
			if !ok {
				break
			}
			root1 = next
			if len(root1.EndIndex) > 0 && isSuffixPalindrome(int32(wi), int32(j+1)) {
				for _, id := range root1.EndIndex {
					if id != int32(wi) {
						res = append(res, []int{int(wi), int(id)})
					}
				}
			}
		}

		// B+AB，枚举后缀B
		root2 := prefixTrie.root
		for j := len(word) - 1; j >= 0; j-- {
			char := word[j]
			next, ok := root2.Children[char]
			if !ok {
				break
			}
			root2 = next
			if len(root2.EndIndex) > 0 && isPrefixPalindrome(int32(wi), int32(j)) {
				for _, id := range root2.EndIndex {
					if id != int32(wi) {
						res = append(res, []int{int(id), int(wi)})
					}
				}
			}
		}
	}

	// 特殊处理空串+回文串的情况
	isPalindrome := make([]bool, len(words))
	for i, word := range words {
		isPalindrome[i] = manachers[i].IsPalindrome(0, int32(len(word)))
	}
	for i, word := range words {
		if word == "" {
			for j := range words {
				if isPalindrome[j] {
					res = append(res, []int{i, j})
					res = append(res, []int{j, i})
				}
			}
		}
	}
	return
}

// Prefix-Suffix Palindrome (Hard version)
// https://www.luogu.com.cn/problem/CF1326D2
// 给定一个字符串。
// 要求选取他的一个前缀(可以为空)和与该前缀不相交的一个后缀(可以为空)拼接成回文串，
// 且该回文串长度最大。求方案.
//
// !先选取互为反串的最长前后缀，然后在剩余字符串中选取最长回文前缀或后缀。
func CF1326D2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(s string) string {
		n := len(s)
		ptr := 0
		for ptr < n && s[ptr] == s[n-1-ptr] {
			ptr++
		}
		if ptr == n {
			return s
		}

		M := NewManacher(s)
		preStart, sufEnd := int32(ptr), int32(n-1-ptr+1)
		preEnd, sufStart := int32(0), int32(0)

		// 回文前缀
		for e := sufEnd; e > preStart; e-- {
			if M.IsPalindrome(preStart, e) {
				preEnd = e
				break
			}
		}
		// 回文后缀
		for s := preStart; s < sufEnd; s++ {
			if M.IsPalindrome(s, sufEnd) {
				sufStart = s
				break
			}
		}

		if len1, len2 := preEnd-preStart, sufEnd-sufStart; len1 >= len2 {
			return s[:preEnd] + s[n-ptr:]
		} else {
			return s[:ptr] + s[sufStart:]
		}

	}

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, solve(s))
	}
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
			k = min32(ma.oddRadius[left+right-i], right-i+1)
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
			k = min32(ma.evenRadius[left+right-i+1], right-i+1)
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
		ma.maxOdd1[start] = max32(ma.maxOdd1[start], length)
		ma.maxOdd2[end] = max32(ma.maxOdd2[end], length)
	}
	for i := int32(0); i < n; i++ {
		if i-1 >= 0 {
			ma.maxOdd1[i] = max32(ma.maxOdd1[i], ma.maxOdd1[i-1]-2)
		}
		if i+1 < n {
			ma.maxOdd2[i] = max32(ma.maxOdd2[i], ma.maxOdd2[i+1]-2)
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
		ma.maxEven1[start] = max32(ma.maxEven1[start], length)
		ma.maxEven2[end] = max32(ma.maxEven2[end], length)
	}
	for i := int32(0); i < n; i++ {
		if i-1 >= 0 {
			ma.maxEven1[i] = max32(ma.maxEven1[i], ma.maxEven1[i-1]-2)
		}
		if i+1 < n {
			ma.maxEven2[i] = max32(ma.maxEven2[i], ma.maxEven2[i+1]-2)
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

type TrieNode[K comparable] struct {
	Children map[K]*TrieNode[K]
	EndIndex []int32
}

func NewTrieNode[K comparable]() *TrieNode[K] {
	return &TrieNode[K]{Children: map[K]*TrieNode[K]{}}
}

type Trie[K comparable] struct {
	root *TrieNode[K]
}

func NewTrie[K comparable]() *Trie[K] {
	return &Trie[K]{root: NewTrieNode[K]()}
}

func (t *Trie[K]) Insert(n int32, f func(int32) K, id int32) {
	cur := t.root
	for i := int32(0); i < n; i++ {
		char := f(i)
		if v, ok := cur.Children[char]; ok {
			cur = v
		} else {
			newNode := NewTrieNode[K]()
			cur.Children[char] = newNode
			cur = newNode
		}
	}
	cur.EndIndex = append(cur.EndIndex, id)
}

func (t *Trie[K]) Search(n int32, f func(int32) K) (ids []int32) {
	cur := t.root
	for i := int32(0); i < n; i++ {
		char := f(i)
		if v, ok := cur.Children[char]; ok {
			cur = v
		} else {
			return
		}
	}
	return cur.EndIndex
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	P3501()
}

// P3501 [POI2010] ANT-Antisymmetry (!对于一个0/1序列，求出其中异或意义下回文的子串数量。)
// https://www.luogu.com.cn/problem/P3501
// 对于一个01字符串，如果将这个字符串0和1取反后，再将整个串反过来和原串一样，就称作“反对称”字符串。
// 比如00001111和010101就是反对称的，1001就不是。
// 现在给出一个长度为N的01字符串，求它有多少个子串是反对称的。
// eg:
// 11001011
// 7个反对称子串分别是：01（出现两次），10（出现两次），0101，1100和001011
//
// 反对称一定是偶数长度的回文串.
func P3501() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s01 string
	fmt.Fscan(in, &s01)

	M := NewManacherWithEqual(s01)
	equal := func(i, j int32) bool { return s01[i] != s01[j] } // 异或意义下的回文
	M.SetEqual(equal)
	evenRadius := M.GetEvenRadius()
	res := 0
	for _, v := range evenRadius {
		res += int(v)
	}
	fmt.Fprintln(out, res)
}

type Sequence = string

// 自定义回文串判定函数的马拉车算法.
type ManacherWithEqual struct {
	n          int32
	seq        Sequence
	oddRadius  []int32
	evenRadius []int32
	maxOdd1    []int32
	maxOdd2    []int32
	maxEven1   []int32
	maxEven2   []int32
	equal      func(i, j int32) bool
}

func NewManacherWithEqual(seq Sequence) *ManacherWithEqual {
	defaultEqual := func(i, j int32) bool { return seq[i] == seq[j] }
	m := &ManacherWithEqual{
		n:     int32(len(seq)),
		seq:   seq,
		equal: defaultEqual,
	}
	return m
}

func (ma *ManacherWithEqual) SetEqual(equal func(i, j int32) bool) {
	ma.equal = equal
}

// 查询切片s[start:end]是否为回文串.
// 空串不为回文串.
func (ma *ManacherWithEqual) IsPalindrome(start, end int32) bool {
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
func (ma *ManacherWithEqual) GetOddRadius() []int32 {
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
		for i-k >= 0 && i+k < n && ma.equal(i-k, i+k) {
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
func (ma *ManacherWithEqual) GetEvenRadius() []int32 {
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
		for i-k-1 >= 0 && i+k < n && ma.equal(i-k-1, i+k) {
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
func (ma *ManacherWithEqual) GetLongestOddStartsAt(index int32) int32 {
	ma._initOdds()
	return ma.maxOdd1[index]
}

// 以s[index]结尾的最长奇回文子串的长度.
func (ma *ManacherWithEqual) GetLongestOddEndsAt(index int32) int32 {
	ma._initOdds()
	return ma.maxOdd2[index]
}

// 以s[index]开头的最长偶回文子串的长度.
func (ma *ManacherWithEqual) GetLongestEvenStartsAt(index int32) int32 {
	ma._initEvens()
	return ma.maxEven1[index]
}

// 以s[index]结尾的最长偶回文子串的长度.
func (ma *ManacherWithEqual) GetLongestEvenEndsAt(index int32) int32 {
	ma._initEvens()
	return ma.maxEven2[index]
}

func (ma *ManacherWithEqual) Len() int32 {
	return ma.n
}

func (ma *ManacherWithEqual) _initOdds() {
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

func (ma *ManacherWithEqual) _initEvens() {
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

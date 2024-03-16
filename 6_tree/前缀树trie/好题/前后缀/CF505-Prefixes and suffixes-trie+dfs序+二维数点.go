package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

// 505. Prefixes and suffixes
// https://codeforces.com/problemsets/acmsguru/problem/99999/505
//
// 给定一个只包含小写字母的字符串集合words.
// 在线处理q个查询：words中有多少个字符串s满足prefix是s的一个前缀, suffix是s的一个后缀.
// n,q,|words[i]|<=1e5
//
// trie+dfs+二维数点:
// !对于给出的每个字符串正着插入字典树A，倒着插入字典树B，
// 对于一个前缀来说，在字典树A上得到的dfs序[st,en]就是所有的匹配串，
// 同理，后缀在字典树B上dfs序[st,en]表示所有的后缀匹配串，
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	prefixTrie, suffixTrie := NewTrie(), NewTrie()
	points := make(map[[2]int32]int32)
	for _, word := range words {
		pos1 := prefixTrie.Insert(word)
		pos2 := suffixTrie.Insert(ReverseString(word))
		points[[2]int32{pos1, pos2}]++
	}
	getLidRid := func(t *Trie) (lid, rid []int32) {
		lid, rid = make([]int32, t.Size()), make([]int32, t.Size())
		dfn := int32(0)
		var dfs func(int32)
		dfs = func(pos int32) {
			lid[pos] = dfn
			dfn++
			for _, next := range t.nodes[pos].children {
				dfs(next)
			}
			rid[pos] = dfn
		}
		dfs(0)
		return
	}
	lid1, rid1 := getLidRid(prefixTrie)
	lid2, rid2 := getLidRid(suffixTrie)
	S := NewRectangleSum()
	for k, v := range points {
		p1, p2 := k[0], k[1]
		S.AddPoint(lid1[p1], lid2[p2], v)
	}
	S.Build()

	// 查询words中有多少个字符串s满足prefix是s的一个前缀, suffix是s的一个后缀.
	query := func(prefix, suffix string) int32 {
		pos1, ok1 := prefixTrie.Search(prefix)
		if !ok1 {
			return 0
		}
		pos2, ok2 := suffixTrie.Search(ReverseString(suffix))
		if !ok2 {
			return 0
		}
		x1, x2 := lid1[pos1], rid1[pos1]
		y1, y2 := lid2[pos2], rid2[pos2]
		return int32(S.Query(x1, x2, y1, y2))
	}

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var prefix, suffix string
		fmt.Fscan(in, &prefix, &suffix)
		fmt.Fprintln(out, query(prefix, suffix))
	}
}

type TrieNode struct {
	children map[byte]int32
	endCount int32
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[byte]int32)}
}

type Trie struct {
	nodes []*TrieNode
}

func NewTrie() *Trie {
	return &Trie{nodes: []*TrieNode{NewTrieNode()}}
}

func (t *Trie) Insert(s string) (pos int32) {
	for i := 0; i < len(s); i++ {
		char := s[i]
		if next, ok := t.nodes[pos].children[char]; ok {
			pos = next
		} else {
			next = int32(len(t.nodes))
			t.nodes[pos].children[char] = next
			t.nodes = append(t.nodes, NewTrieNode())
			pos = next
		}
	}
	t.nodes[pos].endCount++
	return pos
}

func (t *Trie) Search(s string) (pos int32, ok bool) {
	nodes := t.nodes
	for i := 0; i < len(s); i++ {
		char := s[i]
		if next, ok := nodes[pos].children[char]; ok {
			pos = next
		} else {
			return 0, false
		}
	}
	return pos, true
}

// Trie 中的节点数(包含根节点).
func (t *Trie) Size() int32 {
	return int32(len(t.nodes))
}

// RectangleSumStatic
type StaticRectangleSum struct {
	points [][3]int32
	xs     []int32
	ys     []int32
	wm     *waveletMatrix
}

func NewRectangleSum() *StaticRectangleSum {
	return &StaticRectangleSum{
		points: [][3]int32{},
	}
}

func (s *StaticRectangleSum) AddPoint(x, y, w int32) {
	s.points = append(s.points, [3]int32{x, y, w})
}

func (s *StaticRectangleSum) Build() {
	sort.Slice(s.points, func(i, j int) bool {
		return s.points[i][0] < s.points[j][0]
	})
	n := len(s.points)
	xs, ys, ws := make([]int32, n), make([]int32, n), make([]int32, n)
	for i, p := range s.points {
		xs[i], ys[i], ws[i] = p[0], p[1], p[2]
	}
	s.xs = xs

	set := make(map[int32]struct{}, len(ys))
	for _, y := range ys {
		set[y] = struct{}{}
	}
	sortedSet := make([]int32, 0, len(set))
	for y := range set {
		sortedSet = append(sortedSet, y)
	}
	sort.Slice(sortedSet, func(i, j int) bool { return sortedSet[i] < sortedSet[j] })
	s.ys = sortedSet

	comp := make(map[int32]int32, len(sortedSet))
	for i, y := range sortedSet {
		comp[y] = int32(i)
	}

	newYs := make([]int32, len(ys))
	for i, y := range ys {
		newYs[i] = comp[y]
	}

	maxLog := bits.Len(uint(len(sortedSet)))
	s.wm = newWaveletMatrix(newYs, ws, int32(maxLog))
}

// 求矩形x1<=x<x2,y1<=y<y2的权值和 注意是左闭右开
func (s *StaticRectangleSum) Query(x1, x2, y1, y2 int32) int {
	return s.rectSum(x1, x2, y2) - s.rectSum(x1, x2, y1)
}

func (s *StaticRectangleSum) rectSum(left, right, upper int32) int {
	left = int32(sort.Search(len(s.xs), func(i int) bool { return s.xs[i] >= left }))
	right = int32(sort.Search(len(s.xs), func(i int) bool { return s.xs[i] >= right }))
	upper = int32(sort.Search(len(s.ys), func(i int) bool { return s.ys[i] >= upper }))
	return s.wm.RectSum(left, right, upper)
}

func newWaveletMatrix(ys, ws []int32, maxLog int32) *waveletMatrix {
	n := int32(len(ys))
	mat := make([]*bitVector, 0, maxLog)
	zs := make([]int32, 0, maxLog)
	data := make([][]int, maxLog)
	for i := range data {
		data[i] = make([]int, n+1)
	}

	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}

	for d := maxLog - 1; d >= 0; d-- {
		vec := newBitVector(n + 1)
		ls, rs := make([]int32, 0, n), make([]int32, 0, n)
		for i, val := range order {
			if (ys[val]>>uint(d))&1 == 1 {
				rs = append(rs, val)
				vec.Set(int32(i))
			} else {
				ls = append(ls, val)
			}
		}
		vec.Build()
		mat = append(mat, vec)
		zs = append(zs, int32(len(ls)))
		order = append(ls, rs...)
		for i, val := range order {
			data[maxLog-d-1][i+1] = data[maxLog-d-1][i] + int(ws[val])
		}
	}

	return &waveletMatrix{
		n:      n,
		maxLog: maxLog,
		mat:    mat,
		zs:     zs,
		data:   data,
	}
}

type waveletMatrix struct {
	n      int32
	maxLog int32
	mat    []*bitVector
	zs     []int32
	data   [][]int
}

func (w *waveletMatrix) RectSum(left, right, upper int32) int {
	res := 0
	for d := int32(0); d < w.maxLog; d++ {
		if (upper>>(w.maxLog-d-1))&1 == 1 {
			res += int(w.data[d][w.mat[d].Count(0, right)])
			res -= int(w.data[d][w.mat[d].Count(0, left)])
			left = w.mat[d].Count(1, left) + w.zs[d]
			right = w.mat[d].Count(1, right) + w.zs[d]
		} else {
			left = w.mat[d].Count(0, left)
			right = w.mat[d].Count(0, right)
		}
	}
	return res
}

type bitVector struct {
	n     int32
	block []int
	sum   []int32
}

func newBitVector(n int32) *bitVector {
	blockCount := (n + 63) >> 6
	return &bitVector{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int32, blockCount+1),
	}
}

func (f *bitVector) Set(i int32) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *bitVector) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + int32(bits.OnesCount(uint(f.block[i])))
	}
}

// 统计 [0,end) 中 value 的个数
func (f *bitVector) Count(value, end int32) int32 {
	if value == 1 {
		return f.count1(end)
	}
	return end - f.count1(end)
}

func (f *bitVector) count1(k int32) int32 {
	mask := (1 << uint(k&63)) - 1
	return f.sum[k>>6] + int32(bits.OnesCount(uint(f.block[k>>6]&mask)))
}

func ReverseString(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes)
}

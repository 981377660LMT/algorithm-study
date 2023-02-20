// https://atcoder.jp/contests/jag2017autumn/tasks/jag2017autumn_h

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int = 1e9 + 7

// 有一个包含若干个单词的集合patterns.
// !给定一个字符串s,将s分割成若干个单词,使得每个部分都在patterns中.
// !求分割方案数.(mod 1e9+7)
// ∑patterns[i].length≤2e5
//
// dp[i]+快速判断子串是否在模式串中 => AC自动机
func SeparateString(s string, patterns []string) int {
	aho := NewAhoCorasick(26, 'a')
	for i, s := range patterns {
		aho.Add(s, i)
	}
	aho.Build(true)

	n := len(s)
	dp := make([]int, n+1)
	dp[0] = 1
	pos := 0
	for i := 0; i < n; i++ {
		pos = aho.Move(pos, s[i])
		indexes := aho.IndexAll(pos)
		for _, j := range indexes {
			dp[i+1] = (dp[i+1] + dp[i-len(patterns[j])+1]) % MOD
		}
	}

	return dp[n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}
	var s string
	fmt.Fscan(in, &s)
	fmt.Fprintln(out, SeparateString(s, words))
}

type AhoCorasick struct {
	*Trie
	count []int
}

// size: 字符集大小
// margin: 字符集起始字符
func NewAhoCorasick(size int, margin byte) *AhoCorasick {
	res := &AhoCorasick{Trie: NewTrie(size+1, margin)}
	return res
}

// matchAll: 匹配时是否需要每个模式串的信息.
//   如果需要匹配每个模式串的信息, 那么在构建时会对匹配的模式串索引排序(额外开销).
func (ac *AhoCorasick) Build(matchAll bool) {
	n := len(ac.stack)
	ac.count = make([]int, n)
	for i := 0; i < n; i++ {
		if matchAll {
			sort.Ints(ac.stack[i].Indexes)
		}
		ac.count[i] = len(ac.stack[i].Indexes)
	}

	var que []int
	for i := 0; i < ac.size-1; i++ {
		if *ac.Next(0, i) != -1 {
			*ac.Next(*ac.Next(0, i), ac.size-1) = 0
			que = append(que, *ac.Next(0, i))
		} else {
			*ac.Next(0, i) = 0
		}
	}

	for len(que) > 0 {
		x := ac.stack[que[0]]
		fail := x.next[ac.size-1]
		ac.count[que[0]] += ac.count[fail]
		que = que[1:]

		for i := 0; i < ac.size-1; i++ {
			nx := &x.next[i]
			if *nx < 0 {
				*nx = *ac.Next(fail, i)
				continue
			}
			que = append(que, *nx)
			*ac.Next(*nx, ac.size-1) = *ac.Next(fail, i)

			if matchAll {
				idx := ac.stack[*nx].Indexes
				idy := ac.stack[*ac.Next(fail, i)].Indexes
				idz := make([]int, 0, len(idx)+len(idy))

				// set union
				i, j := 0, 0
				for i < len(idx) && j < len(idy) {
					if idx[i] < idy[j] {
						idz = append(idz, idx[i])
						i++
					} else if idx[i] > idy[j] {
						idz = append(idz, idy[j])
						j++
					} else {
						idz = append(idz, idx[i])
						i++
						j++
					}
				}
				for i < len(idx) {
					idz = append(idz, idx[i])
					i++
				}
				for j < len(idy) {
					idz = append(idz, idy[j])
					j++
				}

				ac.stack[*nx].Indexes = idz
			}
		}

	}
}

// matchAll: 是否返回`每个`模式串的匹配次数.
//  返回值: 如果 matchAll 为 false, 则返回值只有一个元素, 为所有模式串的匹配次数之和.
//  如果 matchAll 为 true, 则返回值有 len(模式串集合) 个元素, 第 i 个元素为第 i 个模式串的匹配次数.
func (ac *AhoCorasick) Match(s string, matchAll bool) []int {
	size := 1
	if matchAll {
		size = ac.Size()
	}
	res := make([]int, size)
	pos := 0
	for i := 0; i < len(s); i++ {
		pos = *ac.Next(pos, int(s[i]-ac.margin))
		if matchAll {
			for _, x := range ac.stack[pos].Indexes {
				res[x]++
			}
		} else {
			res[0] += ac.count[pos]
		}
	}
	return res
}

// 从结点pos开始匹配字符串s.
//  返回值: 如果 matchAll 为 false, 则返回值只有一个元素, 为所有模式串的匹配次数之和.
//  如果 matchAll 为 true, 则返回值有 len(模式串集合) 个元素, 第 i 个元素为第 i 个模式串的匹配次数.
func (ac *AhoCorasick) MatchFrom(pos int, s string, matchAll bool) []int {
	size := 1
	if matchAll {
		size = ac.Size()
	}
	res := make([]int, size)
	for i := 0; i < len(s); i++ {
		pos = *ac.Next(pos, int(s[i]-ac.margin))
		if matchAll {
			for _, x := range ac.stack[pos].Indexes {
				res[x]++
			}
		} else {
			res[0] += ac.count[pos]
		}
	}
	return res
}

// 当前节点匹配的完整的模式串的个数.
func (ac *AhoCorasick) Count(pos int) int {
	return ac.count[pos]
}

//
//
//
//
type Trie struct {
	size   int
	margin byte
	stack  []*trieNode
}

type trieNode struct {
	Indexes []int // 存储了哪些模式串的索引
	Key     byte
	index   int   // 最后一次被更新的模式串的索引
	next    []int // children position
}

// size: 字符集大小
// margin: 字符集起始字符
func NewTrie(size int, margin byte) *Trie {
	res := &Trie{size: size, margin: margin}
	root := res.newNode('$')
	res.stack = append(res.stack, root)
	return res
}

// 向字典树中添加一个模式串.
func (t *Trie) Add(s string, index int) {
	if len(s) == 0 {
		return
	}
	pos := 0
	for i := 0; i < len(s); i++ {
		k := int(s[i] - t.margin)
		if *t.Next(pos, k) != -1 {
			pos = *t.Next(pos, k)
			continue
		}
		nextPos := len(t.stack)
		*t.Next(pos, k) = nextPos
		node := t.newNode(s[i])
		t.stack = append(t.stack, node)
		pos = nextPos
	}
	t.stack[pos].index = index
	t.stack[pos].Indexes = append(t.stack[pos].Indexes, index)
}

func (t *Trie) Find(s string) int {
	if len(s) == 0 {
		return 0
	}
	pos := 0
	for i := 0; i < len(s); i++ {
		k := int(s[i] - t.margin)
		if *t.Next(pos, k) == -1 {
			return -1
		}
		pos = *t.Next(pos, k)
	}
	return pos
}

// 从结点pos移动到c指向的下一个结点.
func (t *Trie) Move(pos int, c byte) int {
	if pos < 0 || pos >= len(t.stack) {
		return -1
	}
	return *t.Next(pos, int(c-t.margin))
}

func (t *Trie) Index(pos int) int {
	if pos < 0 {
		return -1
	}
	return t.stack[pos].index
}

func (t *Trie) IndexAll(pos int) []int {
	if pos < 0 {
		return nil
	}
	return t.stack[pos].Indexes
}

// Trie 中的节点数(包含根节点).
func (t *Trie) Size() int { return len(t.stack) }

func (t *Trie) Next(pos, j int) *int {
	return &t.stack[pos].next[j]
}

func (t *Trie) newNode(c byte) *trieNode {
	next := make([]int, t.size)
	for i := range next {
		next[i] = -1
	}
	return &trieNode{
		index: -1,
		Key:   c,
		next:  next,
	}
}

// https://oi-wiki.org/string/ac-automaton/
// pos可以理解为当前节点的位置，pos=0表示根节点
// Trie 中的结点表示的是某个模式串的前缀。
// 我们在后文也将其称作状态。一个结点表示一个状态，Trie 的边就是状态的转移。
package main

import "fmt"

type Trie struct {
	trie *trie
}

func Constructor() Trie {
	return Trie{trie: NewTrie(26, 'a')}
}

func (this *Trie) Insert(word string) {
	this.trie.Add(word, 1)
}

func (this *Trie) Search(word string) bool {
	pos := this.trie.Find(word)
	return pos != -1 && this.trie.Index(pos) != -1
}

func (this *Trie) StartsWith(prefix string) bool {
	pos := this.trie.Find(prefix)
	return pos != -1
}

func main() {
	t := Constructor()
	t.Insert("apple")
	fmt.Println(t.Search("apple"))
	fmt.Println(t.Search("app"))
}

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.Search(word);
 * param_3 := obj.StartsWith(prefix);
 */

type trie struct {
	size   int
	margin byte
	stack  []*trieNode
}

type trieNode struct {
	index   int // 最后一次被更新的字符串的索引
	key     byte
	indexes []int // 存储了哪些字符串的索引
	next    []int // children position
}

// size: 字符集大小
// margin: 字符集起始字符
func NewTrie(size int, margin byte) *trie {
	res := &trie{size: size, margin: margin}
	root := res.newNode('$')
	res.stack = append(res.stack, root)
	return res
}

// 将s添加到Trie中，index为s的索引.
func (t *trie) Add(s string, index int) {
	pos := 0
	for i := 0; i < len(s); i++ {
		k := int(s[i] - t.margin)
		if *t.next(pos, k) != -1 {
			pos = *t.next(pos, k)
			continue
		}
		nextPos := len(t.stack)
		*t.next(pos, k) = nextPos
		node := t.newNode(s[i])
		t.stack = append(t.stack, node)
		pos = nextPos
	}
	t.stack[pos].index = index
	t.stack[pos].indexes = append(t.stack[pos].indexes, index)
}

// 返回字符串s的前缀在Trie中的位置，如果不存在则返回-1.
func (t *trie) Find(s string) (pos int) {
	for i := 0; i < len(s); i++ {
		k := int(s[i] - t.margin)
		if *t.next(pos, k) == -1 {
			return -1
		}
		pos = *t.next(pos, k)
	}
	return pos
}

// 从结点pos出发，沿着字符c移动到下一个结点.
func (t *trie) Move(pos int, c byte) int {
	if pos < 0 || pos >= len(t.stack) {
		return -1
	}
	return *t.next(pos, int(c-t.margin))
}

// 返回结点pos处的字符串的索引.
func (t *trie) Index(pos int) int {
	if pos < 0 {
		return -1
	}
	return t.stack[pos].index
}

// 返回结点pos处存储的所有字符串的索引.
func (t *trie) IndexAll(pos int) []int {
	if pos < 0 {
		return []int{}
	}
	return t.stack[pos].indexes
}

func (t *trie) next(i, j int) *int {
	return &t.stack[i].next[j]
}

func (t *trie) newNode(c byte) *trieNode {
	next := make([]int, t.size)
	for i := range next {
		next[i] = -1
	}
	return &trieNode{
		index: -1,
		key:   c,
		next:  next,
	}
}

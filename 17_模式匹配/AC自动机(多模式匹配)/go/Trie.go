package main

func main() {
	trie := NewTrie(26, 'a')
}

type Trie struct {
	size   int
	margin byte
	stack  []*trieNode
}

type trieNode struct {
	index   int
	key     byte
	next    []int
	indexes []int // 存储了哪些字符串的索引
}

// size: 字符集大小
// margin: 字符集起始字符
func NewTrie(size int, margin byte) *Trie {
	res := &Trie{size: size, margin: margin}
	root := res.newNode('$')
	res.stack = append(res.stack, root)
	return res
}

func (t *Trie) Add(s string, index int) {
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

func (t *Trie) Find(s string) int {
	pos := 0
	for i := 0; i < len(s); i++ {
		k := int(s[i] - t.margin)
		if *t.next(pos, k) == -1 {
			return -1
		}
		pos = *t.next(pos, k)
	}
	return pos
}

func (t *Trie) Move(pos int, c byte) int {
	if pos < 0 || pos >= len(t.stack) {
		return -1
	}
	return *t.next(pos, int(c-t.margin))
}

func (t *Trie) Index(pos int) int {
	if pos < 0 {
		return -1
	}
	return t.stack[pos].index
}

func (t *Trie) IndexAll(pos int) []int {
	if pos < 0 {
		return []int{}
	}
	return t.stack[pos].indexes
}

func (t *Trie) Size() int { return len(t.stack) }

func (t *Trie) next(i, j int) *int {
	return &t.stack[i].next[j]
}

func (t *Trie) newNode(c byte) *trieNode {
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

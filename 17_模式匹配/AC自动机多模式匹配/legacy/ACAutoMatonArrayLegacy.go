// https://nyaannyaan.github.io/library/string/aho-corasick.hpp
// 默认字符集为26个小写字母,可以修改_SIZE和_MARGIN来修改字符集大小和起始字符.
// !Trie由数组实现，比map实现的Trie快很多.
//
// API
//  Insert(id int, s string, didInsert func(pos int)) *ACAutoMatonArrayLegacy
//  Build(heavy bool, dp func(fail, next int))
//  Match(state int, s string) map[int][]int
//  Move(pos, char) int
//  Count(pos) int
//  Accept(pos) bool
//  Size() int

package main

//###############################################################
// 1032. 字符流
// https://leetcode.cn/problems/stream-of-characters/

type StreamChecker struct {
	acm   *ACAutoMatonArrayLegacy
	state int
}

func Constructor(words []string) StreamChecker {
	acm := NewACAutoMatonArrayLegacy()
	for i, word := range words {
		acm.Insert(i, word, nil)
	}
	acm.Build(false, nil)
	return StreamChecker{acm: acm}
}

// !从字符流中接收一个新字符，如果字符流中的任一非空后缀能匹配 words 中的某一字符串，返回 true
func (this *StreamChecker) Query(letter byte) bool {
	this.state = this.acm.Move(this.state, letter)
	return this.acm.Accept(this.state)
}

//###############################################################

const (
	_SIZE   int  = 26  // 字符集大小
	_MARGIN byte = 'a' // 字符集起始字符
)

type ACAutoMatonArrayLegacy struct {
	*Trie // !继承自Trie
	count []int
	heavy bool
}

func NewACAutoMatonArrayLegacy() *ACAutoMatonArrayLegacy {
	res := &ACAutoMatonArrayLegacy{Trie: _newTrie()}
	return res
}

// bfs为字典树的每个结点添加失配指针,结点要跳转到哪里.
//
//	heavy: 是否处理出每个结点匹配到的模式串id.
//	dp: AC自动机构建过程中的回调函数,入参为`(next结点的fail指针, next结点)`.
func (ac *ACAutoMatonArrayLegacy) Build(heavy bool, dp func(fail, next int)) {
	ac.count = make([]int, len(ac.states))
	for i, v := range ac.states {
		ac.count[i] = len(v.matching)
	}
	ac.heavy = heavy

	var queue []int
	for i := 0; i < _SIZE; i++ {
		if next := *ac._next(0, i); next != -1 {
			*ac._next(next, _SIZE) = 0
			queue = append(queue, next)
		} else {
			*ac._next(0, i) = 0
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		node := ac.states[cur]
		fail := node.children[_SIZE]
		ac.count[cur] += ac.count[fail]
		if dp != nil {
			dp(fail, cur)
		}

		for i := 0; i < _SIZE; i++ {
			next := node.children[i]
			if next < 0 {
				node.children[i] = *ac._next(fail, i)
				continue
			}

			queue = append(queue, next)
			move := *ac._next(fail, i)
			*ac._next(next, _SIZE) = move
			if dp != nil {
				dp(move, next)
			}
			if heavy {
				ac.states[next].matching = append(ac.states[next].matching, ac.states[move].matching...)
			}
		}
	}
}

// 从状态`state`开始匹配字符串`s`.
//
//	state : ac自动机的状态.根节点状态为0.
//	s : 待匹配的字符串.
//	返回每个模式串在`s`中出现的下标.
func (ac *ACAutoMatonArrayLegacy) Match(state int, s string) map[int][]int {
	if !ac.heavy {
		panic("需要调用build(heavy=True)构建AC自动机")
	}
	res := make(map[int][]int)
	root := state
	for i := 0; i < len(s); i++ {
		root = ac.Move(root, s[i])
		for _, m := range ac.states[root].matching {
			res[m] = append(res[m], i-len(ac.patterns[m])+1)
		}
	}
	return res
}

// 当前节点匹配的完整的模式串的个数.
func (ac *ACAutoMatonArrayLegacy) Count(pos int) int {
	return ac.count[pos]
}

// 当前节点状态是否为可接受(后缀匹配到了模式串).
func (ac *ACAutoMatonArrayLegacy) Accept(pos int) bool {
	return ac.count[pos] > 0
}

type trieNode struct {
	matching []int          // 存储了哪些模式串的索引(id)
	children [_SIZE + 1]int // 转移到的状态.最后一个位置存储了fail指针.
	id       int            // 属于哪个模式串(覆盖更新)
}

type Trie struct {
	states   []*trieNode
	patterns []string
}

func _newTrie() *Trie {
	res := &Trie{}
	root := res.newNode('$')
	res.states = append(res.states, root)
	return res
}

func (t *Trie) newNode(c byte) *trieNode {
	res := &trieNode{}
	for i := range res.children {
		res.children[i] = -1
	}
	res.id = -1
	return res
}

// 向字典树中添加一个模式串.
func (t *Trie) Insert(id int, pattern string, didInsert func(pos int)) *Trie {
	if len(pattern) == 0 {
		return t
	}
	pos := 0
	for i := 0; i < len(pattern); i++ {
		char := pattern[i]
		offset := int(char - _MARGIN)
		if child := *t._next(pos, offset); child != -1 {
			pos = child
			continue
		}
		newState := len(t.states)
		*t._next(pos, offset) = newState
		node := t.newNode(char)
		t.states = append(t.states, node)
		pos = newState
	}
	t.states[pos].id = id
	t.states[pos].matching = append(t.states[pos].matching, id)
	t.patterns = append(t.patterns, pattern)
	if didInsert != nil {
		didInsert(pos)
	}
	return t
}

// 从结点pos移动到c指向的下一个结点.
func (t *Trie) Move(pos int, c byte) int {
	if pos < 0 || pos >= len(t.states) {
		return -1
	}
	return *t._next(pos, int(c-_MARGIN))
}

// Trie 中的节点数(包含根节点).
func (t *Trie) Size() int { return len(t.states) }

func (t *Trie) _next(pos, j int) *int {
	return &t.states[pos].children[j]
}

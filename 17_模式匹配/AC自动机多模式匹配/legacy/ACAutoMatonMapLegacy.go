// !TODO: 看Nyann的代码并修复，优化复杂度

package main

//###############################################################
// 1032. 字符流
// https://leetcode.cn/problems/stream-of-characters/

type StreamChecker struct {
	acm   *ACAutoMatonLegacy
	state int
}

func Constructor(words []string) StreamChecker {
	acm := NewACAutoMatonLegacy()
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

// ###############################################################
// 面试题 17.17. 多次搜索
// https://leetcode.cn/problems/multi-search-lcci/
func multiSearch(big string, smalls []string) [][]int {
	acm := NewACAutoMatonLegacy()
	for i, word := range smalls {
		acm.Insert(i, word, nil)
	}
	acm.Build(true, nil)
	res := make([][]int, len(smalls))
	matching := acm.Match(0, big)
	for pi, pos := range matching {
		for _, p := range pos {
			res[pi] = append(res[pi], p)
		}
	}
	return res
}

//###############################################################

// ###############################################################
// 2781. 最长合法子字符串的长度
// https://leetcode.cn/problems/length-of-the-longest-valid-substring/
func longestValidSubstring(word string, forbidden []string) int {
	const INF int = 1e18
	acm := NewACAutoMatonLegacy()
	minLen := make(map[int]int)

	for i, word := range forbidden {
		acm.Insert(i, word, func(pos int) {
			minLen[pos] = min(GetOrDefault(minLen, pos, INF), len(word))
		})
	}

	acm.Build(false, func(pre, cur int) {
		minLen[cur] = min(GetOrDefault(minLen, cur, INF), GetOrDefault(minLen, pre, INF))
	})

	res := 0
	left := 0
	state := 0
	for right := 0; right < len(word); right++ {
		state = acm.Move(state, word[right])
		start := right - GetOrDefault(minLen, state, INF) + 1
		left = max(left, start+1)
		res = max(res, right-left+1)

	}
	return res
}

type K = int
type V = int

func GetOrDefault(dict map[K]V, key K, defaultValue V) V {
	if v, ok := dict[key]; ok {
		return v
	}
	return defaultValue
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

//###############################################################

type ACAutoMatonLegacy struct {
	patterns  []string       // 模式串列表
	children  []map[byte]int // trie树,children[i]表示节点(状态)i的所有子节点,0表示虚拟根节点
	wordCount []int          // trie树结点附带的信息.count[i]表示节点(状态)i的匹配个数
	matching  [][]int        // trie树结点附带的信息.matching[i]表示节点(状态)i对应的字符串在patterns中的下标
	fail      []int          // fail[i]表示节点(状态)i的失配指针
	heavy     bool           // 构建时是否在matching中处理出每个结点匹配到的模式串id
}

func NewACAutoMatonLegacy() *ACAutoMatonLegacy {
	res := &ACAutoMatonLegacy{}
	res.children = append(res.children, make(map[byte]int))
	res.matching = append(res.matching, make([]int, 0))
	return res
}

// 将模式串`pattern`插入到Trie树中.模式串一般是`被禁用的单词`.
//
//	pid : 模式串的唯一标识id.
//	pattern : 模式串.
//	didInsert : 模式串插入后的回调函数,入参为结束字符所在的结点(状态).
func (acm *ACAutoMatonLegacy) Insert(pid int, pattern string, didInsert func(state int)) *ACAutoMatonLegacy {
	if len(pattern) == 0 {
		return acm
	}
	root := 0
	for i := 0; i < len(pattern); i++ {
		char := pattern[i]
		nexts := acm.children[root]
		if _, ok := nexts[char]; ok {
			root = nexts[char]
		} else {
			nextState := len(acm.children)
			nexts[char] = nextState
			root = nextState
			acm.children = append(acm.children, make(map[byte]int))
			acm.matching = append(acm.matching, make([]int, 0))
		}
	}
	acm.matching[root] = append(acm.matching[root], pid)
	acm.patterns = append(acm.patterns, pattern)
	if didInsert != nil {
		didInsert(root)
	}
	return acm
}

// 构建失配指针.
// bfs为字典树的每个结点添加失配指针,结点要跳转到哪里.
// AC自动机的失配指针指向的节点所代表的字符串 是 当前节点所代表的字符串的 最长后缀.
//
//	heavy : 是否处理出每个结点匹配到的模式串id.
//	dp : AC自动机构建过程中的回调函数,入参为`(next结点的fail指针, next结点)`.
func (acm *ACAutoMatonLegacy) Build(heavy bool, dp func(move, next int)) {
	acm.wordCount = make([]int, len(acm.children))
	for i := 0; i < len(acm.children); i++ {
		acm.wordCount[i] = len(acm.matching[i])
	}
	acm.fail = make([]int, len(acm.children))
	acm.heavy = heavy
	queue := make([]int, 0)
	for _, v := range acm.children[0] {
		queue = append(queue, v)
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		curFail := acm.fail[cur]
		for input, next := range acm.children[cur] {
			move := acm.Move(curFail, input)
			acm.fail[next] = move
			acm.wordCount[next] += acm.wordCount[move]
			if heavy {
				for _, m := range acm.matching[move] {
					acm.matching[next] = append(acm.matching[next], m)
				}
			}
			if dp != nil {
				dp(move, next)
			}
			queue = append(queue, next)
		}
	}
}

// 从当前状态`state`沿着字符`input`转移到的下一个状态.
// 沿着失配链上跳,找到第一个可由char转移的节点.
func (acm *ACAutoMatonLegacy) Move(state int, input byte) int {
	for {
		nexts := acm.children[state]
		if _, ok := nexts[input]; ok {
			return nexts[input]
		}
		if state == 0 {
			return 0
		}
		state = acm.fail[state]
	}
}

// 从状态`state`开始匹配字符串`s`.
//
//	state : ac自动机的状态.根节点状态为0.
//	s : 待匹配的字符串.
//	返回每个模式串在`s`中出现的下标.
func (acm *ACAutoMatonLegacy) Match(state int, s string) map[int][]int {
	if !acm.heavy {
		panic("需要调用build(heavy=True)构建AC自动机")
	}
	res := make(map[int][]int)
	root := state
	for i := 0; i < len(s); i++ {
		root = acm.Move(root, s[i])
		for _, m := range acm.matching[root] {
			res[m] = append(res[m], i-len(acm.patterns[m])+1)
		}
	}
	return res
}

// 当前状态`state`匹配到的模式串个数.
func (acm *ACAutoMatonLegacy) Count(state int) int {
	return acm.wordCount[state]
}

// 当前状态`state`是否为匹配状态.
func (acm *ACAutoMatonLegacy) Accept(state int) bool {
	return acm.wordCount[state] > 0
}

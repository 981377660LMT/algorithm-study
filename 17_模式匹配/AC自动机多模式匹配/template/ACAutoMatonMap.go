package main

const INF = int(2e18)

// https://leetcode.cn/problems/length-of-the-longest-valid-substring/description/
func longestValidSubstring(word string, forbidden []string) int {
	acm := NewACAutoMatonMap()
	for _, w := range forbidden {
		acm.AddString([]byte(w))
	}
	acm.BuildSuffixLink()

	minBorder := make([]int, acm.Size()) // 每个状态匹配到的模式串的最小长度
	for i := range minBorder {
		minBorder[i] = INF
	}
	for i, pos := range acm.WordPos {
		minBorder[pos] = min(minBorder[pos], len(forbidden[i]))
	}
	acm.Dp(func(from, to int) { minBorder[to] = min(minBorder[to], minBorder[from]) })

	res, left, pos := 0, 0, 0
	for right, char := range word {
		pos = acm.Move(pos, byte(char))
		left = max(left, right-minBorder[pos]+2)
		res = max(res, right-left+1)
	}
	return res
}

type T = byte

type ACAutoMatonMap struct {
	WordPos    []int         // WordPos[i] 表示加入的第i个模式串对应的节点编号.
	children   []map[T]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	suffixLink []int32       // 又叫fail.指向当前节点最长真后缀对应结点.
	bfsOrder   []int32       // 结点的拓扑序,0表示虚拟节点.
}

func NewACAutoMatonMap() *ACAutoMatonMap {
	return &ACAutoMatonMap{
		WordPos:  []int{},
		children: []map[T]int32{{}},
	}
}

func (ac *ACAutoMatonMap) AddString(s []T) int {
	if len(s) == 0 {
		return 0
	}
	pos := 0
	for i := 0; i < len(s); i++ {
		ord := s[i]
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			pos = int(next)
		} else {
			nextState := len(ac.children)
			nexts[ord] = int32(nextState)
			pos = nextState
			ac.children = append(ac.children, map[T]int32{})
		}
	}
	ac.WordPos = append(ac.WordPos, pos)
	return pos
}

// 功能与 AddString 相同.
func (ac *ACAutoMatonMap) AddFrom(n int, getOrd func(i int) T) int {
	if n == 0 {
		return 0
	}
	pos := 0
	for i := 0; i < n; i++ {
		ord := getOrd(i)
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			pos = int(next)
		} else {
			nextState := len(ac.children)
			nexts[ord] = int32(nextState)
			pos = nextState
			ac.children = append(ac.children, map[T]int32{})
		}
	}
	ac.WordPos = append(ac.WordPos, pos)
	return pos
}

func (ac *ACAutoMatonMap) AddChar(pos int, ord T) int {
	nexts := ac.children[pos]
	if next, ok := nexts[ord]; ok {
		return int(next)
	}
	nextState := len(ac.children)
	nexts[ord] = int32(nextState)
	ac.children = append(ac.children, map[T]int32{})
	return nextState
}

func (ac *ACAutoMatonMap) Move(pos int, ord T) int {
	for {
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			return int(next)
		}
		if pos == 0 {
			return 0
		}
		pos = int(ac.suffixLink[pos])
	}
}

func (ac *ACAutoMatonMap) BuildSuffixLink() {
	ac.suffixLink = make([]int32, len(ac.children))
	for i := range ac.suffixLink {
		ac.suffixLink[i] = -1
	}
	ac.bfsOrder = make([]int32, len(ac.children))
	head, tail := 0, 1
	for head < tail {
		v := ac.bfsOrder[head]
		head++
		for char, next := range ac.children[v] {
			ac.bfsOrder[tail] = next
			tail++
			f := ac.suffixLink[v]
			for f != -1 {
				if _, ok := ac.children[f][char]; ok {
					break
				}
				f = ac.suffixLink[f]
			}
			if f == -1 {
				ac.suffixLink[next] = 0
			} else {
				ac.suffixLink[next] = ac.children[f][char]
			}
		}
	}
}

func (ac *ACAutoMatonMap) Empty() bool {
	return len(ac.children) == 1
}

func (ac *ACAutoMatonMap) Clear() {
	ac.WordPos = ac.WordPos[:0]
	ac.children = ac.children[:1]
	ac.children[0] = map[T]int32{}
	ac.suffixLink = ac.suffixLink[:0]
	ac.bfsOrder = ac.bfsOrder[:0]
}

func (ac *ACAutoMatonMap) GetCounter() []int {
	counter := make([]int, len(ac.children))
	for _, pos := range ac.WordPos {
		counter[pos]++
	}
	for _, v := range ac.bfsOrder {
		if v != 0 {
			counter[v] += counter[ac.suffixLink[v]]
		}
	}
	return counter
}

func (ac *ACAutoMatonMap) GetIndexes() [][]int {
	res := make([][]int, len(ac.children))
	for i, pos := range ac.WordPos {
		res[pos] = append(res[pos], i)
	}
	for _, v := range ac.bfsOrder {
		if v != 0 {
			from, to := ac.suffixLink[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int, 0, len(arr1)+len(arr2))
			i, j := 0, 0
			for i < len(arr1) && j < len(arr2) {
				if arr1[i] < arr2[j] {
					arr3 = append(arr3, arr1[i])
					i++
				} else if arr1[i] > arr2[j] {
					arr3 = append(arr3, arr2[j])
					j++
				} else {
					arr3 = append(arr3, arr1[i])
					i++
					j++
				}
			}
			for i < len(arr1) {
				arr3 = append(arr3, arr1[i])
				i++
			}
			for j < len(arr2) {
				arr3 = append(arr3, arr2[j])
				j++
			}
			res[to] = arr3
		}
	}
	return res
}

func (ac *ACAutoMatonMap) Dp(f func(from, to int)) {
	for _, v := range ac.bfsOrder {
		if v != 0 {
			f(int(ac.suffixLink[v]), int(v))
		}
	}
}

func (trie *ACAutoMatonMap) BuildFailTree() [][]int {
	res := make([][]int, trie.Size())
	trie.Dp(func(pre, cur int) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (ac *ACAutoMatonMap) BuildTrieTree() [][]int {
	res := make([][]int, ac.Size())
	var dfs func(int)
	dfs = func(cur int) {
		for _, next := range ac.children[cur] {
			res[cur] = append(res[cur], int(next))
			dfs(int(next))
		}
	}
	dfs(0)
	return res
}

func (ac *ACAutoMatonMap) Size() int {
	return len(ac.children)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

package main

const INF = int(2e18)

func longestValidSubstring(word string, forbidden []string) int {
	acm := NewACAutoMatonMap()
	for _, w := range forbidden {
		acm.AddString(w)
	}
	acm.BuildSuffixLink()

	minLen := make([]int, acm.Size()) // 每个状态匹配到的模式串的最小长度
	for i := range minLen {
		minLen[i] = INF
	}
	for i, pos := range acm.WordPos {
		minLen[pos] = min(minLen[pos], len(forbidden[i]))
	}
	acm.Dp(func(from, to int) { minLen[to] = min(minLen[to], minLen[from]) })

	res, left, pos := 0, 0, 0
	for right, char := range word {
		pos = acm.Move(pos, int(char))
		left = max(left, right-minLen[pos]+2)
		res = max(res, right-left+1)
	}
	return res
}

type ACAutoMatonMap struct {
	WordPos    []int         // WordPos[i] 表示加入的第i个模式串对应的节点编号.
	children   []map[int]int // children[v][c] 表示节点v通过字符c转移到的节点.
	suffixLink []int         // 又叫fail.指向当前节点最长真后缀对应结点.
	bfsOrder   []int         // 结点的拓扑序,0表示虚拟节点.
}

func NewACAutoMatonMap() *ACAutoMatonMap {
	return &ACAutoMatonMap{
		WordPos:  []int{},
		children: []map[int]int{{}},
	}
}

func (ac *ACAutoMatonMap) AddString(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := 0
	for _, char := range str {
		ord := int(char)
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			pos = next
		} else {
			nextState := len(ac.children)
			nexts[ord] = nextState
			pos = nextState
			ac.children = append(ac.children, map[int]int{})
		}
	}
	ac.WordPos = append(ac.WordPos, pos)
	return pos
}

func (ac *ACAutoMatonMap) AddChar(pos int, ord int) int {
	nexts := ac.children[pos]
	if next, ok := nexts[ord]; ok {
		return next
	}
	nextState := len(ac.children)
	nexts[ord] = nextState
	ac.children = append(ac.children, map[int]int{})
	return nextState
}

func (ac *ACAutoMatonMap) Move(pos int, ord int) int {
	for {
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			return next
		}
		if pos == 0 {
			return 0
		}
		pos = ac.suffixLink[pos]
	}
}

func (ac *ACAutoMatonMap) BuildSuffixLink() {
	ac.suffixLink = make([]int, len(ac.children))
	for i := range ac.suffixLink {
		ac.suffixLink[i] = -1
	}
	ac.bfsOrder = make([]int, len(ac.children))
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
			f(ac.suffixLink[v], v)
		}
	}
}

func (ac *ACAutoMatonMap) Size() int {
	return len(ac.children)
}

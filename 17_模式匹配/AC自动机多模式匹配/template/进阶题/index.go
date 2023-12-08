package main

func main() {

}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个字符串.
type ACAutoMatonArray struct {
	WordPos            []int   // WordPos[i] 表示加入的第i个模式串对应的节点编号.
	Parent             []int   // parent[v] 表示节点v的父节点.
	sigma              int     // 字符集大小.
	offset             int     // 字符集的偏移量.
	children           [][]int // children[v][c] 表示节点v通过字符c转移到的节点.
	suffixLink         []int   // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	bfsOrder           []int   // 结点的拓扑序,0表示虚拟节点.
	needUpdateChildren bool    // 是否需要更新children数组.
}

func NewACAutoMatonArray(sigma, offset int) *ACAutoMatonArray {
	res := &ACAutoMatonArray{sigma: sigma, offset: offset}
	res.newNode()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *ACAutoMatonArray) AddString(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := 0
	for _, s := range str {
		ord := int(s) - trie.offset
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
			trie.Parent = append(trie.Parent, pos)
		}
		pos = trie.children[pos][ord]
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// 在pos位置添加一个字符，返回新的节点编号.
func (trie *ACAutoMatonArray) AddChar(pos int, ord int) int {
	ord -= trie.offset
	if trie.children[pos][ord] != -1 {
		return trie.children[pos][ord]
	}
	trie.children[pos][ord] = trie.newNode()
	trie.Parent = append(trie.Parent, pos)
	return trie.children[pos][ord]
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray) Move(pos int, ord int) int {
	ord -= trie.offset
	if trie.needUpdateChildren {
		return trie.children[pos][ord]
	}
	for {
		nexts := trie.children[pos]
		if nexts[ord] != -1 {
			return nexts[ord]
		}
		if pos == 0 {
			return 0
		}
		pos = trie.suffixLink[pos]
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray) Size() int {
	return len(trie.children)
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *ACAutoMatonArray) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.suffixLink = make([]int, len(trie.children))
	for i := range trie.suffixLink {
		trie.suffixLink[i] = -1
	}
	trie.bfsOrder = make([]int, len(trie.children))
	head, tail := 0, 0
	trie.bfsOrder[tail] = 0
	tail++
	for head < tail {
		v := trie.bfsOrder[head]
		head++
		for i, next := range trie.children[v] {
			if next == -1 {
				continue
			}
			trie.bfsOrder[tail] = next
			tail++
			f := trie.suffixLink[v]
			for f != -1 && trie.children[f][i] == -1 {
				f = trie.suffixLink[f]
			}
			trie.suffixLink[next] = f
			if f == -1 {
				trie.suffixLink[next] = 0
			} else {
				trie.suffixLink[next] = trie.children[f][i]
			}
		}
	}
	if !needUpdateChildren {
		return
	}
	for _, v := range trie.bfsOrder {
		for i, next := range trie.children[v] {
			if next == -1 {
				f := trie.suffixLink[v]
				if f == -1 {
					trie.children[v][i] = 0
				} else {
					trie.children[v][i] = trie.children[f][i]
				}
			}
		}
	}
}

// 获取每个状态包含的模式串的个数.
func (trie *ACAutoMatonArray) GetCounter() []int {
	counter := make([]int, len(trie.children))
	for _, pos := range trie.WordPos {
		counter[pos]++
	}
	for _, v := range trie.bfsOrder {
		if v != 0 {
			counter[v] += counter[trie.suffixLink[v]]
		}
	}
	return counter
}

// 获取每个状态包含的模式串的索引.
func (trie *ACAutoMatonArray) GetIndexes() [][]int {
	res := make([][]int, len(trie.children))
	for i, pos := range trie.WordPos {
		res[pos] = append(res[pos], i)
	}
	for _, v := range trie.bfsOrder {
		if v != 0 {
			from, to := trie.suffixLink[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int, 0, len(arr1)+len(arr2))
			i, j := 0, 0
			for i < len(arr1) && j < len(arr2) {
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

// 按照拓扑序进行转移(EnumerateFail).
func (trie *ACAutoMatonArray) Dp(f func(from, to int)) {
	for _, v := range trie.bfsOrder {
		if v != 0 {
			f(trie.suffixLink[v], v)
		}
	}
}

func (trie *ACAutoMatonArray) newNode() int {
	trie.Parent = append(trie.Parent, -1)
	nexts := make([]int, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.children = append(trie.children, nexts)
	return len(trie.children) - 1
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

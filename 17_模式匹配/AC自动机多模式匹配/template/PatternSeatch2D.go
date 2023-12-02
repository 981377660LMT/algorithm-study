// 2D Pattern Search (Baker-Bird Algorithm)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ALDS1_14_C
// row,col<=1000
func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row1, col1 int
	fmt.Fscan(in, &row1, &col1)

	ords1 := make([][]int, row1)
	for i := 0; i < row1; i++ {
		ords1[i] = make([]int, col1)
		var s string
		fmt.Fscan(in, &s)
		for j, v := range s {
			ords1[i][j] = int(v)
		}
	}

	var winRow, winCol int
	fmt.Fscan(in, &winRow, &winCol)
	ords2 := make([][]int, winRow)
	for i := 0; i < winRow; i++ {
		ords2[i] = make([]int, winCol)
		var s string
		fmt.Fscan(in, &s)
		for j, v := range s {
			ords2[i][j] = int(v)
		}
	}

	res := PatternSearch2D(ords1, ords2)
	sort.Slice(res, func(i, j int) bool {
		if res[i][0] == res[j][0] {
			return res[i][1] < res[j][1]
		}
		return res[i][0] < res[j][0]
	})
	for _, p := range res {
		fmt.Fprintln(out, p[0], p[1])
	}
}

// 二维模式串匹配.
// 返回所有匹配的左上角坐标.
// O(row1*col1+row2*col2)
func PatternSearch2D(
	longer [][]int,
	shorter [][]int,
) [][2]int {
	aho := NewACAutoMatonMap()
	for _, s := range shorter {
		aho.AddString(s)
	}
	aho.BuildSuffixLink()
	longerStates := make([][]int, len(longer[0]))
	for i := range longerStates {
		longerStates[i] = make([]int, len(longer))
	}
	for i, row := range longer {
		pos := 0
		for j, ord := range row {
			pos = aho.Move(pos, ord)
			longerStates[j][i] = pos
		}
	}

	pf := PrefixFunction(aho.WordPos)
	res := [][2]int{}
	for j := range longerStates {
		match := KMP(longerStates[j], aho.WordPos, pf)
		for _, i := range match {
			res = append(res, [2]int{i, j - len(shorter[0]) + 1})
		}
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

func (ac *ACAutoMatonMap) AddString(s []int) int {
	if len(s) == 0 {
		return 0
	}
	pos := 0
	for i := 0; i < len(s); i++ {
		ord := s[i]
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

func PrefixFunction(s []int) []int {
	n := len(s)
	res := make([]int, n)
	len := 0
	for i := 1; i < n; i++ {
		if s[i] == s[len] {
			len++
			res[i] = len
		} else {
			if len != 0 {
				len = res[len-1]
				i--
			} else {
				res[i] = 0
			}
		}
	}
	return res
}

func KMP(longer []int, shorter []int, pf []int) []int {
	n, m := len(longer), len(shorter)
	match := []int{}
	i, j := 0, 0
	for i < n {
		if shorter[j] == longer[i] {
			i++
			j++
		}
		if j == m {
			match = append(match, i-j)
			j = pf[j-1]
		} else if i < n && shorter[j] != longer[i] {
			if j != 0 {
				j = pf[j-1]
			} else {
				i++
			}
		}
	}
	return match
}

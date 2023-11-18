package main

import "fmt"

// 参考python

const INF = int(2e18)

type ACAutoMatonMap struct {
	words      []int
	bfsOrder   []int
	children   []map[rune]int
	suffixLink []int
}

func NewACAutoMatonMap() *ACAutoMatonMap {
	return &ACAutoMatonMap{
		words:      []int{},
		bfsOrder:   []int{},
		children:   []map[rune]int{{}},
		suffixLink: []int{},
	}
}

func (ac *ACAutoMatonMap) AddString(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := 0
	for _, char := range str {
		nexts := ac.children[pos]
		if next, ok := nexts[char]; ok {
			pos = next
		} else {
			nextState := len(ac.children)
			nexts[char] = nextState
			pos = nextState
			ac.children = append(ac.children, map[rune]int{})
		}
	}
	ac.words = append(ac.words, pos)
	return pos
}

func (ac *ACAutoMatonMap) AddChar(pos int, char rune) int {
	nexts := ac.children[pos]
	if next, ok := nexts[char]; ok {
		return next
	}
	nextState := len(ac.children)
	nexts[char] = nextState
	ac.children = append(ac.children, map[rune]int{})
	return nextState
}

func (ac *ACAutoMatonMap) Move(pos int, char rune) int {
	for {
		nexts := ac.children[pos]
		if next, ok := nexts[char]; ok {
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
			ac.suffixLink[next] = f
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
	for _, pos := range ac.words {
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
	for i, pos := range ac.words {
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

func main() {
	fmt.Println("Hello, World!")
}

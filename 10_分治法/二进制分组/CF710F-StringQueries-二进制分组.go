// https://codeforces.com/contest/710/submission/187615267
// https://www.cnblogs.com/TianMeng-hyl/p/14989441.html
// https://www.cnblogs.com/Dfkuaid-210/p/bit_divide.html
// https://codeforces.com/contest/710/submission/187615267

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 1 s : 在数据结构中插入 s
// 2 s : 在数据结构中删除 s
// 3 s : 在数据结构中查询 s 出现的次数
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	var q int
	fmt.Fscan(in, &q)

	bg1 := NewBinaryGrouping(func() IPreprocessor {
		return NewACAutoMatonArray(26, 'a')
	})
	bg2 := NewBinaryGrouping(func() IPreprocessor {
		return NewACAutoMatonArray(26, 'a')
	})

	for i := 0; i < q; i++ {
		var op int
		var s string
		fmt.Fscan(in, &op, &s)
		if op == 1 {
			bg1.Add(s)
		} else if op == 2 {
			bg2.Add(s)
		} else {
			res := 0
			bg1.Query(func(p IPreprocessor) {
				pos := 0
				for _, c := range s {
					pos = p.(*SimpleACAutoMatonArray).Move(pos, int(c))
					res += int(p.(*SimpleACAutoMatonArray).count[pos])
				}
			})
			bg2.Query(func(p IPreprocessor) {
				pos := 0
				for _, c := range s {
					pos = p.(*SimpleACAutoMatonArray).Move(pos, int(c))
					res -= int(p.(*SimpleACAutoMatonArray).count[pos])
				}
			})
			out.Flush() // 强制在线，需要刷新缓冲区
		}
	}
}

type V = string

type IPreprocessor interface {
	Add(value V)
	Build()
	Clear()
}

type BinaryGrouping struct {
	groups             [][]V
	preprocessors      []IPreprocessor
	createPreprocessor func() IPreprocessor
}

func NewBinaryGrouping(createPreprocessor func() IPreprocessor) *BinaryGrouping {
	return &BinaryGrouping{
		createPreprocessor: createPreprocessor,
	}
}

func (b *BinaryGrouping) Add(value V) {
	k := 0
	for k < len(b.groups) && len(b.groups[k]) > 0 {
		k++
	}
	if k == len(b.groups) {
		b.groups = append(b.groups, []V{})
		b.preprocessors = append(b.preprocessors, b.createPreprocessor())
	}
	b.groups[k] = append(b.groups[k], value)
	b.preprocessors[k].Add(value)
	for i := 0; i < k; i++ {
		for _, v := range b.groups[i] {
			b.preprocessors[k].Add(v)
		}
		b.groups[k] = append(b.groups[k], b.groups[i]...)
		b.preprocessors[i].Clear()
		b.groups[i] = b.groups[i][:0]
	}
}

func (b *BinaryGrouping) Query(onQuery func(p IPreprocessor)) {
	for i := 0; i < len(b.preprocessors); i++ {
		onQuery(b.preprocessors[i])
	}
}

// 满足IPreprocessor接口的AC自动机.
type SimpleACAutoMatonArray struct {
	sigma    int32     // 字符集大小.
	offset   int32     // 字符集的偏移量.
	count    []int32   // count[i] 表示第i个状态匹配到的个数.
	children [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
}

func NewACAutoMatonArray(sigma, offset int) *SimpleACAutoMatonArray {
	res := &SimpleACAutoMatonArray{sigma: int32(sigma), offset: int32(offset)}
	res.newNode()
	return res
}

// 添加一个字符串.
func (trie *SimpleACAutoMatonArray) Add(str string) {
	pos := int32(0)
	for _, s := range str {
		ord := int32(s) - trie.offset
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
		}
		pos = (trie.children[pos][ord])
	}
	trie.count[pos]++
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *SimpleACAutoMatonArray) Move(pos int, ord int) int {
	ord -= int(trie.offset)
	return int(trie.children[pos][ord])
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *SimpleACAutoMatonArray) Size() int {
	return len(trie.children)
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *SimpleACAutoMatonArray) Build() {
	suffixLink := make([]int32, len(trie.children))
	for i := range suffixLink {
		suffixLink[i] = -1
	}
	bfsOrder := make([]int32, len(trie.children))
	head, tail := 0, 0
	bfsOrder[tail] = 0
	tail++
	for head < tail {
		v := bfsOrder[head]
		head++
		for i, next := range trie.children[v] {
			if next == -1 {
				continue
			}
			bfsOrder[tail] = next
			tail++
			f := suffixLink[v]
			for f != -1 && trie.children[f][i] == -1 {
				f = suffixLink[f]
			}
			suffixLink[next] = f
			if f == -1 {
				suffixLink[next] = 0
			} else {
				suffixLink[next] = trie.children[f][i]
			}
		}
	}
	for _, v := range bfsOrder {
		for i, next := range trie.children[v] {
			if next == -1 {
				f := suffixLink[v]
				if f == -1 {
					trie.children[v][i] = 0
				} else {
					trie.children[v][i] = trie.children[f][i]
				}
			}
		}
	}
}

func (trie *SimpleACAutoMatonArray) Clear() {
	trie.count = trie.count[:1]
	trie.children = trie.children[:1]
}

func (trie *SimpleACAutoMatonArray) newNode() int32 {
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.children = append(trie.children, nexts)
	trie.count = append(trie.count, 0)
	return int32(len(trie.children) - 1)
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

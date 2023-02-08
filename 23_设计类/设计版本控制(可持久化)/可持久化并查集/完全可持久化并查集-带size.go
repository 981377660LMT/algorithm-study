// # https://atcoder.jp/contests/code-thanks-festival-2017-open/tasks/code_thanks_festival_2017_h
// # 给定n个集合,初始时第i个集合只有一个元素i (i=1,2,...,n)
// # 之后进行m次合并操作,每次合并ai和bi所在的集合
// # 如果ai和bi在同一个集合,则无事发生
// # 给定q个询问,问ai和bi是在第几次操作后第一次连通的,如果不连通则输出-1

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	uf := NewPersistentUnionfind()
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		uf.Union(uf.CurVersion, a, b)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		if !uf.IsConnected(uf.CurVersion, a, b) {
			fmt.Fprintln(out, -1)
			continue
		}

		left, right := 0, uf.CurVersion
		for left <= right {
			mid := (left + right) / 2
			if uf.IsConnected(mid, a, b) {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		fmt.Fprintln(out, left)
	}
}

type PersistentUnionfind struct {
	CurVersion int
	par        *PersistentArray
	savePoints []*AryNode
}

func NewPersistentUnionfind() *PersistentUnionfind {
	savePoints := []*AryNode{nil}
	return &PersistentUnionfind{
		par:        NewPersistentArray(4, -1), // 3或者4
		savePoints: savePoints,
	}
}

// 合并x和y所在的集合,返回当前版本号。
//  0<=version<=curVersion
//  每进行一次合并操作,版本号+1;
//  Union(curVersion,0,0)表示不进行任何操作,但是会在savePoints中增加一个版本号。
func (p *PersistentUnionfind) Union(version, x, y int) int {
	x, y = p.Find(version, x), p.Find(version, y)
	ptr := p.savePoints[version]
	if x != y {
		sizeX := -p.par.Get(p.savePoints[version], x)
		sizeY := -p.par.Get(p.savePoints[version], y)
		if sizeX > sizeY {
			p.par.ch(&ptr, x, -sizeX-sizeY)
			p.par.ch(&ptr, y, x)
		} else {
			p.par.ch(&ptr, y, -sizeX-sizeY)
			p.par.ch(&ptr, x, y)
		}
	}

	p.savePoints = append(p.savePoints, ptr)
	p.CurVersion++
	return p.CurVersion
}

//  0<=version<=curVersion
func (p *PersistentUnionfind) Find(version, x int) int {
	y := p.par.Get(p.savePoints[version], x)
	if y < 0 {
		return x
	}
	return p.Find(version, y)
}

//  0<=version<=curVersion
func (p *PersistentUnionfind) IsConnected(version, x, y int) bool {
	return p.Find(version, x) == p.Find(version, y)
}

//  0<=version<=curVersion
func (p *PersistentUnionfind) GetSize(version, x int) int {
	return -p.par.Get(p.savePoints[version], p.Find(version, x))
}

type E = int

type AryNode struct {
	data     E
	children []*AryNode
}

type PersistentArray struct {
	log  int
	null E
	Root *AryNode // 初始版本
}

// aryLog 一般取3或4
func NewPersistentArray(aryLog int, null E) *PersistentArray {
	res := &PersistentArray{log: aryLog, null: null}
	return res
}

func (p *PersistentArray) Build(nums []E) *AryNode {
	for i := 0; i < len(nums); i++ {
		p.Root = p.Set(p.Root, i, nums[i])
	}
	return p.Root
}

func (p *PersistentArray) Get(version *AryNode, index int) E {
	if version == nil {
		return p.null
	}
	if index == 0 {
		return version.data
	}
	return p.Get(version.children[index&((1<<p.log)-1)], index>>p.log)
}

func (p *PersistentArray) Set(version *AryNode, index int, data E) *AryNode {
	newNode := AryNode{data: p.null, children: make([]*AryNode, 1<<p.log)}
	if version != nil {
		copy(newNode.children, version.children)
		newNode.data = version.data
	}

	version = &newNode
	if index == 0 {
		version.data = data
	} else {
		ptr := p.Set(version.children[index&((1<<p.log)-1)], index>>p.log, data)
		version.children[index&((1<<p.log)-1)] = ptr
	}
	return version
}

func (p *PersistentArray) ch(version **AryNode, index int, data E) {
	*version = p.Set(*version, index, data)
}

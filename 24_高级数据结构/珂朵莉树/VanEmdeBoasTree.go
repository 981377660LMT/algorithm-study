// Van Emde Boas Tree (梵峨眉大悲寺树)

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	demo := func() {
		van := NewVanEmdeBoasTree()
		fmt.Println(van.Min(), van.Max(), van.Size(), van.Prev(0), van.Next(0))
		van.Insert(1)
		van.Insert(2)
		fmt.Println(van.Has(1), van.Has(2), van.Has(3))
		fmt.Println(van.Prev(1), van.Prev(2), van.Prev(3), van.Prev(100), van.Prev(-1000000000000))
		van.Erase(1)
		fmt.Println(van.Min(), van.Max(), van.Size())
		van.Insert(-111)
		fmt.Println(van.Min(), van.Max(), van.Size())
		fmt.Println(van.Has(-111), van.Prev(-1), van.Min())

		n := int(5e5)
		fs := NewVanEmdeBoasTree()
		time1 := time.Now()
		for i := 0; i < n; i++ {
			fs.Insert(i)
			fs.Next(i)
			fs.Prev(i)
			fs.Has(i)
			fs.Erase(i)
			fs.Insert(i)
		}
		fmt.Println(time.Since(time1)) // !5e5 => 234ms(depth=8)
	}
	_ = demo
	// demo()

	// https://judge.yosupo.jp/problem/predecessor_problem
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	set := NewVanEmdeBoasTree()
	var s string
	fmt.Fscan(in, &s)
	for i, v := range s {
		if v == '1' {
			set.Insert(i)
		}
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		switch op {
		case 0:
			var k int
			fmt.Fscan(in, &k)
			set.Insert(k)
		case 1:
			var k int
			fmt.Fscan(in, &k)
			set.Erase(k)
		case 2:
			var k int
			fmt.Fscan(in, &k)
			if set.Has(k) {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		case 3:
			var k int
			fmt.Fscan(in, &k)
			ceiling := set.Next(k)
			if ceiling != INF {
				fmt.Fprintln(out, ceiling)
			} else {
				fmt.Fprintln(out, -1)
			}
		case 4:
			var k int
			fmt.Fscan(in, &k)
			floor := set.Prev(k)
			if floor != -INF {
				fmt.Fprintln(out, floor)
			} else {
				fmt.Fprintln(out, -1)
			}

		}
	}
}

const INF int = 1e18

type VanEmdeBoasTree struct {
	root *VNode
	size int
}

// !建立一个元素范围为(-INF,INF)的VanEmdeBoasTree.
func NewVanEmdeBoasTree() *VanEmdeBoasTree {
	return &VanEmdeBoasTree{root: NewVNode(32)} // 16/32/64
}

func (van *VanEmdeBoasTree) Has(x int) bool {
	return van.root.Has(x)
}

func (van *VanEmdeBoasTree) Insert(x int) bool {
	if van.Has(x) {
		return false
	}
	van.size++
	van.root.Insert(x)
	return true
}

func (van *VanEmdeBoasTree) Erase(x int) bool {
	if !van.Has(x) {
		return false
	}
	van.size--
	van.root.Erase(x)
	return true
}

// 返回小于等于i的最大元素.如果不存在,返回-INF.
func (van *VanEmdeBoasTree) Prev(x int) int {
	return van.root.Prev(x)
}

// 返回大于等于i的最小元素.如果不存在,返回INF.
func (van *VanEmdeBoasTree) Next(x int) int {
	return van.root.Next(x)
}

func (van *VanEmdeBoasTree) Size() int {
	return van.size
}

// 如果没有元素,返回INF.
func (van *VanEmdeBoasTree) Min() int {
	return van.root.min
}

// 如果没有元素,返回-INF.
func (van *VanEmdeBoasTree) Max() int {
	return van.root.max
}

type VNode struct {
	min, max, dep int
	aux           *VNode
	son           map[int]*VNode
}

func NewVNode(dep int) *VNode {
	return &VNode{min: INF, max: -INF, dep: dep, son: map[int]*VNode{}}
}

func (v *VNode) Has(x int) bool {
	vMin, vMax := v.min, v.max
	if x == vMin || x == vMax {
		return true
	}
	vDep := v.dep
	if x < vMin || x > vMax || vDep == 0 {
		return false
	}
	i := x >> vDep
	soni, ok := v.son[i]
	if !ok {
		return false
	}
	return soni.Has(x - (i << vDep))
}

func (v *VNode) Insert(x int) {
	vMin, vMax := v.min, v.max
	if vMin > vMax {
		v.min, v.max = x, x
		return
	}
	if min, max := vMin, vMax; min == max {
		if x < min {
			v.min = x
			return
		} else if x > max {
			v.max = x
			return
		}
	}
	if x < vMin {
		x, v.min = vMin, x
	}
	if x > vMax {
		x, v.max = vMax, x
	}
	vDep := v.dep
	i := x >> vDep
	soni, ok := v.son[i]
	if !ok {
		soni = NewVNode(vDep >> 1)
		v.son[i] = soni
	}
	if soni.Empty() {
		if v.aux == nil {
			v.aux = NewVNode(vDep >> 1)
		}
		v.aux.Insert(i)
	}
	soni.Insert(x - (i << vDep))
}

func (v *VNode) Erase(x int) {
	vMin, vMax := v.min, v.max
	if vMin == x && vMax == x {
		v.min, v.max = INF, -INF
		return
	}

	vDep := v.dep
	aux := v.aux
	if x == v.min {
		if aux == nil || aux.Empty() {
			v.min = vMax
			return
		}
		auxMin := aux.min
		x = (auxMin << vDep) + v.son[auxMin].min
		v.min = x
	}
	if x == vMax {
		if aux == nil || aux.Empty() {
			v.max = vMin
			return
		}
		auxMax := aux.max
		x = (auxMax << vDep) + v.son[auxMax].max
		v.max = x
	}

	i := x >> vDep
	soni := v.son[i]
	soni.Erase(x - (i << vDep))
	if soni.Empty() {
		aux.Erase(i)
	}
}

func (v *VNode) Prev(x int) int {
	vMin := v.min
	if x < vMin {
		return -INF
	}
	vMax := v.max
	if x >= vMax {
		return vMax
	}
	vDep := v.dep
	i := x >> vDep
	hi := i << vDep
	lo := x - hi
	soni, ok := v.son[i]
	if ok && lo >= soni.min {
		return hi + soni.Prev(lo)
	}
	var y int
	if v.aux != nil && i > 0 {
		y = v.aux.Prev(i - 1)
	} else {
		y = -INF
	}
	if y == -INF {
		return vMin
	}
	return (y << vDep) + v.son[y].max
}

func (v *VNode) Next(x int) int {
	vMin := v.min
	if x <= vMin {
		return vMin
	}
	vMax := v.max
	if x > vMax {
		return INF
	}
	vDep := v.dep
	i := x >> vDep
	hi := i << vDep
	lo := x - hi
	soni, ok := v.son[i]
	if ok && lo <= soni.max {
		return hi + soni.Next(lo)
	}
	var y int
	if v.aux != nil {
		y = v.aux.Next(i + 1)
	} else {
		y = INF
	}
	if y == INF {
		return vMax
	}
	return (y << vDep) + v.son[y].min
}

func (v *VNode) Empty() bool {
	return v.min > v.max
}

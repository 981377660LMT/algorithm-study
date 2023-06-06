// 珂朵莉树(ODT)/Intervals

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const INF int = 1e18

func demo() {
	odt := NewODTVan(-INF)
	odt.Set(0, 3, 1)
	odt.Set(3, 5, 2)
	fmt.Println(odt.Len, odt.Count, odt)
}

func UnionOfInterval() {
	// https://atcoder.jp/contests/abc256/tasks/abc256_d
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	odt := NewODTVan(-INF)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		odt.Set(l, r, 1)
	}
	odt.EnumerateAll(func(l, r int, x int) {
		if x == 1 {
			fmt.Fprintln(out, l, r)
		}
	})
}

func main() {
	// demo()
	UnionOfInterval()
}

type Value = int

type ODTVan struct {
	Len        int // 区间数
	Count      int // 区间元素个数之和
	llim, rlim int
	noneValue  Value
	data       map[int]Value
	ss         *_Van
}

// 指定哨兵 noneValue 建立一个 ODT.
//  区间为[-INF,INF).
func NewODTVan(noneValue Value) *ODTVan {
	res := &ODTVan{
		rlim:      INF,
		llim:      -INF,
		noneValue: noneValue,
		data:      make(map[int]Value),
		ss:        _NewVan(),
	}
	return res
}

// 返回包含 x 的区间的信息.
func (odt *ODTVan) Get(x int, erase bool) (start, end int, value Value) {
	start, end = odt.ss.Prev(x), odt.ss.Next(x+1)
	v := odt._getOrNone(start)
	if erase && v != odt.noneValue {
		odt.Len--
		odt.Count -= end - start
		odt.data[start] = odt.noneValue
		odt.mergeAt(start)
		odt.mergeAt(end)
	}
	return
}

func (odt *ODTVan) Set(start, end int, value Value) {
	odt.EnumerateRange(start, end, func(l, r int, x Value) {}, true)
	odt.ss.Insert(start)
	odt.data[start] = value
	if value != odt.noneValue {
		odt.Len++
		odt.Count += end - start
	}
	odt.mergeAt(start)
	odt.mergeAt(end)
}

func (odt *ODTVan) EnumerateAll(f func(start, end int, value Value)) {
	odt.EnumerateRange(odt.llim, odt.rlim, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *ODTVan) EnumerateRange(start, end int, f func(start, end int, value Value), erase bool) {
	if !(odt.llim <= start && start <= end && end <= odt.rlim) {
		panic(fmt.Sprintf("invalid range [%d, %d)", start, end))
	}

	NONE := odt.noneValue
	if !erase {
		l := odt.ss.Prev(start)
		for l < end {
			r := odt.ss.Next(l + 1)
			f(max(l, start), min(r, end), odt._getOrNone(l))
			l = r
		}
		return
	}

	// 分割区间
	p := odt.ss.Prev(start)
	if p < start {
		odt.ss.Insert(start)
		v := odt._getOrNone(p)
		odt.data[start] = v
		if v != NONE {
			odt.Len++
		}
	}
	p = odt.ss.Next(end)
	if end < p {
		v := odt._getOrNone(odt.ss.Prev(end))
		odt.data[end] = v
		odt.ss.Insert(end)
		if v != NONE {
			odt.Len++
		}
	}
	p = start
	for p < end {
		q := odt.ss.Next(p + 1)
		v := odt._getOrNone(p)
		f(p, q, v)
		if v != NONE {
			odt.Len--
			odt.Count -= q - p
		}
		odt.ss.Erase(p)
		p = q
	}
	odt.ss.Insert(start)
	odt.data[start] = NONE
}

func (odt *ODTVan) String() string {
	sb := []string{}
	odt.EnumerateAll(func(start, end int, value Value) {
		var v interface{} = value
		if value == odt.noneValue {
			v = "nil"
		}
		sb = append(sb, fmt.Sprintf("[%d,%d):%v", start, end, v))
	})
	return fmt.Sprintf("ODT{%v}", strings.Join(sb, ", "))
}

func (odt *ODTVan) mergeAt(p int) {
	if p <= 0 || odt.rlim <= p {
		return
	}
	q := odt.ss.Prev(p - 1)
	if dataP, dataQ := odt._getOrNone(p), odt._getOrNone(q); dataP == dataQ {
		if dataP != odt.noneValue {
			odt.Len--
		}
		odt.ss.Erase(p)
	}
}

func (odt *ODTVan) _getOrNone(key int) Value {
	if value, ok := odt.data[key]; ok {
		return value
	}
	return odt.noneValue
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

type _Van struct {
	root *VNode
	size int
}

func _NewVan() *_Van {
	return &_Van{root: NewVNode(32)} // 16/32
}

func (van *_Van) Has(x int) bool {
	return van.root.Has(x)
}

func (van *_Van) Insert(x int) bool {
	if van.Has(x) {
		return false
	}
	van.size++
	van.root.Insert(x)
	return true
}

func (van *_Van) Erase(x int) bool {
	if !van.Has(x) {
		return false
	}
	van.size--
	van.root.Erase(x)
	return true
}

// 返回小于等于i的最大元素.如果不存在,返回-INF.
func (van *_Van) Prev(x int) int {
	return van.root.Prev(x)
}

// 返回大于等于i的最小元素.如果不存在,返回INF.
func (van *_Van) Next(x int) int {
	return van.root.Next(x)
}

func (van *_Van) Size() int {
	return van.size
}

// 如果没有元素,返回INF.
func (van *_Van) Min() int {
	return van.root.min
}

// 如果没有元素,返回-INF.
func (van *_Van) Max() int {
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

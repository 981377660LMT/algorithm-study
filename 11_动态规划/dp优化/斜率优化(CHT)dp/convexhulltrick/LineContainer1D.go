// 注意cpp里的迭代器:
// !Begin指向第一个元素,
// !End指向最后一个元素的下一个位置,
// 这里的迭代器设计为:
// !Begin指向第一个元素的前一个位置,First指向第一个元素
// !Last指向最后一个元素,End指向最后一个元素的下一个位置

// 这里删除元素会引起迭代器失效

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func lineAddGetMin() {
	// https://judge.yosupo.jp/problem/line_add_get_min
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	L := NewLineContainer1D(true)
	for i := 0; i < n; i++ {
		var k, m int
		fmt.Fscan(in, &k, &m)
		L.Add(k, m)
	}
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			L.Add(a, b)
		} else {
			var x int
			fmt.Fscan(in, &x)
			fmt.Fprintln(out, L.Query(x))
		}
	}
}

func main() {
	lineAddGetMin()
}

const INF int = 2e18

type Line struct {
	k, m, p int
}

type LineContainer1D struct {
	minimize bool
	sl       *SpecializedSortedList
}

func NewLineContainer1D(minimize bool) *LineContainer1D {
	return &LineContainer1D{
		minimize: minimize,
		sl:       NewSpecializedSortedList(func(a, b *Line) bool { return a.k < b.k }),
	}
}

// 向集合中添加一条线，表示为y = kx + m
func (lc *LineContainer1D) Add(k, m int) {
	if lc.minimize {
		k, m = -k, -m
	}
	newLine := &Line{k: k, m: m}
	lc.sl.Add(newLine)
	it1 := lc.sl.BisectRightByK(k) - 1
	it2 := it1
	it1++
	for lc.insect(it2, it1) {
		lc.sl.Pop(it1)
	}

	it3 := it2
	if it3 != 0 {
		it3--
		if lc.insect(it3, it2) {
			lc.sl.Pop(it2)
			lc.insect(it3, it2)
		}
	}

	for it2 := it3; it2 != 0 && lc.sl.At(it3-1).p >= lc.sl.At(it2).p; it2 = it3 {
		it3--
		lc.sl.Pop(it2)
		lc.insect(it3, it2)
	}
}

// 查询 kx + m 的最小值（或最大值).
func (lc *LineContainer1D) Query(x int) int {
	// !这里有一个关键点：尽管Line<T>结构体中的operator<按k值对线性函数进行排序，
	// !但LineContainer类在维护这些线性函数时，确保了它们的交点的x坐标（p值）是有序的。
	// 这使得query函数可以通过调用lower_bound(x)来找到给定x值对应的最大（或最小）y值。
	if lc.sl.Len() == 0 {
		panic("empty container")
	}
	pos := lc.sl.BisectLeftByP(x)
	line := lc.sl.At(pos)
	v := line.k*x + line.m
	if lc.minimize {
		return -v
	}
	return v
}

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *LineContainer1D) insect(posX, posY int32) bool {
	if posY == lc.sl.Len() {
		line := lc.sl.At(posX)
		line.p = INF
		return false
	}

	line1, line2 := lc.sl.At(posX), lc.sl.At(posY)
	if line1.k == line2.k {
		if line1.m > line2.m {
			line1.p = INF
		} else {
			line1.p = -INF
		}
	} else {
		// lc_div
		a, b := line2.m-line1.m, line1.k-line2.k
		tmp := 0
		if (a^b) < 0 && a%b != 0 {
			tmp = 1
		}
		line1.p = a/b - tmp
	}
	return line1.p >= line2.p
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const _LOAD int32 = 50 // 75/100/150/200

type S = *Line

type SpecializedSortedList struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSpecializedSortedList(less func(a, b S) bool, elements ...S) *SpecializedSortedList {
	elements = append(elements[:0:0], elements...)
	res := &SpecializedSortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := int32(len(elements))
	blocks := [][]S{}
	for start := int32(0); start < n; start += _LOAD {
		end := min32(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]S, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *SpecializedSortedList) Add(value S) *SpecializedSortedList {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}

	pos, index := sl._locRight(value)

	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SpecializedSortedList) Pop(index int32) {
	pos, startIndex := sl._findKth(index)
	sl._delete(pos, startIndex)
}

func (sl *SpecializedSortedList) At(index int32) S {
	if index < 0 || index >= sl.size {
		return nil
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SpecializedSortedList) BisectRightByK(k int) int32 {
	pos, index := sl._locRightByK(k)
	return sl._queryTree(pos) + index
}

func (sl *SpecializedSortedList) BisectLeftByP(p int) int32 {
	pos, index := sl._locLeftByP(p)
	return sl._queryTree(pos) + index
}

func (sl *SpecializedSortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SpecializedSortedList) Len() int32 {
	return sl.size
}

func (sl *SpecializedSortedList) _delete(pos, index int32) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], sl.blocks[pos][index+1:]...)
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
	sl.shouldRebuildTree = true
}

func (sl *SpecializedSortedList) _locLeftByP(p int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(-1)
	right := int32(len(sl.blocks) - 1)
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.mins[mid].p >= p {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		last := block[len(block)-1]
		if last.p >= p {
			right--
		}
	}
	pos = right

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if cur[mid].p >= p {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _locRight(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, sl.mins[mid]) {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, cur[mid]) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _locRightByK(k int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if k < sl.mins[mid].k {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if k < cur[mid].k {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _buildTree() {
	sl.tree = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = int32(len(sl.blocks[i]))
	}
	tree := sl.tree
	for i := 0; i < len(tree); i++ {
		j := i | (i + 1)
		if j < len(tree) {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *SpecializedSortedList) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < int32(len(tree)); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SpecializedSortedList) _queryTree(end int32) int32 {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := int32(0)
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SpecializedSortedList) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	pos = -1
	m := int32(len(tree))
	bitLength := bits.Len32(uint32(m))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

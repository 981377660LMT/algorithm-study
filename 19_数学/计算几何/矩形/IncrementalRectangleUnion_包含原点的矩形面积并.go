package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	KarenAndCards()
}

// https://www.luogu.com.cn/problem/CF815D
// 给定n张卡片，每张卡片有(a,b,c)三个属性.
// !一张卡片可以打败另一张卡片当且仅当，至少有两个属性严格大于另一张卡片的对应属性.
// !求出a<=x,b<=y,c<=z的卡片中，可以打败这n张卡片的卡片数量.
//
// 考虑无法打败的情况.
func KarenAndCards() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, x, y, z int
	fmt.Fscan(in, &n, &x, &y, &z)
	cards := make([][3]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cards[i][0], &cards[i][1], &cards[i][2])
	}

	points := make([][3]int, 0, n*3)
	for _, card := range cards {
		a, b, c := card[0], card[1], card[2]
		// 无法打败：固定一个维度，另外两个维度的值都不超过该维度的值
		points = append(points, [3]int{x, b, c}, [3]int{a, y, c}, [3]int{a, b, z})
	}
	volumn := CuboidUnionVolumn(points)
	res := x*y*z - volumn
	fmt.Println(res)
}

func demo() {
	R := NewIncrementalRectangleUnionRange()
	R.Add(-2, 1, -2, 1)
	fmt.Println(R.Query())
	R.Add(-1, 2, -1, 2)
	fmt.Println(R.Query())
}

const INF int = 2e9

// 求出所有左下角为原点的立方体的体积并.
//
//	points: [x, y, z] 每个点的坐标，表示一个[0, x] * [0, y] * [0, z]的立方体.
func CuboidUnionVolumn(points [][3]int) int {
	points = append(points[:0:0], points...)
	sort.Slice(points, func(i, j int) bool { return points[i][2] > points[j][2] })
	preZ := INF
	res := 0
	area := 0
	manager := NewIncrementalRectangleUnion()
	for _, p := range points {
		x, y, z := p[0], p[1], p[2]
		res += (preZ - z) * area
		manager.Add(x, y)
		area = manager.Query()
		preZ = z
	}
	res += preZ * area
	return res
}

type IncrementalRectangleUnionRange struct {
	ur *IncrementalRectangleUnion
	ul *IncrementalRectangleUnion
	dr *IncrementalRectangleUnion
	dl *IncrementalRectangleUnion
}

func NewIncrementalRectangleUnionRange() *IncrementalRectangleUnionRange {
	return &IncrementalRectangleUnionRange{
		ur: NewIncrementalRectangleUnion(),
		ul: NewIncrementalRectangleUnion(),
		dr: NewIncrementalRectangleUnion(),
		dl: NewIncrementalRectangleUnion(),
	}
}

// Add [x1, x2] * [y1, y2].
// x1 <= 0 <= x2.
// y1 <= 0 <= y2.
func (irur *IncrementalRectangleUnionRange) Add(x1, x2, y1, y2 int) {
	irur.ur.Add(x2, y2)
	irur.ul.Add(-x1, y2)
	irur.dr.Add(x2, -y1)
	irur.dl.Add(-x1, -y1)
}

func (irur *IncrementalRectangleUnionRange) Query() int {
	return irur.ur.Query() + irur.ul.Query() + irur.dr.Query() + irur.dl.Query()
}

// 包含原点的矩形合并/矩形面积并.
type IncrementalRectangleUnion struct {
	sum int
	sl  *_SL
}

func NewIncrementalRectangleUnion() *IncrementalRectangleUnion {
	return &IncrementalRectangleUnion{
		sl: _NewSL(func(a, b S) bool { return a.start < b.start }, S{0, INF}, S{INF, 0}),
	}
}

// Add [0, x] * [0, y].
// x >= 0.
// y >= 0.
func (iru *IncrementalRectangleUnion) Add(x, y int) {
	pos := iru.sl.BisectLeft(S{x, -INF})
	item := iru.sl.At(pos)
	if item.end >= y {
		return
	}
	nextY := item.end
	pos--
	pre := iru.sl.At(pos)
	for pre.end <= y {
		x1 := pre.start
		y1 := pre.end
		iru.sl.Pop(pos)
		pos--
		pre = iru.sl.At(pos)
		iru.sum -= (x1 - pre.start) * (y1 - nextY)
	}
	iru.sum += (x - iru.sl.At(pos).start) * (y - nextY)
	pos = iru.sl.BisectLeft(S{x, -INF})
	if iru.sl.At(pos).start == x {
		iru.sl.Pop(pos)
	}
	iru.sl.Add(S{x, y})
}

func (iru *IncrementalRectangleUnion) Query() int {
	return iru.sum
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type S = struct{ start, end int }

type _SL struct {
	less              func(a, b S) bool
	size              int
	blocks            [][]S
	mins              []S
	tree              []int
	shouldRebuildTree bool
}

func _NewSL(less func(a, b S) bool, elements ...S) *_SL {
	elements = append(elements[:0:0], elements...)
	res := &_SL{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]S{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
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

func (sl *_SL) Add(value S) *_SL {
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
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *_SL) Pop(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}

func (sl *_SL) At(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *_SL) BisectLeft(value S) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

func (sl *_SL) _delete(pos, index int) {
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

func (sl *_SL) _locLeft(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := -1
	right := len(sl.blocks) - 1
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	pos = right

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = len(cur)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(cur[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *_SL) _locRight(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := 0
	right := len(sl.blocks)
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
	right = len(cur)
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

func (sl *_SL) _buildTree() {
	sl.tree = make([]int, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = len(sl.blocks[i])
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

func (sl *_SL) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *_SL) _queryTree(end int) int {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := 0
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *_SL) _findKth(k int) (pos, index int) {
	if k < len(sl.blocks[0]) {
		return 0, k
	}
	last := len(sl.blocks) - 1
	lastLen := len(sl.blocks[last])
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	pos = -1
	bitLength := bits.Len32(uint32(len(tree)))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < len(tree) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
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

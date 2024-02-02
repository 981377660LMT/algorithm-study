// 不删除莫队,复杂度和普通莫队一样
// !删除操作很麻烦的时候使用
// 发明者:https://snuke.hatenablog.com/entry/2016/07/01/000000

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	StaticRangeInversionsQuery()
	// AT1219()
	// Luogu5906()
}

// Static Range Inversions Query - 静态区间逆序对查询
// https://judge.yosupo.jp/problem/static_range_inversions_query
func StaticRangeInversionsQuery() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	mo := NewMoRollback(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		mo.AddQuery(l, r)
	}

	// 离散化
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	mp := make(map[int]int, len(allNums))
	for i, v := range allNums {
		mp[v] = i
	}
	newNums := make([]int, n)
	for i, v := range nums {
		newNums[i] = mp[v]
	}

	bit := NewBitArray(len(allNums))
	inv, snap, snapInv := 0, 0, 0 // inv: 当前逆序对数, snap: 当前快照状态, snapInv: 当前快照逆序对数
	history := make([]int, 0, n)  // history: 当前操作历史便于undo
	res := make([]int, q)

	add := func(index, delta int) {
		if delta == 1 { // add_right ->
			x := newNums[index]
			inv += bit.QueryRange(x+1, len(allNums))
			bit.Add(x, 1)
			history = append(history, x)
		} else { // add_left <-
			x := newNums[index]
			inv += bit.QueryPrefix(x)
			bit.Add(x, 1)
			history = append(history, x)
		}
	}
	_move := func(state int) {
		for len(history) > state {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			bit.Add(x, -1)
		}
	}
	reset := func() {
		_move(0)
		inv = 0
	}
	snapshot := func() {
		snap = len(history)
		snapInv = inv
	}
	rollback := func() {
		_move(snap)
		inv = snapInv
	}
	query := func(qi int) {
		res[qi] = inv
	}

	mo.Run(add, reset, snapshot, rollback, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 歴史の研究
// https://www.luogu.com.cn/problem/AT_joisc2014_c
// 给定一个数组nums和q个查询(l,r)
// 每次查询[l,r]区间内的`重要度`,一个数字num的重要度定义为`num乘以区间内num的个数`
func AT1219() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	mo := NewMoRollback(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		mo.AddQuery(l, r)
	}

	res := make([]int, q)
	cur := 0                     // 当前区间的答案
	snap, snapCur := 0, 0        // 当前快照状态,当前快照答案
	history := make([]int, 0, n) // x
	counter := make(map[int]int)

	add := func(index, _ int) { // TODO
		x := nums[index]
		counter[x] += 1
		cur = max(cur, x*counter[x])
		history = append(history, x)
	}
	_move := func(state int) {
		for len(history) > state {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			counter[x]-- // TODO
		}
	}

	reset := func() {
		_move(0)
		cur = 0
	}
	snapshot := func() {
		snap = len(history)
		snapCur = cur
	}
	rollback := func() {
		_move(snap)
		cur = snapCur
	}
	query := func(qi int) {
		res[qi] = cur
	}

	mo.Run(add, reset, snapshot, rollback, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://www.luogu.com.cn/problem/P5906
// 给定一个序列，多次询问一段区间 [l,r]，求区间中相同的数的最远间隔距离。
// 如果区间内不存在两个数相同，则输出 0。
// 序列中两个元素的间隔距离指的是两个元素下标差的绝对值。
//
// !维护每个数在区间内索引的最大值和最小值.
func Luogu5906() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	// 离散化
	_pool := make(map[int]int)
	id := func(o int) int {
		if v, ok := _pool[o]; ok {
			return v
		}
		v := len(_pool)
		_pool[o] = v
		return v
	}
	for i, v := range nums {
		nums[i] = id(v)
	}

	var q int
	fmt.Fscan(in, &q)
	mo := NewMoRollback(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		mo.AddQuery(l, r)
	}

	res := make([]int, q)
	cur := 0                        // 当前区间的答案
	snap, snapCur := 0, 0           // 当前快照状态,当前快照答案
	history := make([][3]int, 0, n) // (x,preMinPos,preMaxPos)
	minPos := make([]int, n)
	maxPos := make([]int, n)
	for i := range minPos {
		minPos[i] = n
		maxPos[i] = -1
	}

	add := func(index, _ int) { // TODO
		x := nums[index]
		preMinPos := minPos[x]
		preMaxPos := maxPos[x]
		minPos[x] = min(preMinPos, index)
		maxPos[x] = max(preMaxPos, index)
		cur = max(cur, maxPos[x]-minPos[x])
		history = append(history, [3]int{x, preMinPos, preMaxPos})
	}
	_move := func(state int) {
		for len(history) > state {
			item := history[len(history)-1]
			x, preMinPos, preMaxPos := item[0], item[1], item[2]
			history = history[:len(history)-1]
			minPos[x] = preMinPos
			maxPos[x] = preMaxPos
		}
	}

	reset := func() {
		_move(0)
		cur = 0
	}
	snapshot := func() {
		snap = len(history)
		snapCur = cur
	}
	rollback := func() {
		_move(snap)
		cur = snapCur
	}
	query := func(qi int) {
		if cur > 0 {
			res[qi] = cur
		}
	}

	mo.Run(add, reset, snapshot, rollback, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoRollback struct {
	chunkSize          int
	left, right, order []int
}

func NewMoRollback(n, q int) *MoRollback {
	chunkSize := max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	order := make([]int, q)
	for i := range order {
		order[i] = i
	}
	return &MoRollback{chunkSize: chunkSize, order: order}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//
//	0 <= left <= right <= n
func (mo *MoRollback) AddQuery(left, right int) {
	mo.left = append(mo.left, left)
	mo.right = append(mo.right, right)
}

// 返回每个查询的结果.
//
//	add: 将数据添加到窗口.
//	reset: 将窗口重置为初始状态.
//	snapshot: 保存当前窗口的状态.
//	rollback: 恢复窗口的状态.
//	query: 查询窗口的状态.
func (mo *MoRollback) Run(
	add func(index, delta int),
	reset func(),
	snapshot func(),
	rollback func(),
	query func(qi int),
) {
	left, right, order := mo.left, mo.right, mo.order
	chunkSize := mo.chunkSize
	sort.Slice(order, func(i, j int) bool {
		ii, jj := order[i], order[j]
		iblock, jblock := left[ii]/chunkSize, left[jj]/chunkSize
		if iblock != jblock {
			return iblock < jblock
		}
		return right[ii] < right[jj]
	})

	reset()
	for _, idx := range order {
		if right[idx]-left[idx] < chunkSize {
			for i := left[idx]; i < right[idx]; i++ {
				add(i, 1)
			}
			query(idx)
			rollback()
		}
	}

	nr := 0
	lastBlock := -1
	for _, idx := range order {
		if right[idx]-left[idx] < chunkSize {
			continue
		}
		block := left[idx] / chunkSize
		if lastBlock != block {
			reset()
			lastBlock = block
			nr = (block + 1) * chunkSize
		}
		for nr < right[idx] {
			add(nr, 1)
			nr++
		}
		snapshot()
		for j := (block+1)*chunkSize - 1; j >= left[idx]; j-- {
			add(j, -1)
		}
		query(idx)
		rollback()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) MaxRight(check func(index, preSum int) bool) int {
	i := 0
	s := 0
	k := 1
	for 2*k <= b.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < b.n {
			t := s + b.data[i+k-1]
			if check(i+k, t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

// 0/1 树状数组查找第 k(0-based) 个1的位置.
// UpperBound.
func (b *BITArray) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

type FastSet struct {
	n, lg int
	seg   [][]int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.n
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

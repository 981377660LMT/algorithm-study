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
	Cf765F()
	// StaticRangeInversionsQuery()
	// StaticRangeModeQuery()

}

const INF int = 1e9 + 10

func SolveCf765F(nums []int, queries [][]int) []int {
	n, q := len(nums), len(queries)
	S := NewFastSet(n)
	order := argSort(nums)
	keys := reArrage(nums, order)
	nums = argSort(order) // 不带相同值的离散化，转换为 0-n-1

	M := NewMoRollback(n)
	for i := 0; i < q; i++ {
		M.AddQuery(queries[i][0], queries[i][1])
	}

	res := make([]int, q)
	curRes := INF
	snapState, snapRes := 0, 0
	var history []int32

	add := func(index int) {
		index = nums[index]
		left := S.Prev(index)
		right := S.Next(index)
		S.Insert(index)
		history = append(history, int32(index))
		if left != -1 {
			curRes = min(curRes, keys[index]-keys[left])
		}
		if right != n {
			curRes = min(curRes, keys[right]-keys[index])
		}
	}

	reset := func() {
		S.Enumerate(0, n, func(p int) { S.Erase(p) })
		curRes = INF
	}

	snapShot := func() {
		snapState = len(history)
		snapRes = curRes
	}

	rollback := func() {
		for len(history) > snapState {
			index := history[len(history)-1]
			history = history[:len(history)-1]
			S.Erase(int(index))
		}
		curRes = snapRes
	}

	query := func(qi int) {
		res[qi] = curRes
	}

	M.Run(add, add, reset, snapShot, rollback, query, -1)
	return res
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
	getRank, _, size := DiscretizeCompressed(nums, 0)
	for i := 0; i < n; i++ {
		nums[i] = getRank(nums[i])
	}

	mo := NewMoRollback(n)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		mo.AddQuery(l, r)
	}

	bit := NewBitArray(size)
	inv, snapState, snapInv := 0, 0, 0 // inv: 当前逆序对数, snap: 当前快照状态, snapInv: 当前快照逆序对数
	history := make([]int, 0, n)       // history: 当前操作历史便于undo
	res := make([]int, q)

	addLeft := func(index int) {
		x := nums[index]
		inv += bit.QueryPrefix(x)
		bit.Add(x, 1)
		history = append(history, x)
	}
	addRight := func(index int) {
		x := nums[index]
		inv += bit.QueryAll() - bit.QueryPrefix(x+1)
		bit.Add(x, 1)
		history = append(history, x)
	}

	_rollback := func(state int) {
		for len(history) > state {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			bit.Add(x, -1)
		}
	}
	reset := func() {
		_rollback(0)
		inv = 0
	}
	snapshot := func() {
		snapState = len(history)
		snapInv = inv
	}
	rollback := func() {
		_rollback(snapState)
		inv = snapInv
	}
	query := func(qi int) {
		res[qi] = inv
	}

	mo.Run(addLeft, addRight, reset, snapshot, rollback, query, -1)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 区间众数查询，存在多个众数时，返回任意一个.
// https://judge.yosupo.jp/problem/static_range_mode_query
func StaticRangeModeQuery() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	D := NewDictionary()
	for i := 0; i < n; i++ {
		nums[i] = D.Id(nums[i])
	}

	M := NewMoRollback(n)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		M.AddQuery(l, r)
	}
	res := make([][2]int, q) // (众数, 众数出现次数)

	counter := make([]int, D.Size())
	history := make([]int, 0, n)

	maxCount, maxKey := 0, 0
	snapState := 0
	snapCount, snapKey := 0, 0

	add := func(index int) {
		x := nums[index]
		history = append(history, x)
		counter[x]++
		if counter[x] > maxCount {
			maxCount = counter[x]
			maxKey = x
		}
	}

	reset := func() {
		for _, v := range history {
			counter[v] = 0
		}
		history = history[:0]
		maxCount, maxKey = 0, 0
	}

	snapshot := func() {
		snapState = len(history)
		snapCount, snapKey = maxCount, maxKey
	}

	rollback := func() {
		for len(history) > snapState {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			counter[x]--
		}
		maxCount, maxKey = snapCount, snapKey
	}

	query := func(qi int) {
		res[qi] = [2]int{D.Value(maxKey), maxCount}
	}

	M.Run(add, add, reset, snapshot, rollback, query, -1)
	for _, v := range res {
		fmt.Fprintln(out, v[0], v[1])
	}
}

// https://codeforces.com/problemset/submission/765/240821486
// https://www.luogu.com.cn/problem/CF765F
// https://leetcode.cn/problems/minimum-absolute-difference-queries/description/
// 给定一个数组和q组查询，每组查询包含两个整数start,end，求出[start,end)区间内的 abs(a[i]-a[j])的最小值(i!=j)。
func Cf765F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][]int, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		queries[i] = []int{l, r}
	}

	res := SolveCf765F(nums, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoRollback struct {
	left, right []int32
}

func NewMoRollback(n int) *MoRollback {
	return &MoRollback{}
}

func (mo *MoRollback) AddQuery(start, end int) {
	mo.left = append(mo.left, int32(start))
	mo.right = append(mo.right, int32(end))
}

// addLeft : 将index位置的元素加入到区间左端.
// addRight: 将index位置的元素加入到区间右端.
// reset: 重置区间.
// snapShot: 快照当前状态.
// rollback: 回滚到快照状态.
// query: 查询当前区间.
// blockSize: 分块大小.-1表示使用默认值.
func (mo *MoRollback) Run(
	addLeft func(index int),
	addRight func(index int),
	reset func(),
	snapShot func(),
	rollback func(),
	query func(qid int),
	blockSize int,
) {
	q := int32(len(mo.left))
	if q == 0 {
		return
	}
	n := int32(0)
	for _, v := range mo.right {
		n = max32(n, v)
	}
	blockSize32 := int32(blockSize)
	if blockSize32 == -1 {
		blockSize32 = int32(max32(1, n/max32(1, int32(math.Sqrt(float64(q*2/3))))))
	}
	queryId := make([][]int32, (n-1)/blockSize32+1)
	naive := func(qid int32) {
		snapShot()
		for i := mo.left[qid]; i < mo.right[qid]; i++ {
			addRight(int(i))
		}
		query(int(qid))
		rollback()
	}

	for qid := int32(0); qid < q; qid++ {
		l, r := mo.left[qid], mo.right[qid]
		iL, iR := l/blockSize32, r/blockSize32
		if iL == iR {
			naive(qid)
			continue
		}
		queryId[iL] = append(queryId[iL], qid)
	}

	for _, order := range queryId {
		if len(order) == 0 {
			continue
		}
		sort.Slice(order, func(i, j int) bool {
			return mo.right[order[i]] < mo.right[order[j]]
		})
		lMax := int32(0)
		for _, qid := range order {
			lMax = max32(lMax, mo.left[qid])
		}
		reset()
		l, r := lMax, lMax
		for _, qid := range order {
			L, R := mo.left[qid], mo.right[qid]
			for r < R {
				addRight(int(r))
				r++
			}
			snapShot()
			for L < l {
				l--
				addLeft(int(l))
			}
			query(int(qid))
			rollback()
			l = lMax
		}
	}
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
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
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

// (紧)离散化.
//
//	offset: 离散化的起始值偏移量.
//
//	getRank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//	count: 离散化(去重)后的元素个数.
func DiscretizeCompressed(nums []int, offset int) (getRank func(value int) int, getValue func(rank int) int, count int) {
	set := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	count = len(set)
	rank := make([]int, 0, count)
	for v := range set {
		rank = append(rank, v)
	}
	sort.Ints(rank)
	mp := make(map[int]int, count)
	for i, v := range rank {
		mp[v] = i + offset
	}
	getRank = func(v int) int { return mp[v] }
	getValue = func(r int) int { return rank[r-offset] }
	count = len(nums)
	return
}

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}

func argSort(nums []int) []int {
	order := make([]int, len(nums))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func reArrage(nums []int, order []int) []int {
	res := make([]int, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

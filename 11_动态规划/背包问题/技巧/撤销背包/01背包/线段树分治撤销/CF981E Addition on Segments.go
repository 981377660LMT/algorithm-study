package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

// https://www.luogu.com.cn/problem/solution/CF981E
// 给一个长度为n的序列(初始全为0)和q条操作(以(left,right,x)表示将[left,right]中的每个数都加上x(1<=x<=n)
// 对于1≤k≤n,求哪些k满足:选出若干条操作后序列最大值为k.(对于一个k,每条操作至多用一次)
// n,q<=1e4
//
// 序列最大值为k 等价于 `能取到k`, 变为01背包判定问题，可以bitset优化.
// !线段树分治, 求出每个时间点的答案(可以取到哪些值)将它们并起来.
// !时间复杂度O(nqlogn/w)
func AdditionOnSegments(n int, operations [][3]int) []int {
	initState := NewBitset(n)
	initState.Set(0)
	resBit := NewBitset(n)
	seg := NewSegmentTreeDivideAndConquerCopy(
		&initState,
		func(state *State, mutationId int) {
			dp := *state
			add := operations[mutationId][2]
			dp.IOr(dp.Copy().Lsh(add))
		},
		func(state *State) *State {
			dp := *state
			res := dp.Copy()
			return &res
		},
		func(state *State, queryId int) {
			dp := *state
			resBit.IOr(dp)
		},
	)

	for i, e := range operations {
		seg.AddMutation(e[0], e[1]+1, i)
	}
	for i := 1; i <= n; i++ {
		seg.AddQuery(i, i)
	}
	seg.Run()

	res := []int{}
	resBit.ForEach(func(p int) bool {
		if 1 <= p && p <= n {
			res = append(res, p)
		}
		return false
	})

	return res
}

type State = Bitset

// 线段树分治copy版.
// 如果修改操作难以撤销，可以在每个节点处保存一份副本.
type SegmentTreeDivideAndConquerCopy struct {
	initState *State
	mutate    func(state *State, mutationId int)
	copy      func(state *State) *State
	query     func(state *State, queryId int)
	mutations []struct{ start, end, id int }
	queries   []struct{ time, id int }
	nodes     [][]int
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	  initState: 数据结构的初始状态.
//		mutate: 添加编号为`mutationId`的变更后产生的副作用.
//		copy: 拷贝一份数据结构的副本.
//		query: 响应编号为`queryId`的查询.
//		一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` 和 **O(q)** 次`query`.
func NewSegmentTreeDivideAndConquerCopy(
	initState *State,
	mutate func(state *State, mutationId int),
	copy func(state *State) *State,
	query func(state *State, queryId int),
) *SegmentTreeDivideAndConquerCopy {
	return &SegmentTreeDivideAndConquerCopy{
		initState: initState,
		mutate:    mutate, copy: copy, query: query,
	}
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
func (o *SegmentTreeDivideAndConquerCopy) AddMutation(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.mutations = append(o.mutations, struct{ start, end, id int }{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`id`的查询.
func (o *SegmentTreeDivideAndConquerCopy) AddQuery(time int, id int) {
	o.queries = append(o.queries, struct{ time, id int }{time, id})
}

func (o *SegmentTreeDivideAndConquerCopy) Run() {
	if len(o.queries) == 0 {
		return
	}
	times := make([]int, len(o.queries))
	for i := range o.queries {
		times[i] = o.queries[i].time
	}
	sort.Ints(times)
	uniqueInplace(&times)
	usedTimes := make([]bool, len(times)+1)
	usedTimes[0] = true
	for _, e := range o.mutations {
		usedTimes[lowerBound(times, e.start)] = true
		usedTimes[lowerBound(times, e.end)] = true
	}
	for i := 1; i < len(times); i++ {
		if !usedTimes[i] {
			times[i] = times[i-1]
		}
	}
	uniqueInplace(&times)

	n := len(times)
	offset := 1
	for offset < n {
		offset <<= 1
	}
	o.nodes = make([][]int, offset+n)
	for _, e := range o.mutations {
		left := offset + lowerBound(times, e.start)
		right := offset + lowerBound(times, e.end)
		eid := e.id << 1
		for left < right {
			if left&1 == 1 {
				o.nodes[left] = append(o.nodes[left], eid)
				left++
			}
			if right&1 == 1 {
				right--
				o.nodes[right] = append(o.nodes[right], eid)
			}
			left >>= 1
			right >>= 1
		}
	}

	for _, q := range o.queries {
		pos := offset + upperBound(times, q.time) - 1
		o.nodes[pos] = append(o.nodes[pos], (q.id<<1)|1)
	}

	o.dfs(1, o.initState)
}

func (o *SegmentTreeDivideAndConquerCopy) dfs(now int, state *State) {
	curNodes := o.nodes[now]
	for _, id := range curNodes {
		if id&1 == 1 {
			o.query(state, id>>1)
		} else {
			o.mutate(state, id>>1)
		}
	}
	if now<<1 < len(o.nodes) {
		o.dfs(now<<1, o.copy(state))
	}
	if (now<<1)|1 < len(o.nodes) {
		o.dfs((now<<1)|1, o.copy(state))
	}
}

func uniqueInplace(sorted *[]int) {
	tmp := *sorted
	slow := 0
	for fast := 0; fast < len(tmp); fast++ {
		if tmp[fast] != tmp[slow] {
			slow++
			tmp[slow] = tmp[fast]
		}
	}
	*sorted = tmp[:slow+1]
}

func lowerBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func upperBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

//

type Bitset []uint

func NewBitset(n int) Bitset { return make(Bitset, n>>6+1) } // (n+64-1)>>6

func (b Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

func (b Bitset) Copy() Bitset {
	res := make(Bitset, len(b))
	copy(res, b)
	return res
}

func (bs Bitset) Clear() {
	for i := range bs {
		bs[i] = 0
	}
}

// 遍历所有 1 的位置
// 如果对范围有要求，可在 f 中 return p < n
func (b Bitset) ForEach(f func(p int) (shouldBreak bool)) {
	for i, v := range b {
		for ; v != 0; v &= v - 1 {
			j := i<<6 | bits.TrailingZeros(v)
			if f(j) {
				return
			}
		}
	}
}

// 返回第一个 0 的下标，若不存在则返回-1。
func (b Bitset) Index0() int {
	for i, v := range b {
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
	}
	return -1
}

// 返回第一个 1 的下标，若不存在则返回-1.
func (b Bitset) Index1() int {
	for i, v := range b {
		if v != 0 {
			return i<<6 | bits.TrailingZeros(v)
		}
	}
	return -1
}

// 返回下标 >= p 的第一个 1 的下标，若不存在则返回-1.
func (b Bitset) Next1(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i] & (^uint(0) << (p & 63)) // mask off bits below bound
		if v != 0 {
			return i<<6 | bits.TrailingZeros(v)
		}
		for i++; i < len(b); i++ {
			if b[i] != 0 {
				return i<<6 | bits.TrailingZeros(b[i])
			}
		}
	}
	return -1
}

// 返回下标 >= p 的第一个 0 的下标，若不存在则返回-1.
func (b Bitset) Next0(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i]
		if p&63 != 0 {
			v |= ^(^uint(0) << (p & 63))
		}
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
		for i++; i < len(b); i++ {
			if ^b[i] != 0 {
				return i<<6 | bits.TrailingZeros(^b[i])
			}
		}
	}
	return -1
}

// 返回最后第一个 1 的下标，若不存在则返回 -1
func (b Bitset) LastIndex1() int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0 {
			return i<<6 | (bits.Len(b[i]) - 1) // 如果再 +1，需要改成 i<<6 + bits.Len(b[i])
		}
	}
	return -1
}

// += 1 << i，模拟进位
func (b Bitset) Add(i int) { b.FlipRange(i, b.Next0(i)) }

// -= 1 << i，模拟借位
func (b Bitset) Sub(i int) { b.FlipRange(i, b.Next1(i)) }

// 判断 [l,r) 范围内的数是否全为 0
// https://codeforces.com/contest/1107/problem/D（标准做法是二维前缀和）
func (b Bitset) All0(l, r int) bool {
	i := l >> 6
	if i == r>>6 {
		mask := ^uint(0)<<(l&63) ^ ^uint(0)<<(r&63)
		return b[i]&mask == 0
	}
	if b[i]>>(l&63) != 0 {
		return false
	}
	for i++; i < r>>6; i++ {
		if b[i] != 0 {
			return false
		}
	}
	mask := ^uint(0) << (r & 63)
	return b[r>>6]&^mask == 0
}

// 判断 [l,r) 范围内的数是否全为 1
func (b Bitset) All1(l, r int) bool {
	i := l >> 6
	if i == r>>6 {
		mask := ^uint(0)<<(l&63) ^ ^uint(0)<<(r&63)
		return b[i]&mask == mask
	}
	mask := ^uint(0) << (l & 63)
	if b[i]&mask != mask {
		return false
	}
	for i++; i < r>>6; i++ {
		if ^b[i] != 0 {
			return false
		}
	}
	mask = ^uint(0) << (r & 63)
	return ^(b[r>>6] | mask) == 0
}

// 反转 [l,r) 范围内的比特
// https://codeforces.com/contest/1705/problem/E
func (b Bitset) FlipRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] ^= maskL ^ maskR
		return
	}
	b[i] ^= maskL
	for i++; i < r>>6; i++ {
		b[i] = ^b[i]
	}
	b[i] ^= ^maskR
}

// 将 [l,r) 范围内的比特全部置 1
func (b Bitset) SetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] |= maskL ^ maskR
		return
	}
	b[i] |= maskL
	for i++; i < r>>6; i++ {
		b[i] = ^uint(0)
	}
	b[i] |= ^maskR
}

// 将 [l,r) 范围内的比特全部置 0
func (b Bitset) ResetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] &= ^maskL | maskR
		return
	}
	b[i] &= ^maskL
	for i++; i < r>>6; i++ {
		b[i] = 0
	}
	b[i] &= maskR
}

// 左移 k 位
// LC1981 https://leetcode-cn.com/problems/minimize-the-difference-between-target-and-chosen-elements/
func (b Bitset) Lsh(k int) Bitset {
	if k == 0 {
		return b
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return b
	}
	if offset == 0 {
		// Fast path
		copy(b[shift:], b)
	} else {
		for i := len(b) - 1; i > shift; i-- {
			b[i] = b[i-shift]<<offset | b[i-shift-1]>>(64-offset)
		}
		b[shift] = b[0] << offset
	}
	for i := 0; i < shift; i++ {
		b[i] = 0
	}
	return b
}

// 右移 k 位
func (b Bitset) Rsh(k int) Bitset {
	if k == 0 {
		return b
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return b
	}
	lim := len(b) - 1 - shift
	if offset == 0 {
		// Fast path
		copy(b, b[shift:])
	} else {
		for i := 0; i < lim; i++ {
			b[i] = b[i+shift]>>offset | b[i+shift+1]<<(64-offset)
		}
		// 注意：若前后调用 lsh 和 rsh，需要注意超出 n 的范围的 1 对结果的影响（如果需要，可以把范围开大点）
		b[lim] = b[len(b)-1] >> offset
	}
	for i := lim + 1; i < len(b); i++ {
		b[i] = 0
	}
	return b
}

// 借用 bits 库中的一些方法的名字
func (b Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}
func (b Bitset) TrailingZeros() int { return b.Index1() }
func (b Bitset) Len() int           { return b.LastIndex1() + 1 }

// 下面几个方法均需保证长度相同
func (b Bitset) Equals(c Bitset) bool {
	for i, v := range b {
		if v != c[i] {
			return false
		}
	}
	return true
}

func (b Bitset) HasSubset(c Bitset) bool {
	for i, v := range b {
		if v|c[i] != v {
			return false
		}
	}
	return true
}

// 将 c 的元素合并进 b
func (b Bitset) IOr(c Bitset) {
	for i, v := range c {
		b[i] |= v
	}
}

func (b Bitset) Or(c Bitset) Bitset {
	res := make(Bitset, len(b))
	for i, v := range b {
		res[i] = v | c[i]
	}
	return res
}

func (b Bitset) IAnd(c Bitset) {
	for i, v := range c {
		b[i] &= v
	}
}

func (b Bitset) And(c Bitset) Bitset {
	res := make(Bitset, len(b))
	for i, v := range b {
		res[i] = v & c[i]
	}
	return res
}

func (b Bitset) IXor(c Bitset) Bitset {
	for i, v := range c {
		b[i] ^= v
	}
	return b
}

func (b Bitset) Xor(c Bitset) Bitset {
	res := make(Bitset, len(b))
	for i, v := range b {
		res[i] = v ^ c[i]
	}
	return res
}

func (b Bitset) String() string {
	sb := strings.Builder{}
	sb.WriteString("BitSet{")
	nums := []string{}
	b.ForEach(func(pos int) bool {
		nums = append(nums, fmt.Sprintf("%d", pos))
		return false
	})
	sb.WriteString(strings.Join(nums, ","))
	sb.WriteString("}")
	return sb.String()
}

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	operations := make([][3]int, q)
	for i := range operations {
		var left, right, x int
		fmt.Fscan(in, &left, &right, &x)
		operations[i] = [3]int{left, right, x}
	}

	res := AdditionOnSegments(n, operations)
	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

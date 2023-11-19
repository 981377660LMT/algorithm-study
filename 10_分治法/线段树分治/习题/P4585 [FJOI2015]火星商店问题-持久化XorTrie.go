// P4585 [FJOI2015]火星商店问题
// https://www.luogu.com.cn/problem/P4585
// 给定 n 个集合，每个集合元素有两个值，一个是价值，一个是存在时间，每个集合初始有一个存在时间无限的物品。
// 每天都有一个 1 操作和若干个 2 操作
// 操作 0 index v : 在编号为index的集合中加入一个物品v。
// 操作 1 left right x day : 在 left 到 right 集合内查询未过期的物品(day天之内)，使 value xor x 最大，输出最大值。
//
// 用可持久化Trie辅助维护，避免了树套树的毒瘤做法。注意这种树形数据结构都是支持撤销的。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const INF int = 1e18

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n,q int
	fmt.Fscan(in, &n, &q)
	values:=make([]int,n)  // 每个商店的不下架物品的价值
	for i:=range values{
		fmt.Fscan(in,&values[i])
	}

	operations := make([][5]int, q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var shop, price int
			fmt.Fscan(in, &shop, &price)
			shop--
			operations[i] = [5]int{0, shop, price, 0, 0}
		} else {
			var start, end, x, day int
			fmt.Fscan(in, &start, &end, &x, &day)
			start--
			operations[i] = [5]int{1, start, end, x, day}
		}
	}

	seg := NewSegmentTreeDivideAndConquerCopy()
	queryCount := 0
	for i, operation := range operations {
		kind := operation[0]
		if kind == 0 
			// seg.AddMutation(i, i+
		} else {
			x := operation[1]
			queries = append(queries, x)
			seg.AddQuery(i, len(queries)-1)
		}
	}

	res := make([]int, len(queries))
	xorTrie := NewBinaryTriePersistent(1<<20, true)
	initState := xorTrie.NewRoot()
	seg.Run(
		&initState,
		func(state *State, pos int) {
			operation := operations[pos]
			k, b := operation[1], operation[2]
			*state = xorTrie.AddLine(*state, Line{k: k, b: b})
		},
		func(state *State) *State {
			copy_ := xorTrie.Copy(*state)
			return &copy_
		},
		func(state *State, queryId int) {
			x := queries[queryId]
			res[queryId] = xorTrie.Query(*state, x).value
		},
	)

	for _, v := range res {
		if v == -INF {
			fmt.Fprintln(out, "EMPTY SET")
		} else {
			fmt.Fprintln(out, v)
		}
	}
}

type State = *Node

type segMutation struct{ start, end, id int }
type segQuery struct{ time, id int }

// 线段树分治copy版.
// 如果修改操作难以撤销，可以在每个节点处保存一份副本.
// !调用O(n)次拷贝注意不要超出内存.
type SegmentTreeDivideAndConquerCopy struct {
	initState *State
	mutate    func(state *State, mutationId int)
	copy      func(state *State) *State
	query     func(state *State, queryId int)
	mutations []segMutation
	queries   []segQuery
	nodes     [][]int
}

func NewSegmentTreeDivideAndConquerCopy() *SegmentTreeDivideAndConquerCopy {
	return &SegmentTreeDivideAndConquerCopy{}
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
func (o *SegmentTreeDivideAndConquerCopy) AddMutation(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.mutations = append(o.mutations, segMutation{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`id`的查询.
func (o *SegmentTreeDivideAndConquerCopy) AddQuery(time int, id int) {
	o.queries = append(o.queries, segQuery{time, id})
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	 initState: 数据结构的初始状态.
//		mutate: 添加编号为`mutationId`的变更后产生的副作用.
//		copy: 拷贝一份数据结构的副本.
//		query: 响应编号为`queryId`的查询.
//		一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` 和 **O(q)** 次`query`.
func (o *SegmentTreeDivideAndConquerCopy) Run(
	initState *State,
	mutate func(state *State, mutationId int),
	copy func(state *State) *State,
	query func(state *State, queryId int),
) {
	if len(o.queries) == 0 {
		return
	}
	o.initState = initState
	o.mutate, o.copy, o.query = mutate, copy, query
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

// XorTriePersistent
type BinaryTriePersistent struct {
	maxLog     int
	persistent bool
	xorLazy    int
}

type Node struct {
	width       int
	value       int
	count       int
	left, right *Node
}

func NewBinaryTriePersistent(max int, persistent bool) *BinaryTriePersistent {
	return &BinaryTriePersistent{maxLog: bits.Len(uint(max)), persistent: persistent}
}

func (bt *BinaryTriePersistent) NewRoot() *Node {
	return nil
}

func (bt *BinaryTriePersistent) Add(root *Node, num, count int) *Node {
	if root == nil {
		root = _newNode(0, 0)
	}
	return bt._add(root, bt.maxLog, num, count)
}

func (bt *BinaryTriePersistent) Remove(root *Node, num, count int) *Node {
	return bt.Add(root, num, -count)
}

// 0<=k<exist.
//
//	如果不存在,返回(*,false).
func (bt *BinaryTriePersistent) Kth(root *Node, k int, xor int) (res int, ok bool) {
	if root == nil || k < 0 || k >= root.count {
		return
	}
	return bt._kth(root, 0, k, bt.maxLog, xor) ^ xor, true
}

func (bt *BinaryTriePersistent) Max(root *Node, xor int) (res int, ok bool) {
	return bt.Kth(root, root.count-1, xor)
}

func (bt *BinaryTriePersistent) Min(root *Node, xor int) (res int, ok bool) {
	return bt.Kth(root, 0, xor)
}

func (bt *BinaryTriePersistent) CountLess(root *Node, num int, xor int) int {
	if root == nil {
		return 0
	}
	return bt._prefixCount(root, bt.maxLog, num, xor, 0)
}

func (bt *BinaryTriePersistent) BisectLeft(root *Node, num int, xor int) int {
	return bt.CountLess(root, num, xor)
}

func (bt *BinaryTriePersistent) CountLessOrEqual(root *Node, num int, xor int) int {
	return bt.CountLess(root, num+1, xor)
}

func (bt *BinaryTriePersistent) BisectRight(root *Node, num int, xor int) int {
	return bt.CountLess(root, num+1, xor)
}

func (bt *BinaryTriePersistent) CountRange(root *Node, start, end int, xor int) int {
	return bt.CountLess(root, end, xor) - bt.CountLess(root, start, xor)
}

func (bt *BinaryTriePersistent) Size(root *Node) int {
	if root == nil {
		return 0
	}
	return root.count
}

func (bt *BinaryTriePersistent) Enumerate(root *Node, f func(value int, count int)) {
	if root == nil {
		return
	}
	var dfs func(root *Node, val int, height int)
	dfs = func(root *Node, val int, height int) {
		if height == 0 {
			f(val, root.count)
			return
		}
		if c := root.left; c != nil {
			dfs(c, (val<<c.width)|c.value, height-c.width)
		}
		if c := root.right; c != nil {
			dfs(c, (val<<c.width)|c.value, height-c.width)
		}
	}
	dfs(root, 0, bt.maxLog)
}

func (bt *BinaryTriePersistent) Copy(node *Node) *Node {
	if node == nil || !bt.persistent {
		return node
	}
	return &Node{
		width: node.width,
		value: node.value,
		count: node.count,
		left:  node.left,
		right: node.right,
	}
}

func (bt *BinaryTriePersistent) _add(root *Node, height int, val, count int) *Node {
	root = bt.Copy(root)
	root.count += count
	if height == 0 {
		return root
	}
	goRight := (val>>(height-1))&1 == 1
	var c *Node
	if goRight {
		c = root.right
	} else {
		c = root.left
	}
	if c == nil {
		c = _newNode(height, val)
		c.count = count
		if goRight {
			root.right = c
		} else {
			root.left = c
		}
		return root
	}
	w := c.width
	if (val >> (height - w)) == c.value {
		c = bt._add(c, height-w, val&((1<<(height-w))-1), count)
		if goRight {
			root.right = c
		} else {
			root.left = c
		}
		return root
	}
	same := w - 1 - _topbit((val>>(height-w))^(c.value))
	n := _newNode(same, c.value>>(w-same))
	n.count = c.count + count
	c = bt.Copy(c)
	c.width = w - same
	c.value &= (1 << (w - same)) - 1
	if (val>>(height-same-1))&1 == 1 {
		n.left = c
		n.right = _newNode(height-same, val&((1<<(height-same))-1))
		n.right.count = count
	} else {
		n.right = c
		n.left = _newNode(height-same, val&((1<<(height-same))-1))
		n.left.count = count
	}
	if goRight {
		root.right = n
	} else {
		root.left = n
	}
	return root
}

func (bt *BinaryTriePersistent) _prefixCount(root *Node, height int, limit int, xor int, val int) int {
	now := (val << height) ^ xor
	a, b := limit>>height, now>>height
	if a > b {
		return root.count
	}
	if height == 0 || a < b {
		return 0
	}
	res := 0
	if c := root.left; c != nil {
		w := c.width
		res += bt._prefixCount(c, height-w, limit, xor, (val<<w)|c.value)
	}
	if c := root.right; c != nil {
		w := c.width
		res += bt._prefixCount(c, height-w, limit, xor, (val<<w)|c.value)
	}
	return res
}

func (bt *BinaryTriePersistent) _kth(root *Node, val, k, height, xor int) int {
	if height == 0 {
		return val
	}
	left, right := root.left, root.right
	if (xor>>(height-1))&1 == 1 {
		left, right = right, left
	}
	leftSize := 0
	if left != nil {
		leftSize = left.count
	}
	var c *Node
	if k < leftSize {
		c = left
	} else {
		c = right
		k -= leftSize
	}
	w := c.width
	return bt._kth(c, (val<<w)|c.value, k, height-w, xor)
}

func _newNode(width int, value int) *Node {
	return &Node{width: width, value: value}
}

func _topbit(x int) int {
	if x == 0 {
		return -1
	}
	return 31 - bits.LeadingZeros32(uint32(x))
}

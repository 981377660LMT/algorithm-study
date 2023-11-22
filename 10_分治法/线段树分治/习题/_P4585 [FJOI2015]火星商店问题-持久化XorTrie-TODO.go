// P4585 [FJOI2015]火星商店问题
// https://www.luogu.com.cn/problem/P4585
// https://blog.csdn.net/weixin_43823767/article/details/101099572
// 给定 n 个集合，每个集合元素有两个值，一个是价值，一个是存在时间，每个集合初始有一个存在时间无限的物品。
// 每天都有一个 1 操作和若干个 2 操作
// 操作 0 index v : 在编号为index的集合中加入一个物品v。
// 操作 1 left right x day : 在 left 到 right 集合内查询未过期的物品(day天之内)，使 value xor x 最大，输出最大值。
//
// 用可持久化Trie辅助维护，避免了树套树的毒瘤做法。注意这种树形数据结构都是支持撤销的。
//
// !不同于以往的线段树分治，这道题是单点修改，(时间)区间查询
// divideInterval??
// 将查询丢到多个线段树结点上，然后每个结点上开一个Trie.
// TODO 有问题

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

	var n, q int
	fmt.Fscan(in, &n, &q)
	values := make([]int, n) // 每个商店的不下架物品的价值
	for i := range values {
		fmt.Fscan(in, &values[i])
	}

	operations := make([][4]int, q)
	mutations := [][2]int{} // (id, time)
	queries := [][3]int{}   // (id, startTime, endTime)
	curTime := 0
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var shop, price int
			fmt.Fscan(in, &shop, &price)
			shop--
			curTime += 1
			operations[i] = [4]int{shop, price, 0, 0}
			mutations = append(mutations, [2]int{i, curTime})
		} else {
			var start, end, x, day int
			fmt.Fscan(in, &start, &end, &x, &day)
			start--
			operations[i] = [4]int{start, end, x, day}
			queries = append(queries, [3]int{i, curTime - day, curTime + 1})
		}
	}

	res := make([]int, q)
	for i := range res {
		res[i] = -1
	}

	bisectRight := func(nums []*Node, target int) int {

		left, right := 0, len(nums)-1
		for left <= right {
			mid := (left + right) >> 1
			if nums[mid].lastIndex >= target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return left
	}

	seg := NewSegmentTreeDivideAndConquerCopyRangeQuery()
	// 按照商店编号降序，保证查询结点时可以二分
	sort.Slice(mutations, func(i, j int) bool {
		return operations[mutations[i][0]][0] > operations[mutations[j][0]][0]
	})
	for _, mutation := range mutations {
		time, id := mutation[1], mutation[0]
		seg.AddMutation(time, id)
	}
	for _, query := range queries {
		id, start, end := query[0], query[1], query[2]
		seg.AddQuery(start, end, id)
		fmt.Println("query", id, start, end)
	}

	xorTrie := NewXorTrieSimplePersistent(int(1e5 + 10))
	root := xorTrie.NewRoot()
	for _, v := range values {
		root = xorTrie.Insert(root, v, -1) // INF表示永不下架
	}
	initState := []*Node{root}
	seg.Run(
		&initState,
		func(state *State, mutationId int) {
			nums := *state
			mutation := operations[mutationId]
			shop, price := mutation[1], mutation[2]
			nums = append(nums, xorTrie.Insert(nums[len(nums)-1], price, shop))
			*state = nums
		},
		func(state *State) *State {
			nums := *state
			copy_ := xorTrie.Copy(nums[0])
			return &[]*Node{copy_}
		},
		func(state *State, queryId int) {
			nums := *state
			start, end, value := operations[queryId][0], operations[queryId][1], operations[queryId][2]
			// 二分 lastIndex 找到 end 对应
			pos := bisectRight(nums, end-1) - 1
			if pos < 0 {
				return
			}
			root := nums[pos]
			fmt.Println(queryId, start, end, value, root)
			maxXor, _ := xorTrie.Query(root, value, start)
			res[queryId] = max(res[queryId], maxXor)
		},
	)

	for _, v := range res {
		if v != -1 {
			fmt.Fprintln(out, v)
		}
	}
}

type State = []*Node

type segMutation struct{ time, id int }
type segQuery struct{ start, end, id int }

// 线段树分治copy版.
// 如果修改操作难以撤销，可以在每个节点处保存一份副本.
// !调用O(n)次拷贝注意不要超出内存.
type SegmentTreeDivideAndConquerCopyRangeQuery struct {
	initState *State
	mutate    func(state *State, mutationId int)
	copy      func(state *State) *State
	query     func(state *State, queryId int)
	queries   []segQuery
	mutations []segMutation
	nodes     [][]int
}

func NewSegmentTreeDivideAndConquerCopyRangeQuery() *SegmentTreeDivideAndConquerCopyRangeQuery {
	return &SegmentTreeDivideAndConquerCopyRangeQuery{}
}

// 在时间`time`时添加一个编号为`id`的变更.
func (o *SegmentTreeDivideAndConquerCopyRangeQuery) AddMutation(time int, id int) {
	o.mutations = append(o.mutations, segMutation{time, id})
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的查询.
func (o *SegmentTreeDivideAndConquerCopyRangeQuery) AddQuery(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.queries = append(o.queries, segQuery{startTime, endTime, id})
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	 initState: 数据结构的初始状态.
//		mutate: 添加编号为`mutationId`的变更后产生的副作用.
//		copy: 拷贝一份数据结构的副本.
//		query: 响应编号为`queryId`的查询.
//		一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` 和 **O(q)** 次`query`.
func (o *SegmentTreeDivideAndConquerCopyRangeQuery) Run(
	initState *State,
	mutate func(state *State, mutationId int),
	copy func(state *State) *State,
	query func(state *State, queryId int),
) {
	if len(o.mutations) == 0 {
		return
	}
	o.initState = initState
	o.mutate, o.copy, o.query = mutate, copy, query
	set := make(map[int]struct{})
	for i := range o.mutations {
		set[o.mutations[i].time] = struct{}{}
	}
	for i := range o.queries {
		set[o.queries[i].start] = struct{}{}
		set[o.queries[i].end] = struct{}{}
	}
	times := make([]int, 0, len(set))
	for k := range set {
		times = append(times, k)
	}
	sort.Ints(times)

	n := len(times)
	offset := 1
	for offset < n {
		offset <<= 1
	}
	o.nodes = make([][]int, offset+n)

	for _, q := range o.mutations {
		id := (q.id << 1) | 1
		pos := offset + upperBound(times, q.time) - 1
		for pos > 0 {
			o.nodes[pos] = append(o.nodes[pos], id)
			pos >>= 1
		}
	}

	for _, e := range o.queries {
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

	o.dfs(1, o.initState)
}

func (o *SegmentTreeDivideAndConquerCopyRangeQuery) dfs(now int, state *State) {
	if now<<1 < len(o.nodes) {
		o.dfs(now<<1, o.copy(state))
	}
	if (now<<1)|1 < len(o.nodes) {
		o.dfs((now<<1)|1, o.copy(state))
	}
	curNodes := o.nodes[now]
	for _, id := range curNodes {
		if id&1 == 1 {
			o.mutate(state, id>>1)
		} else {
			o.query(state, id>>1)
		}
	}
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

// 区间异或最大值.
type XorRangeMax struct {
	trie  *XorTrieSimplePersistent
	roots []*Node
	size  int
}

func NewXorRangeMax(max int) *XorRangeMax {
	trie := NewXorTrieSimplePersistent(max)
	root := trie.NewRoot()
	return &XorRangeMax{
		trie:  trie,
		roots: []*Node{root},
	}
}

// 查询x与区间[start,end)中的数异或的最大值以及最大值对应的下标.
// 如果不存在,返回0,-1.
func (xm *XorRangeMax) Query(start, end int, x int) (maxXor, maxIndex int) {
	if start < 0 {
		start = 0
	}
	if size := xm.Size(); end > size {
		end = size
	}
	if start >= end {
		return 0, -1
	}
	maxXor, node := xm.trie.Query(xm.roots[end], x, start)
	if node == nil {
		return 0, -1
	}
	return maxXor, node.lastIndex
}

func (xm *XorRangeMax) Push(num int) {
	xm.roots = append(xm.roots, xm.trie.Insert(xm.roots[xm.size], num, xm.size))
	xm.size += 1
}

func (xm *XorRangeMax) Size() int {
	return xm.size
}

type Node struct {
	lastIndex int // 最后一次被更新的时间
	chidlren  [2]*Node
}

type XorTrieSimplePersistent struct {
	bit int
}

func NewXorTrieSimplePersistent(upper int) *XorTrieSimplePersistent {
	return &XorTrieSimplePersistent{bit: bits.Len(uint(upper))}
}

func (trie *XorTrieSimplePersistent) NewRoot() *Node {
	return nil
}

func (trie *XorTrieSimplePersistent) Copy(node *Node) *Node {
	if node == nil {
		return node
	}
	return &Node{
		lastIndex: node.lastIndex,
		chidlren:  node.chidlren,
	}
}

func (trie *XorTrieSimplePersistent) Insert(root *Node, num int, lastIndex int) *Node {
	if root == nil {
		root = &Node{}
	}
	return trie._insert(root, num, trie.bit-1, lastIndex)
}

// 查询num与root中的数异或的最大值以及最大值对应的结点.
// !如果root为nil,返回0.
func (trie *XorTrieSimplePersistent) Query(root *Node, num int, leftIndex int) (maxXor int, node *Node) {
	if root == nil {
		return
	}
	for k := trie.bit - 1; k >= 0; k-- {
		bit := (num >> k) & 1
		if root.chidlren[bit^1] != nil && root.chidlren[bit^1].lastIndex >= leftIndex {
			bit ^= 1
			maxXor |= 1 << k
		}
		root = root.chidlren[bit]
	}
	return maxXor, root
}

func (trie *XorTrieSimplePersistent) _insert(root *Node, num int, depth int, lastIndex int) *Node {
	root = trie.Copy(root)
	root.lastIndex = lastIndex
	if depth < 0 {
		return root
	}
	bit := (num >> depth) & 1
	if root.chidlren[bit] == nil {
		root.chidlren[bit] = &Node{}
	}
	root.chidlren[bit] = trie._insert(root.chidlren[bit], num, depth-1, lastIndex)
	return root
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(nums ...int) int {
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

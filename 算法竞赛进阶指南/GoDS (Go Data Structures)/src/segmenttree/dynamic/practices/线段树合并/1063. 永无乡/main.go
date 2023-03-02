// https://www.acwing.com/problem/content/description/1065/

// 永无乡包含 n座岛，编号从 1到 n，每座岛都有自己的独一无二的重要度，
// 按照重要度可以将这 n座岛排名，名次用 1到 n来表示.
// 现在有两种操作：
// B x y 在x和y之间建立一座桥，使得x和y之间可以互相到达
// Q x y 询问当前与岛 x连通的所有岛中第 k重要的是哪座岛,如果该岛屿不存在，则输出 −1.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, brigdeNum int
	fmt.Fscan(in, &n, &brigdeNum)
	rank := make([]int, n) // 每个岛的重要度
	rankToId := make([]int, n)
	for i := range rank {
		var r int
		fmt.Fscan(in, &r)
		r--
		rank[i] = r
		rankToId[r] = i
	}

	uf := NewUnionFindArray(n)
	roots := make([]*LazyNode, n) // 每个岛的线段树
	for i := range roots {
		roots[i] = CreateSegmentTree(0, n-1)
		roots[i].Set(rank[i], 1)
	}

	// 在x和y之间建立一座桥，使得x和y之间可以互相到达
	addEdge := func(u, v int) {
		if uf.IsConnected(u, v) {
			return
		}
		uf.UnionWithCallback(u, v, func(big, small int) {
			roots[big] = roots[big].Merge(roots[small])
		})
	}

	for i := 0; i < brigdeNum; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		addEdge(u, v)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "B" {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u, v = u-1, v-1
			addEdge(u, v)
		} else {
			var u, k int
			fmt.Fscan(in, &u, &k)
			u--
			leader := uf.Find(u)
			if roots[leader].QueryAll() < k {
				fmt.Fprintln(out, -1)
			} else {
				num := roots[leader].kth(k)
				fmt.Fprintln(out, rankToId[num]+1) // !输出岛的编号
			}
		}
	}
}

// 指定区间上下界建立权值线段树.
func CreateSegmentTree(lower, upper int) *LazyNode {
	root := &LazyNode{left: lower, right: upper}
	return root
}

type LazyNode struct {
	left, right           int
	sum                   int
	lazy                  int
	leftChild, rightChild *LazyNode
}

func (LazyNode) op(a, b int) int {
	return a + b
}

func (o *LazyNode) propagate(add int) {
	o.lazy += add                         // % mod
	o.sum += (o.right - o.left + 1) * add // % mod
}

func (o *LazyNode) pushDown() {
	m := (o.left + o.right) >> 1
	if o.leftChild == nil {
		o.leftChild = &LazyNode{left: o.left, right: m}
	}
	if o.rightChild == nil {
		o.rightChild = &LazyNode{left: m + 1, right: o.right}
	}
	if add := o.lazy; add != 0 {
		o.leftChild.propagate(add)
		o.rightChild.propagate(add)
		o.lazy = 0
	}
}

func (o *LazyNode) pushUp() {
	o.sum = o.op(o.leftChild.QueryAll(), o.rightChild.QueryAll())
}

// Build from array. [1,len(nums))]
func (o *LazyNode) Build(nums []int) {
	o.build(nums, 1, len(nums))
}

func (o *LazyNode) build(nums []int, left, right int) {
	o.left, o.right = left, right
	if left == right {
		o.sum = int(nums[left-1])
		return
	}
	m := (left + right) >> 1
	o.leftChild = &LazyNode{}
	o.leftChild.build(nums, left, m)
	o.rightChild = &LazyNode{}
	o.rightChild.build(nums, m+1, right)
	o.pushUp()
}

// [left, right]
func (o *LazyNode) Update(left, right int, add int) {
	if left <= o.left && o.right <= right {
		o.propagate(add)
		return
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if left <= m {
		o.leftChild.Update(left, right, add)
	}
	if m < right {
		o.rightChild.Update(left, right, add)
	}
	o.pushUp()
}

// [left, right]
func (o *LazyNode) Query(left, right int) int {
	if o == nil || left > o.right || right < o.left {
		return 0
	}
	if left <= o.left && o.right <= right {
		return o.sum
	}
	o.pushDown()
	return o.op(o.leftChild.Query(left, right), o.rightChild.Query(left, right))
}

func (o *LazyNode) Set(pos int, val int) {
	if o.left == o.right {
		o.sum = val
		return
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if pos <= m {
		o.leftChild.Set(pos, val)
	} else {
		o.rightChild.Set(pos, val)
	}
	o.pushUp()
}

func (o *LazyNode) Get(pos int) int {
	if o.left == o.right {
		return o.sum
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if pos <= m {
		return o.leftChild.Get(pos)
	}
	return o.rightChild.Get(pos)
}

func (o *LazyNode) QueryAll() int {
	if o != nil {
		return o.sum
	}
	return 0
}

// 线段树合并
func (o *LazyNode) Merge(b *LazyNode) *LazyNode {
	if o == nil {
		return b
	}
	if b == nil {
		return o
	}
	if o.left == o.right {
		o.sum += b.sum
		return o
	}
	o.leftChild = o.leftChild.Merge(b.leftChild)
	o.rightChild = o.rightChild.Merge(b.rightChild)
	o.pushUp()
	return o
}

// 线段树分裂
//  将区间 [l,r] 从原树分离到 other 上, this 为原树的剩余部分
func (o *LazyNode) Split(left, right int) (this, other *LazyNode) {
	this, other = o.split(nil, left, right)
	return
}

func (o *LazyNode) split(b *LazyNode, l, r int) (*LazyNode, *LazyNode) {
	if o == nil || l > o.right || r < o.left {
		return o, nil
	}
	if l <= o.left && o.right <= r {
		return nil, o
	}
	if b == nil {
		b = &LazyNode{left: o.left, right: o.right}
	}
	o.leftChild, b.leftChild = o.leftChild.split(b.leftChild, l, r)
	o.rightChild, b.rightChild = o.rightChild.split(b.rightChild, l, r)
	o.pushUp()
	b.pushUp()
	return o, b
}

// 权值线段树求第 k 小
// 调用前需保证 1 <= k <= root.QueryAll()
func (o *LazyNode) kth(k int) int {
	if o.left == o.right {
		return o.left
	}
	if lc := o.leftChild.QueryAll(); k <= lc {
		return o.leftChild.kth(k)
	} else {
		return o.rightChild.kth(k - lc)
	}
}

// NewUnionFindWithCallback ...
func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

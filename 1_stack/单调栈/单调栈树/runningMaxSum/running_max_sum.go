// 数字 = [3, 1, 4, 1, 5]
// runningMax(nums) = [3, 3, 4, 4, 5], runningMaxSum(nums) = 19
// runningMin(nums) = [3, 1, 1, 1, 1], runningMinSum(nums) = 7

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	P5648()
}

func demo() {
	arr := []int{3, 1, 4, 1, 5}
	runningMaxSum := NewRunningMaxSum(arr, true)
	runningMinSum := NewRunningMaxSum(arr, false)
	fmt.Println(runningMaxSum.Query(1, 3)) // 5
	fmt.Println(runningMinSum.Query(0, 5)) // 7
}

// P5648 Mivik的神力
// https://www.luogu.com.cn/problem/P5648
func P5648() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	runningMaxSum := NewRunningMaxSum(nums, true)
	lastRes := 0
	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		start := 1 + (u^lastRes)%n
		len := 1 + (v^(lastRes+1))%(n-start+1)
		start--
		lastRes = runningMaxSum.Query(int32(start), int32(start+len))
		fmt.Fprintln(out, lastRes)
	}
}

// 3420. 统计 K 次操作以内得到非递减子数组的数目
// https://leetcode.cn/problems/count-non-decreasing-subarrays-after-k-operations/description/
func countNonDecreasingSubarrays(nums []int, k int) int64 {
	n := int32(len(nums))
	rs := NewRunningMaxSum(nums, true)
	presum := make([]int, n+1)
	for i := int32(0); i < n; i++ {
		presum[i+1] = presum[i] + nums[i]
	}
	cost := func(l, r int32) int {
		return rs.Query(l, r) - (presum[r] - presum[l])
	}

	res := 0
	for left := int32(0); left < n; left++ {
		right := MaxRight32(left, func(r int32) bool { return cost(left, r) <= k }, n)
		res += int(right - left)
	}
	return int64(res)
}

type RunningMaxSum struct {
	n    int32
	nums []int
	tree *CompressedBinaryLiftWithSum[int]
}

// isMax: true => running max, false => running min.
func NewRunningMaxSum(nums []int, isMax bool) *RunningMaxSum {
	copy_ := append(nums[:0:0], nums...)
	n := int32(len(copy_))
	var checkParent func(i, j int32) bool
	if isMax {
		checkParent = func(i, j int32) bool { return copy_[i] < copy_[j] }
	} else {
		checkParent = func(i, j int32) bool { return copy_[i] > copy_[j] }
	}
	weight := func(start, end int32) int { return copy_[start] * int(end-start) }
	e := func() int { return 0 }
	op := func(e1, e2 int) int { return e1 + e2 }
	tree := BuildMonoStackTree(n, checkParent, weight, e, op)
	return &RunningMaxSum{n: n, nums: copy_, tree: tree}
}

// O(logn) 查询区间[start,end)的 running max/min sum.
func (s *RunningMaxSum) Query(start, end int32) int {
	if start < 0 {
		start = 0
	}
	if end > s.n {
		end = s.n
	}
	if start >= end {
		return 0
	}

	segStart, runningSum := s.tree.LastTrueWithSum(start, func(e int32, _ int) bool { return e < end }, true)
	segSum := s.nums[segStart] * int(end-segStart)
	res := runningSum + segSum
	return res
}

// 单调栈树.
// 将每个数的下标与右侧第一个(严格)大于它的数的下标连接，构成的树.
//
//		n: 数组长度.
//		checkParent: j是否可以是i的父节点(0 <= i < j < n).
//		weight: 权重函数(边权或点权, 0 <= start < end <= n).
//		e: 幺元.
//		op: 结合运算.
//
//	 e.g.:
//	 nums = [2, 1, 4, 3, 5]
//
//	 	      5(虚拟根节点)
//	 	     /
//	 	    4
//	 	   /|
//	 	  2 3
//	 	 /|
//	 	0 1
func BuildMonoStackTree[S any](
	n int32, checkParent func(i, j int32) bool, /** 0 <= i < j < n */
	weight func(start, end int32) S, /** 0 <= start < end <= n */
	e func() S, op func(e1, e2 S) S,
) *CompressedBinaryLiftWithSum[S] {
	stack := []int32{}
	depths := make([]int32, n+1)
	parents := make([]int32, n+1)
	parents[n] = -1
	values := make([]S, n+1)
	values[n] = e()
	for i := int32(n - 1); i >= 0; i-- {
		for len(stack) > 0 && !checkParent(i, stack[len(stack)-1]) {
			stack = stack[:len(stack)-1]
		}

		var p int32
		if len(stack) == 0 {
			p = n
		} else {
			p = stack[len(stack)-1]
		}
		depths[i] = depths[p] + 1
		parents[i] = p
		values[i] = weight(i, p)
		stack = append(stack, i)
	}

	tree := NewCompressedBinaryLiftWithSum(n+1, depths, parents, func(i int32) S { return values[i] }, e, op)
	return tree
}

type CompressedBinaryLiftWithSum[S any] struct {
	Depth       []int32
	Parent      []int32
	jump        []int32 // 指向当前节点的某个祖先节点.
	attachments []S     // 从当前结点到`jump`结点的路径上的聚合值(不包含`jump`结点).
	singles     []S     // 当前结点的聚合值.
	e           func() S
	op          func(e1, e2 S) S
}

// values: 每个点的`点权`.
// 如果需要查询边权，则每个点的`点权`设为`该点与其父亲结点的边权`, 根节点的`点权`设为`幺元`.
func NewCompressedBinaryLiftWithSum[S any](
	n int32, depthOnTree, parentOnTree []int32, values func(i int32) S,
	e func() S, op func(e1, e2 S) S,
) *CompressedBinaryLiftWithSum[S] {
	res := &CompressedBinaryLiftWithSum[S]{
		Depth:       depthOnTree,
		Parent:      parentOnTree,
		jump:        make([]int32, n),
		attachments: make([]S, n),
		singles:     make([]S, n),
		e:           e,
		op:          op,
	}
	for i := int32(0); i < n; i++ {
		res.jump[i] = -1
		res.attachments[i] = res.e()
		res.singles[i] = values(i)
	}
	for i := int32(0); i < n; i++ {
		res._consider(i)
	}
	return res
}

// root:-1表示无根.
func NewCompressedBinaryLiftWithSumFromTree[S any](
	tree [][]int32, root int32, values func(i int32) S,
	e func() S, op func(e1, e2 S) S,
) *CompressedBinaryLiftWithSum[S] {
	n := int32(len(tree))
	res := &CompressedBinaryLiftWithSum[S]{
		Depth:       make([]int32, n),
		Parent:      make([]int32, n),
		jump:        make([]int32, n),
		attachments: make([]S, n),
		singles:     make([]S, n),
		e:           e,
		op:          op,
	}
	for i := int32(0); i < n; i++ {
		res.attachments[i] = res.e()
		res.singles[i] = values(i)
	}
	if root != -1 {
		res.Parent[root] = -1
		res.jump[root] = root
		res._setUp(tree, root)
	} else {
		for i := int32(0); i < n; i++ {
			res.Parent[i] = -1
		}
		for i := int32(0); i < n; i++ {
			if res.Parent[i] == -1 {
				res.jump[i] = i
				res._setUp(tree, i)
			}
		}
	}
	return res
}

func (bl *CompressedBinaryLiftWithSum[S]) FirstTrue(start int32, predicate func(end int32) bool) int32 {
	for !predicate(start) {
		if predicate(bl.jump[start]) {
			start = bl.Parent[start]
		} else {
			if start == bl.jump[start] {
				return -1
			}
			start = bl.jump[start]
		}
	}
	return start
}

func (bl *CompressedBinaryLiftWithSum[S]) FirstTrueWithSum(start int32, predicate func(end int32, sum S) bool, isEdge bool) (int32, S) {
	if isEdge {
		sum := bl.e() // 不包含_singles[start]
		for {
			if predicate(start, sum) {
				return start, sum
			}
			jumpStart, jumpSum := bl.jump[start], bl.op(sum, bl.attachments[start])
			if predicate(jumpStart, jumpSum) {
				sum = bl.op(sum, bl.singles[start])
				start = bl.Parent[start]
			} else {
				if start == jumpStart {
					return -1, jumpSum
				}
				sum = jumpSum
				start = jumpStart
			}
		}
	} else {
		sum := bl.e() // 不包含_singles[start]
		for {
			sumWithSingle := bl.op(sum, bl.singles[start])
			if predicate(start, sumWithSingle) {
				return start, sumWithSingle
			}
			jumpStart, jumpSum1 := bl.jump[start], bl.op(sum, bl.attachments[start])
			jumpSum2 := bl.op(jumpSum1, bl.singles[jumpStart])
			if predicate(jumpStart, jumpSum2) {
				sum = sumWithSingle
				start = bl.Parent[start]
			} else {
				if start == jumpStart {
					return -1, jumpSum2
				}
				sum = jumpSum1
				start = jumpStart
			}
		}
	}
}

func (bl *CompressedBinaryLiftWithSum[S]) LastTrue(start int32, predicate func(end int32) bool) int32 {
	if !predicate(start) {
		return -1
	}
	for {
		if predicate(bl.jump[start]) {
			if start == bl.jump[start] {
				return start
			}
			start = bl.jump[start]
		} else if predicate(bl.Parent[start]) {
			start = bl.Parent[start]
		} else {
			return start
		}
	}
}

func (bl *CompressedBinaryLiftWithSum[S]) LastTrueWithSum(start int32, predicate func(end int32, sum S) bool, isEdge bool) (int32, S) {
	if isEdge {
		sum := bl.e() // 不包含_singles[start]
		if !predicate(start, sum) {
			return -1, sum
		}
		for {
			jumpStart, jumpSum := bl.jump[start], bl.op(sum, bl.attachments[start])
			if predicate(jumpStart, jumpSum) {
				if start == jumpStart {
					return start, sum
				}
				sum = jumpSum
				start = jumpStart
			} else {
				parentStart, parentSum := bl.Parent[start], bl.op(sum, bl.singles[start])
				if predicate(parentStart, parentSum) {
					sum = parentSum
					start = parentStart
				} else {
					return start, sum
				}
			}
		}
	} else {
		if !predicate(start, bl.singles[start]) {
			return -1, bl.singles[start]
		}
		sum := bl.e() // 不包含_singles[start]
		for {
			jumpStart, jumpSum1 := bl.jump[start], bl.op(sum, bl.attachments[start])
			jumpSum2 := bl.op(jumpSum1, bl.singles[jumpStart])
			if predicate(jumpStart, jumpSum2) {
				if start == jumpStart {
					return start, jumpSum2
				}
				sum = jumpSum1
				start = jumpStart
			} else {
				parentStart, parentSum1 := bl.Parent[start], bl.op(sum, bl.singles[start])
				parentSum2 := bl.op(parentSum1, bl.singles[parentStart])
				if predicate(parentStart, parentSum2) {
					sum = parentSum1
					start = parentStart
				} else {
					return start, parentSum1
				}
			}
		}
	}
}

func (bl *CompressedBinaryLiftWithSum[S]) UpToDepth(root int32, toDepth int32) int32 {
	if !(0 <= toDepth && toDepth <= bl.Depth[root]) {
		return -1
	}
	for bl.Depth[root] > toDepth {
		if bl.Depth[bl.jump[root]] < toDepth {
			root = bl.Parent[root]
		} else {
			root = bl.jump[root]
		}
	}
	return root
}

func (bl *CompressedBinaryLiftWithSum[S]) UpToDepthWithSum(root int32, toDepth int32, isEdge bool) (int32, S) {
	sum := bl.e() // 不包含_singles[root]
	if !(0 <= toDepth && toDepth <= bl.Depth[root]) {
		return -1, sum
	}
	for bl.Depth[root] > toDepth {
		if bl.Depth[bl.jump[root]] < toDepth {
			sum = bl.op(sum, bl.singles[root])
			root = bl.Parent[root]
		} else {
			sum = bl.op(sum, bl.attachments[root])
			root = bl.jump[root]
		}
	}
	if !isEdge {
		sum = bl.op(sum, bl.singles[root])
	}
	return root, sum
}

func (bl *CompressedBinaryLiftWithSum[S]) KthAncestor(node, k int32) int32 {
	targetDepth := bl.Depth[node] - k
	return bl.UpToDepth(node, targetDepth)
}

func (bl *CompressedBinaryLiftWithSum[S]) KthAncestorWithSum(node, k int32, isEdge bool) (int32, S) {
	targetDepth := bl.Depth[node] - k
	return bl.UpToDepthWithSum(node, targetDepth, isEdge)
}

func (bl *CompressedBinaryLiftWithSum[S]) Lca(a, b int32) int32 {
	if bl.Depth[a] > bl.Depth[b] {
		a = bl.KthAncestor(a, bl.Depth[a]-bl.Depth[b])
	} else if bl.Depth[a] < bl.Depth[b] {
		b = bl.KthAncestor(b, bl.Depth[b]-bl.Depth[a])
	}
	for a != b {
		if bl.jump[a] == bl.jump[b] {
			a = bl.Parent[a]
			b = bl.Parent[b]
		} else {
			a = bl.jump[a]
			b = bl.jump[b]
		}
	}
	return a
}

// 查询路径`a`到`b`的聚合值.
// isEdge 是否是边权.
func (bl *CompressedBinaryLiftWithSum[S]) LcaWithSum(a, b int32, isEdge bool) (int32, S) {
	var e S // 不包含_singles[a]和_singles[b]
	if bl.Depth[a] > bl.Depth[b] {
		end, sum := bl.UpToDepthWithSum(a, bl.Depth[b], true)
		a, e = end, sum
	} else if bl.Depth[a] < bl.Depth[b] {
		end, sum := bl.UpToDepthWithSum(b, bl.Depth[a], true)
		b, e = end, sum
	} else {
		e = bl.e()
	}
	for a != b {
		if bl.jump[a] == bl.jump[b] {
			e = bl.op(e, bl.singles[a])
			e = bl.op(e, bl.singles[b])
			a = bl.Parent[a]
			b = bl.Parent[b]
		} else {
			e = bl.op(e, bl.attachments[a])
			e = bl.op(e, bl.attachments[b])
			a = bl.jump[a]
			b = bl.jump[b]
		}
	}
	if !isEdge {
		e = bl.op(e, bl.singles[a])
	}
	return a, e
}

func (bl *CompressedBinaryLiftWithSum[S]) Jump(start, target, step int32) int32 {
	lca := bl.Lca(start, target)
	dep1, dep2, deplca := bl.Depth[start], bl.Depth[target], bl.Depth[lca]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return bl.KthAncestor(start, step)
	}
	return bl.KthAncestor(target, dist-step)
}

func (bl *CompressedBinaryLiftWithSum[S]) InSubtree(maybeChild, maybeAncestor int32) bool {
	return bl.Depth[maybeChild] >= bl.Depth[maybeAncestor] &&
		bl.KthAncestor(maybeChild, bl.Depth[maybeChild]-bl.Depth[maybeAncestor]) == maybeAncestor
}

func (bl *CompressedBinaryLiftWithSum[S]) Dist(a, b int32) int32 {
	return bl.Depth[a] + bl.Depth[b] - 2*bl.Depth[bl.Lca(a, b)]
}

func (bl *CompressedBinaryLiftWithSum[S]) _consider(root int32) {
	if root == -1 || bl.jump[root] != -1 {
		return
	}
	p := bl.Parent[root]
	bl._consider(p)
	bl._addLeaf(root, p)
}

func (bl *CompressedBinaryLiftWithSum[S]) _addLeaf(leaf, parent int32) {
	if parent == -1 {
		bl.jump[leaf] = leaf
	} else if tmp := bl.jump[parent]; bl.Depth[parent]-bl.Depth[tmp] == bl.Depth[tmp]-bl.Depth[bl.jump[tmp]] {
		bl.jump[leaf] = bl.jump[tmp]
		bl.attachments[leaf] = bl.op(bl.singles[leaf], bl.attachments[parent])
		bl.attachments[leaf] = bl.op(bl.attachments[leaf], bl.attachments[tmp])
	} else {
		bl.jump[leaf] = parent
		bl.attachments[leaf] = bl.singles[leaf] // copy
	}
}

func (bl *CompressedBinaryLiftWithSum[S]) _setUp(tree [][]int32, root int32) {
	queue := []int32{root}
	head := 0
	for head < len(queue) {
		cur := queue[head]
		head++
		nexts := tree[cur]
		for _, next := range nexts {
			if next == bl.Parent[cur] {
				continue
			}
			bl.Depth[next] = bl.Depth[cur] + 1
			bl.Parent[next] = cur
			queue = append(queue, next)
			bl._addLeaf(next, cur)
		}
	}
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含,使用时需要right-1.
// right<=upper.
func MaxRight32(left int32, check func(right int32) bool, upper int32) int32 {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func test() {
	runningMaxSum := func(nums []int, start, end int) int {
		if start < 0 {
			start = 0
		}
		if end > len(nums) {
			end = len(nums)
		}
		if start >= end {
			return 0
		}
		res := nums[start]
		curMax := nums[start]
		for i := start + 1; i < end; i++ {
			if nums[i] > curMax {
				curMax = nums[i]
			}
			res += curMax
		}
		return res
	}

	runningMinSum := func(nums []int, start, end int) int {
		if start < 0 {
			start = 0
		}
		if end > len(nums) {
			end = len(nums)
		}
		if start >= end {
			return 0
		}
		res := nums[start]
		curMin := nums[start]
		for i := start + 1; i < end; i++ {
			if nums[i] < curMin {
				curMin = nums[i]
			}
			res += curMin
		}
		return res
	}

	for i := 0; i < 5; i++ {
		n := rand.Intn(1000) + 1
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = rand.Intn(1000)
		}
		runningMaxSum1 := NewRunningMaxSum(nums, true)
		runningMinSum1 := NewRunningMaxSum(nums, false)

		for s := int32(0); s < int32(n); s++ {
			for e := int32(0); e < int32(n); e++ {
				if runningMaxSum(nums, int(s), int(e)) != runningMaxSum1.Query(s, e) {
					fmt.Println("error")
				}
				if runningMinSum(nums, int(s), int(e)) != runningMinSum1.Query(s, e) {
					fmt.Println("error")
				}
			}
		}
	}

	fmt.Println("pass")
}

func testTime() {
	n, q := int(2e5), int(2e5)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Intn(1e9)
	}

	starts, ends := make([]int32, q), make([]int32, q)
	for i := 0; i < n; i++ {
		starts[i] = int32(rand.Intn(n))
		ends[i] = int32(rand.Intn(n))
		if starts[i] > ends[i] {
			starts[i], ends[i] = ends[i], starts[i]
		}
	}

	time1 := time.Now()
	R := NewRunningMaxSum(nums, true)
	for i := 0; i < q; i++ {
		R.Query(starts[i], ends[i])
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1))
}

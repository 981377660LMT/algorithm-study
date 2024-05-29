// 可能的子集和,子集和输出方案
// SubsetSum 求方案有O(n*max(nums)),O(n*target/w),O(2^(n/2)三种解法
// 一般的dp，bitset优化dp, 折半搜索，最优算法取决于数据范围

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	JumpingSequence()
	// testTime()
}

func demo() {

	fmt.Println(SubsetSumTargetDp1([]int{2, 3, 4, 5, 6, 7, 8, 9}, 10))
	fmt.Println(SubsetSumTargetDp4([]int{2, 3, 4, 5, 6, 7, 8, 9}, 10))
	fmt.Println(SubsetSumTargetDp3([]int{3, 4, 5, 6, 7, 8, 9}, 11))
	fmt.Println(SubsetSumTargetDp2([]int{300}, 300))
	fmt.Println(SubsetSumTargetDp5([]int{2, 3, 4, 5, 6, 7, 8, 9}, 10))

}

// G - Jumping sequence
// https://atcoder.jp/contests/abc221/tasks/abc221_g
// 有一个无限大的二维平面，开始时位于原点(0,0).
// 给定一个长为n的序列nums.
// 每一步可以向上、下、左、右四个方向移动nums[i].
// 想要在n步后到达(A,B).问是否存在并输出方案(L/R/U/D).
// n<=2000,di<=1800.
//
//  1. 坐标系旋转45度下考虑问题, (x,y)=>(x-y,x+y).
//  2. 问题变为 d1 ± d2 ± ... ± dn = C.两边同时加∑di,系数变为0/2，除以二后转化为01背包问题.
//     x和y维度分别独立处理.
func JumpingSequence() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, targetX, targetY int
	fmt.Fscan(in, &n, &targetX, &targetY)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	newTargetX := targetX + targetY
	newTargetY := targetX - targetY
	sum := 0
	for _, v := range nums {
		sum += v
	}

	if (newTargetX+sum)%2 == 1 {
		fmt.Fprintln(out, "No")
		return
	}
	newTargetX = (newTargetX + sum) / 2
	newTargetY = (newTargetY + sum) / 2
	if newTargetX < 0 || sum < newTargetX {
		fmt.Fprintln(out, "No")
		return
	}
	if newTargetY < 0 || sum < newTargetY {
		fmt.Fprintln(out, "No")
		return
	}

	res1, ok1 := SubsetSumTarget(nums, newTargetX)
	res2, ok2 := SubsetSumTarget(nums, newTargetY)
	if !ok1 && newTargetX != 0 {
		fmt.Fprintln(out, "No")
		return
	}
	if !ok2 && newTargetY != 0 {
		fmt.Fprintln(out, "No")
		return
	}

	fmt.Fprintln(out, "Yes")
	state := make([]int, n)
	for _, i := range res1 {
		state[i] |= 1
	}
	for _, i := range res2 {
		state[i] |= 2
	}

	cmd := "LUDR"
	for _, v := range state {
		fmt.Fprint(out, string(cmd[v]))
	}
}

func yuki04() {
	// !能否分成两个和相等的子集(天平称重问题)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
		sum += nums[i]
	}

	_, ok := SubsetSumTarget(nums, sum/2)
	can := ok || sum == 0
	if can && sum%2 == 0 {
		fmt.Fprintln(out, "possible")
	} else {
		fmt.Fprintln(out, "impossible")
	}
}

func testTime() {
	max := int(1e7)
	n := 40
	nums := make([]int, n)
	for i := range nums {
		nums[i] = max
	}
	sumNums := 0
	for _, v := range nums {
		sumNums += v
	}

	target := int(1e7)
	time1 := time.Now()
	SubsetSumTargetDp1(nums, target)
	time2 := time.Now()
	SubsetSumTargetDp2(nums, target)
	time3 := time.Now()
	SubsetSumTargetDp3(nums, target)
	time4 := time.Now()
	SubsetSumTargetDp4(nums, target)
	time5 := time.Now()
	SubsetSumTargetDp5(nums, target)
	time6 := time.Now()

	cost1 := n * max
	cost2 := int(math.Pow(float64(sumNums), 1.5)) / 100
	cost3 := n * target / 500 // TODO, 待调整
	cost4 := int(math.Pow(float64(sumNums), 1.5) / 500)
	cost5 := 1 << (n / 2)

	fmt.Println(time2.Sub(time1)) // 3.05312s
	fmt.Println(time3.Sub(time2)) // 1.752335916s
	fmt.Println(time4.Sub(time3)) // 63.696375ms
	fmt.Println(time5.Sub(time4)) // 4.00797975s
	fmt.Println(time6.Sub(time5)) // 24.724042ms
	fmt.Println(cost1, cost2, cost3, cost4, cost5)
}

const INF int = 1e18

// 能否用nums中的若干个数凑出和为target.
func SubsetSumTarget(nums []int, target int) (solution []int, ok bool) {
	if target <= 0 {
		return
	}
	n := len(nums)
	if n == 0 {
		return
	}
	max_ := 0
	sum := 0
	for _, v := range nums {
		max_ = max(max_, v)
		sum += v
	}

	cost1 := n * max_
	cost2 := int(math.Pow(float64(sum), 1.5) / 100)
	cost3 := n * target / 500 // 经验值
	cost4 := INF
	if n <= 60 {
		cost4 = 1 << (n / 2)
	}
	minCost := min(cost1, min(cost2, min(cost3, cost4)))
	if minCost == cost1 {
		return SubsetSumTargetDp1(nums, target)
	}
	if minCost == cost2 {
		return SubsetSumTargetDp2(nums, target)
	}
	if minCost == cost3 {
		return SubsetSumTargetDp3(nums, target)
	}
	return SubsetSumTargetDp5(nums, target)
}

// 能否用nums中的若干个数凑出和为target(KnapsackProblemWithBoundedWeights).
// "Linear Time Algorithms for Knapsack Problems with Bounded Weights" by David Pisinger
//
//	O(n*max(nums)))
func SubsetSumTargetDp1(nums []int, target int) (solution []int, ok bool) {
	if target <= 0 {
		return
	}

	n := len(nums)
	max_ := 0
	for _, v := range nums {
		max_ = max(max_, v)
	}
	right, curSum := 0, 0
	for right < n && curSum+nums[right] <= target {
		curSum += nums[right]
		right++
	}
	if right == n && curSum != target {
		return
	}

	offset := target - max_ + 1
	dp := make([]int, 2*max_)
	for i := range dp {
		dp[i] = -1
	}
	pre := make([][]int, n)
	for i := range pre {
		pre[i] = make([]int, 2*max_)
		for j := range pre[i] {
			pre[i][j] = -1
		}
	}

	dp[curSum-offset] = right
	for i := right; i < n; i++ {
		ndp := make([]int, len(dp))
		copy(ndp, dp)
		p := pre[i]
		a := nums[i]
		for j := 0; j < max_; j++ {
			if ndp[j+a] < dp[j] {
				ndp[j+a] = dp[j]
				p[j+a] = -2
			}
		}
		for j := 2*max_ - 1; j >= max_; j-- {
			for k := ndp[j] - 1; k >= max(dp[j], 0); k-- {
				if ndp[j-nums[k]] < k {
					ndp[j-nums[k]] = k
					p[j-nums[k]] = k
				}
			}
		}
		dp = ndp
	}

	if dp[max_-1] == -1 {
		return
	}

	used := make([]bool, n)
	i, j := n-1, max_-1
	for i >= right {
		p := pre[i][j]
		if p == -2 {
			used[i] = !used[i]
			j -= nums[i]
			i--
		} else if p == -1 {
			i--
		} else {
			used[p] = !used[p]
			j += nums[p]
		}
	}

	for i >= 0 {
		used[i] = !used[i]
		i--
	}

	for i := 0; i < n; i++ {
		if used[i] {
			solution = append(solution, i)
		}
	}

	ok = true
	return
}

// 能否用nums中的若干个数凑出和为target.
//
//	O(sum(nums)^1.5)
func SubsetSumTargetDp2(nums []int, target int) (solution []int, ok bool) {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if target > sum {
		return
	}
	counter := make([]int, sum+1)
	dp := make([]int, sum+1)
	last := make([]int, sum+1)
	id := make(map[int][]int)
	for i, v := range nums {
		if v <= sum {
			id[v] = append(id[v], i)
			counter[v]++
		}
	}

	dp[0] = 1
	for i := 1; i <= sum; i++ {
		if counter[i] == 0 {
			continue
		}
		for j := 0; j < i; j++ {
			c := 0
			for k := j; k <= sum; k += i {
				if dp[k] == 1 {
					c = counter[i]
				} else if c > 0 {
					dp[k] = 1
					c--
					last[k] = id[i][c]
				}
			}
		}
	}

	if dp[target] == 0 {
		return
	}

	for target > 0 {
		solution = append(solution, last[target])
		target -= nums[last[target]]
	}

	ok = true
	return
}

// Bitset优化dp.能否用nums中的若干个数凑出和为target.
//
//	O(n*target/w)
func SubsetSumTargetDp3(nums []int, target int) (solution []int, ok bool) {
	_enumerateBits64 := func(s uint64, f func(bit int)) {
		for s != 0 {
			i := bits.TrailingZeros64(s)
			f(i)
			s ^= 1 << i
		}
	}

	_argSort := func(nums []int) []int {
		order := make([]int, len(nums))
		for i := range order {
			order[i] = i
		}
		sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
		return order
	}

	n := len(nums)
	order := _argSort(nums)
	dp := NewBitsetDynamic(1, 1)
	last := make([]int, target+1)
	for i := range last {
		last[i] = -1
	}

	for k := 0; k < n; k++ {
		v := nums[order[k]]
		if v > target {
			continue
		}

		newSize := dp.Size() + v
		ndp := dp.CopyAndResize(newSize)
		ndp.IOrRange(v, newSize, dp)
		if ndp.Size() > target+1 {
			ndp.Resize(target + 1)
		}
		for i := 0; i < len(ndp.data); i++ {
			var updatedBits uint64
			if i < len(dp.data) {
				updatedBits = dp.data[i] ^ ndp.data[i]
			} else {
				updatedBits = ndp.data[i]
			}
			_enumerateBits64(updatedBits, func(p int) {
				last[(i<<6)|p] = order[k]
			})
		}

		dp = ndp
	}

	if target >= dp.Size() || !dp.Has(target) {
		return
	}

	for target > 0 {
		i := last[target]
		solution = append(solution, i)
		target -= nums[i]
	}

	ok = true
	return
}

// 能否用nums中的若干个数凑出和为target.
//
//	O(sum(nums)^1.5/w).常数较大.
func SubsetSumTargetDp4(nums []int, target int) (solution []int, ok bool) {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	n := len(nums)
	ids := make([][]int, sum+1)
	for i := 0; i < n; i++ {
		ids[nums[i]] = append(ids[nums[i]], i)
	}
	pre := make([][2]int, n)
	for i := range pre {
		pre[i] = [2]int{-1, -1}
	}
	var grpVals []int
	var rawIdx []int
	for x := 1; x <= sum; x++ {
		I := ids[x]
		for len(I) >= 3 {
			a, b := I[len(I)-1], I[len(I)-2]
			I = I[:len(I)-2]
			c := len(pre)
			pre = append(pre, [2]int{a, b})
			ids[2*x] = append(ids[2*x], c)
		}
		for _, i := range I {
			grpVals = append(grpVals, x)
			rawIdx = append(rawIdx, i)
		}
	}
	I, tmp := SubsetSumTargetDp3(grpVals, target)
	if !tmp {
		return
	}
	for _, i := range I {
		st := []int{rawIdx[i]}
		for len(st) > 0 {
			c := st[len(st)-1]
			st = st[:len(st)-1]
			if c < n {
				solution = append(solution, c)
				continue
			}
			a, b := pre[c][0], pre[c][1]
			st = append(st, a, b)
		}
	}
	ok = true
	return
}

// SubsetSumTargetMeetInTheMiddle
// 折半搜索.能否用nums中的若干个数凑出和为target.
//
//	O(2^(n/2))
func SubsetSumTargetDp5(nums []int, target int) (solution []int, ok bool) {
	_merge := func(a, b [][2]int) [][2]int {
		n1, n2 := len(a), len(b)
		res := make([][2]int, n1+n2)
		i, j, k := 0, 0, 0
		for i < n1 && j < n2 {
			if a[i][0] < b[j][0] {
				res[k] = a[i]
				i++
			} else {
				res[k] = b[j]
				j++
			}
			k++
		}
		for i < n1 {
			res[k] = a[i]
			i++
			k++
		}
		for j < n2 {
			res[k] = b[j]
			j++
			k++
		}
		return res
	}

	_subsetSumSortedWithState := func(nums []int) [][2]int {
		dp := [][2]int{{0, 0}} // (sum, state)
		for i, x := range nums {
			ndp := make([][2]int, len(dp))
			for j, p := range dp {
				ndp[j][0] = p[0] + x
				ndp[j][1] = p[1] | 1<<i
			}
			dp = _merge(dp, ndp)
		}
		return dp
	}

	_resolveState := func(leftSize, leftState, rightSize, rightState int) []int {
		res := []int{}
		for i := 0; i < leftSize; i++ {
			if leftState&(1<<i) != 0 {
				res = append(res, i)
			}
		}
		for i := 0; i < rightSize; i++ {
			if rightState&(1<<i) != 0 {
				res = append(res, i+leftSize)
			}
		}
		return res
	}

	n := len(nums)
	mid := n / 2
	dp1 := _subsetSumSortedWithState(nums[:mid])
	dp2 := _subsetSumSortedWithState(nums[mid:])
	left, right := 0, len(dp2)-1
	for left < len(dp1) && right >= 0 {
		sum := dp1[left][0] + dp2[right][0]
		if sum == target {
			solution = _resolveState(mid, dp1[left][1], n-mid, dp2[right][1])
			ok = true
			return
		}
		if sum < target {
			left++
		} else {
			right--
		}
	}

	return
}

// 动态bitset，支持切片操作.
type BitSetDynamic struct {
	n    int
	data []uint64
}

// 建立一个大小为 n 的 bitset，初始值为 filledValue.
func NewBitsetDynamic(n int, filledValue int) *BitSetDynamic {
	if !(filledValue == 0 || filledValue == 1) {
		panic("filledValue should be 0 or 1")
	}
	data := make([]uint64, n>>6+1)
	if filledValue == 1 {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= (len(data) << 6) - n
		}
	}
	return &BitSetDynamic{n: n, data: data}
}

func (bs *BitSetDynamic) Add(i int) *BitSetDynamic {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BitSetDynamic) Has(i int) bool {
	return bs.data[i>>6]>>(i&63)&1 == 1
}

func (bs *BitSetDynamic) Discard(i int) {
	bs.data[i>>6] &^= 1 << (i & 63)
}

func (bs *BitSetDynamic) Flip(i int) {
	bs.data[i>>6] ^= 1 << (i & 63)
}

func (bs *BitSetDynamic) AddRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] |= maskL ^ maskR
		return
	}
	bs.data[i] |= maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = ^uint64(0)
	}
	bs.data[i] |= ^maskR
}

func (bs *BitSetDynamic) DiscardRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] &= ^maskL | maskR
		return
	}
	bs.data[i] &= ^maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = 0
	}
	bs.data[i] &= maskR
}

func (bs *BitSetDynamic) FlipRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] ^= maskL ^ maskR
		return
	}
	bs.data[i] ^= maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = ^bs.data[i]
	}
	bs.data[i] ^= ^maskR
}

func (bs *BitSetDynamic) Fill(zeroOrOne int) {
	if zeroOrOne == 0 {
		for i := range bs.data {
			bs.data[i] = 0
		}
	} else {
		for i := range bs.data {
			bs.data[i] = ^uint64(0)
		}
		if bs.n != 0 {
			bs.data[len(bs.data)-1] >>= (len(bs.data) << 6) - bs.n
		}
	}
}

func (bs *BitSetDynamic) Clear() {
	for i := range bs.data {
		bs.data[i] = 0
	}
}

func (bs *BitSetDynamic) OnesCount(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > bs.n {
		end = bs.n
	}
	if start == 0 && end == bs.n {
		res := 0
		for _, v := range bs.data {
			res += bits.OnesCount64(v)
		}
		return res
	}
	pos1 := start >> 6
	pos2 := end >> 6
	if pos1 == pos2 {
		return bits.OnesCount64(bs.data[pos1] & (^uint64(0) << (start & 63)) & ((1 << (end & 63)) - 1))
	}
	count := 0
	if (start & 63) > 0 {
		count += bits.OnesCount64(bs.data[pos1] & (^uint64(0) << (start & 63)))
		pos1++
	}
	for i := pos1; i < pos2; i++ {
		count += bits.OnesCount64(bs.data[i])
	}
	if (end & 63) > 0 {
		count += bits.OnesCount64(bs.data[pos2] & ((1 << (end & 63)) - 1))
	}
	return count
}

func (bs *BitSetDynamic) AllOne(start, end int) bool {
	i := start >> 6
	if i == end>>6 {
		mask := ^uint64(0)<<(start&63) ^ ^uint64(0)<<(end&63)
		return (bs.data[i] & mask) == mask
	}
	mask := ^uint64(0) << (start & 63)
	if (bs.data[i] & mask) != mask {
		return false
	}
	for i++; i < end>>6; i++ {
		if bs.data[i] != ^uint64(0) {
			return false
		}
	}
	mask = ^uint64(0) << (end & 63)
	return ^(bs.data[end>>6] | mask) == 0
}

func (bs *BitSetDynamic) AllZero(start, end int) bool {
	i := start >> 6
	if i == end>>6 {
		mask := ^uint64(0)<<(start&63) ^ ^uint64(0)<<(end&63)
		return (bs.data[i] & mask) == 0
	}
	if (bs.data[i] >> (start & 63)) != 0 {
		return false
	}
	for i++; i < end>>6; i++ {
		if bs.data[i] != 0 {
			return false
		}
	}
	mask := ^uint64(0) << (end & 63)
	return (bs.data[end>>6] & ^mask) == 0
}

// 返回第一个 1 的下标，若不存在则返回-1.
func (bs *BitSetDynamic) IndexOfOne(position int) int {
	if position == 0 {
		for i, v := range bs.data {
			if v != 0 {
				return i<<6 | bs._lowbit(v)
			}
		}
		return -1
	}
	for i := position >> 6; i < len(bs.data); i++ {
		v := bs.data[i] & (^uint64(0) << (position & 63))
		if v != 0 {
			return i<<6 | bs._lowbit(v)
		}
		for i++; i < len(bs.data); i++ {
			if bs.data[i] != 0 {
				return i<<6 | bs._lowbit(bs.data[i])
			}
		}
	}
	return -1
}

// 返回第一个 0 的下标，若不存在则返回-1。
func (bs *BitSetDynamic) IndexOfZero(position int) int {
	if position == 0 {
		for i, v := range bs.data {
			if v != ^uint64(0) {
				return i<<6 | bs._lowbit(^v)
			}
		}
		return -1
	}

	i := position >> 6
	if i < len(bs.data) {
		v := bs.data[i]
		if position&63 != 0 {
			v |= ^((^uint64(0)) << (position & 63))
		}
		if ^v != 0 {
			res := i<<6 | bs._lowbit(^v)
			if res < bs.n {
				return res
			}
			return -1
		}
		for i++; i < len(bs.data); i++ {
			if ^bs.data[i] != 0 {
				res := i<<6 | bs._lowbit(^bs.data[i])
				if res < bs.n {
					return res
				}
				return -1
			}
		}
	}
	return -1
}

// 返回右侧第一个 1 的位置(`包含`当前位置).
//
//	如果不存在, 返回 n.
func (bs *BitSetDynamic) Next(index int) int {
	if index < 0 {
		index = 0
	}
	if index >= bs.n {
		return bs.n
	}
	k := index >> 6
	x := bs.data[k]
	s := index & 63
	x = (x >> s) << s
	if x != 0 {
		return (k << 6) | bs._lowbit(x)
	}
	for i := k + 1; i < len(bs.data); i++ {
		if bs.data[i] == 0 {
			continue
		}
		return (i << 6) | bs._lowbit(bs.data[i])
	}
	return bs.n
}

// 返回左侧第一个 1 的位置(`包含`当前位置).
//
//	如果不存在, 返回 -1.
func (bs *BitSetDynamic) Prev(index int) int {
	if index >= bs.n-1 {
		index = bs.n - 1
	}
	if index < 0 {
		return -1
	}
	k := index >> 6
	if (index & 63) < 63 {
		x := bs.data[k]
		x &= (1 << ((index & 63) + 1)) - 1
		if x != 0 {
			return (k << 6) | bs._topbit(x)
		}
		k--
	}
	for i := k; i >= 0; i-- {
		if bs.data[i] == 0 {
			continue
		}
		return (i << 6) | bs._topbit(bs.data[i])
	}
	return -1
}

func (bs *BitSetDynamic) Equals(other *BitSetDynamic) bool {
	if len(bs.data) != len(other.data) {
		return false
	}
	for i := range bs.data {
		if bs.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) IsSubset(other *BitSetDynamic) bool {
	if bs.n > other.n {
		return false
	}
	for i, v := range bs.data {
		if (v & other.data[i]) != v {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) IsSuperset(other *BitSetDynamic) bool {
	if bs.n < other.n {
		return false
	}
	for i, v := range other.data {
		if (v & bs.data[i]) != v {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) Ior(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] |= v
	}
	return bs
}

func (bs *BitSetDynamic) Iand(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] &= v
	}
	return bs
}

func (bs *BitSetDynamic) Ixor(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] ^= v
	}
	return bs
}

func (bs *BitSetDynamic) Or(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] | v
	}
	return res
}

func (bs *BitSetDynamic) And(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] & v
	}
	return res
}

func (bs *BitSetDynamic) Xor(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] ^ v
	}
	return res
}

func (bs *BitSetDynamic) IOrRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		bs.data[start>>6] |= other._get(a) << (start & 63)
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		bs.data[end>>6] |= other._get(b) << (end & 63)
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] |= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] |= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic) IAndRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		if other._get(a) == 0 {
			bs.data[start>>6] &^= 1 << (start & 63)
		}
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		if other._get(b) == 0 {
			bs.data[end>>6] &^= 1 << (end & 63)
		}
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] &= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] &= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic) IXorRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		bs.data[start>>6] ^= other._get(a) << (start & 63)
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		bs.data[end>>6] ^= other._get(b) << (end & 63)
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] ^= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] ^= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

// 类似js中类型数组的set操作.如果超出赋值范围，抛出异常.
//
//	other: 要赋值的bitset.
//	offset: 赋值的起始元素下标.
func (bs *BitSetDynamic) Set(other *BitSetDynamic, offset int) {
	left, right := offset, offset+other.n
	if right > bs.n {
		panic("out of range")
	}
	a, b := 0, other.n
	for left < right && (left&63) != 0 {
		if other.Has(a) {
			bs.Add(left)
		} else {
			bs.Discard(left)
		}
		a++
		left++
	}
	for left < right && (right&63) != 0 {
		right--
		b--
		if other.Has(b) {
			bs.Add(right)
		} else {
			bs.Discard(right)
		}
	}

	// other[a:b] -> this[start:end]
	l, r := left>>6, right>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] = other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] = (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic) Slice(start, end int) *BitSetDynamic {
	if start < 0 {
		start += bs.n
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += bs.n
	}
	if end > bs.n {
		end = bs.n
	}
	if start >= end {
		return NewBitsetDynamic(0, 0)
	}
	if start == 0 && end == bs.n {
		return bs.Copy()
	}
	res := NewBitsetDynamic(end-start, 0)
	remain := (end - start) & 63
	for i := 0; i < remain; i++ {
		if bs.Has(end - 1) {
			res.Add(end - start - 1)
		}
		end--
	}
	n := (end - start) >> 6
	hi := start & 63
	lo := 64 - hi
	s := start >> 6
	if hi == 0 {
		for i := 0; i < n; i++ {
			res.data[i] ^= bs.data[s+i]
		}
	} else {
		for i := 0; i < n; i++ {
			res.data[i] ^= (bs.data[s+i] >> hi) ^ (bs.data[s+i+1] << lo)
		}
	}
	return res
}

func (bs *BitSetDynamic) Copy() *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	copy(res.data, bs.data)
	return res
}

func (bs *BitSetDynamic) CopyAndResize(size int) *BitSetDynamic {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (uint64(1) << remainingBits) - 1
		newBits[len(newBits)-1] &= mask
	}
	return &BitSetDynamic{data: newBits, n: size}
}

func (bs *BitSetDynamic) Resize(size int) {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (uint64(1) << remainingBits) - 1
		newBits[len(newBits)-1] &= mask
	}
	bs.data = newBits
	bs.n = size
}

func (bs *BitSetDynamic) Expand(size int) {
	if size <= bs.n {
		return
	}
	bs.Resize(size)
}

func (bs *BitSetDynamic) BitLength() int {
	return bs._lastIndexOfOne() + 1
}

// 遍历所有 1 的位置.
func (bs *BitSetDynamic) ForEach(f func(pos int)) {
	for i, v := range bs.data {
		for v != 0 {
			j := (i << 6) | bs._lowbit(v)
			f(j)
			v &= v - 1
		}
	}
}

func (bs *BitSetDynamic) Size() int {
	return bs.n
}

func (bs *BitSetDynamic) String() string {
	sb := strings.Builder{}
	sb.WriteString("BitSetDynamic{")
	nums := []string{}
	bs.ForEach(func(pos int) {
		nums = append(nums, fmt.Sprintf("%d", pos))
	})
	sb.WriteString(strings.Join(nums, ","))
	sb.WriteString("}")
	return sb.String()
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 1, 2)
func (bs *BitSetDynamic) _topbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(x)
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 0, 2)
func (bs *BitSetDynamic) _lowbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return bits.TrailingZeros64(x)
}

func (bs *BitSetDynamic) _get(i int) uint64 {
	return bs.data[i>>6] >> (i & 63) & 1
}

func (bs *BitSetDynamic) _lastIndexOfOne() int {
	for i := len(bs.data) - 1; i >= 0; i-- {
		x := bs.data[i]
		if x != 0 {
			return (i << 6) | (bs._topbit(x))
		}
	}
	return -1
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

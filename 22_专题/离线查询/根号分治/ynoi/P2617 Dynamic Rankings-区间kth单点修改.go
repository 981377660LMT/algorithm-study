// PointSetRangeKth
// P2617 Dynamic Rankings
// https://www.luogu.com.cn/problem/P2617
//
// 序列分块+值域分块的做法,离线，时空复杂度均为O(n*sqrt(n)).
// 维护两个信息：
// sum1[i][j] -> 前i个序列块中有多少个数出现在第j个值域块中.(用于维护<个数)
// sum2[i][x] -> 前i个序列块中有多少个数恰好等于x.(用于维护==个数)
//
// 二分法要求我们快速回答这个问题:
// !在这个数字集合当中有几个数字比x大
// 难以维护，与分块不太搭
// 而基于值域分块求解区间第k小仅仅需要回答这个问题：
// !这个数字集合当中有几个数字恰好落在了第i个值域块里，在这个数字集合当中有几个数字恰好为x。
// !先确定答案在哪个值域块，再在那个值域块中求出精确的答案
//
// 序列分块+值域分块的好处是可以支持区间修改.缺点是离线.
// 为了减少空间复杂度O(nsqrt(n))，可以把块的大小增大一些，调成4-10倍.
// 这样1e5的数据大约需要1e5*32/10=3e6个int,接近O(nlogn)的空间复杂度了.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// 0 start end kth: 查询[start, end)内排名为k(k>=0)的数.
// 1 pos target: 将pos位置的数修改为target.
func PointSetRangeKth(nums []int, operations [][4]int) []int {
	allNums := append(nums[:0:0], nums...)
	for i := range operations {
		kind := operations[i][0]
		if kind == 1 {
			allNums = append(allNums, operations[i][2])
		}
	}
	set := make(map[int]struct{})
	for _, v := range allNums {
		set[v] = struct{}{}
	}
	sorted := make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	valueToId := make(map[int]int, len(sorted))
	idToValue := make([]int, len(sorted))
	for i, v := range sorted {
		valueToId[v] = i
		idToValue[i] = v
	}

	for i := range nums {
		nums[i] = valueToId[nums[i]]
	}

	blockSize1 := int(math.Sqrt(float64(len(nums)))+1) * 4 // !增加块的大小以减少块的数量，减少内存占用.
	block1 := UseBlock(len(nums), blockSize1)
	belong1, blockStart1, blockEnd1, blockCount1 := block1.belong, block1.blockStart, block1.blockEnd, block1.blockCount
	block2 := UseBlock(len(sorted), int(math.Sqrt(float64(len(sorted)))+1))
	belong2, blockStart2, blockEnd2, blockCount2 := block2.belong, block2.blockStart, block2.blockEnd, block2.blockCount
	sum1 := make([][]int, blockCount1+1)
	for i := range sum1 {
		sum1[i] = make([]int, blockCount2)
	}
	sum2 := make([][]int, blockCount1+1)
	for i := range sum2 {
		sum2[i] = make([]int, len(sorted))
	}
	tmpCounter1 := make([]int, blockCount2) // 查询时临时保存散块内出现在第i个值域块中的数的个数.
	tmpCounter2 := make([]int, len(sorted)) // 查询时临时保存散块内等于i的数的个数.

	// 初始化前缀和.
	init := func() {
		for bid := 0; bid < blockCount1; bid++ {
			copy(sum1[bid+1], sum1[bid])
			copy(sum2[bid+1], sum2[bid])
			s1, s2 := sum1[bid+1], sum2[bid+1]
			for i := blockStart1[bid]; i < blockEnd1[bid]; i++ {
				num := nums[i]
				vid := belong2[num]
				s1[vid]++
				s2[num]++
			}
		}
	}
	init()

	// 更新散块信息.
	updateTmp := func(start, end, delta int) {
		for i := start; i < end; i++ {
			num := nums[i]
			vid := belong2[num]
			tmpCounter1[vid] += delta
			tmpCounter2[num] += delta
		}
	}

	// 查询区间第k小(k>=0).
	queryKth := func(start, end, kth int) int {
		bid1, bid2 := belong1[start], belong1[end-1]

		// 散块.
		if bid1 == bid2 {
			updateTmp(start, end, 1)
			todo := kth + 1
			for vid := 0; vid < blockCount2; vid++ {
				if tmp := todo - tmpCounter1[vid]; tmp > 0 {
					todo = tmp
					continue
				}
				// 答案在这个值域中.
				for j := blockStart2[vid]; j < blockEnd2[vid]; j++ {
					todo -= tmpCounter2[j]
					if todo <= 0 {
						updateTmp(start, end, -1)
						return j
					}
				}
			}
		} else {
			// 完整块.
			updateTmp(start, blockEnd1[bid1], 1)
			updateTmp(blockStart1[bid2], end, 1)
			todo := kth + 1
			for vid := 0; vid < blockCount2; vid++ {
				curCount := tmpCounter1[vid] + sum1[bid2][vid] - sum1[bid1+1][vid]
				if tmp := todo - curCount; tmp > 0 {
					todo = tmp
					continue
				}
				// 答案在这个值域中.
				for j := blockStart2[vid]; j < blockEnd2[vid]; j++ {
					curCount := tmpCounter2[j] + sum2[bid2][j] - sum2[bid1+1][j]
					todo -= curCount
					if todo <= 0 {
						updateTmp(start, blockEnd1[bid1], -1)
						updateTmp(blockStart1[bid2], end, -1)
						return j
					}
				}
			}
		}

		panic("unreachable")
	}

	// 单点修改.
	pointSet := func(pos, target int) {
		if nums[pos] == target {
			return
		}
		preValue, curValue := nums[pos], target
		preVid, curVid := belong2[preValue], belong2[curValue]
		for bid := belong1[pos]; bid < blockCount1; bid++ {
			s1, s2 := sum1[bid+1], sum2[bid+1]
			s1[preVid]--
			s1[curVid]++
			s2[preValue]--
			s2[curValue]++
		}
		nums[pos] = target
	}

	res := []int{}
	for _, op := range operations {
		kind := op[0]
		if kind == 0 {
			start, end, k := op[1], op[2], op[3]
			kth := queryKth(start, end, k)
			res = append(res, idToValue[kth])
		} else {
			pos, target := op[1], op[2]
			target = valueToId[target]
			pointSet(pos, target)
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	operations := make([][4]int, q)
	for i := range operations {
		var op string
		fmt.Fscan(in, &op)
		if op == "Q" {
			var start, end, k int
			fmt.Fscan(in, &start, &end, &k)
			start--
			k--
			operations[i] = [4]int{0, start, end, k}
		} else {
			var pos, target int
			fmt.Fscan(in, &pos, &target)
			pos--
			operations[i] = [4]int{1, pos, target}
		}
	}

	res := PointSetRangeKth(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
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

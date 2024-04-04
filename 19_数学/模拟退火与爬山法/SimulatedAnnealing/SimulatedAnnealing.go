// 模拟退火优化求解最小值(最优解).
// https://www.cnblogs.com/shenben/p/11342308.html
// https://vlight.me/2018/06/08/Simulated-Annealing/
// https://oi-wiki.org//misc/simulated-annealing/
// https://www.luogu.com/article/b2recz8n
//
// Atcoder Heuristic Contest：
// https://www.youtube.com/watch?v=hxic3DVMPTg&list=PLLeJZg4opYKY6yCPd7j0b5NS4b7krV2yF
//
// !技巧：可以在时限内重复跑 SA 取最优值，防止脸黑
//
// api:
//
//	SetK(k float64) *SimulatedAnnealing[X] (影响接受较差解的概率)
//	SetReduce(reduce float64) *SimulatedAnnealing[X] (温度衰减率)
//	SetTimeLimitMs(timeLimitMs float64) *SimulatedAnnealing[X] (时间限制)
//	Optimize(initX X) (求解最优解)
//	GetBestX() X (获取最优解)
//	GetBestY() float64 (获取最优值)
//	SetInitTemperature(initTemperature float64) *SimulatedAnnealing[X]
//	SetThresholdTemperature(thresholdTemperature float64) *SimulatedAnnealing[X]

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// 参数说明:
//
//	initTemperature: 初始温度.一般设置为 1e5.
//	thresholdTemperature: 目标温度阈值.一般设置为 1e-8.
//	!k: 玻尔兹曼常数.一般设置为 1.越大表示越难接受较差的解.对答案影响较大.
//	!reduce: 温度衰减率.一般设置为 0.99 到 0.999.对答案影响较大.
//	initX: 初始的自变量.
//
// 如何生成新解：
//
//	最小化函数值：温度越低，自变量随机偏移越小。
//	坐标系内：随机生成一个点，或者生成一个向量。
//	序列问题：random.shuffle()或者随机交换两个元素。
//	网格问题：可以看做二维序列，每次交换两个格子即可。
type SimulatedAnnealing[X any] struct {
	bestX X       // 最优解.
	bestY float64 // 最优解的分数.

	reduce               float64 // 温度衰减率.默认值为 0.99.
	initTemperature      float64 // 初始温度.默认值为 1e5.
	thresholdTemperature float64 // 温度阈值.默认值为 1e-8.
	k                    float64 // 玻尔兹曼常数.默认值为 1.
	timeLimitMs          float64 // 时间限制(毫秒).默认为-1，表示不限制时间.

	eval      func(x X) float64                       // 函数值.
	next      func(oldX X, temperature float64) X     // 生成新的自变量.
	summarize func(newX X, newY float64, accept bool) // 下一次迭代的x和y、本轮是否接受了新解.

	calculated bool
}

func NewSimulatedAnnealing[X any](
	eval func(x X) float64,
	next func(oldX X, temperature float64) X,
	summarize func(newX X, newY float64, accept bool),
) *SimulatedAnnealing[X] {
	return &SimulatedAnnealing[X]{
		bestY:                math.MaxFloat64,
		reduce:               0.99,
		initTemperature:      1e5,
		thresholdTemperature: 1e-8,
		k:                    1,
		timeLimitMs:          -1,
		eval:                 eval,
		next:                 next,
		summarize:            summarize,
	}
}

func (sa *SimulatedAnnealing[X]) Optimize(initX X) {
	if sa.timeLimitMs == -1 {
		sa._run(initX)
	} else {
		sa._runWithinTimeLimitMs(initX, sa.timeLimitMs)
	}
}

func (sa *SimulatedAnnealing[X]) GetBestX() X       { return sa.bestX }
func (sa *SimulatedAnnealing[X]) GetBestY() float64 { return sa.bestY }

func (sa *SimulatedAnnealing[X]) SetReduce(reduce float64) *SimulatedAnnealing[X] {
	sa.reduce = reduce
	return sa
}

func (sa *SimulatedAnnealing[X]) SetInitTemperature(initTemperature float64) *SimulatedAnnealing[X] {
	sa.initTemperature = initTemperature
	return sa
}

func (sa *SimulatedAnnealing[X]) SetThresholdTemperature(thresholdTemperature float64) *SimulatedAnnealing[X] {
	sa.thresholdTemperature = thresholdTemperature
	return sa
}

func (sa *SimulatedAnnealing[X]) SetK(k float64) *SimulatedAnnealing[X] {
	sa.k = k
	return sa
}

func (sa *SimulatedAnnealing[X]) SetTimeLimitMs(timeLimitMs float64) *SimulatedAnnealing[X] {
	sa.timeLimitMs = timeLimitMs
	return sa
}

func (sa *SimulatedAnnealing[X]) _run(initX X) {
	x := initX
	y := sa.eval(x)
	t := sa.initTemperature
	k, threshold := sa.k, sa.thresholdTemperature
	for t > threshold {
		nextX := sa.next(x, t)
		nextY := sa.eval(nextX)
		accpet := false
		// 最小直接取，或者以一定概率接受较大的值
		if nextY < y || math.Exp((y-nextY)/(k*t)) > rand.Float64() {
			accpet = true
			x = nextX
			y = nextY
		}
		sa.summarize(nextX, nextY, accpet)
		t *= sa.reduce
	}

	if !sa.calculated || sa.bestY > y {
		sa.bestX = x
		sa.bestY = y
		sa.calculated = true
	}
}

func (sa *SimulatedAnnealing[X]) _runWithinTimeLimitMs(initX X, timeLimitMs float64) {
	timeLimitMs64 := int64(timeLimitMs)
	x := initX
	y := sa.eval(x)
	t := sa.initTemperature
	k := sa.k
	startTime := time.Now()
	for time.Since(startTime).Milliseconds() < timeLimitMs64 {
		nextX := sa.next(x, t)
		nextY := sa.eval(nextX)
		accept := false
		if nextY < y || math.Exp((y-nextY)/(k*t)) > rand.Float64() {
			accept = true
			x = nextX
			y = nextY
		}
		sa.summarize(nextX, nextY, accept)
		t *= sa.reduce
	}

	if !sa.calculated || sa.bestY > y {
		sa.bestX = x
		sa.bestY = y
		sa.calculated = true
	}
}

func min64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func abs64(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
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

func main() {
	// demo()
	// [2999,3914,1064,927,64,1130,2048,235,159,3549,241,987,972,976,279,1004]
	fmt.Println(canPartitionKSubsets([]int{2999, 3914, 1064, 927, 64, 1130, 2048, 235, 159, 3549, 241, 987, 972, 976, 279, 1004}, 4))
}

func demo() {
	// 求解 f(x) = x^4 + x 的最小值
	eval := func(x float64) float64 {
		return (x*x*x*x + x)
	}
	next := func(old float64, temperature float64) float64 {
		return old + (2*rand.Float64()-1)*temperature
	}
	summarize := func(nextX float64, nextY float64, accept bool) {}

	sa := NewSimulatedAnnealing[float64](eval, next, summarize)
	sa.Optimize(0)
	fmt.Println(sa.GetBestY())
}

// 698. 划分为k个相等的子集 (状压dp，判定性)
// https://leetcode.cn/problems/partition-to-k-equal-sum-subsets/description/
// 给定一个整数数组  nums 和一个正整数 k，找出是否有可能把这个数组分成 k 个非空子集，其总和都相等。
func canPartitionKSubsets(nums []int, k int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%k != 0 {
		return false
	}
	target := sum / k

	n := len(nums)

	type Arg = []int

	eval := func(x Arg) float64 {
		diff := sum
		vi, gi := 0, 0
		for vi < n && gi < k {
			curSum := 0
			for vi < n && curSum+x[vi] <= target {
				curSum += x[vi]
				vi++
			}
			diff -= curSum
			gi++
		}
		return float64(diff)
	}

	swapI, swapJ := 0, 0
	next := func(old Arg, _ float64) Arg {
		swapI, swapJ = rand.Intn(n), rand.Intn(n)
		if swapI == swapJ {
			swapJ = (swapJ + 1) % n
			return old
		}
		old[swapI], old[swapJ] = old[swapJ], old[swapI]
		return old
	}

	ok := false
	summarize := func(nextX Arg, nextY float64, accept bool) {
		if nextY == 0 {
			ok = true
		}
		if !accept {
			nextX[swapI], nextX[swapJ] = nextX[swapJ], nextX[swapI]
		}
	}

	sa := NewSimulatedAnnealing(eval, next, summarize)
	sa.SetK(5)
	for i := 0; i < 100; i++ {
		rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
		sa.Optimize(nums)
		if ok {
			return true
		}
	}
	return false
}

// 1515. 服务中心的最佳位置 (凸优化, reduce=0.99, k=0.01)
// https://leetcode.cn/problems/best-position-for-a-service-centre/description/
// 给你一个数组 positions ，其中 positions[i] = [xi, yi] 表示第 i 个客户在二维地图上的位置，返回到所有客户的 欧几里得距离的最小总和 。
func getMinDistSum(positions [][]int) float64 {
	type Arg = [2]float64
	eval := func(arg Arg) float64 {
		x, y := arg[0], arg[1]
		res := 0.0
		for _, p := range positions {
			px, py := float64(p[0]), float64(p[1])
			res += math.Sqrt((px-x)*(px-x) + (py-y)*(py-y))
		}
		return res
	}
	next := func(oldArg Arg, temperature float64) Arg {
		x, y := oldArg[0], oldArg[1]
		return Arg{x + (2*rand.Float64()-1)*temperature, y + (2*rand.Float64()-1)*temperature}
	}
	summarize := func(nextArg Arg, nextY float64, accept bool) {}

	sa := NewSimulatedAnnealing[Arg](eval, next, summarize)
	sa.SetK(0.01) // !k很小，表示接受较差的解的概率较大
	res := math.MaxFloat64
	for i := 0; i < 10; i++ {
		sa.Optimize(Arg{0, 0})
		res = min64(res, sa.GetBestY())
	}
	return res
}

// 1723. 完成所有工作的最短时间 (状压dp，最小化)
// https://leetcode.cn/problems/find-minimum-time-to-finish-all-jobs/solutions/1/gong-shui-san-xie-yi-ti-shuang-jie-jian-4epdd/
func minimumTimeRequired(jobs []int, k int) int {
	n := len(jobs)

	type Arg = []int
	eval := func(arg Arg) float64 {
		// 分配任务到最小的组
		groupSum := make([]int, k)
		for _, job := range arg {
			best, bestArg := groupSum[0], 0
			for i, sum := range groupSum {
				if sum < best {
					best, bestArg = sum, i
				}
			}
			groupSum[bestArg] += job
		}
		max_ := 0
		for _, sum := range groupSum {
			max_ = max(max_, sum)
		}
		return float64(max_)
	}

	swapI, swapJ := 0, 0
	next := func(old Arg, _ float64) Arg {
		swapI, swapJ = rand.Intn(n), rand.Intn(n)
		if swapI == swapJ {
			swapJ = (swapJ + 1) % n
			return old
		}
		old[swapI], old[swapJ] = old[swapJ], old[swapI]
		return old
	}

	res := math.MaxFloat64
	summarize := func(nextArg Arg, nextY float64, accept bool) {
		res = min64(res, nextY)
		if !accept {
			nextArg[swapI], nextArg[swapJ] = nextArg[swapJ], nextArg[swapI]
		}
	}

	sa := NewSimulatedAnnealing[Arg](eval, next, summarize)
	for i := 0; i < 20; i++ {
		rand.Shuffle(len(jobs), func(i, j int) { jobs[i], jobs[j] = jobs[j], jobs[i] })
		sa.Optimize(jobs)
	}
	return int(res)
}

// 1815. 得到新鲜甜甜圈的最多组数 (状压dp，最大化)
// https://leetcode.cn/problems/maximum-number-of-groups-getting-fresh-donuts/
// https://zhuanlan.zhihu.com/p/600471525
func maxHappyGroups(batchSize int, groups []int) int {
	n := len(groups)

	type Arg = []int
	eval := func(x Arg) float64 {
		remain := 0
		happyCount := 0
		for _, need := range x {
			if remain == 0 {
				happyCount++
			}
			remain = (remain + need) % batchSize
		}
		return -float64(happyCount) // 因为要最大化，所以取负数
	}

	swapI, swapJ := 0, 0
	next := func(oldX Arg, _ float64) Arg {
		swapI, swapJ = rand.Intn(n), rand.Intn(n)
		if swapI == swapJ {
			swapJ = (swapJ + 1) % n
			return oldX
		}
		oldX[swapI], oldX[swapJ] = oldX[swapJ], oldX[swapI]
		return oldX
	}

	res := float64(0)
	summarize := func(newX Arg, newY float64, accept bool) {
		res = max64(res, -newY) // 因为要最大化，所以取负数
		if !accept {
			newX[swapI], newX[swapJ] = newX[swapJ], newX[swapI]
		}
	}

	sa := NewSimulatedAnnealing[Arg](eval, next, summarize)
	for i := 0; i < 30; i++ {
		rand.Shuffle(len(groups), func(i, j int) { groups[i], groups[j] = groups[j], groups[i] })
		sa.Optimize(groups)
	}
	return int(res)
}

// 1879. 两个数组最小的异或值之和 (状压dp)
// https://leetcode.cn/problems/minimum-xor-sum-of-two-arrays/description/
func minimumXORSum(nums1 []int, nums2 []int) int {
	n := len(nums1)

	type Arg = []int
	eval := func(x Arg) float64 {
		res := 0
		for i := 0; i < n; i++ {
			res += x[i] ^ nums2[i]
		}
		return float64(res)
	}

	swapI, swapJ := 0, 0
	next := func(old Arg, _ float64) Arg {
		swapI, swapJ = rand.Intn(n), rand.Intn(n)
		if swapI == swapJ {
			swapJ = (swapJ + 1) % n
			return old
		}
		old[swapI], old[swapJ] = old[swapJ], old[swapI]
		return old
	}

	res := math.MaxFloat64
	summarize := func(nextX Arg, nextY float64, accept bool) {
		res = min64(res, nextY)
		if !accept {
			nextX[swapI], nextX[swapJ] = nextX[swapJ], nextX[swapI]
		}
	}

	sa := NewSimulatedAnnealing[Arg](eval, next, summarize)
	for i := 0; i < 20; i++ {
		rand.Shuffle(len(nums1), func(i, j int) { nums1[i], nums1[j] = nums1[j], nums1[i] })
		sa.Optimize(nums1)
	}
	return int(res)
}

// 模板题 https://www.luogu.com.cn/problem/P1337
// todo 教学题 https://atcoder.jp/contests/intro-heuristics/tasks/intro_heuristics_a
//  https://atcoder.jp/contests/ahc001/tasks/ahc001_a
//  https://atcoder.jp/contests/ahc002/tasks/ahc002_a

func P1337() {}

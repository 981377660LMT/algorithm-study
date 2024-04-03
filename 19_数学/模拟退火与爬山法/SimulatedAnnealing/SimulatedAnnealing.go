// 模拟退火优化求解最小值(最优解).
// https://www.cnblogs.com/shenben/p/11342308.html
// https://vlight.me/2018/06/08/Simulated-Annealing/
// https://oi-wiki.org//misc/simulated-annealing/
// https://www.luogu.com/article/b2recz8n
// !技巧：可以在时限内重复跑 SA 取最优值，防止脸黑

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	demo()
}

// 1515. 服务中心的最佳位置
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
	abandon := func(old Arg) {}

	sa := NewSimulatedAnnealing[Arg](eval, next, abandon)
	sa.SetK(0.01) // !k很小，表示接受较差的解的概率较大
	res := math.MaxFloat64
	for i := 0; i < 10; i++ {
		sa.Optimize(Arg{0, 0})
		res = min64(res, sa.GetBestY())
	}
	return res
}

func demo() {
	// 求解 f(x) = x^4 + x 的最小值
	eval := func(x float64) float64 {
		return (x*x*x*x + x)
	}
	next := func(old float64, temperature float64) float64 {
		return old + (2*rand.Float64()-1)*temperature
	}
	abandon := func(old float64) {}

	sa := NewSimulatedAnnealing[float64](eval, next, abandon)
	sa.Optimize(0)
	fmt.Println(sa.GetBestY())
}

// 模板题 https://www.luogu.com.cn/problem/P1337
// LC1515 https://leetcode.cn/problems/best-position-for-a-service-centre/
// http://poj.org/problem?id=2420
// UVa 10228 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=14&page=show_problem&problem=1169
// todo 教学题 https://atcoder.jp/contests/intro-heuristics/tasks/intro_heuristics_a
//  https://atcoder.jp/contests/ahc001/tasks/ahc001_a
//  https://atcoder.jp/contests/ahc002/tasks/ahc002_a

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

	eval    func(x X) float64                   // 函数值.
	next    func(oldX X, temperature float64) X // 生成新的自变量.
	abandon func(oldX X)                        // 处理舍弃的自变量.

	calculated bool
}

func NewSimulatedAnnealing[X any](
	eval func(x X) float64,
	next func(oldX X, temperature float64) X,
	abandon func(oldX X),
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
		abandon:              abandon,
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
		// 最小直接取，或者以一定概率接受较大的值
		if nextY < y || math.Exp((y-nextY)/(k*t)) > rand.Float64() {
			sa.abandon(x)
			x = nextX
			y = nextY
		}
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
		if nextY < y || math.Exp((y-nextY)/(k*t)) > rand.Float64() {
			sa.abandon(x)
			x = nextX
			y = nextY
		}
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

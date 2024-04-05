// 模拟退火优化求解最小值(最优解).
// https://www.cnblogs.com/shenben/p/11342308.html
// https://vlight.me/2018/06/08/Simulated-Annealing/
// https://oi-wiki.org//misc/simulated-annealing/
// https://www.luogu.com/article/b2recz8n
// 教学题 https://atcoder.jp/contests/intro-heuristics/tasks/intro_heuristics_a
//        https://atcoder.jp/contests/ahc001/tasks/ahc001_a
//        https://atcoder.jp/contests/ahc002/tasks/ahc002_a

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
//	SetInitTemperature(initTemperature float64) *SimulatedAnnealing[X] (取决于数据量范围设定)
//	Optimize(initX X) (求解最优解)
//	GetBestX() X (获取最优解)
//	GetBestY() float64 (获取最优值)
//	SetThresholdTemperature(thresholdTemperature float64) *SimulatedAnnealing[X]
//
// 可以考虑模拟退火的：
// 1. 状压题、搜索题(n非常小) -> 最优序列
// 2. 凸优化问题(爬山法)
// 3. 计算几何问题

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

// 参数说明:
//
//	initTemperature: 初始温度.一般设置为 1e5.
//	thresholdTemperature: 目标温度阈值.一般设置为 1e-8.
//	!k: 玻尔兹曼常数.一般设置为 1.越小表示越难接受较差的解(爬山法).对答案影响较大.
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

	evaluate  func(x X) float64                       // 估值函数.
	next      func(oldX X, temperature float64) X     // 生成新的自变量.
	summarize func(newX X, newY float64, accept bool) // 下一次迭代的x和y、本轮是否接受了新解.

	calculated bool
}

func NewSimulatedAnnealing[X any](
	evaluate func(x X) float64, next func(oldX X, temperature float64) X,
	summarize func(newX X, newY float64, accept bool), // 可选
) *SimulatedAnnealing[X] {
	return &SimulatedAnnealing[X]{
		bestY:                math.MaxFloat64,
		reduce:               0.99,
		initTemperature:      1e5,
		thresholdTemperature: 1e-8,
		k:                    1,
		timeLimitMs:          -1,
		evaluate:             evaluate,
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
	y := sa.evaluate(x)
	t := sa.initTemperature
	k, threshold := sa.k, sa.thresholdTemperature
	for t > threshold {
		nextX := sa.next(x, t)
		nextY := sa.evaluate(nextX)
		accept := false
		// 最小直接取，或者以一定概率接受较大的值
		if nextY < y || math.Exp((y-nextY)/(k*t)) > rand.Float64() {
			accept = true
			x = nextX
			y = nextY
		}
		if sa.summarize != nil {
			sa.summarize(nextX, nextY, accept)
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
	y := sa.evaluate(x)
	t := sa.initTemperature
	k := sa.k
	startTime := time.Now()
	for time.Since(startTime).Milliseconds() < timeLimitMs64 {
		nextX := sa.next(x, t)
		nextY := sa.evaluate(nextX)
		accept := false
		if nextY < y || math.Exp((y-nextY)/(k*t)) > rand.Float64() {
			accept = true
			x = nextX
			y = nextY
		}
		if sa.summarize != nil {
			sa.summarize(nextX, nextY, accept)
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
	// SP34()
	// P1337()
	// P2503()
	// P2538()
	// P3878()
	P3936()

	// P5544()
}

func demo() {
	// 求解 f(x) = x^4 + x 的最小值
	evaluate := func(x float64) float64 {
		return (x*x*x*x + x)
	}
	next := func(old float64, temperature float64) float64 {
		return old + (2*rand.Float64()-1)*temperature
	}
	summarize := func(nextX float64, nextY float64, accept bool) {}

	sa := NewSimulatedAnnealing[float64](evaluate, next, summarize)
	sa.Optimize(0)
	fmt.Println(sa.GetBestY())
}

// RUNAWAY - Run Away (最大化)
// https://www.luogu.com.cn/problem/SP34
// 在给定范围内找一个点，使得距离所有点的最小值最大。
func SP34() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	// [0,maxX] * [0,maxY]
	solve := func(points [][2]int, maxX, maxY int) (bestX, bestY float64) {
		type Arg = [2]float64
		res := 0.0
		evaluate := func(arg Arg) float64 {
			x, y := arg[0], arg[1]
			minDist := math.MaxFloat64
			for _, p := range points {
				px, py := float64(p[0]), float64(p[1])
				minDist = min64(minDist, math.Sqrt((px-x)*(px-x)+(py-y)*(py-y)))
			}
			if minDist > res {
				res = minDist
				bestX, bestY = x, y
			}
			return -minDist // 最大化最小值
		}

		next := func(oldArg Arg, temperature float64) Arg {
			x, y := oldArg[0], oldArg[1]
			nextX, nextY := x+(2*rand.Float64()-1)*temperature, y+(2*rand.Float64()-1)*temperature
			if nextX < 0 {
				nextX = 0
			}
			if nextX > float64(maxX) {
				nextX = float64(maxX)
			}
			if nextY < 0 {
				nextY = 0
			}
			if nextY > float64(maxY) {
				nextY = float64(maxY)
			}
			return Arg{nextX, nextY}
		}

		sa := NewSimulatedAnnealing[Arg](evaluate, next, nil)
		sa.SetK(0.01) // !k很小，表示接受较差的解的概率较小，接近爬山法
		sa.SetTimeLimitMs(100)
		for i := 0; i < 8; i++ { // 跑8轮，每轮100ms
			sa.Optimize(Arg{float64(maxX) / 2, float64(maxY) / 2})
		}

		return
	}

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		var x, y, m int
		fmt.Fscan(in, &x, &y, &m)
		points := make([][2]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &points[j][0], &points[j][1])
		}
		bestX, bestY := solve(points, x, y)
		fmt.Fprintf(out, "The safest point is (%.1f, %.1f)\n", bestX, bestY)
	}
}

// P1337 [JSOI2004] 平衡点 / 吊打XXX (凸优化，最小化)
// https://www.luogu.com.cn/problem/P1337
// 给定若干个点(x,y,weight)，求这些点的重心坐标.
// 输出两个浮点数（保留小数点后三位）表示横坐标和纵坐标，两个数以一个空格隔开.
func P1337() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	points := make([][3]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1], &points[i][2])
	}

	type Arg = [2]float64
	evaluate := func(arg Arg) float64 {
		x, y := arg[0], arg[1]
		res := 0.0
		for _, p := range points {
			px, py, w := float64(p[0]), float64(p[1]), float64(p[2])
			res += w * math.Sqrt((px-x)*(px-x)+(py-y)*(py-y))
		}
		return res
	}

	next := func(oldArg Arg, temperature float64) Arg {
		x, y := oldArg[0], oldArg[1]
		return Arg{x + (2*rand.Float64()-1)*temperature, y + (2*rand.Float64()-1)*temperature}
	}

	sa := NewSimulatedAnnealing[Arg](evaluate, next, nil)
	sa.SetK(0.001) // !k很小，表示接受较差的解的概率较小，接近爬山法
	sa.SetTimeLimitMs(100)
	for i := 0; i < 8; i++ { // 跑8轮，每轮100ms
		sa.Optimize(Arg{0, 0})
	}

	res := sa.GetBestX()
	fmt.Fprintf(out, "%.3f %.3f\n", res[0], res[1])
}

// P2503 [HAOI2006] 均分数据 (状压dp，最小化)
// https://www.luogu.com.cn/problem/P2503
// 把n个数分成k组，每组单独求和，使这k个数的方差最小
// 输出一行一个实数，表示最小均方差的值(保留小数点后两位数字)。
// !类似工人分配，把每个数加到当前和最小的组里.
func P2503() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	totalSum := 0
	for _, v := range nums {
		totalSum += v
	}
	avg := float64(totalSum) / float64(k)

	type Arg = []int
	evaluate := func(x Arg) float64 {
		groupSum := make([]int, k)
		for _, v := range x {
			bestArg, best := 0, groupSum[0]
			for i, sum := range groupSum {
				if sum < best {
					bestArg, best = i, sum
				}
			}
			groupSum[bestArg] += v
		}

		res := 0.0
		for _, sum := range groupSum {
			res += (float64(sum) - avg) * (float64(sum) - avg)
		}
		res /= float64(k)
		return res
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

	sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
	for i := 0; i < 20; i++ {
		rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
		sa.Optimize(nums)
	}
	fmt.Fprintf(out, "%.2f\n", math.Sqrt(res))
}

// P2538 [SCOI2008] 城堡 (dp，最小化)
// https://www.luogu.com.cn/problem/P2538
// 给定一个无向带权基环树，图中有m个特殊点.
// !你需要将不超过k个非特殊点变为特殊点，
// !使得所有非特殊点到最近的特殊点的距离的最大值最小。
// 求出这个最小值.
// n<=50.
// 对所有的非特殊点，分为两类：一类是变为特殊点A，一类是不变B.
// 估值函数：dijkstra求出非特殊点到最近的特殊点的距离.
// !生成新解：随机交换A和B种任意两个点.
// TO:WA
func P2538() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	nexts := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nexts[i])
	}
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weights[i])
	}
	isSpecial := make([]bool, n)
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(in, &x)
		isSpecial[x] = true
	}

	if m+k == n {
		fmt.Fprintln(out, 0)
		return
	}

	adjList := make([][][2]int, n) // (to, weight)
	for i := 0; i < n; i++ {
		next, weight := nexts[i], weights[i]
		adjList[i] = append(adjList[i], [2]int{next, weight})
		adjList[next] = append(adjList[next], [2]int{i, weight})
	}
	distToNearestSpecial := func(start int, state []int8, dist []int, pq *Heap[[2]int]) int {
		for i := 0; i < n; i++ {
			dist[i] = math.MaxInt32
		}
		dist[start] = 0
		pq.Clear()
		pq.Push([2]int{start, 0})
		for pq.Len() > 0 {
			top := pq.Pop()
			node, d := top[0], top[1]
			if state[node] == 0 || state[node] == 2 {
				return d
			}
			if dist[node] < d {
				continue
			}
			for _, edge := range adjList[node] {
				to, weight := edge[0], edge[1]
				if newDist := d + weight; newDist < dist[to] {
					dist[to] = newDist
					pq.Push([2]int{to, newDist})
				}
			}
		}
		return math.MaxInt32
	}

	getInitState := func() []int8 {
		remain := k
		notSpecial := make([]int, 0, n-m)
		initState := make([]int8, n) // 0: A, 1: B, 2: special
		for i := 0; i < n; i++ {
			if isSpecial[i] {
				initState[i] = 2
			} else {
				notSpecial = append(notSpecial, i)
			}
		}
		rand.Shuffle(len(notSpecial), func(i, j int) { notSpecial[i], notSpecial[j] = notSpecial[j], notSpecial[i] })
		for _, v := range notSpecial {
			if remain > 0 {
				initState[v] = 0
				remain--
			} else {
				initState[v] = 1
			}
		}
		return initState
	}

	type Arg = []int8
	evaluate := func(arg Arg) float64 {
		res := 0
		dist := make([]int, n)
		pq := NewHeap[[2]int](func(a, b [2]int) bool { return a[1] < b[1] }, nil)
		for i := 0; i < n; i++ {
			if arg[i] == 1 {
				res = max(res, distToNearestSpecial(i, arg, dist, pq))
			}
		}
		return float64(res)
	}

	swapI, swapJ := 0, 0
	next := func(old Arg, _ float64) Arg {
		A, B := make([]int, 0), make([]int, 0)
		for i := 0; i < n; i++ {
			if old[i] == 0 {
				A = append(A, i)
			} else if old[i] == 1 {
				B = append(B, i)
			}
		}
		if len(A) == 0 || len(B) == 0 {
			return old
		}
		swapI, swapJ = A[rand.Intn(len(A))], B[rand.Intn(len(B))]
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

	sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
	sa.SetTimeLimitMs(130)
	sa.SetK(1)
	for i := 0; i < 6; i++ {
		sa.Optimize(getInitState())
	}
	fmt.Fprintln(out, res)
}

// P3878 [TJOI2010] 分金币 (折半枚举)
// 现在有 n 枚金币，它们可能会有不同的价值，第 i 枚金币的价值为 vi。
// 现在要把它们分成两部分，要求这两部分金币数目之差不超过 1，问这样分成的两部分金币的价值之差最小是多少？
func P3878() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(nums []int) int {
		n := len(nums)

		totalSum := 0
		for _, v := range nums {
			totalSum += v
		}
		mid := n / 2

		type Arg = []int
		evaluate := func(arg Arg) float64 {
			sum1 := 0
			for i := 0; i < mid; i++ {
				sum1 += arg[i]
			}
			return abs64(float64(sum1 - (totalSum - sum1)))
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

		sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
		for i := 0; i < 20; i++ {
			rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
			sa.Optimize(nums)
		}
		return int(res)
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var n int32
		fmt.Fscan(in, &n)
		nums := make([]int, n)
		for j := int32(0); j < n; j++ {
			fmt.Fscan(in, &nums[j])
		}
		fmt.Println(solve(nums))
	}
}

// P3936 Coloring (二维，最小化)
// https://www.luogu.com.cn/problem/P3936
// 将一个n*m的网格图用c种不同的颜色染色，规定每种颜色的格子的数量.
// !求一个方案，让相邻格子不同颜色的边的数量F尽量小。
// 输出共n行，每行m个数，表示你构造出的n∗m的F尽量少的染色方案。
//
// 首先按顺序把1−c这c种数全部填进表格里
// 然后每次随机选两个颜色不同的块交换
// n,m<=20,c<=50
func P3936() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row, col, c int16
	fmt.Fscan(in, &row, &col, &c)
	limits := make([]int16, c)
	for i := int16(0); i < c; i++ {
		fmt.Fscan(in, &limits[i])
	}

	all := int(row * col)

	if c == 1 {
		for i := int16(0); i < row; i++ {
			for j := int16(0); j < col; j++ {
				fmt.Fprint(out, 1, " ")
			}
			fmt.Fprintln(out)
		}
		return
	}

	type Arg = []int
	evaluate := func(arg Arg) float64 {
		diff := 0
		for i := int16(0); i < row; i++ {
			for j := int16(0); j < col; j++ {
				cur := arg[i*col+j]
				if i > 0 && arg[(i-1)*col+j] != cur {
					diff++
				}
				if j > 0 && arg[i*col+j-1] != cur {
					diff++
				}
			}
		}
		return float64(diff)
	}

	swapI, swapJ := 0, 0
	next := func(old Arg, _ float64) Arg {
		for {
			swapI, swapJ = rand.Intn(all), rand.Intn(all)
			if old[swapI] != old[swapJ] {
				break
			}
		}
		old[swapI], old[swapJ] = old[swapJ], old[swapI]
		return old
	}

	var bestArg Arg
	bestY := math.MaxFloat64
	summarize := func(nextArg Arg, nextY float64, accept bool) {
		if nextY < bestY {
			bestY = nextY
			bestArg = append(nextArg[:0:0], nextArg...)
		}
		if !accept {
			nextArg[swapI], nextArg[swapJ] = nextArg[swapJ], nextArg[swapI]
		}
	}

	// 把相同颜色的摆一起，能得到一个初始的较优解。
	getInitArg := func() Arg {
		res := make([]int, all)
		ptr := 0
		for color, count := range limits {
			for i := int16(0); i < count; i++ {
				res[ptr] = color + 1
				ptr++
			}
		}
		return res
	}

	sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
	sa.SetTimeLimitMs(995)
	for i := 0; i < 5; i++ {
		sa.SetK(0.4 * float64(i+1))
		initArg := getInitArg()
		sa.Optimize(initArg)
	}
	// fmt.Println(bestY)
	for i := int16(0); i < row; i++ {
		for j := int16(0); j < col; j++ {
			fmt.Fprint(out, bestArg[i*col+j], " ")
		}
		fmt.Fprintln(out)
	}
}

// P5544 [JSOI2016] 炸弹攻击1 (凸优化，最大化)
// https://www.luogu.com.cn/problem/P5544
// 在一个平面内有几个圆以及一些点.
// 使用一个半径不超过R的圆，尽可能多的覆盖平面上的点，并且不与别的圆重合。
// 求最多能覆盖多少个点.
//
// !如果直接把这个杀死敌人个数当作参考的话，这是个整值，
// 导致整个二维函数很不平滑，模拟退火的效果会非常不好（形象理解一下，地图上大量充斥着 0，很可能走不出去）
// !考虑设一个返回值为实数的平滑的函数，能对「即使当前点杀死敌人数量是 0，那么它离 1 有多近」有良好的参考
// !「当前点对应的最大半径还需再增加多少能碰到第一个敌人」是一个好的选择，即为deltaR
// !同时杀死敌人数量又不能不考虑，于是设这样一个估价函数，即为 count
// f(x,y) = c*deltaR - count (c 为一个常数)，最小化这个函数即可.
// 注意让重心作为初始值.
func P5544() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, R int
	fmt.Fscan(in, &n, &m, &R)
	circles := make([][3]float64, n) // (x, y, r)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &circles[i][0], &circles[i][1], &circles[i][2])
	}
	points := make([][2]float64, m) // (x, y)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1])
	}

	sumX, sumY := 0.0, 0.0
	for _, p := range points {
		sumX += p[0]
		sumY += p[1]
	}
	centerX, centerY := sumX/float64(m), sumY/float64(m)

	res := 0

	type Arg = [2]float64
	evaluate := func(arg Arg) float64 {
		x, y := arg[0], arg[1]
		maxRaduis := float64(R)
		for _, c := range circles {
			dist := math.Sqrt((c[0]-x)*(c[0]-x)+(c[1]-y)*(c[1]-y)) - c[2]
			if dist < maxRaduis {
				maxRaduis = dist
			}
		}
		maxRaduis = max64(maxRaduis, 0)

		deltaR, count := math.MaxFloat64, 0.0
		for _, p := range points {
			dist := math.Sqrt((p[0]-x)*(p[0]-x) + (p[1]-y)*(p[1]-y))
			deltaR = min64(deltaR, dist-maxRaduis)
			if dist <= maxRaduis {
				count++
			}
		}
		res = max(res, int(count)) // 将每次的答案都记录下来取 max

		return max64(0, deltaR)*14 - count
	}

	next := func(oldArg Arg, temperature float64) Arg {
		x, y := oldArg[0], oldArg[1]
		nextX := x + (rand.Float64()*2-1)*temperature
		nextY := y + (rand.Float64()*2-1)*temperature
		return Arg{nextX, nextY}
	}

	sa := NewSimulatedAnnealing[Arg](evaluate, next, nil)

	sa.SetK(0.01) // !k很小，表示接受较差的解的概率较小，接近爬山法
	sa.SetTimeLimitMs(49)
	for i := 0; i < 20; i++ { // 跑20轮，每轮49ms
		sa.Optimize(Arg{centerX, centerY})
	}

	fmt.Fprintln(out, res)
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

	evaluate := func(x Arg) float64 {
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

	sa := NewSimulatedAnnealing(evaluate, next, summarize)
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
// 就是费马点.
func getMinDistSum(positions [][]int) float64 {
	type Arg = [2]float64
	evaluate := func(arg Arg) float64 {
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

	sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
	sa.SetK(0.01) // !k很小，表示接受较差的解的概率较小，接近爬山法
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
	evaluate := func(arg Arg) float64 {
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

	sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
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
	evaluate := func(x Arg) float64 {
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

	sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
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
	evaluate := func(x Arg) float64 {
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

	sa := NewSimulatedAnnealing[Arg](evaluate, next, summarize)
	for i := 0; i < 20; i++ {
		rand.Shuffle(len(nums1), func(i, j int) { nums1[i], nums1[j] = nums1[j], nums1[i] })
		sa.Optimize(nums1)
	}
	return int(res)
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) Clear() { h.data = h.data[:0] }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}

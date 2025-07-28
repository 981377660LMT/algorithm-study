package main

import "math"

// 3629. 通过质数传送到达终点的最少跳跃次数
// https://leetcode.cn/problems/minimum-jumps-to-reach-end-via-prime-teleportation/description/
func minJumps(nums []int) int {
	n := len(nums)

	primeToIndexes := make(map[int][]int)
	visited := make([]bool, n)
	for i, num := range nums {
		for p := range E.GetPrimeFactors(num) {
			if _, ok := primeToIndexes[p]; !ok {
				primeToIndexes[p] = []int{}
			}
			primeToIndexes[p] = append(primeToIndexes[p], i)
		}
	}

	setUsed := func(cur int) {
		visited[cur] = true
	}

	findUnused := func(cur int) (next int) {
		if cur-1 >= 0 && !visited[cur-1] {
			return cur - 1
		}
		if cur+1 < n && !visited[cur+1] {
			return cur + 1
		}

		if E.IsPrime(nums[cur]) {
			// !延迟删除
			nexts := primeToIndexes[nums[cur]]
			for len(nexts) > 0 && visited[nexts[len(nexts)-1]] {
				nexts = nexts[:len(nexts)-1]
			}
			primeToIndexes[nums[cur]] = nexts

			if len(nexts) > 0 {
				return nexts[len(nexts)-1]
			}
		}

		return -1
	}

	dist := OnlineBfs(n, 0, setUsed, findUnused)
	return dist[n-1]
}

var E *EratosthenesSieve = NewEratosthenesSieve(1e6 + 10)

// 埃氏筛
type EratosthenesSieve struct {
	minPrime []int
}

func NewEratosthenesSieve(maxN int) *EratosthenesSieve {
	minPrime := make([]int, maxN+1)
	for i := range minPrime {
		minPrime[i] = i
	}
	upper := int(math.Sqrt(float64(maxN))) + 1
	for i := 2; i < upper; i++ {
		if minPrime[i] < i {
			continue
		}
		for j := i * i; j <= maxN; j += i {
			if minPrime[j] == j {
				minPrime[j] = i
			}
		}
	}
	return &EratosthenesSieve{minPrime}
}

func (es *EratosthenesSieve) IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	return es.minPrime[n] == n
}

func (es *EratosthenesSieve) GetPrimeFactors(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		m := es.minPrime[n]
		res[m]++
		n /= m
	}
	return res
}

func (es *EratosthenesSieve) GetPrimes() []int {
	res := []int{}
	for i, x := range es.minPrime {
		if i >= 2 && i == x {
			res = append(res, x)
		}
	}
	return res
}

const INF int = 1e18

// 在线bfs.
//
//	不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
//	setUsed(u)：将 u 标记为已访问。
//	findUnused(u)：找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `-1`。
//
// https://leetcode.cn/problems/minimum-reverse-operations/solution/python-zai-xian-bfs-jie-jue-bian-shu-hen-y58m/
func OnlineBfs(
	n int, start int,
	setUsed func(u int), findUnused func(cur int) (next int),
) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
	setUsed(start)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for {
			next := findUnused(cur)
			if next == -1 {
				break
			}
			dist[next] = dist[cur] + 1 // weight
			queue = append(queue, next)
			setUsed(next)
		}
	}

	return
}

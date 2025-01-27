package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	abc390_d()
}

func demo() {
	t1 := time.Now()
	ptr := 0
	EnumerateSetPartition(12 /* n */, func(groups [][]int) {
		ptr++
	})
	fmt.Println(time.Since(t1))

	t2 := time.Now()
	ptr = 0
	EnumerateSetPartitionSum(
		12,
		func(presum int, mask int) int {
			return presum + mask
		}, 0,
		func(sum int) {
			ptr++
		})
	fmt.Println(time.Since(t2))
}

func abc390_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	subsum := make([]int, 1<<n)
	for i := 0; i < n; i++ {
		for s := 0; s < 1<<i; s++ {
			subsum[s|1<<i] = subsum[s] + nums[i]
		}
	}

	set := make(map[int]struct{})
	reducer := func(presum int, groupMask int) int {
		return presum ^ subsum[groupMask]
	}
	initial := 0
	consumer := func(totalSum int) {
		set[totalSum] = struct{}{}
	}

	EnumerateSetPartitionSum(n, reducer, initial, consumer)

	fmt.Fprintln(out, len(set))
}

// 遍历所有的子集分割方式(贝尔数).
// n = 12 => 4213597, 50ms.
func EnumerateSetPartition(n int, consumer func(groups [][]int)) {
	groups := [][]int{} // 或者用一个 roots 数组表示集合的根节点（代表元）
	var f func(int)
	f = func(p int) {
		if p == n {
			consumer(groups)
			return
		}

		groups = append(groups, []int{p})
		f(p + 1)
		groups = groups[:len(groups)-1]
		for i := range groups {
			groups[i] = append(groups[i], p)
			f(p + 1)
			groups[i] = groups[i][:len(groups[i])-1]
		}
	}

	f(0)
}

// 遍历所有的子集分割方式，同时计算划分的和.
func EnumerateSetPartitionSum[S any](
	n int,
	reducer func(presum S, groupMask int) S, initial S,
	consumer func(totalSum S),
) {
	var dfs func(int, S)
	dfs = func(remain int, sum S) {
		if remain == 0 {
			consumer(sum)
			return
		}

		lb := remain & -remain
		next := remain ^ lb
		for s := next; ; s = (s - 1) & next {
			subset := s | lb
			dfs(remain^subset, reducer(sum, subset))
			if s == 0 {
				break
			}
		}
	}

	dfs((1<<n)-1, initial)
}

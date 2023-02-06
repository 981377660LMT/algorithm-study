// https://www.luogu.com.cn/problem/P5569
// GarsiaWachs 算法 nlogn
// 现要将石子有次序地合并成一堆。规定每次只能选`相邻`的 2堆石子合并成新的一堆，
// 并将新的一堆石子数记为该次合并的得分。
package main

import (
	"bufio"
	"fmt"
	"os"
)

func mergeStone(costs []int) int {
	res := 0
	t := 1
	var combine func(k int)
	combine = func(k int) {
		tmp := costs[k] + costs[k-1]
		res += tmp
		for i := k; i < t-1; i++ {
			costs[i] = costs[i+1]
		}
		t--
		j := 0
		for j = k - 1; j > 0 && costs[j-1] < tmp; j-- {
			costs[j] = costs[j-1]
		}
		costs[j] = tmp
		for j >= 2 && costs[j] >= costs[j-2] {
			d := t - j
			combine(j - 1)
			j = t - d
		}
	}

	for i := 1; i < len(costs); i++ {
		costs[t] = costs[i]
		t++
		for t >= 3 && costs[t-1] >= costs[t-3] {
			combine(t - 2)
		}
	}

	for t > 1 {
		combine(t - 1)
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	fmt.Fprintln(out, mergeStone(nums))
}

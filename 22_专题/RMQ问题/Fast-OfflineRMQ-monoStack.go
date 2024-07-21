// 离线RMQ (区间最值查询)
// 单调栈+二分(比并查集稍慢一些)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	demo()
}

func yosupo() {
	// https://judge.yosupo.jp/problem/staticrmq
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	queries := make([][2]int32, q)
	for i := int32(0); i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
	}
	minIndexes := OfflineRmqMonoStack(queries, func(i, j int32) bool {
		return nums[i] < nums[j]
	})

	for i := int32(0); i < q; i++ {
		fmt.Fprintln(out, nums[minIndexes[i]])
	}
}

func demo() {
	type pair struct{ left, right int32 }
	nums := []pair{{1, 10}, {2, 5}, {1, 5}, {3, 7}, {2, 6}}
	queries := [][2]int32{{0, 5}, {1, 4}, {2, 3}, {3, 5}, {1, 3}}
	minIndexes := OfflineRmqMonoStack(queries, func(i, j int32) bool {
		if nums[i].left == nums[j].left {
			return nums[i].right < nums[j].right
		}
		return nums[i].left < nums[j].left
	})
	for _, i := range minIndexes {
		fmt.Println(nums[i])
	}

}

// 每个查询是左闭右开区间[start, end).返回每个查询的最小值的下标.
//
//	0<=start<end<=n.
func OfflineRmqMonoStack(queries [][2]int32, less func(i, j int32) bool) []int32 {
	n := int32(0)
	for _, query := range queries {
		n = max32(n, query[1])
	}

	// 离线询问
	type pair struct{ left, qid int32 }
	qs := make([][]pair, n)
	for i, q := range queries {
		left, right := q[0], q[1]-1
		qs[right] = append(qs[right], pair{left, int32(i)})
	}

	res := make([]int32, len(queries))
	minStack := []int32{}
	for right := int32(0); right < n; right++ {
		for len(minStack) > 0 && less(right, minStack[len(minStack)-1]) {
			minStack = minStack[:len(minStack)-1]
		}
		minStack = append(minStack, right)
		for _, p := range qs[right] {
			index := sort.Search(len(minStack), func(i int) bool { return minStack[i] >= p.left })
			res[p.qid] = minStack[index]
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
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

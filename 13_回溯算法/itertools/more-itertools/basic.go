package main

import "fmt"

func main() {
	Batched(10, 3, func(start, end int) { fmt.Println(start, end) })    // [0, 3), [3, 6), [6, 9), [9, 10)
	Distribute(10, 3, func(start, end int) { fmt.Println(start, end) }) // [0, 4), [4, 7), [7, 10)
}

// 将 n 个元素分成大小为 size 的批次，最后一个批次可能小于 size.
// 适用于分片处理大量数据的场景.
// 返回两个组的大小和组的个数.
func Batched(n, size int, f func(start, end int)) (size1, count1, size2, count2 int) {
	if f != nil {
		for i := 0; i < n; i += size {
			f(i, min(i+size, n))
		}
	}
	size1 = size
	count1 = n / size
	size2 = n % size
	count2 = 0
	if size2 > 0 {
		count2 = 1
	}
	return
}

// 将 n 个元素分成 groupCount 个组.每个组的大小尽可能均等分配，使得每个组的大小差距不超过1.
// 适用于将任务分配给多个工作线程的场景.
// 返回两个组的大小和组的个数.
func Distribute(n, groupCount int, f func(start, end int)) (size1, count1, size2, count2 int) {
	q := n / groupCount
	r := n % groupCount
	if f != nil {
		start := 0
		for i := 0; i < groupCount; i++ {
			end := start + q
			if i < r {
				end++
			}
			f(start, end)
			start = end
		}
	}
	size1 = q + 1
	count1 = r
	size2 = q
	count2 = groupCount - r
	return
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

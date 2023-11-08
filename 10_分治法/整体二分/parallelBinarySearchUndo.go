package main

import (
	"fmt"
)

func main() {
	demo()
}

func demo() {
	// n个操作，第i个操作将curSum增加i+1.
	// q个查询，第i个查询形如：curSum是否大于等于i+1.
	// 对于每个查询，输出第一个满足条件的操作的编号.

	curSum := 0
	history := make([]int, 0, 10)
	res := ParallelBinarySearchUndo(
		10, 10,
		func(mutationId int) {
			history = append(history, curSum)
			curSum += mutationId + 1
		},
		func() {
			curSum = history[len(history)-1]
			history = history[:len(history)-1]
		},
		func(queryId int) bool {
			return curSum >= 56
		},
	)

	fmt.Println(res)
}

// !如果undo操作(nlong次)比reset(logn次)更优，则可以使用这个版本.
//
// 整体二分解决这样一类问题:
//   - 给定一个长度为n的操作序列, 按顺序执行这些操作;
//   - 给定q个查询,每个查询形如:"条件qi为真(满足条件)是在第几次操作之后?".
//     !要求对条件为真的判定具有单调性，即某个操作后qi为真,后续操作都会满足qi为真.
//
// 返回:
//   - -1 => 不需要操作就满足条件的查询.
//   - [0, n) => 满足条件的最早的操作的编号(0-based).
//   - n => 执行完所有操作后都不满足条件的查询.
//
// https://betrue12.hateblo.jp/entry/2019/08/14/152227
func ParallelBinarySearchUndo(
	n, q int,
	mutate func(mutationId int), // 执行第 mutationId 次操作，一共调用 nlogn 次.
	undo func(), // 撤销上一次`mutate`操作，一共调用 nlogn 次.
	predicate func(queryId int) bool // 判断第 queryId 次查询是否满足条件，一共调用 qlogn 次.
) []int {
	left, right := make([]int, q), make([]int, q)
	for i := 0; i < q; i++ {
		left[i], right[i] = 0, n
	}

	// 不需要操作就满足条件的查询
	for i := 0; i < q; i++ {
		if predicate(i) {
			right[i] = -1
		}
	}

	for {
		mids := make([]int, q)
		for i := range mids {
			mids[i] = -1
		}
		for i := 0; i < q; i++ {
			if left[i] <= right[i] {
				mids[i] = (left[i] + right[i]) >> 1
			}
		}

		// csr 数组保存二元对 (qi,mid).
		indeg := make([]int, n+2)
		for i := 0; i < q; i++ {
			mid := mids[i]
			if mid != -1 {
				indeg[mid+1]++
			}
		}
		for i := 0; i < n+1; i++ {
			indeg[i+1] += indeg[i]
		}
		total := indeg[n+1]
		if total == 0 {
			break
		}
		counter := append(indeg[:0:0], indeg...)
		csr := make([]int, total)
		for i := 0; i < q; i++ {
			mid := mids[i]
			if mid != -1 {
				csr[counter[mid]] = i
				counter[mid]++
			}
		}

		times := 0
		for _, pos := range csr {
			for times < mids[pos] {
				mutate(times)
				times++
			}
			if predicate(pos) {
				right[pos] = times - 1
			} else {
				left[pos] = times + 1
			}
		}
		for i := 0; i < times; i++ {
			undo()
		}
	}

	return right
}

// 二维滑动窗口最小值

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(SlidingWindowMinimum2D([][]int{
		{1, 2, 3, 4, 5},
		{2, 3, 4, 5, 6},
		{3, 4, 5, 6, 7},
		{4, 5, 6, 7, 8},
		{5, 6, 7, 8, 9},
	}, 3, 3, true))

	fmt.Println(QuerySlidingWindowMinimum2D([][]int{
		{1, 2, 3, 4, 5},
		{2, 3, 4, 5, 6},
		{3, 4, 5, 6, 7},
		{4, 5, 6, 7, 8},
		{5, 6, 7, 8, 9},
	}, 5, 5, 3, 3, true))

}

// 二维滑动窗口最小值.
// 求出每个rowSize*colSize窗口(和)的最值.
func SlidingWindowMinimum2D(grid [][]int, windowRow, windowCol int, isMax bool) [][]int {
	ROW, COL := len(grid), len(grid[0])
	cols := make([]*MonoQueue, COL)
	for c := 0; c < COL; c++ {
		if isMax {
			cols[c] = NewMonoQueue(func(a, b int) bool { return a > b })
		} else {
			cols[c] = NewMonoQueue(func(a, b int) bool { return a < b })
		}
	}
	for r := 0; r < windowRow; r++ {
		for c := 0; c < COL; c++ {
			cols[c].Append(grid[r][c])
		}
	}

	res := [][]int{}
	for r := 0; r < ROW-windowRow+1; r++ {
		res = append(res, []int{})
		// 维护这一行的滑动窗口最值.
		var window *MonoQueue
		if isMax {
			window = NewMonoQueue(func(a, b int) bool { return a > b })
		} else {
			window = NewMonoQueue(func(a, b int) bool { return a < b })
		}
		for c := 0; c < COL; c++ {
			// 每列的最值.
			window.Append(cols[c].Head())
			if c >= windowCol-1 {
				res[r] = append(res[r], window.Head())
				window.Popleft()
			}
		}

		// 下一行进入窗口.
		if r+windowRow < ROW {
			for c := 0; c < COL; c++ {
				cols[c].Append(grid[r+windowRow][c])
				cols[c].Popleft()
			}
		}
	}

	return res
}

// 在每个bigRow*bigCol窗口范围内所有的smallRow*smallCol窗口和的最值.
func QuerySlidingWindowMinimum2D(grid [][]int, bigRow, bigCol, smallRow, smallCol int, isMax bool) [][]int {
	ROW, COL := len(grid), len(grid[0])
	preSum := make([][]int, ROW+1)
	for r := range preSum {
		preSum[r] = make([]int, COL+1)
	}
	for r := 0; r < ROW; r++ {
		for c := 0; c < COL; c++ {
			preSum[r+1][c+1] = preSum[r][c+1] + preSum[r+1][c] - preSum[r][c] + grid[r][c]
		}
	}

	windowSum := make([][]int, ROW-smallRow+1)
	for r := 0; r < ROW-smallRow+1; r++ {
		for c := 0; c < COL-smallCol+1; c++ {
			windowSum[r] = append(windowSum[r],
				preSum[r+smallRow][c+smallCol]-
					preSum[r][c+smallCol]-
					preSum[r+smallRow][c]+
					preSum[r][c])
		}
	}
	windowMax := SlidingWindowMinimum2D(windowSum, bigRow-smallRow+1, bigCol-smallCol+1, isMax)
	return windowMax
}

type V = int

// 单调队列维护滑动窗口最小值.
// 单调队列队头元素为当前窗口最小值，队尾元素为当前窗口最大值.
type MonoQueue struct {
	MinQueue       []V
	_minQueueCount []int
	_less          func(a, b V) bool
	_len           int
}

func NewMonoQueue(less func(a, b V) bool) *MonoQueue {
	return &MonoQueue{
		_less: less,
	}
}

func (q *MonoQueue) Append(value V) *MonoQueue {
	count := 1
	for len(q.MinQueue) > 0 && q._less(value, q.MinQueue[len(q.MinQueue)-1]) {
		q.MinQueue = q.MinQueue[:len(q.MinQueue)-1]
		count += q._minQueueCount[len(q._minQueueCount)-1]
		q._minQueueCount = q._minQueueCount[:len(q._minQueueCount)-1]
	}
	q.MinQueue = append(q.MinQueue, value)
	q._minQueueCount = append(q._minQueueCount, count)
	q._len++
	return q
}

func (q *MonoQueue) Popleft() {
	q._minQueueCount[0]--
	if q._minQueueCount[0] == 0 {
		q.MinQueue = q.MinQueue[1:]
		q._minQueueCount = q._minQueueCount[1:]
	}
	q._len--
}

func (q *MonoQueue) Head() V {
	return q.MinQueue[0]
}

func (q *MonoQueue) Min() V {
	return q.MinQueue[0]
}

func (q *MonoQueue) Len() int {
	return q._len
}

func (q *MonoQueue) String() string {
	sb := []string{}
	for i := 0; i < len(q.MinQueue); i++ {
		sb = append(sb, fmt.Sprintf("%v", pair{q.MinQueue[i], q._minQueueCount[i]}))
	}
	return fmt.Sprintf("MonoQueue{%v}", strings.Join(sb, ", "))
}

type pair struct {
	value V
	count int
}

func (p pair) String() string {
	return fmt.Sprintf("(value: %v, count: %v)", p.value, p.count)
}

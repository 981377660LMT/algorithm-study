package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	positions := make([]int, m)
	jumps := make([]int, m)
	for i := range positions {
		var p, j int
		fmt.Fscan(in, &p, &j)
		positions[i] = p
		jumps[i] = j
	}

	fmt.Fprintln(out, 雅加达的摩天楼(n, m, positions, jumps))
}

// https://www.luogu.com.cn/problem/P3645
// P3645 [APIO2015] 雅加达的摩天楼
//
// 有m只doge分布在n个摩天大楼上。楼和doge都是从0开始编号。
// 每只doge初始位置positions[i]，弹跳力jumps[i]。
// 它每一次跳会恰好跳jumps[i]个大楼。比如从x可以到 x±jumps[i]。
// 现在，0号doge 要把某信息传给1号doge。
// 对于一只doge，若它尚未知道信息，就不能动。
// 对于一只doge，若它已经知道信息，可以选择把信息告诉处于同一位置的doge们，或者跳去别的位置。
// !求最少跳的步数。如果无法传递，输出-1。
// n,m<=3e4
//
// !带两种权的最短路考虑01bfs.
// !注意到bfs的状态由(pos,jump)两个维度决定,即在第pos个位置,当前doge的弹跳力为jump.
// jump <= sqrt(n)时，只有O(n*sqrt(n))个状态.
// jump > sqrt(n)时，只有O(m*sqrt(n))个状态(最多m只doge，每只doge只有n/jump个可行位置).
// 由于状态数较多，采用bitset来存储状态.
func 雅加达的摩天楼(n int, m int, positions []int, jumps []int) int {
	target := positions[1]
	groups := make([][]int, n)
	for i, p := range positions {
		groups[p] = append(groups[p], jumps[i])
	}

	maxJump := 0
	for _, j := range jumps {
		maxJump = max(maxJump, j)
	}

	visited := make([]_BS, n)
	for i := range visited {
		visited[i] = _NewBS(maxJump + 1)
	}
	queue := &Deque{} // [pos, jump, dist]
	queue.Append([3]int{positions[0], jumps[0], 0})
	visited[positions[0]].Set(jumps[0])

	for queue.Size() > 0 {
		item := queue.PopLeft()
		curPos, curJump, curDist := item[0], item[1], item[2]
		if curPos == target {
			return curDist
		}

		// doge 接力
		for _, nextJump := range groups[curPos] {
			if !visited[curPos].Has(nextJump) {
				visited[curPos].Set(nextJump)
				queue.AppendLeft([3]int{curPos, nextJump, curDist})
			}
		}

		// doge 不接力
		if curPos-curJump >= 0 && !visited[curPos-curJump].Has(curJump) {
			visited[curPos-curJump].Set(curJump)
			queue.Append([3]int{curPos - curJump, curJump, curDist + 1})
		}

		if curPos+curJump < n && !visited[curPos+curJump].Has(curJump) {
			visited[curPos+curJump].Set(curJump)
			queue.Append([3]int{curPos + curJump, curJump, curDist + 1})
		}
	}

	return -1
}

type _BS []uint64

func _NewBS(n int) _BS       { return make(_BS, n>>6+1) }
func (b _BS) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b _BS) Set(p int)      { b[p>>6] |= 1 << (p & 63) }

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

type D = [3]int
type Deque struct{ l, r []D }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

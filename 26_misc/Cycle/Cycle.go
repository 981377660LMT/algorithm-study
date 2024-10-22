// 与环相关的工具函数.

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Cycle struct {
	n int32
}

func NewCycle(n int32) *Cycle {
	return &Cycle{n: n}
}

func (c *Cycle) Dist(u, v int32) int32 {
	d := abs32(u - v)
	return min32(d, c.n-d)
}

func (c *Cycle) Path(u, v int32) []int32 {
	if c.DistLeft(u, v) <= c.DistRight(u, v) {
		return c.PathLeft(u, v)
	}
	return c.PathRight(u, v)
}

func (c *Cycle) Segment(u, v int32, consumer func(from, to int32)) {
	if c.DistLeft(u, v) <= c.DistRight(u, v) {
		c.SegmentLeft(u, v, consumer)
	} else {
		c.SegmentRight(u, v, consumer)
	}
}

func (c *Cycle) SegmentLeft(from, to int32, consumer func(from, to int32)) {
	if from >= to {
		consumer(from, to)
		return
	}
	consumer(from, 0)
	consumer(c.n-1, to)
}

func (c *Cycle) SegmentRight(from, to int32, consumer func(from, to int32)) {
	if to >= from {
		consumer(from, to)
		return
	}
	consumer(from, c.n-1)
	consumer(0, to)
}

func (c *Cycle) PathLeft(from, to int32) []int32 {
	if from >= to {
		return c.makeRange(from, to-1, -1)
	}
	return append(c.makeRange(from, -1, -1), c.makeRange(c.n-1, to-1, -1)...)
}

func (c *Cycle) PathRight(from, to int32) []int32 {
	if to >= from {
		return c.makeRange(from, to+1, 1)
	}
	return append(c.makeRange(from, c.n, 1), c.makeRange(0, to+1, 1)...)
}

func (c *Cycle) DistLeft(from, to int32) int32 {
	if from >= to {
		return from - to
	}
	return from + c.n - to
}

func (c *Cycle) DistRight(from, to int32) int32 {
	if to >= from {
		return to - from
	}
	return to + c.n - from
}

// x 是否在 from 到 to 之间的逆时针路径上.
func (c *Cycle) OnPathLeft(from, to, x int32) bool {
	if x < to {
		x += c.n
	}
	if from < to {
		from += c.n
	}
	return to <= x && x <= from
}

// x 是否在 from 到 to 之间的顺时针路径上.
func (c *Cycle) OnPathRight(from, to, x int32) bool {
	if from > to {
		to += c.n
	}
	if from > x {
		x += c.n
	}
	return from <= x && x <= to
}

func (c *Cycle) JumpLeft(from, steps int32) int32 {
	res := (from - steps) % c.n
	if res < 0 {
		res += c.n
	}
	return res
}

func (c *Cycle) JumpRight(from, steps int32) int32 {
	res := (from + steps) % c.n
	if res < 0 {
		res += c.n
	}
	return res
}

func (c *Cycle) makeRange(from, to, step int32) (res []int32) {
	if step > 0 {
		for i := from; i < to; i += step {
			res = append(res, i)
		}
	} else {
		for i := from; i > to; i += step {
			res = append(res, i)
		}
	}
	return
}

func abs32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func main() {
	abc376f()
}

// F - Hands on Ring (Hard) (abc376 F)
//
// n的环形格子。两个棋子，初始位于0,1。
// 给定q个指令，每个指令指定一个棋子移动到某个格子上，期间可以移动另外一个棋子。
// 依次执行这些指令，问最小的移动次数。
// n,q<=3000.
// !dp[index][posL][posR]表示执行到第index个指令，左边的棋子在posL，右边的棋子在posR时的最小移动次数。
func abc376f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF32 int32 = 1e9 + 10

	var N, Q int32
	fmt.Fscan(in, &N, &Q)
	H, T := make([]byte, Q), make([]int32, Q)
	for i := int32(0); i < Q; i++ {
		var s string
		var t int32
		fmt.Fscan(in, &s, &t)
		H[i] = s[0]
		T[i] = int32(t - 1)
	}

	C := NewCycle(N)

	// Returns (distance, new position of the other point).
	moveLeft := func(cur, to, other int32) (int32, int32) {
		if !C.OnPathLeft(cur, to, other) {
			return C.DistLeft(cur, to), other
		}
		otherTo := to - 1
		if to == 0 {
			otherTo = N - 1
		}
		return C.DistLeft(cur, to) + C.DistLeft(other, otherTo), otherTo
	}

	// Returns (distance, new position of the other point).
	moveRight := func(cur, to, other int32) (int32, int32) {
		if !C.OnPathRight(cur, to, other) {
			return C.DistRight(cur, to), other
		}
		otherTo := to + 1
		if to == N-1 {
			otherTo = 0
		}
		return C.DistRight(cur, to) + C.DistRight(other, otherTo), otherTo
	}

	memo := make([]map[int32]int32, Q)
	for i := range memo {
		memo[i] = make(map[int32]int32)
	}

	var dp func(index, posL, posR int32) int32
	dp = func(index, posL, posR int32) int32 {
		if index == Q {
			return 0
		}
		hash := posL<<12 | posR
		if res, ok := memo[index][hash]; ok {
			return res
		}
		to := T[index]
		res := INF32
		if H[index] == 'L' {
			d1, r1 := moveLeft(posL, to, posR)
			res = min32(res, d1+dp(index+1, to, r1))
			d2, r2 := moveRight(posL, to, posR)
			res = min32(res, d2+dp(index+1, to, r2))
		} else {
			d1, l1 := moveLeft(posR, to, posL)
			res = min32(res, d1+dp(index+1, l1, to))
			d2, l2 := moveRight(posR, to, posL)
			res = min32(res, d2+dp(index+1, l2, to))
		}
		memo[index][hash] = res
		return res
	}

	res := dp(0, 0, 1)
	fmt.Fprintln(out, res)
}

// https://ei1333.github.io/library/other/connected-grid-states.hpp
// 二维网格连通性dp问题,维护`并查集的状态`(每个状态编号).
// 什么都不存在的状态表示为0.
// !墙壁 => '#', 通路 => '.'.
//
// !列数(1<=col<=11): dp的状态数(并查集的状态数)
// 1: 3
// 2: 11
// 3: 34
// 4: 102
// 5: 306
// 6: 914
// 7: 2722
// 8: 8080
// 9: 23929
// 10: 70747
// 11: 208944
// 12: 616706
// 13: 1819663

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

// Ex - Unite
// https://atcoder.jp/contests/abc296/tasks/abc296_h
// 给定一个R行C列的网格,每个格子有'#'和'.'两种状态(墙壁/通路)
// 现在需要将网格中的某些格子变成墙壁,使得网格中的所有墙壁都连通.
// !问最少需要变成墙壁的格子数.
// R<=100,C<=7
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)

	grid := make([][]byte, ROW)
	for i := 0; i < ROW; i++ {
		grid[i] = make([]byte, COL)
		fmt.Fscan(in, &grid[i])
	}

	S := NewConnectedGridStates(COL, func(state []int) uint64 {
		hash := uint64(0)
		for _, num := range state {
			hash = hash*13331 + uint64(num+2)
		}
		return hash
	})

	dp := []int{0}
	for r := 0; r < ROW; r++ {
		for c := 0; c < COL; c++ {
			ndp := make([]int, len(S.GetStates(c+1)))
			for i := range ndp {
				ndp[i] = INF
			}

			for s, preRes := range dp {
				{ // '#'
					to := S.SetWall(c, s)
					if to != -1 {
						opt := 0
						if grid[r][c] == '.' {
							opt = 1
						}
						ndp[to] = min(ndp[to], preRes+opt)
					}
				}

				// '.'
				if grid[r][c] != '#' {
					to := S.SetGround(c, s)
					if to != -1 {
						ndp[to] = min(ndp[to], preRes)
					}
				}
			}

			dp = ndp
		}
	}

	res := INF
	for s := range dp {
		ok := true
		for _, j := range S.states[0][s] {
			if j >= 1 { // 有连通的组
				ok = false
			}
		}
		if ok {
			res = min(res, dp[s])
		}
	}

	fmt.Fprintln(out, res)
}

type ConnectedGridStates struct {
	states               [][][]int //! -2:消失,-1:不连通,0:初始都不连通的状态, >0:并查集的状态编号(存在连通的组)
	nextWall, nextGround [][]int
	width                int
}

// 预处理列数为col的二维网格的连接状态转移关系.
//  1<=col<=11.
func NewConnectedGridStates(col int, hash func(state []int) uint64) *ConnectedGridStates {
	if col < 1 {
		panic("invalid col")
	}

	states := make([][][]int, col)   // 当前在第col列,并查集状态为s,当前状态
	nextWall := make([][]int, col)   // 当前在第col列,并查集状态为s,置为墙壁后的并查集状态
	nextGround := make([][]int, col) // 当前在第col列,并查集状态为s,置为通路后的并查集状态
	width := col
	res := &ConnectedGridStates{
		states:     states,
		nextWall:   nextWall,
		nextGround: nextGround,
		width:      width,
	}

	stateId := make([]map[uint64]int, width) // stateId[col][hash(state)] => id
	for i := 0; i < width; i++ {
		stateId[i] = map[uint64]int{}
	}
	state := make([]int, width) // 并查集状态
	for i := 0; i < width; i++ {
		state[i] = -1
	}
	states[0] = append(states[0], state)
	stateId[0][hash(state)] = 0 // 初始状态的并查集对应编号为0
	nextWall[0] = append(nextWall[0], -1)
	nextGround[0] = append(nextGround[0], -1)

	type item struct {
		stateId int
		state   []int
		curCol  int
	}
	queue := []item{}
	queue = append(queue, item{0, state, 0})

	modify := func(state []int) []int {
		id := make([]int, width)
		for i := 0; i < width; i++ {
			id[i] = -1
		}
		now := 0
		for i := 0; i < width; i++ {
			if preS := state[i]; preS != -1 {
				if id[preS] == -1 {
					id[preS] = now
					now++
				}
				state[i] = id[preS]
			}
		}
		return state
	}

	// 将新的状态压入队列,返回新状态编号.
	push := func(nextCol int, state []int) int {
		hash_ := hash(state)
		if id, ok := stateId[nextCol][hash_]; ok {
			return id
		}
		newId := len(states[nextCol])
		stateId[nextCol][hash_] = newId
		states[nextCol] = append(states[nextCol], state)
		nextWall[nextCol] = append(nextWall[nextCol], -1)
		nextGround[nextCol] = append(nextGround[nextCol], -1)
		queue = append(queue, item{newId, state, nextCol})
		return newId
	}

	// bfs代替dp.
	for len(queue) > 0 {
		sid, state, curCol := queue[0].stateId, queue[0].state, queue[0].curCol
		queue = queue[1:]
		nextCol := 0
		if curCol+1 != width {
			nextCol = curCol + 1
		}
		max_ := state[0]
		for k := 1; k < width; k++ {
			if max_ < state[k] {
				max_ = state[k]
			}
		}

		{
			// '.', 将当前格子置为通路
			if state[curCol] == -1 {
				nextGround[curCol][sid] = push(nextCol, state)
			} else {
				ok := false
				for k := 0; k < width; k++ {
					if curCol != k && state[k] == state[curCol] {
						ok = true
					}
				}
				if ok {
					to := make([]int, width)
					copy(to, state)
					to[curCol] = -1
					nextGround[curCol][sid] = push(nextCol, modify(to))
				} else if max_ >= 1 { // disconnected,当前列与其他列不连通
					nextGround[curCol][sid] = -1
				} else { // disappeared,该状态已经消失
					nextGround[curCol][sid] = -2
				}
			}
		}

		{
			// '#', 将当前格子置为墙壁
			to := make([]int, width)
			copy(to, state)
			if curCol == 0 {
				if state[curCol] == -1 {
					to[curCol] = max_ + 1
				}
			} else {
				if state[curCol] == -1 {
					if state[curCol-1] == -1 {
						to[curCol] = max_ + 1
					} else {
						to[curCol] = state[curCol-1]
					}
				} else if state[curCol-1] != -1 {
					for k := 0; k < width; k++ {
						if state[k] == state[curCol] {
							to[k] = state[curCol-1]
						}
					}
				}
			}
			nextWall[curCol][sid] = push(nextCol, modify(to))
		}
	}

	for c := 0; c < width; c++ {
		nc := 0
		if c+1 != width {
			nc = c + 1
		}
		len_ := len(states[nc])
		for k := 0; k < len(nextGround[c]); k++ {
			if nextGround[c][k] == -2 {
				nextGround[c][k] = len_
			}
		}
		for k := 0; k < len(nextWall[c]); k++ {
			if nextWall[c][k] == -2 {
				nextWall[c][k] = len_
			}
		}
		nextGround[c] = append(nextGround[c], len_)
		nextWall[c] = append(nextWall[c], -1)
	}

	for i := 0; i < width; i++ {
		cur := make([]int, width)
		for j := 0; j < width; j++ {
			cur[j] = -1
		}
		states[i] = append(states[i], cur)
	}

	return res
}

// 当前在第c列,状态为state时,将该列的格子置为墙壁后的(转移)状态.
//  0<=state<len(States[c]).
//  0<=c<col.
//  O(1).
func (cgs *ConnectedGridStates) SetWall(c, state int) int {
	return cgs.nextWall[c][state]
}

// 当前在第c列,状态为state时,将该列的格子置为通路后的(转移)状态.
//  0<=state<len(States[c]).
//  0<=c<col.
//  O(1).
func (cgs *ConnectedGridStates) SetGround(c, state int) int {
	return cgs.nextGround[c][state]
}

// 返回第c列的所有状态.States[c][state]表示第c列的状态为state时的并查集状态.
//  0<=c<2*col.当c大于col时,返回第c-col列的所有状态.
//  O(1).
func (cgs *ConnectedGridStates) GetStates(c int) [][]int {
	if c >= cgs.width {
		c -= cgs.width
	}
	return cgs.states[c]
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

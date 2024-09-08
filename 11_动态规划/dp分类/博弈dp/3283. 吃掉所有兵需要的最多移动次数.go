// 3283. 吃掉所有兵需要的最多移动次数
// https://leetcode.cn/problems/maximum-number-of-moves-to-kill-all-pawns/description/
// 给你一个 50 x 50 的国际象棋棋盘，棋盘上有 一个 马和一些兵。
// 给你两个整数 kx 和 ky ，其中 (kx, ky) 表示马所在的位置，同时还有一个二维数组 positions ，其中 positions[i] = [xi, yi] 表示第 i 个兵在棋盘上的位置。
// Alice 和 Bob 玩一个回合制游戏，Alice 先手。玩家的一次操作中，可以执行以下操作：
// 玩家选择一个仍然在棋盘上的兵，然后移动马，通过 最少 的 步数 吃掉这个兵。
// 注意 ，玩家可以选择 任意 一个兵，不一定 要选择从马的位置出发 最少 移动步数的兵。
// 在马吃兵的过程中，马 可能 会经过一些其他兵的位置，但这些兵 不会 被吃掉。只有 选中的兵在这个回合中被吃掉。
// Alice 的目标是 最大化 两名玩家的 总 移动次数，直到棋盘上不再存在兵，而 Bob 的目标是 最小化 总移动次数。
// 假设两名玩家都采用 最优 策略，请你返回 Alice 可以达到的 最大 总移动次数。
// 在一次 移动 中，如下图所示，马有 8 个可以移动到的位置，每个移动位置都是沿着坐标轴的一个方向前进 2 格，然后沿着垂直的方向前进 1 格。

package main

const INF32 int32 = 1e9

var DIR8 = [][]int32{{-2, -1}, {-2, 1}, {2, -1}, {2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}}

func maxMoves(kx int, ky int, positions [][]int) int {
	ROW, COL := int32(50), int32(50)

	// 预处理兵+初始的马到各个格子的最短距离.
	bfs := func(sx, sy int32) [][]int32 {
		dist := make([][]int32, ROW)
		for i := range dist {
			dist[i] = make([]int32, COL)
			for j := range dist[i] {
				dist[i][j] = -1
			}
		}
		dist[sx][sy] = 0
		queue := [][2]int32{{sx, sy}}
		for len(queue) > 0 {
			x, y := queue[0][0], queue[0][1]
			queue = queue[1:]
			for _, d := range DIR8 {
				nx, ny := x+d[0], y+d[1]
				if nx >= 0 && nx < ROW && ny >= 0 && ny < COL && dist[nx][ny] == -1 {
					dist[nx][ny] = dist[x][y] + 1
					queue = append(queue, [2]int32{nx, ny})
				}
			}
		}
		return dist
	}
	dists := make([][][]int32, len(positions)+1)
	for i, pos := range positions {
		dists[i] = bfs(int32(pos[0]), int32(pos[1]))
	}
	dists[len(positions)] = bfs(int32(kx), int32(ky))
	dist := func(fromIndex, toIndex int32) int32 {
		tx, ty := int32(positions[toIndex][0]), int32(positions[toIndex][1])
		return dists[fromIndex][tx][ty]
	}

	m := int32(len(positions))
	mask := int32((1 << m) - 1)
	memo := make([]int32, (m+1)*mask*2)
	for i := range memo {
		memo[i] = -1
	}

	var dfs func(index, state, player int32) int32
	dfs = func(index, state, player int32) int32 {
		if state == mask {
			return 0
		}
		hash := index*mask*2 + state*2 + player
		if memo[hash] != -1 {
			return memo[hash]
		}
		res, op := -INF32, max32
		if player == 1 {
			res, op = INF32, min32
		}
		for i := int32(0); i < m; i++ {
			if state&(1<<i) == 0 {
				tmp := dist(index, i) + dfs(i, state|(1<<i), 1^player)
				res = op(res, tmp)
			}
		}
		memo[hash] = res
		return res
	}

	res := dfs(m, 0, 0)
	return int(res)
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

func min32(a, b int32) int32 {
	if a < b {
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

// 3283. 吃掉所有兵需要的最多移动次数
// https://leetcode.cn/problems/maximum-number-of-moves-to-kill-all-pawns/description/
// 给你一个 50 x 50 的国际象棋棋盘，棋盘上有 一个 马和一些兵。
// 给你两个整数 kx 和 ky ，其中 (kx, ky) 表示马所在的位置，同时还有一个二维数组 positions ，
// 其中 positions[i] = [xi, yi] 表示第 i 个兵在棋盘上的位置。
// Alice 和 Bob 玩一个回合制游戏，Alice 先手。玩家的一次操作中，可以执行以下操作：
// 玩家选择一个仍然在棋盘上的兵，然后移动马，通过 最少 的 步数 吃掉这个兵。
// 注意 ，玩家可以选择 任意 一个兵，不一定 要选择从马的位置出发 最少 移动步数的兵。
// 在马吃兵的过程中，马 可能 会经过一些其他兵的位置，但这些兵 不会 被吃掉。只有 选中的兵在这个回合中被吃掉。
// Alice 的目标是 最大化 两名玩家的 总 移动次数，直到棋盘上不再存在兵，而 Bob 的目标是 最小化 总移动次数。
// 假设两名玩家都采用 最优 策略，请你返回 Alice 可以达到的 最大 总移动次数。
// 在一次 移动 中，如下图所示，马有 8 个可以移动到的位置，每个移动位置都是沿着坐标轴的一个方向前进 2 格，然后沿着垂直的方向前进 1 格。

const DIR8 = [
  [-2, -1],
  [-2, 1],
  [2, -1],
  [2, 1],
  [-1, -2],
  [-1, 2],
  [1, -2],
  [1, 2]
]

const INF = 1e9
function maxMoves(kx: number, ky: number, positions: number[][]): number {
  const ROW = 50
  const COL = 50
  const bfs = (sx: number, sy: number): number[][] => {
    const dist: number[][] = Array(ROW)
    for (let i = 0; i < ROW; i++) dist[i] = Array(COL).fill(-1)
    dist[sx][sy] = 0
    let queue: number[][] = [[sx, sy]]
    while (queue.length) {
      const nextQueue: number[][] = []
      for (let i = 0; i < queue.length; i++) {
        const [x, y] = queue[i]
        DIR8.forEach(([dx, dy]) => {
          const nx = x + dx
          const ny = y + dy
          if (nx >= 0 && nx < ROW && ny >= 0 && ny < COL && dist[nx][ny] === -1) {
            dist[nx][ny] = dist[x][y] + 1
            nextQueue.push([nx, ny])
          }
        })
      }
      queue = nextQueue
    }
    return dist
  }

  const dists: number[][][] = Array(positions.length + 1)
  for (let i = 0; i < positions.length; i++) {
    dists[i] = bfs(positions[i][0], positions[i][1])
  }
  dists[positions.length] = bfs(kx, ky)

  const dist = (fromIndex: number, toIndex: number): number => {
    const tx = positions[toIndex][0]
    const ty = positions[toIndex][1]
    return dists[fromIndex][tx][ty]
  }

  const m = positions.length
  const mask = (1 << m) - 1
  const memo: Int32Array = new Int32Array((m + 1) * mask * 2).fill(-1)
  const dfs = (index: number, state: number, player: number): number => {
    if (state === mask) return 0
    const hash = index * mask * 2 + state * 2 + player
    if (memo[hash] !== -1) return memo[hash]
    let res = player === 0 ? -INF : INF
    const op = player === 0 ? Math.max : Math.min
    for (let i = 0; i < m; i++) {
      if ((state & (1 << i)) === 0) {
        const tmp = dist(index, i) + dfs(i, state | (1 << i), player ^ 1)
        res = op(res, tmp)
      }
    }
    memo[hash] = res
    return res
  }

  return dfs(m, 0, 0)
}

export {}

// 2258. 逃离火灾
// https://leetcode.cn/problems/escape-the-spreading-fire/description/

// 0 表示草地。
// 1 表示着火的格子。
// 2 表示一座墙，你跟火都不能通过这个格子。
// 2 <= m, n <= 300
// 4 <= m * n <= 2e4

// 一开始你在最左上角的格子 (0, 0) ，你想要到达最右下角的安全屋格子 (m - 1, n - 1) 。
// 每一分钟，你可以移动到 相邻 的草地格子。
// 每次你移动 之后 ，着火的格子会扩散到所有不是墙的 相邻 格子。
// !请你返回你在初始位置可以停留的 最多 分钟数，且停留完这段时间后你还能安全到达安全屋。
// 如果无法实现，请你返回 -1 。如果不管你在初始位置停留多久，你 总是 能到达安全屋，请你返回 1e9 。

// 首先通过 BFS 处理出人到每个格子的最短时间 manTime，
// 以及火到每个格子的最短时间 fireTime。

const INF = 1e9
const DIR4 = [
  [0, 1],
  [0, -1],
  [1, 0],
  [-1, 0]
]

enum State {
  Unvisited = 0,
  PersonVisited = 1,
  FireVisited = 2,
  Wall = 3
}

function maximumMinutes(grid: number[][]): number {
  const ROW = grid.length
  const COL = grid[0].length
  let left = 0
  let right = ROW * COL

  if (!check(left)) return -1 // 无法实现
  if (check(right)) return INF // 不管你在初始位置停留多久，你 总是 能到达安全屋

  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }

  return right

  function check(mid: number): boolean {
    const startRow = 0
    const startCol = 0
    const targetRow = ROW - 1
    const targetCol = COL - 1

    const visited = new Uint8Array(ROW * COL).fill(State.Unvisited)

    let queueMan: [r: number, c: number][] = [[startRow, startCol]]
    visited[startRow * COL + startCol] = State.PersonVisited
    let queueFire: [r: number, c: number][] = []
    for (let r = 0; r < ROW; r++) {
      const row = grid[r]
      for (let c = 0; c < COL; c++) {
        const cell = row[c]
        if (cell === 1) {
          queueFire.push([r, c])
          visited[r * COL + c] = State.FireVisited
        } else if (cell === 2) {
          visited[r * COL + c] = State.Wall
        }
      }
    }

    for (let _ = 0; _ < mid; _++) spreadFire()

    while (queueMan.length) {
      spreadMan()
      if (visited[targetRow * COL + targetCol] === State.PersonVisited) {
        return true
      }
      spreadFire()
    }

    return false

    function spreadMan(): void {
      const nextQueue: [r: number, c: number][] = []
      queueMan.forEach(([curR, curC]) => {
        // !只能由人访问过的点才能继续扩散
        if (visited[curR * COL + curC] !== State.PersonVisited) return
        for (const [dr, dc] of DIR4) {
          const nextR = curR + dr
          const nextC = curC + dc
          if (
            nextR >= 0 &&
            nextR < ROW &&
            nextC >= 0 &&
            nextC < COL &&
            visited[nextR * COL + nextC] === State.Unvisited
          ) {
            visited[nextR * COL + nextC] = State.PersonVisited
            nextQueue.push([nextR, nextC])
          }
        }
      })

      queueMan = nextQueue
    }

    function spreadFire(): void {
      const nextQueue: [r: number, c: number][] = []
      queueFire.forEach(([curR, curC]) => {
        for (const [dr, dc] of DIR4) {
          const nextR = curR + dr
          const nextC = curC + dc
          if (
            nextR >= 0 &&
            nextR < ROW &&
            nextC >= 0 &&
            nextC < COL &&
            // !只能继续扩散到未访问/人访问过的点
            (visited[nextR * COL + nextC] === State.Unvisited ||
              visited[nextR * COL + nextC] === State.PersonVisited)
          ) {
            visited[nextR * COL + nextC] = State.FireVisited
            nextQueue.push([nextR, nextC])
          }
        }
      })

      queueFire = nextQueue
    }
  }
}

export {}

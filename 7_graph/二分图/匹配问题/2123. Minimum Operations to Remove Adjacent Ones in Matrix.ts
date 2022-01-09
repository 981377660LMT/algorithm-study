// 相邻两个1组成一条边，每条边都要去掉一个端点，
// 其实是找最小点覆盖，即求二分图的最大匹配，跑匈牙利算法
function minimumOperations(grid: number[][]): number {
  const [row, col] = [grid.length, grid[0].length]
  const adjMap = new Map<number, Set<number>>()

  for (let r = 0; r < row; r++) {
    for (let c = 0; c < col; c++) {
      if (grid[r][c] === 0) continue

      const cur = r * col + c

      if (r - 1 >= 0 && grid[r - 1][c] === 1) {
        const up = (r - 1) * col + c
        !adjMap.has(up) && adjMap.set(up, new Set())
        adjMap.get(up)!.add(cur)
        !adjMap.has(cur) && adjMap.set(cur, new Set())
        adjMap.get(cur)!.add(up)
      }

      if (c - 1 >= 0 && grid[r][c - 1] === 1) {
        const left = r * col + c - 1
        !adjMap.has(left) && adjMap.set(left, new Set())
        adjMap.get(left)!.add(cur)
        !adjMap.has(cur) && adjMap.set(cur, new Set())
        adjMap.get(cur)!.add(left)
      }
    }
  }

  return hungarian(adjMap)
}

function hungarian(adjMap: Map<number, Set<number>>): number {
  let maxMatching = 0
  const visited = new Set<number>()
  const matching = new Map<number, number>()

  const colors = bisect(adjMap)
  for (const cur of adjMap.keys()) {
    // 从左侧还没有匹配到的男生出发，并重置visited
    if (colors.get(cur) === 0 && !matching.has(cur)) {
      visited.clear()
      if (dfs(cur)) maxMatching++
    }
  }

  return maxMatching

  // 匈牙利算法核心:寻找增广路径 找到的话最大匹配加一
  // dfs(cur) 表示给cur找匹配
  function dfs(cur: number): boolean {
    if (visited.has(cur)) return false
    visited.add(cur)

    for (const next of adjMap.get(cur) || []) {
      // 是增广路径或者dfs找到增广路径
      if (!matching.has(next) || dfs(matching.get(next)!)) {
        matching.set(cur, next)
        matching.set(next, cur)
        return true
      }
    }

    return false
  }

  // 二分图检测、获取colors
  function bisect(adjMap: Map<number, Set<number>>): Map<number, number> {
    const colors = new Map<number, number>()

    const dfs = (cur: number, color: number): void => {
      colors.set(cur, color)
      for (const next of adjMap.get(cur) || []) {
        if (!colors.has(next)) {
          dfs(next, color ^ 1)
        } else if (colors.get(next) === colors.get(cur)) {
          throw new Error('不是二分图')
        }
      }
    }

    for (const cur of adjMap.keys()) {
      if (!colors.has(cur)) dfs(cur, 0)
    }

    return colors
  }
}

console.log(
  minimumOperations([
    [1, 1, 0],
    [0, 1, 1],
    [1, 1, 1],
  ])
)

export {}

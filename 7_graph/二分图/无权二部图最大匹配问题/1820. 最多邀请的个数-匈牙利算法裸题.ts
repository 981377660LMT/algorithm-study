// # 给定一个 m x n 的整数矩阵 grid ，其中 grid[i][j] 等于 0 或 1 。
// # 若 grid[i][j] == 1 ，则表示第 i 个男孩可以邀请第 j 个女孩参加派对。
// # 一个男孩最多可以邀请一个女孩，一个女孩最多可以接受一个男孩的一个邀请。
// # 返回可能的最多邀请的个数。
// 男孩是行,女孩是列

function maximumInvitations(grid: number[][]): number {
  const [row, col] = [grid.length, grid[0].length]
  const matched = new Map<number, number>()
  let res = 0

  for (let boy = 0; boy < row; boy++) {
    const visited = new Set<number>()
    if (dfs(boy, visited)) res++
  }

  return res

  // 当前男生能否匹配一个女生或者找到增广路径
  function dfs(boy: number, visited: Set<number>): boolean {
    for (let girl = 0; girl < col; girl++) {
      if (grid[boy][girl] === 0) continue
      if (visited.has(girl)) continue
      visited.add(girl)
      // 女生未配对或可找到增广路径
      if (!matched.has(girl) || dfs(matched.get(girl)!, visited)) {
        matched.set(girl, boy)
        return true
      }
    }

    return false
  }
}

console.log(
  maximumInvitations([
    [1, 1, 1],
    [1, 0, 1],
    [0, 0, 1],
  ])
)

// # 输出: 3
// # 解释: 按下列方式邀请：
// # - 第 1 个男孩邀请第 3 个女孩。
// # - 第 2 个男孩邀请第 1 个女孩。
// # - 第 3 个男孩未邀请任何人。
// # - 第 4 个男孩邀请第 2 个女孩。
export {}

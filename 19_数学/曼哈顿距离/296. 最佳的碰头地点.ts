// 求解多个绝对值之和最小值:奇尖偶平中间最小
function minTotalDistance(grid: number[][]): number {
  const xNums: number[] = []
  const yNums: number[] = []
  for (let i = 0; i < grid.length; i++) {
    for (let j = 0; j < grid[i].length; j++) {
      if (grid[i][j] === 1) {
        xNums.push(i)
        yNums.push(j)
      }
    }
  }

  xNums.sort((a, b) => a - b)
  yNums.sort((a, b) => a - b)
  const xMiddle = xNums[xNums.length >> 1]
  const yMiddle = yNums[yNums.length >> 1]

  let res = 0
  for (const x of xNums) {
    res += Math.abs(x - xMiddle)
  }

  for (const y of yNums) {
    res += Math.abs(y - yMiddle)
  }

  return res
}

console.log(
  minTotalDistance([
    [1, 0, 0, 0, 1],
    [0, 0, 0, 0, 0],
    [0, 0, 1, 0, 0],
  ])
)
// 有一队人（两人或以上）想要在一个地方碰面，他们希望能够最小化他们的总行走距离。
// 输出: 6
// 解析:
// 给定的三个人分别住在(0,0)，(0,4) 和 (2,2):
// (0,2) 是一个最佳的碰面点，其总行走距离为 2 + 2 + 2 = 6，
// 最小，因此返回 6。

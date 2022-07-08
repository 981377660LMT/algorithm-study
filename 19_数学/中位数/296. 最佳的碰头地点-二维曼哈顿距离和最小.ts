// 二维的求中位数
function minTotalDistance(grid: number[][]): number {
  const px: number[] = []
  const py: number[] = []
  for (let i = 0; i < grid.length; i++) {
    for (let j = 0; j < grid[i].length; j++) {
      if (grid[i][j] === 1) {
        px.push(i)
        py.push(j)
      }
    }
  }

  px.sort((a, b) => a - b)
  py.sort((a, b) => a - b)
  const midX = px[px.length >> 1]
  const midY = py[py.length >> 1]

  let res = 0
  for (const x of px) res += Math.abs(x - midX)
  for (const y of py) res += Math.abs(y - midY)
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

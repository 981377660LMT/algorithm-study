// 注意降维使逻辑更清晰
function nearestValidPoint(x: number, y: number, points: number[][]): number {
  const validPoints = points
    .map<[x: number, y: number, index: number]>((point, index) => [point[0], point[1], index])
    .filter(([row, col]) => row === x || col === y)

  if (validPoints.length === 0) return -1

  return validPoints
    .map(([row, col, index]) => [Math.hypot(row - x, col - y), index])
    .sort((a, b) => a[0] - b[0] || a[1] - b[1])[0][1]
}

// 请返回距离你当前位置 曼哈顿距离 最近的 有效 点的下标（下标从 0 开始）。
// 如果有多个最近的有效点，请返回下标 最小 的一个。如果没有有效点，请返回 -1 。

console.log(
  nearestValidPoint(3, 4, [
    [1, 2],
    [3, 1],
    [2, 4],
    [2, 3],
    [4, 4],
  ])
)

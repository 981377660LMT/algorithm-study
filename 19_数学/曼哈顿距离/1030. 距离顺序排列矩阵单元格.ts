/**
 * @param {number} rows
 * @param {number} cols
 * @param {number} rCenter
 * @param {number} cCenter
 * @return {number[][]}
 * 返回矩阵中的所有单元格的坐标，
 * 并按到 (r0, c0) 的距离从最小到最大的顺序排，
 * 其中，两单元格(r1, c1) 和 (r2, c2) 之间的距离是曼哈顿距离|r1 - r2| + |c1 - c2|。

 */
var allCellsDistOrder = function (
  rows: number,
  cols: number,
  rCenter: number,
  cCenter: number
): number[][] {
  const bucket = Array.from<unknown, [row: number, col: number][]>({ length: 200 }, () => [])
  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      const dis = Math.abs(i - rCenter) + Math.abs(j - cCenter)
      bucket[dis].push([i, j])
    }
  }
  console.log(bucket)
  return bucket.flat()
}

console.log(allCellsDistOrder(1, 2, 0, 0))
console.log(allCellsDistOrder(2, 2, 0, 1))

/**
 * @param {number[][]} lines  1 <= lines.length <= 10^5 且 lines[i].length == 2
 * @return {number}
 * 所有直线以 [k,b] 的形式存于二维数组 lines 中，不存在重合的两条直线
 * @summary
 * 算一下交点坐标，很容易发现可以对应到(k, b)坐标系的N个点找斜率最大值问题，然后写斜率最大值的算法就好了
 *
 */
function minRecSize(lines: number[][]): number {
  if (lines.length <= 2) return 0

  lines.sort((a, b) => a[0] - b[0] || b[1] - a[1])

  const nlines: number[][] = []
  let minX = Infinity,
    minY = Infinity,
    maxX = -Infinity,
    maxY = -Infinity

  for (let i = 0; i < lines.length; i++) {
    if (
      i === 0 ||
      i === lines.length - 1 ||
      lines[i][0] !== lines[i - 1][0] ||
      lines[i][0] !== lines[i + 1][0]
    ) {
      nlines.push(lines[i])
    }
  }

  for (let i = 0; i < nlines.length - 1; i++) {
    const [k1, b1] = nlines[i]
    for (let j = i + 1; j < nlines.length; j++) {
      const [k2, b2] = nlines[j]

      if (k1 == k2) continue

      const x = -(b1 - b2) / (k1 - k2)
      const y = k1 * x + b1

      minX = Math.min(minX, x)
      minY = Math.min(minY, y)
      maxX = Math.max(maxX, x)
      maxY = Math.max(maxY, y)
    }
  }

  const res = (maxX - minX) * (maxY - minY)

  return res === Infinity ? 0 : res
}

console.log(
  minRecSize([
    [2, 3],
    [3, 0],
    [4, 1],
  ])
)

// 类似于 LCP 29. 乐团站位.ts

// m 和 n 都是 偶数
// 1 <= k <= 10^9
function rotateGrid(grid: number[][], k: number): number[][] {
  const [row, col] = [grid.length, grid[0].length]
  // 维护每层的位置
  let [top, bottom] = [0, row - 1]
  let [left, right] = [0, col - 1]

  const res = Array.from<unknown, number[]>({ length: row }, () => Array(col).fill(0))

  // 从外向内处理每一层
  while (top < bottom && left < right) {
    const curLayer: [x: number, y: number][] = []

    // 从左上角开始，左下右上放入，旋转之后看变到了这层的哪一个
    for (let r = top; r < bottom; r++) curLayer.push([r, left])
    for (let c = left; c < right; c++) curLayer.push([bottom, c])
    for (let r = bottom; r > top; r--) curLayer.push([r, right])
    for (let c = right; c > left; c--) curLayer.push([top, c])

    const rotate = k % curLayer.length
    curLayer.forEach(([preR, preC], index) => {
      const [nextR, nextC] = curLayer[(index + rotate) % curLayer.length]
      res[nextR][nextC] = grid[preR][preC]
    })

    top++
    bottom--
    left++
    right--
  }

  return res
}

console.log(
  rotateGrid(
    [
      [1, 2, 3, 4],
      [5, 6, 7, 8],
      [9, 10, 11, 12],
      [13, 14, 15, 16],
    ],
    2
  )
)
// 输出：[[3,4,8,12],[2,11,10,16],[1,7,6,15],[5,9,13,14]]

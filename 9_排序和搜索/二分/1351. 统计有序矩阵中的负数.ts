// 你可以设计一个时间复杂度为 O(n + m) 的解决方案吗？
function countNegatives(grid: number[][]): number {
  const [m, n] = [grid.length, grid[0].length]
  let row = m - 1
  let col = 0
  let res = 0

  while (row >= 0 && col < n) {
    if (grid[row][col] < 0) {
      res += n - col
      row--
    } else {
      col++
    }
  }
  return res
}

console.log(
  countNegatives2([
    [4, 3, 2, -1],
    [3, 2, 1, -1],
    [1, 1, -1, -2],
    [-1, -1, -2, -3],
  ])
)
// 给你一个 m * n 的矩阵 grid，矩阵中的元素无论是按行还是按列，都以非递增顺序排列。

// 请你统计并返回 grid 中 负数 的数目。

// mlogn
function countNegatives2(grid: number[][]): number {
  let res = 0

  for (const row of grid) {
    res += count(row)
  }

  return res

  // 最左能力二分
  function count(row: number[]): number {
    let [l, r] = [0, row.length - 1]
    while (l <= r) {
      const mid = (l + r) >> 1
      if (row[mid] >= 0) {
        l = mid + 1
      } else {
        r = mid - 1
      }
    }

    return row.length - l
  }
}

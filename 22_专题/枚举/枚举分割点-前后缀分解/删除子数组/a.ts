// 给你一个下标从 0 开始、大小为 n * m 的二维整数矩阵 grid ，定义一个下标从 0 开始、大小为 n * m 的的二维矩阵 p。如果满足以下条件，则称 p 为 grid 的 乘积矩阵 ：

// 对于每个元素 p[i][j] ，它的值等于除了 grid[i][j] 外所有元素的乘积。乘积对 12345 取余数。
// 返回 grid 的乘积矩阵。

const MOD = 12345

function constructProductMatrix(grid: number[][]): number[][] {
  const ROW = grid.length
  const COL = grid[0].length

  // 每一行的乘积
  const rowProduct = Array(ROW).fill(0)
  for (let i = 0; i < ROW; i++) {
    rowProduct[i] = grid[i].reduce((acc, cur) => (acc * cur) % MOD, 1)
  }
  const rowWithOutOne = productWithoutOne(
    rowProduct,
    () => 1,
    (a, b) => (a * b) % MOD
  )

  const res = Array(ROW)
  for (let i = 0; i < ROW; i++) {
    res[i] = Array(COL)
    const curRowWithoutOne = rowWithOutOne[i]
    const tmp = productWithoutOne(
      grid[i],
      () => 1,
      (a, b) => (a * b) % MOD
    )
    for (let j = 0; j < COL; j++) {
      res[i][j] = (tmp[j] * curRowWithoutOne) % MOD
    }
  }
  return res
}

function productWithoutOne<E>(nums: E[], e: () => E, op: (a: E, b: E) => E): E[] {
  const n = nums.length
  const res: E[] = Array(n)
  for (let i = 0; i < n; i++) {
    res[i] = e()
  }
  for (let i = 0; i < n - 1; i++) {
    res[i + 1] = op(res[i], nums[i])
  }
  let x = e()
  for (let i = n - 1; ~i; i--) {
    res[i] = op(res[i], x)
    x = op(nums[i], x)
  }
  return res
}

export {}

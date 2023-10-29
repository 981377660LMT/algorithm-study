/* eslint-disable no-inner-declarations */
/**
 * 除自身以外数组的乘积.
 */
function productWithoutOne<E>(nums: ArrayLike<E>, e: () => E, op: (a: E, b: E) => E): E[] {
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

/**
 * 除自身以外数组的乘积二维版本.
 */
function productWithoutOne2D<E>(
  nums: ArrayLike<ArrayLike<E>>,
  e: () => E,
  op: (a: E, b: E) => E
): E[][] {
  const ROW = nums.length
  const COL = nums[0].length

  let unit = e()
  const res: E[][] = Array(ROW)
  for (let i = 0; i < ROW; i++) {
    const row = Array(COL)
    for (let j = 0; j < COL; j++) row[j] = unit
    res[i] = row
  }

  for (let i = ROW - 1; ~i; i--) {
    const tmp1 = res[i]
    const tmp2 = nums[i]
    for (let j = COL - 1; ~j; j--) {
      tmp1[j] = unit
      unit = op(tmp2[j], unit)
    }
  }

  unit = e()
  for (let i = 0; i < ROW; i++) {
    const tmp1 = res[i]
    const tmp2 = nums[i]
    for (let j = 0; j < COL; j++) {
      tmp1[j] = op(unit, tmp1[j])
      unit = op(unit, tmp2[j])
    }
  }

  return res
}

export { productWithoutOne, productWithoutOne2D }

if (require.main === module) {
  function productExceptSelf(nums: number[]): number[] {
    return productWithoutOne(
      nums,
      () => 1,
      (a, b) => a * b
    )
  }

  // 2906. 构造乘积矩阵
  // https://leetcode.cn/problems/construct-product-matrix/description/
  function constructProductMatrix(grid: number[][]): number[][] {
    return productWithoutOne2D(
      grid,
      () => 1,
      (a, b) => (a * b) % 12345
    )
  }
}

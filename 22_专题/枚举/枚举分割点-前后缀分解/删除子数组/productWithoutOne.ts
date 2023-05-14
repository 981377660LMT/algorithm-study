/**
 * 除自身以外数组的乘积.
 */
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

export { productWithoutOne }

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function productExceptSelf(nums: number[]): number[] {
    return productWithoutOne(
      nums,
      () => 1,
      (a, b) => a * b
    )
  }
}

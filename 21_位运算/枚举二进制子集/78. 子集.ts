/**
 *
 * @param nums
 * @returns 子集
 */
function subsets<T>(nums: T[]): T[][] {
  const n = 1 << nums.length
  const res: T[][] = []

  for (let state = 0; state < n; state++) {
    const tmp: T[] = []
    for (let j = 0; j < nums.length; j++) {
      if (state & (1 << j)) tmp.push(nums[j])
    }
    res.push(tmp)
  }

  return res
}

export { subsets }

if (require.main === module) {
  console.log(subsets([1, 2, 3]))
}

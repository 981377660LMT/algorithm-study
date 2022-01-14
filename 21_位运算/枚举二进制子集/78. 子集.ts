/**
 *
 * @param nums
 * @returns 2^n 个子集
 */
function subsets<T>(nums: T[]): T[][] {
  const n = 1 << nums.length
  const res: T[][] = []

  for (let state = 0; state < n; state++) {
    const cands: T[] = []
    for (let j = 0; j < nums.length; j++) {
      if (state & (1 << j)) cands.push(nums[j])
    }
    res.push(cands)
  }

  return res
}

function GosperHack(n: number, k: number) {
  let x = (1 << k) - 1
  const limit = 1 << n
  while (x < limit) {
    console.log(x.toString(2).padStart(5, '0'))
    const lowbit = x & -x
    const r = x + lowbit
    // xor
    x = r | (((x ^ r) >> 2) / lowbit)
  }
}

export { subsets }

if (require.main === module) {
  console.log(subsets([1, 2, 3]))
}

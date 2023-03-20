/**
 * 遍历子集.
 */
function enumerateSubset<T>(nums: ArrayLike<T>, callback: (subset: T[]) => void): void {
  const n = nums.length
  for (let state = 0; state < 1 << n; state++) {
    const cands: T[] = []
    for (let j = 0; j < nums.length; j++) {
      if (state & (1 << j)) cands.push(nums[j])
    }
    callback(cands)
  }
}

export { enumerateSubset }

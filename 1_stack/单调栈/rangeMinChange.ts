interface IChange {
  /** 区间左端点开始位置 */
  lStart: number
  /** 区间左端点结束位置 */
  lEnd: number
  /** 旧的最小值 */
  oldMin: number
  /** 新的最小值 */
  newMin: number
}

/**
 * 维护区间最小值的变化历史。
 * 返回：res[i]，表示右端点r=i+1时，所有受影响区间[l,r)的最小值变化记录：(l, r, old_min, new_min)
 * 每次右端点推进，所有被当前元素“刷新”最小值的区间都会被记录下来，适用于区间DP、单调栈优化等场景。
 */
function rangeMinChange(nums: ArrayLike<number>, defaultValue = Infinity): IChange[][] {
  const n = nums.length
  const res: IChange[][] = Array(n)
  for (let i = 0; i < n; i++) res[i] = []

  type Triple = {
    l: number
    r: number
    min: number
  }

  const stack: Triple[] = []
  for (let i = 0; i < n; i++) {
    const v = nums[i]
    res[i].push({ lStart: i, lEnd: i + 1, oldMin: defaultValue, newMin: v })
    let ptr = i
    while (stack.length > 0) {
      const lc = stack[stack.length - 1]
      if (lc.min <= v) break
      res[i].push({ lStart: lc.l, lEnd: lc.r, oldMin: lc.min, newMin: v })
      ptr = lc.l
      stack.pop()
    }
    stack.push({ l: ptr, r: i + 1, min: v })
    res[i].reverse()
  }

  return res
}

export { rangeMinChange }

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  // Test cases
  const nums = [3, 1, 4, 1, 5, 9, 2, 6]
  const changes = rangeMinChange(nums)

  for (let i = 0; i < changes.length; i++) {
    console.log(`Right endpoint r=${i + 1}:`, changes[i])
  }
}

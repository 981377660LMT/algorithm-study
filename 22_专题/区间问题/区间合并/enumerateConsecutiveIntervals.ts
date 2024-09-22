/* eslint-disable no-inner-declarations */

/** 遍历连续区间. */
function enumerateConsecutiveIntervals(
  n: number,
  supplier: (i: number) => number,
  consumer: (min: number, max: number, isIn: boolean) => boolean | void
): void {
  if (!n) return
  let i = 0
  while (i < n) {
    const start = i
    while (i < n - 1 && supplier(i) + 1 === supplier(i + 1)) {
      i++
    }
    if (consumer(supplier(start), supplier(i), true)) return
    if (i + 1 < n) {
      if (consumer(supplier(i) + 1, supplier(i + 1) - 1, false)) return
    }
    i++
  }
}

export { enumerateConsecutiveIntervals }

if (require.main === module) {
  // 228. 汇总区间
  // https://leetcode.cn/problems/summary-ranges/description/
  function summaryRanges(nums: number[]): string[] {
    const res: string[] = []
    enumerateConsecutiveIntervals(
      nums.length,
      i => nums[i],
      (min, max, isIn) => {
        if (!isIn) return
        res.push(min === max ? `${min}` : `${min}->${max}`)
      }
    )
    return res
  }
}

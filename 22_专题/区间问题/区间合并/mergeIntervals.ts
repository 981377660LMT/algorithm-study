/** 合并所有重叠的闭区间，返回一个不重叠的区间列表. */
function mergeIntervals(
  n: number,
  supplier: (i: number) => { left: number; right: number },
  consumer: (newLeft: number, newRight: number) => void
): void {
  if (n === 0) {
    return
  }
  const order = Array<number>(n)
  for (let i = 0; i < n; i++) order[i] = i
  order.sort((i, j) => {
    const l1 = supplier(i).left
    const l2 = supplier(j).left
    return l1 - l2
  })
  let { left: preL, right: preR } = supplier(order[0])
  for (let i = 1; i < n; i++) {
    const { left: curL, right: curR } = supplier(order[i])
    if (curL <= preR) {
      preR = Math.max(preR, curR)
    } else {
      consumer(preL, preR)
      preL = curL
      preR = curR
    }
  }
  consumer(preL, preR)
}

export { mergeIntervals }

if (require.main === module) {
  // 56. 合并区间
  // https://leetcode.cn/problems/merge-intervals/description/
  function merge(intervals: number[][]): number[][] {
    const res: number[][] = []
    mergeIntervals(
      intervals.length,
      i => ({ left: intervals[i][0], right: intervals[i][1] }),
      (newLeft, newRight) => res.push([newLeft, newRight])
    )
    return res
  }
}

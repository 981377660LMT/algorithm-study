/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */

/**
 * 删除重叠覆盖区间.
 * @param intervals 左闭右开区间.
 * @param removeIncluded 是删除包含的区间还是被包含的区间.默认为删除被包含的区间.
 * @returns 按照区间的起点排序的剩余的区间索引(相同的区间会保留).
 */
function reduceIntervals(intervals: [start: number, end: number][] | number[][], removeIncluded = true): number[] {
  const n = intervals.length
  const res: number[] = []
  const order = Array(n)
  for (let i = 0; i < n; i++) order[i] = i
  if (removeIncluded) {
    order.sort((a, b) => intervals[a][0] - intervals[b][0] || intervals[b][1] - intervals[a][1])
    for (let i = 0; i < n; i++) {
      const cur = order[i]
      if (res.length) {
        const pre = res[res.length - 1]
        const { 0: curStart, 1: curEnd } = intervals[cur]
        const { 0: preStart, 1: preEnd } = intervals[pre]
        if (curEnd <= preEnd && curEnd - curStart < preEnd - preStart) continue
      }
      res.push(cur)
    }
  } else {
    order.sort((a, b) => intervals[a][1] - intervals[b][1] || intervals[b][0] - intervals[a][0])
    for (let i = 0; i < n; i++) {
      const cur = order[i]
      if (res.length) {
        const pre = res[res.length - 1]
        const { 0: curStart, 1: curEnd } = intervals[cur]
        const { 0: preStart, 1: preEnd } = intervals[pre]
        if (curStart <= preStart && curEnd - curStart > preEnd - preStart) continue
      }
      res.push(cur)
    }
  }

  return res
}

export { reduceIntervals }

if (require.main === module) {
  // https://leetcode.cn/problems/remove-covered-intervals/
  // 1288. 删除被覆盖区间
  function removeCoveredIntervals(intervals: number[][]): number {
    return reduceIntervals(intervals).length
  }

  console.log(
    reduceIntervals([
      [1, 2],
      [1, 2]
    ])
  )
}

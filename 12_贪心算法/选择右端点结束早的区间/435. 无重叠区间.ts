// 给定一个区间的集合 intervals ，其中 intervals[i] = [starti, endi] 。
// 返回 需要移除区间的最小数量，使剩余区间互不重叠 。

function eraseOverlapIntervals(intervals: number[][]): number {
  if (intervals.length <= 1) return 0
  intervals.sort((a, b) => a[1] - b[1])

  let res = 0
  let preEnd = -2e15

  for (let i = 0; i < intervals.length; i++) {
    const [start, end] = intervals[i]
    if (start < preEnd) {
      res++
    } else {
      preEnd = end
    }
  }

  return res
}

console.log(
  eraseOverlapIntervals([
    [1, 100],
    [11, 22],
    [1, 11],
    [2, 12]
  ])
)
export {}

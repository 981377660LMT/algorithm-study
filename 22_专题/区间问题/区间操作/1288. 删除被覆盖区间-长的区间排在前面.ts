// 区间覆盖

/**
 *
 * @param intervals 删除列表中被其他区间所覆盖的区间。
 * 返回列表中剩余区间的数目。
 */
function removeCoveredIntervals(intervals: number[][]): number {
  if (intervals.length <= 1) return intervals.length
  // 需要让长的区间排在前面
  intervals.sort((a, b) => a[0] - b[0] || b[1] - a[1])

  let merge = 0

  let preEnd = -1
  for (let index = 0; index < intervals.length; index++) {
    const [_, end] = intervals[index]

    // 三种关系:包含，相交，相离
    // 判断包含：需要让长的区间排在前面
    if (end <= preEnd) {
      merge++
    } else {
      preEnd = end
    }
  }

  return intervals.length - merge
}

console.log(
  removeCoveredIntervals([
    [1, 4],
    [3, 6],
    [2, 8]
  ])
)

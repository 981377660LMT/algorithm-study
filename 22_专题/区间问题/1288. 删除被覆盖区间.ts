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
  let preLeft = intervals[0][0]
  let preRight = intervals[0][1]
  for (let index = 1; index < intervals.length; index++) {
    const [curLeft, curRight] = intervals[index]

    // 三种关系:包含，相交，相离
    // 判断包含：需要让长的区间排在前面
    if (curRight <= preRight) {
      merge++
    } else if (curLeft <= preRight && curRight >= preRight) {
      preRight = curRight
    } else {
      preLeft = curLeft
      preRight = curRight
    }
  }

  return intervals.length - merge
}

console.log(
  removeCoveredIntervals([
    [1, 4],
    [3, 6],
    [2, 8],
  ])
)

/**
 * @param {number[][]} points
 * @return {number}
 * 有多少不重叠的区间
 */
function findMinArrowShots(points: number[][]): number {
  // sort by the earliest finish time
  points.sort((a, b) => a[1] - b[1])

  let preEnd = -Infinity
  let res = 1

  for (let i = 0; i < points.length; i++) {
    const [curStart, curEnd] = points[i]
    if (preEnd < curStart) {
      preEnd = curEnd
      res++
    }
  }

  return res
}

console.log(
  findMinArrowShots([
    [10, 16],
    [2, 8],
    [1, 6],
    [7, 12],
  ])
)
// 输出：2
// 解释：对于该样例，x = 6 可以射爆 [2,8],[1,6] 两个气球，以及 x = 11 射爆另外两个气球

export default 1

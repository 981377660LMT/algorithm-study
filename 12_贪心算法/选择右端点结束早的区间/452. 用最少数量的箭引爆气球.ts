/**
 * @param {number[][]} points
 * @return {number}
 * 返回引爆所有气球所必须射出的 最小 弓箭数 。
 * !有多少不重叠的区间
 */
function findMinArrowShots(points: number[][]): number {
  // 区间结束时间/区间右端点排序
  points.sort((a, b) => a[1] - b[1])

  let res = 0
  let preEnd = -2e15

  for (let i = 0; i < points.length; i++) {
    const [start, end] = points[i]
    if (start > preEnd) {
      preEnd = end
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
    [7, 12]
  ])
)
// 输出：2
// 解释：对于该样例，x = 6 可以射爆 [2,8],[1,6] 两个气球，以及 x = 11 射爆另外两个气球

export default 1

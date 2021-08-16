/**
 * @param {number[][]} points
 * @return {number}
 * 有多少不重叠的区间
 */
var findMinArrowShots = function (points: number[][]): number {
  // sort by the earliest finish time
  points.sort((a, b) => a[1] - b[1])
  let prev = points[0],
    chain = 1

  for (let i = 1; i < points.length; i++) {
    const [prevS, prevE] = prev
    const [currS, currE] = points[i]
    if (prevE < currS) {
      prev = points[i]
      chain++
    }
  }
  console.log(points)
  return chain
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

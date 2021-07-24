/**
 * @param {number[][]} intervals
 * @return {number}
 * @description 给定一个区间的集合，找到需要移除区间的最小数量，使剩余区间互不重叠。
 * @description 思路：按照结尾排序，不能和前一个区间有重叠
 * @description Time Complexity: O(N LogN); Space Complexity: O(1);
 */
var eraseOverlapIntervals = function (intervals) {
  intervals.sort((a, b) => a[1] - b[1])
  let res = 0
  let prev = intervals[0]
  const len = intervals.length

  for (let i = 1; i < len; i++) {
    if (intervals[i][0] < prev[1]) res++
    else prev = intervals[i]
  }

  return res
}

console.log(
  eraseOverlapIntervals([
    [1, 2],
    [2, 3],
    [3, 4],
    [1, 3],
  ])
)
// 输出: 1

// 解释: 移除 [1,3] 后，剩下的区间没有重叠。

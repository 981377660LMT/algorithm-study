/**
 * @param {number[][]} intervals
 * @return {number[]}
 */
const findRightInterval = function (intervals: number[][]): number[] {
  const res: number[] = []
  const startMap = new Map<number, number>() // number index
  const start: number[] = []
  intervals.forEach((value, index) => {
    startMap.set(value[0], index)
    start.push(value[0])
  })
  start.sort((a, b) => a - b)

  for (let i = 0; i < intervals.length; i++) {
    const nearestStart = search(start, intervals[i][1])
    if (nearestStart === Infinity) res.push(-1)
    else res.push(startMap.get(nearestStart)!)
  }

  return res

  /**
   *
   * @param nums
   * @param target
   * @returns 寻找不小于target的最小的start值
   */
  function search(nums: number[], target: number): number {
    let res = Infinity
    let l = 0
    let r = nums.length - 1
    while (l <= r) {
      let mid = (l + r) >> 1
      if (nums[mid] === target) return target
      else if (nums[mid] > target) {
        res = nums[mid]
        r = mid - 1
      } else {
        l = mid + 1
      }
    }

    return res
  }
}

console.log(
  findRightInterval([
    [3, 4],
    [2, 3],
    [1, 2],
  ])
)

// 输出：[-1, 0, 1]
// 解释：对于 [3,4] ，没有满足条件的“右侧”区间。
// 对于 [2,3] ，区间[3,4]具有最小的“右”起点;
// 对于 [1,2] ，区间[2,3]具有最小的“右”起点。

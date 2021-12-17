/**
 * @param {number[][]} intervals  给你一个 无重叠的 ，按照区间起始端点排序的区间列表。
 * @param {number[]} newInterval
 * @return {number[][]}
 * 你需要确保列表中的区间仍然有序且不重叠（如果有必要的话，可以合并区间）。
 */
const insert = function (intervals: number[][], newInterval: number[]): number[][] {
  bisectInsort(intervals, newInterval)
  return mergeSortedArray(intervals)

  function bisectInsort(target: number[][], inserted: number[]) {
    const num = inserted[0]
    let l = 0
    let r = target.length - 1
    while (l <= r) {
      const mid = (l + r) >> 1
      if (target[mid][0] >= num) r = mid - 1
      else l = mid + 1
    }
    target.splice(l, 0, inserted)
  }

  function mergeSortedArray(nums: number[][]) {
    const res: number[][] = [nums[0]]
    for (let index = 1; index < nums.length; index++) {
      const interval = nums[index]
      const [preLeft, preRight] = res[res.length - 1]
      const [curLeft, curRight] = interval

      // 三种关系:包含，相交，相离
      if (curRight <= preRight) {
        continue
      } else if (curLeft <= preRight && curRight >= preRight) {
        res.pop()
        res.push([preLeft, curRight])
      } else {
        res.push(interval)
      }
    }
    return res
  }
}

console.log(
  insert(
    [
      [1, 3],
      [6, 9],
    ],
    [2, 5]
  )
)

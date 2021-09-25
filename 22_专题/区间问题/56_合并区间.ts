// 合并所有重叠的区间，并返回一个不重叠的区间数组

const merge = (intervals: number[][]): number[][] => {
  if (intervals.length <= 1) return intervals
  intervals.sort((a, b) => a[0] - b[0] || b[1] - a[1])

  const mergeSortedArray = (nums: number[][]) => {
    const res: number[][] = [nums[0]]
    for (let index = 1; index < nums.length; index++) {
      const [preLeft, preRight] = res[res.length - 1]
      const [curLeft, curRight] = nums[index]

      // 三种关系:包含，相交，相离
      // 判断包含：需要让长的区间排在前面
      if (curRight <= preRight) {
        continue
      } else if (curLeft <= preRight && curRight >= preRight) {
        res.pop()
        res.push([preLeft, curRight])
      } else {
        res.push(nums[index])
      }
    }
    return res
  }

  return mergeSortedArray(intervals)
}

console.log(
  merge([
    [1, 3],
    [2, 6],
    [8, 10],
    [15, 18],
  ])
)

export {}

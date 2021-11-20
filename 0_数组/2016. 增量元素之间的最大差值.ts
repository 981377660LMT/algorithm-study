// 计算 nums[j] - nums[i] 能求得的 最大差值 ，其中 0 <= i < j < n 且 nums[i] < nums[j] 。

// 思路：维护最小值
function maximumDifference(nums: number[]): number {
  let min = nums[0]
  let res = -1

  for (const num of nums) {
    if (num > min) {
      res = Math.max(res, num - min)
    } else {
      min = Math.min(min, num)
    }
  }

  return res
}

console.log(maximumDifference([7, 1, 5, 4]))

export {}

import { bisectLeft } from '../../9_排序和搜索/二分/bisect'

/**
 * !求严格LIS的长度
 * 维护一个LIS数组
   LIS[i] 表示以LIS[i]结尾的最大上升子序列长度为i+1
 */
function lengthOfLIS(nums: number[]): number {
  if (nums.length <= 1) {
    return nums.length
  }

  const LIS: number[] = [nums[0]]

  for (let i = 1; i < nums.length; i++) {
    if (nums[i] > LIS[LIS.length - 1]) {
      LIS.push(nums[i])
    } else {
      LIS[bisectLeft(LIS, nums[i])] = nums[i]
    }
  }

  return LIS.length
}

console.log(lengthOfLIS([7, 7, 7, 7, 7, 7, 7, 8, 9, 4, 1, 2, 3, 7]))
// 输出：4
// 解释：最长递增子序列是 [2,3,7,101]，因此长度为 4

export { lengthOfLIS }

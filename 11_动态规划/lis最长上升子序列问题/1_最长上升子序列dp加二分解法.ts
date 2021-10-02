import { bisectRight } from '../../9_排序和搜索/二分/7_二分搜索寻找最插右入位置'

/**
 * @param {number[]} nums
 * @return {number}
 * @description 求LIS的长度
 * @summary 维护一个LIS数组
 * 时间复杂度 O(NlogN) ： 遍历 nums 列表需 O(N)，在每个 nums[i] 二分法需 O(logN)。
   空间复杂度 O(N) ： LIS (单增数组)占用线性大小额外空间。
   LIS[i] 表示以LIS[i]结尾的最大上升子序列长度为i+1
 */
const lengthOfLIS = function (nums: number[]): number {
  if (nums.length <= 1) return nums.length
  const LIS: number[] = [nums[0]]
  const bisectLeft = (arr: number[], target: number) => {
    let l = 0
    let r = arr.length - 1
    while (l <= r) {
      const mid = (r + l) >> 1
      if (arr[mid] === target) {
        r--
      } else if (arr[mid] > target) {
        r = mid - 1
      } else {
        l = mid + 1
      }
    }

    return l
  }

  console.log(bisectLeft(nums, 7))
  for (let i = 1; i < nums.length; i++) {
    if (nums[i] > LIS[LIS.length - 1]) {
      LIS.push(nums[i])
    } else {
      LIS[bisectLeft(LIS, nums[i])] = nums[i]
    }
  }
  console.log(LIS)
  return LIS.length
}

console.log(lengthOfLIS([7, 7, 7, 7, 7, 7, 7, 8, 9, 4, 1, 2, 3, 7]))
// 输出：4
// 解释：最长递增子序列是 [2,3,7,101]，因此长度为 4

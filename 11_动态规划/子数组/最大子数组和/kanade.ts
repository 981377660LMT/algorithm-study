/* eslint-disable prefer-destructuring */

function maxSubarraySum(arr: ArrayLike<number>): [max: number, start: number, end: number] {
  const n = arr.length
  if (n === 0) return [0, 0, 0] // !根据题意返回

  let maxSum = arr[0]
  let curSum = arr[0]
  let curStart = 0
  let start = 0
  let end = 0
  for (let i = 1; i < n; i++) {
    if (curSum < 0) {
      curSum = 0
      curStart = i // 重新开始
    }
    curSum += arr[i]
    if (curSum > maxSum) {
      maxSum = curSum
      start = curStart
      end = i + 1
    }
  }

  if (maxSum < 0) return [0, 0, 0] // !根据题意返回

  return [maxSum, start, end]
}

/**
 * 最大两段子段和（两段必须间隔至少 gap 个数）.
 */
function maxSubarraySumTwoSum(arr: ArrayLike<number>, gap: number): number {
  const n = arr.length
  const sufSumMax = Array<number>(n).fill(0)
  sufSumMax[n - 1] = arr[n - 1]
  let curSumMax = arr[n - 1]
  for (let i = n - 2; i >= 0; i--) {
    const v = arr[i]
    curSumMax = Math.max(curSumMax + v, v)
    sufSumMax[i] = Math.max(sufSumMax[i + 1], curSumMax)
  }
  curSumMax = arr[0]
  let preSumMax = arr[0]
  let res = preSumMax + sufSumMax[1 + gap]
  for (let i = 1; i + 1 + gap < n; i++) {
    const v = arr[i]
    curSumMax = Math.max(curSumMax + v, v)
    preSumMax = Math.max(preSumMax, curSumMax)
    res = Math.max(res, preSumMax + sufSumMax[i + 1 + gap])
  }
  return res
}

/**
 * @description
 * 求最大/最小子数组和
 */
function kanade(nums: number[], getMax = true): number {
  if (nums.length === 0) return 0
  if (nums.length === 1) return nums[0]

  let res = nums[0]
  let dp = nums[0]
  for (let i = 1; i < nums.length; i++) {
    const num = nums[i]
    if (getMax) {
      dp = Math.max(dp + num, num)
      res = Math.max(res, dp)
    } else {
      dp = Math.min(dp + num, num)
      res = Math.min(res, dp)
    }
  }

  return res
}

export { kanade, maxSubarraySum, maxSubarraySumTwoSum }

if (require.main === module) {
  console.log(kanade([-2, 1, -3, 4, -1, 2, 1, -5, 4]))
  console.log(maxSubarraySum([-2, 1, -3, 4, -1, 2, 1, -5, 4]))
  console.log(maxSubarraySumTwoSum([1, 2, 3, 4, 5], 2))
}

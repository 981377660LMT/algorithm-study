/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * @description
 * 元素的 频数 是该元素在一个数组中出现的次数。
 * 在一步操作中，你可以选择 nums 的一个下标，并将该下标对应元素的值增加 1 。
 * 执行最多 k 次操作后，返回数组中最高频元素的 最大可能频数 。
 * @summary
 * 排序 再滑窗看每个元素作为矩形高度，可以取到多少的长度
 */
function maxFrequency(nums: number[], k: number): number {
  nums.sort((a, b) => a - b)
  let curSum = 0
  let left = 0
  let res = 1

  for (let right = 1; right < nums.length; right++) {
    curSum += (nums[right] - nums[right - 1]) * (right - left)

    while (left <= right && curSum > k) {
      curSum -= nums[right] - nums[left]
      left++
    }

    res = Math.max(res, right - left + 1)
  }

  return res
}
export {}

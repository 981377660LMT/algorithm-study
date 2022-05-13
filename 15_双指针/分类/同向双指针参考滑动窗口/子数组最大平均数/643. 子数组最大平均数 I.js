/**
 * @param {number[]} nums  1 <= k <= n <= 10**5
 * @param {number} k
 * @return {number}
 * 请你找出平均数最大且 长度为 k 的连续子数组，并输出该最大平均数
 */
const findMaxAverage = function (nums, k) {
  let sum = 0
  let maxSum = -Infinity

  for (let i = 0; i < nums.length; i++) {
    sum += nums[i]
    if (i >= k) sum -= nums[i - k]
    if (i >= k - 1) maxSum = Math.max(maxSum, sum)
  }

  return maxSum / k
}

console.log(findMaxAverage([1, 12, -5, -6, 50, 3], 4))
// 输出：12.75
// 解释：最大平均数 (12-5-6+50)/4 = 51/4 = 12.75

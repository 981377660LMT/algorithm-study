/**
 * @param {number[]} nums  1 <= k <= n <= 10**5
 * @param {number} k
 * @return {number}
 * 请你找出平均数最大且 长度为 k 的连续子数组，并输出该最大平均数
 */
const findMaxAverage = function (nums, k) {
  let sum = 0
  let l = 0
  let r = k - 1
  for (let i = 0; i < k; i++) {
    sum += nums[i]
  }
  let res = sum

  while (r < nums.length - 1) {
    l++
    r++
    sum += nums[r] - nums[l - 1]
    res = Math.max(res, sum)
  }

  return res / k
}

console.log(findMaxAverage([1, 12, -5, -6, 50, 3], 4))
// 输出：12.75
// 解释：最大平均数 (12-5-6+50)/4 = 51/4 = 12.75
const findMaxAverage2 = function (nums, k) {
  let pre = [0]
  for (let i = 1; i <= nums.length; i++) {
    pre[i] = pre[i - 1] + nums[i - 1]
  }
  let res = -Infinity

  for (let i = k; i < pre.length; i++) {
    res = Math.max(res, pre[i] - pre[i - k])
  }

  return res / k
}
console.log(findMaxAverage2([5], 1))

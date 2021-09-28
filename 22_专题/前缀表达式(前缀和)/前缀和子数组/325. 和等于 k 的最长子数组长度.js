/**
 * @param {number[]} nums  1 <= nums.length <= 2 * 10**5
 * @param {number} k
 * @return {number}
 */
var maxSubArrayLen = function (nums, k) {
  let res = 0
  let sum = 0
  const pre = new Map([[0, -1]])

  for (let i = 0; i < nums.length; i++) {
    sum += nums[i]
    !pre.has(sum) && pre.set(sum, i) // 只存一次，存最早出现的那次
    if (pre.has(sum - k)) res = Math.max(res, i - pre.get(sum - k))
  }

  return res
}

console.log(maxSubArrayLen([1, -1, 5, -2, 3], 3))

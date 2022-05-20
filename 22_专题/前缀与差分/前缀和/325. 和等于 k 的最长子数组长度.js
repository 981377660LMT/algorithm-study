/**
 * @param {number[]} nums  1 <= nums.length <= 2 * 10**5
 * @param {number} k
 * @return {number}
 */
function maxSubArrayLen(nums, k) {
  let res = 0
  let sum = 0
  const first = new Map([[0, -1]])

  for (let i = 0; i < nums.length; i++) {
    sum += nums[i]
    if (first.has(sum - k)) res = Math.max(res, i - first.get(sum - k))

    !first.has(sum) && first.set(sum, i) // 只存一次，存最早出现的那次
  }

  return res
}

console.log(maxSubArrayLen([1, -1, 5, -2, 3], 3))

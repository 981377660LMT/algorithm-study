/**
 * @param {number[]} nums  -104 <= nums[i] <= 10
 * @param {number} k  1 <= k <= n <= 104
 * @return {number}
 * 找出 长度大于等于 k 且含最大平均值的连续子数组 并输出这个最大平均值
 * 暴力 O(n^2)
 * @description 暴力法 1884 ms
 */
const findMaxAverage = function (nums: number[], k: number): number {
  let res = -Infinity

  for (let i = 0; i < nums.length - k + 1; i++) {
    let sum = 0
    for (let j = i; j < nums.length; j++) {
      sum += nums[j]
      if (j - i + 1 >= k) {
        res = Math.max(res, sum / (j - i + 1))
      }
    }
  }

  return res
}

console.log(findMaxAverage([1, 12, -5, -6, 50, 3], 4))

export {}

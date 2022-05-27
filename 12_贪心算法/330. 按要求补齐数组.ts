/**
 * @param {number[]} nums
 * @param {number} n
 * @return {number}
 * 从 [1, n] 区间内选取任意个数字补充到 nums 中，
 * 使得 [1, n] 区间内的任何数字都可以用 nums 中某几个数字的和来表示
 * 请输出满足上述要求的最少需要补充的数字个数。
 * 如果当前区间是 [1,x]，我们应该添加数字 x + 1，这样可以覆盖的区间为 [1,2*x+1]
 */
function minPatches(nums: number[], n: number): number {
  let furthest = 0
  let i = 0
  let res = 0
  while (furthest < n) {
    if (i < nums.length && nums[i] <= furthest + 1) {
      furthest += nums[i]
      i++
    } else {
      res++
      furthest = 2 * furthest + 1
    }
  }

  return res
}

console.log(minPatches([1, 5, 10], 20))
// 解释: 我们需要添加 [2, 4]。

// Append Numbers to List to Create Range
// 向数组添加最少的个数
// 最远距离

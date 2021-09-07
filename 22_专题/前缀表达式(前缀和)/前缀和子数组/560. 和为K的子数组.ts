/**
 *
 * @param nums 给定一个整数数组和一个整数 k，你需要找到该数组中和为 k 的连续的子数组的个数。
 * @param k
 * @description 连续子数组的和:前缀和
 */
const subarraySum = (nums: number[], k: number): number => {
  // 两数之差问题:
  // map表示前缀和为某个数的个数
  const pre = new Map<number, number>([[0, 1]]) // 保证一个数等于k时也成立
  let sum = 0
  let res = 0

  for (let i = 0; i < nums.length; i++) {
    sum += nums[i]
    res += pre.get(sum - k) || 0
    pre.set(sum, (pre.get(sum) || 0) + 1)
  }

  console.log(pre, pre)

  return res
}

console.log(subarraySum([1, 2, 1, 2, 1], 3))

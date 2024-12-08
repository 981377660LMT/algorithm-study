export {}

function maxSubarraySum(nums: number[], k: number): number {
  const n = nums.length
  const presum = Array(n + 1)
  presum[0] = 0
  for (let i = 0; i < n; i++) presum[i + 1] = presum[i] + nums[i]

  const premin = Array(k).fill(Infinity)
  premin[0] = 0
  let res = -Infinity
  for (let i = 1; i <= n; i++) {
    const m = i % k
    res = Math.max(res, presum[i] - premin[m])
    premin[m] = Math.min(premin[m], presum[i])
  }
  return res
}

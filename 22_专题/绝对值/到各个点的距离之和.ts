/**
 *
 * @param nums 单调不减的数组
 * @returns 到各个点距离之和
 */
function getDistances(nums: number[]): number[] {
  const n = nums.length
  const preSum = Array<number>(n + 1).fill(0)
  for (let i = 1; i <= n; i++) {
    preSum[i] = preSum[i - 1] + nums[i - 1]
  }

  const res = Array<number>(n).fill(0)
  for (const [i, num] of nums.entries()) {
    const pre = num * (i + 1) - preSum[i + 1]
    const post = preSum[n] - preSum[i] - num * (n - i)
    res[i] = pre + post
  }

  return res
}

if (require.main === module) {
  console.log(getDistances([2, 5, 6]))
}

export { getDistances }

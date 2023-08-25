// https://leetcode.cn/problems/sum-of-special-evenly-spaced-elements-in-array/
// !每个查询[start,step]要计算nums[start:n:step]的和
// n<=5e4,q<=1e5
// 分块思想:
// 1. 如果step比较大(大于根号n，那么查询只需根号n次运算，没问题)
// 2. 如果step比较小(小于根号n,那么只需要在时间复杂度O(n*根号n)内预处理答案，然后O(1)查询)

const MOD = 1e9 + 7

function solve(nums: number[], queries: [start: number, step: number][]): number[] {
  const n = nums.length
  const res = Array(queries.length).fill(0)
  const threshold = ~~Math.sqrt(n / 4) + 1
  const hash = (step: number, start: number): number => step * (n + 1) + start

  // dp[step][start]表示步长为step,起点为start的所有元素的和
  const dp = new Uint32Array((threshold + 1) * (n + 1))

  for (let step = 1; step <= threshold; step++) {
    for (let start = n - 1; ~start; start--) {
      dp[hash(step, start)] = (dp[hash(step, Math.min(n, start + step))] + nums[start]) % MOD
    }
  }

  for (let i = 0; i < queries.length; i++) {
    const [start, step] = queries[i]

    if (step <= threshold) {
      res[i] = dp[hash(step, start)]
    } else {
      let sum = 0
      for (let j = start; j < n; j += step) {
        sum += nums[j]
        sum %= MOD
      }
      res[i] = sum
    }
  }

  return res
}

export {}

if (require.main === module) {
  console.log(
    solve(
      [0, 1, 2, 3, 4, 5, 6, 7],
      [
        [0, 3],
        [5, 1],
        [4, 2]
      ]
    )
  )
}

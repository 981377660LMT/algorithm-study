const MOD = 1e9 + 7

function solve(nums: number[], queries: number[][]): number[] {
  const n = nums.length
  const res = Array(queries.length).fill(0)
  const sqrt = Math.floor(Math.sqrt(n))
  // 注意不要用普通数组 会MLE
  const dp = Array.from({ length: n + 1 }, () => new Uint32Array(sqrt + 1))
  // const dp = Array.from({ length: n + 1 }, () => Array(sqrt + 1).fill(0))

  for (let step = 1; step <= sqrt; step++) {
    for (let start = n - 1; ~start; start--) {
      dp[start][step] = dp[Math.min(n, start + step)][step] + nums[start]
      dp[start][step] %= MOD
    }
  }

  for (let i = 0; i < queries.length; i++) {
    const [start, step] = queries[i]

    if (step <= sqrt) {
      res[i] = dp[start][step]
    } else {
      let sum = 0
      for (let j = start; j < n; j += step) {
        sum += nums[j]
        sum %= MOD
      }
      res[i] = sum
    }
  }

  for (let i = 0; i < n + 10; i++) {
    for (let j = 0; j < sqrt + 10; j++) {
      dp[i][j] = 0
    }
  }

  return res
}

console.log(
  solve(
    [0, 1, 2, 3, 4, 5, 6, 7],
    [
      [0, 3],
      [5, 1],
      [4, 2],
    ]
  )
)

export {}

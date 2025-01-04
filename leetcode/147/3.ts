function longestSubsequence(nums: number[]): number {
  const n = nums.length
  if (n <= 1) {
    return n
  }

  const maxDiff = Math.max(...nums) - Math.min(...nums)
  const dp: Uint16Array[] = Array.from({ length: n }, () => new Uint16Array(maxDiff + 1))
  const prefix: Uint16Array[] = Array.from({ length: n }, () => new Uint16Array(maxDiff + 1))

  let res = 1
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < i; j++) {
      const d = Math.abs(nums[i] - nums[j])
      dp[i][d] = Math.max(dp[i][d], 1 + prefix[j][d])
    }
    prefix[i][maxDiff] = dp[i][maxDiff]
    for (let d = maxDiff - 1; d >= 0; d--) {
      prefix[i][d] = Math.max(dp[i][d], prefix[i][d + 1])
    }
    res = Math.max(res, prefix[i][0])
  }

  return res + 1
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const n = 1e4
  const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 1e4))
  console.time('longestSubsequence')
  console.log(longestSubsequence(nums))
  console.timeEnd('longestSubsequence')
}

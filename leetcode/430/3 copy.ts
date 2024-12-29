export {}
function numberOfSubsequences(nums: number[]): number {
  const n = nums.length
  // 由于题目保证 nums.length >= 7，且 nums[i] >= 1 且 <= 1000

  // 1. 构建 freq 数组: freq[i][val] = 在区间 [i..n-1] 中，val 出现的次数
  //   因为 nums[i] ≤ 1000，因此 val 范围 [1..1000]
  const maxVal = 1000
  const freq: number[][] = Array.from({ length: n + 1 }, () =>
    new Array<number>(maxVal + 1).fill(0)
  )

  // 初始化最末尾一行
  for (let val = 1; val <= maxVal; val++) {
    freq[n][val] = 0
  }

  // 从后往前填充
  for (let i = n - 1; i >= 0; i--) {
    // 先继承下一行
    for (let val = 1; val <= maxVal; val++) {
      freq[i][val] = freq[i + 1][val]
    }
    // 当前 nums[i] 计数 +1
    freq[i][nums[i]]++
  }

  let ans = 0

  // 2. 枚举中间两位 (q, r)；q < r 且 r >= q+2
  for (let q = 0; q < n; q++) {
    for (let r = q + 2; r < n; r++) {
      const X = nums[q],
        Y = nums[r]
      // 3. 枚举左侧 p，满足 p <= q-2
      for (let p = 0; p + 2 <= q; p++) {
        const T = nums[p] * Y
        // 需要 X * nums[s] = T => nums[s] = T / X
        if (T % X === 0) {
          const neededVal = T / X
          if (neededVal >= 1 && neededVal <= maxVal) {
            // 4. 右侧 s 的范围 [r+2..n-1]
            if (r + 2 <= n - 1) {
              ans += freq[r + 2][neededVal]
            }
          }
        }
      }
    }
  }

  return ans
}

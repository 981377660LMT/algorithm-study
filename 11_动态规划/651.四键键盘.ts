// 可以敲击四种:A/ctrl+A/ctrl+C/ctrl+V
const maxA = (n: number) => {
  // dp[i] 表示 i 次操作后最多能显示多少个 A
  const dp = Array.from<unknown, number>({ length: n + 1 }, (_, i) => i)

  for (let i = 1; i <= n; i++) {
    for (let j = 2; j < i; j++) {
      // j是按ctrl+C的地方 从2开始是因为前面至少有一个全选和复制的操作
      // 例如从第三步开始复制,那么现在有dp[j-2]个a
      // i-j是连续按了多少次ctrl+V，+1是因为原来就有dp[j-2]的A在那里
      dp[i] = Math.max(dp[i], dp[j - 2] * (i - (j + 1) + 1 + 1))
    }
  }

  return dp[n]
}

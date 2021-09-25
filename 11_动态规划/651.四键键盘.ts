// 可以敲击四种:A/ctrl+A/ctrl+C/ctrl+V
const maxA = (n: number) => {
  // dp[i] 表示 i 次操作后最多能显示多少个 A
  const dp = Array.from<unknown, number>({ length: n + 1 }, (_, i) => i)

  for (let i = 1; i <= n; i++) {
    for (let j = 2; j < i; j++) {
      // 全选 & 复制 dp[j-2]，连续粘贴 i - j 次  i-1=(j-2)+(i-j+1)
      // 其中 j 变量减 2 是给 ctrl+A/ctrl+C 留下操作数
      dp[i] = Math.max(dp[i], dp[j - 2] * (i - j + 1))
    }
  }

  return dp[n]
}

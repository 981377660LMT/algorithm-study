// 模拟就是翻译题目要求
function getMaximumGenerated(n: number): number {
  const memo = [0, 1]

  if (n === 0 || n === 1) return memo[n]
  for (let i = 2; i <= n; i++) {
    if (i % 2 === 0) {
      memo.push(memo[i / 2])
    } else {
      let v = memo[(i - 1) / 2] + memo[(i + 1) / 2]
      memo.push(v)
    }
  }

  return Math.max(...memo)
}

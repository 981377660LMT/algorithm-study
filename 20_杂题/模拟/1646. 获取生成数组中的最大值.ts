// 模拟就是翻译题目要求
function getMaximumGenerated(n: number): number {
  const res = [0, 1]

  if (n === 0 || n === 1) return res[n]
  for (let i = 2; i <= n; i++) {
    if (i % 2 === 0) {
      res.push(res[i / 2])
    } else {
      res.push(res[(i - 1) / 2] + res[(i + 1) / 2])
    }
  }

  return Math.max(...res)
}

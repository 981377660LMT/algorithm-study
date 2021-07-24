const intergerBreak = (n: number): number => {
  const memo = new Map<number, number>()

  const dp = (n: number): number => {
    if (n === 1 || n === 2) return 1
    if (memo.has(n)) return memo.get(n)!

    let max = -Infinity
    for (let i = 1; i <= n / 2; i++) {
      max = Math.max(max, i * Math.max(n - i, dp(n - i)))
    }

    memo.set(n, max)
    return max
  }

  return dp(n)
}

console.dir(intergerBreak(10), { depth: null })
// 输出36
export {}

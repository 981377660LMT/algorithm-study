/**
 * @param {number} n
 * @return {number}
 * 给定一个正整数 n，将其拆分为至少两个正整数的和，并使这些整数的乘积最大化。 返回你可以获得的最大乘积。
 * f(n) 等价于 max(1 * f(n - 1), 2 * f(n - 2), …, (n - 1) * f(1))。
 * The time complexity is O(n^2).
 */
const intergerBreak = (n: number): number => {
  const memo = new Map<number, number>()

  // 自顶向下记忆化搜索 易于理解
  const dp = (n: number): number => {
    if (n === 1 || n === 2) return 1
    if (memo.has(n)) return memo.get(n)!

    let max = -Infinity
    for (let i = 1; i <= n / 2; i++) {
      // 注意 不拆3比拆3大
      max = Math.max(max, i * Math.max(n - i, dp(n - i)))
    }
    console.log(memo)
    memo.set(n, max)
    return max
  }

  return dp(n)
}

console.dir(intergerBreak(10), { depth: null })
// 输出36
export {}

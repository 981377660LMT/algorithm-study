function translateNum(num: number): number {
  const str = num.toString()
  const dp = Array(str.length).fill(0)
  dp[0] = 1

  for (let i = 1; i < dp.length; i++) {
    const num = parseInt(str.slice(i - 1, i + 1))
    if (num >= 10 && num <= 25) dp[i] = dp[i - 1] + (dp[i - 2] || 1)
    else dp[i] = dp[i - 1]
  }

  return dp[dp.length - 1]
}

console.log(translateNum(12258))
// 输出: 5
// 解释: 12258有5种不同的翻译，分别是"bccfi", "bwfi", "bczi", "mcfi"和"mzi"
// 0 翻译成 “a” ，1 翻译成 “b”，……，11 翻译成 “l”，……，
// 25 翻译成 “z”。一个数字可能有多个翻译。请编程实现一个函数，用来计算一个数字有多少种不同的翻译方法。

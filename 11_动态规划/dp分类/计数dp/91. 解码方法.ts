/**
 * @param {string} s
 * @return {number}
 * @description 请计算并返回 解码 方法的 总数
 */
const numDecodings = function (s: string): number {
  if (s.length === 0 || s[0] === '0') return 0

  const n = s.length
  const dp = new Uint32Array(n + 1)
  dp[0] = 1
  dp[1] = 1

  for (let i = 2; i < n + 1; i++) {
    const cand1 = Number(s.slice(i - 1, i)) // last one digit
    if (cand1 >= 1 && cand1 <= 9) {
      dp[i] += dp[i - 1]
    }

    const cand2 = Number(s.slice(i - 2, i)) // last two digits
    if (cand2 >= 10 && cand2 <= 26) {
      dp[i] += dp[i - 2]
    }
  }

  return dp[n]
}

console.log(numDecodings('226'))
// 它可以解码为 "BZ" (2 26), "VF" (22 6), 或者 "BBF" (2 2 6) 。

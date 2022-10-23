/* eslint-disable no-param-reassign */
// 中心扩展法求回文子串 (进一步,可以dp[i][j]表示s[i:j+1]是否为回文串)

function countSubstrings(s: string): number {
  const helper = (s: string, left: number, right: number) => {
    let count = 0
    while (left >= 0 && right < s.length && s[left] === s[right]) {
      left--
      right++
      count++
    }
    return count
  }

  let res = 0
  for (let i = 0; i < s.length; i++) {
    res += helper(s, i, i) + helper(s, i, i + 1)
  }

  return res
}

export {}

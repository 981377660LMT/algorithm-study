/* eslint-disable no-param-reassign */
// 中心扩展法求回文子串 (进一步,可以dp[i][j]表示s[i:j+1]是否为回文串)

function countSubstrings(s: string): number {
  const expand = (s: string, left: number, right: number) => {
    let count = 0
    while (left >= 0 && right < s.length && s[left] === s[right]) {
      count++ // s[left:right+1]是回文串
      left--
      right++
    }
    return count
  }

  let res = 0
  for (let i = 0; i < s.length; i++) {
    res += expand(s, i, i) + expand(s, i, i + 1)
  }

  return res
}

export {}

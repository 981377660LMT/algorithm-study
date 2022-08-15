/* eslint-disable no-param-reassign */
// 给你一个字符串 s，找到 s 中最长的回文子串。
// O(n^2) 朴素中心扩展
// 回文串一定是对称的，所以我们可以每次循环选择一个中心，进行左右扩展，判断左右字符是否相等即可。
function longestPalindrome2(str: string): string {
  let res = ''
  const n = str.length
  for (let i = 0; i < n; i++) {
    const cand1 = expand(i, i)
    if (cand1.length > res.length) res = cand1
    const cand2 = expand(i, i + 1)
    if (cand2.length > res.length) res = cand2
  }

  return res

  function expand(left: number, right: number): string {
    while (left >= 0 && right < n && str[left] === str[right]) {
      left--
      right++
    }
    return str.slice(left + 1, right)
  }
}

if (require.main === module) {
  console.log(longestPalindrome2('abccccdd'))
}

export default 1

import { maxExpand } from './母题_中心扩展寻找最长的回文串'

// 给你一个字符串 s，找到 s 中最长的回文子串。
// O(n^2) 朴素中心扩展
// 回文串一定是对称的，所以我们可以每次循环选择一个中心，进行左右扩展，判断左右字符是否相等即可。
const longestPalindrome2 = (str: string): string => {
  if (str.length <= 1) return str
  let res = ''
  let tmp = ''

  for (let i = 0; i < str.length; i++) {
    // odd case, like "aba"
    tmp = maxExpand(str, i, i)
    tmp.length > res.length && (res = tmp)
    // even case, like "abba"
    tmp = maxExpand(str, i, i + 1)
    tmp.length > res.length && (res = tmp)
  }

  return res
}
console.log(longestPalindrome2('abccccdd'))

export default 1

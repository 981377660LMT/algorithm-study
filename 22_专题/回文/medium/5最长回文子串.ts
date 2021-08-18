import { helper } from './母题_中心扩展寻找最长的回文串'

// 给定一个包含大写字母和小写字母的字符串，找到通过这些字母构造成的最长的回文串的长度。
// 注意是构造成的
const longestPalindrome = (str: string): number => {
  let res = 0
  // 对于数据范围有限的计数，直接用数组就行了
  const map = Array(58).fill(0)
  for (const letter of str) {
    const index = letter.codePointAt(0)! - 65
    map[index]++
    if (map[index] % 2 === 0) res += 2
  }
  // 说明存在某个字符出现了奇数次，那可以在最中间再补一个字母
  return str.length > res ? res + 1 : res
}
console.log(longestPalindrome('abccccdd'))
////////////////////////////////////////////////////////////////////
// 给你一个字符串 s，找到 s 中最长的回文子串。
const longestPalindrome2 = (str: string): string => {
  if (str.length <= 1) return str
  let res = ''
  let tmp = ''

  for (let i = 0; i < str.length; i++) {
    // odd case, like "aba"
    tmp = helper(str, i, i)
    tmp.length > res.length && (res = tmp)
    // even case, like "abba"
    tmp = helper(str, i, i + 1)
    tmp.length > res.length && (res = tmp)
  }
  return res
}
console.log(longestPalindrome2('abccccdd'))

export default 1

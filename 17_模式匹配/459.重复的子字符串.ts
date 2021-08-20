import { getLPS } from './2_最长公共前后缀lps'

/**
 * @param {string} s
 * @return {boolean}
 * @description
 */
const repeatedSubstringPattern = function (s: string): boolean {
  const lps = getLPS(s)
  // (数组长度-最长相等前后缀的长度) 正好可以被 数组的长度整除(一个重复周期)，说明有该字符串有重复的子字符串。
  return lps[lps.length - 1] !== 0 && s.length % (s.length - lps[lps.length - 1]) === 0
}

console.log(repeatedSubstringPattern('abab'))
console.log(repeatedSubstringPattern('aba'))
console.log(repeatedSubstringPattern('abcabcabcabc'))

export default 1

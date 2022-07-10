import { longestPalindromeSubseq } from '../../../../22_专题/回文/516. 最长回文子序列'

/**
 * @param {string} s
 * @param {number} k
 * @return {boolean}
 * 如果可以通过从字符串中删去最多 k 个字符将其转换为回文，那么这个字符串就是一个「K 回文」
 * 寻找最长回文子串
 */
const isValidPalindrome = function (s: string, k: number): boolean {
  return longestPalindromeSubseq(s) + k >= s.length
}

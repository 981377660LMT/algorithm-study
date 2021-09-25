import { longestPalindromeSubseq } from './516. 最长回文子序列'

// 求的就是不为最长回文子序列的字符个数，所以我们只需要用字符串长度减去最长回文子序列的长度就是最少插入次数了。
function minInsertions(s: string): number {
  return s.length - longestPalindromeSubseq(s)
}

console.log(minInsertions('leetcode'))

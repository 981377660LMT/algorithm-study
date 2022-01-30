// 「快乐前缀」是在原字符串中既是 非空 前缀也是后缀（不包括原字符串自身）的字符串。
// 给你一个字符串 s，请你返回它的 最长快乐前缀。

import { BigIntHasher } from '../../BigIntHasher'

function longestPrefix(s: string): string {
  const leftHasher = new BigIntHasher(s)
  let res = ''

  for (let r = 1; r <= s.length - 1; r++) {
    const leftSliceHash = leftHasher.getHashOfRange(1, r)
    const rightSliceHash = leftHasher.getHashOfRange(s.length - (r - 1), s.length)
    if (leftSliceHash === rightSliceHash) res = s.slice(0, r)
  }

  return res
}

console.log(longestPrefix('level'))
// 输出："l"
// 解释：不包括 s 自己，一共有 4 个前缀（"l", "le", "lev", "leve"）和 4 个后缀（"l", "el", "vel", "evel"）。最长的既是前缀也是后缀的字符串是 "l" 。
console.log(longestPrefix('ababab'))

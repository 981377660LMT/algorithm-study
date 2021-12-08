import { fix } from './最长连续 1 模型'

/**
 * @param {string} s
 * @param {number} k
 * @return {number}
 * 你可以将任意位置上的字符替换成另外的字符，总共可最多替换 k 次
 * @summary
 * 如果当前字符串中的出现次数最多的字母个数 +K大于等于串长度，那么这个串就是满足条件的
 */
const characterReplacement = function (s: string, k: number): number {
  if (!s.length) return 0
  const BASE = 65
  const counter = Array<number>(26).fill(0)

  let l = 0
  let res = 0
  let preMaxCount = 0

  for (let r = 0; r < s.length; r++) {
    const curChar = s[r].codePointAt(0)! - BASE
    counter[curChar]++
    preMaxCount = Math.max(preMaxCount, counter[curChar])

    // 超出 如果maxCount不涨 就不会更新结果 所以只需关注maxCount历史最大值
    if (r - l + 1 > preMaxCount + k) {
      counter[s[l].codePointAt(0)! - BASE]--
      l++
    }

    res = Math.max(res, r - l + 1)
  }

  return res
}

// 不太好 没有做到空间换时间
// const characterReplacement2 = function (s: string, k: number): number {
//   let res = 0

//   for (let i = 0; i <= 25; i++) {
//     res = Math.max(res, fix(s, String.fromCodePoint(65 + i), k))
//   }

//   return res
// }

console.log(characterReplacement('AABABBA', 1))

export {}

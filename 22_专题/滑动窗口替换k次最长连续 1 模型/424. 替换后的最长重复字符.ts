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
  const map = Array<number>(26).fill(0)

  let l = 0
  let r = 0
  let maxTimes = 0

  while (r < s.length) {
    const cur = s[r].codePointAt(0)! - BASE
    map[cur]++
    maxTimes = Math.max(maxTimes, map[cur])
    if (r - l + 1 > maxTimes + k) {
      map[s[l].codePointAt(0)! - BASE]--
      l++
    }
    r++
  }

  return r - l
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

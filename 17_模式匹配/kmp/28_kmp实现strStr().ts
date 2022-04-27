/**
 * @param {string} pattern
 * @param {string} needle
 * @return {number}
 * @description 在 pattern 字符串中找出 needle 字符串出现的第一个位置（下标从 0 开始）。
 * 如果不存在，则返回  -1 。如果needle是空字符串，则返回0。
 * 相当于实现indexOf
 * v8引擎中,indexOf使用了kmp和bm两种算法,在主串长度小于7时使用kmp,大于7的时候使用bm
 * @summary kmp比暴力解法好
 */
function strStr(pattern: string, needle: string): number {
  if (needle.length === 0) return 0
  if (pattern.length < needle.length) return -1

  const next = getNext(needle)
  let hitJ = 0
  for (let i = 0; i < pattern.length; i++) {
    while (hitJ > 0 && pattern[i] !== needle[hitJ]) {
      hitJ = next[hitJ - 1]
    }

    if (pattern[i] === needle[hitJ]) hitJ++

    // 找到头了
    if (hitJ === needle.length) {
      return i - needle.length + 1
    }
  }

  return -1
}

// 求next数组，kmp的核心
function getNext(pattern: string): number[] {
  // next[i]表示[0,i]这一段字符串中最长公共前后缀的长度
  const next = Array<number>(pattern.length).fill(0)
  let j = 0

  for (let i = 1; i < pattern.length; i++) {
    while (j > 0 && pattern[i] !== pattern[j]) {
      //  前进到最长公共后缀结尾处
      j = next[j - 1]
    }

    if (pattern[i] === pattern[j]) j++
    next[i] = j
  }

  return next
}

console.log(strStr('abcdaabcdfabcdababcdg', 'abcdab'))

// 10
export {}

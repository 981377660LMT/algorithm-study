/**
 * @param {string} haystack
 * @param {string} needle
 * @return {number}
 * @description 在 haystack 字符串中找出 needle 字符串出现的第一个位置（下标从 0 开始）。
 * 如果不存在，则返回  -1 。如果needle是空字符串，则返回0。
 * 相当于实现indexOf
 * v8引擎中,indexOf使用了kmp和bm两种算法,在主串长度小于7时使用kmp,大于7的时候使用bm
 * @summary kmp比暴力解法好
 */
const strStr = function (haystack: string, needle: string): number {
  if (needle.length === 0) return 0
  if (haystack.length < needle.length) return -1

  const lps = getLPS(needle)
  let j = 0
  for (let i = 0; i < haystack.length; i++) {
    while (j > 0 && haystack[i] !== needle[j]) {
      j = lps[j - 1]
    }
    if (haystack[i] === needle[j]) j++
    // 找到头了
    if (j === needle.length) {
      return i - needle.length + 1
    }
  }
  return -1
}

// 求lps数组
const getLPS = (pattern: string): number[] => {
  // lps[i]表示[0,i]这一段字符串中lps的长度
  const lps: number[] = []
  let j = 0
  lps.push(j)

  for (let i = 1; i < pattern.length; i++) {
    while (j > 0 && pattern[i] !== pattern[j]) {
      //  回退到最长公共前缀结尾处
      j = lps[j - 1]
    }
    if (pattern[i] === pattern[j]) j++
    lps.push(j)
  }

  return lps
}

console.log(strStr('abcdaabcdfabcdababcdg', 'abcdab'))

// 10
export {}

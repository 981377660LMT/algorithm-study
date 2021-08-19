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
  // haystack指针
  let i = 0
  // needle指针
  let j = 0
  const lps = calculateLPS(needle)

  while (i < haystack.length && j < needle.length) {
    if (j === 0 || haystack[i] === needle[j]) {
      i++
      j++
    } else {
      // j回到指定位置
      j = lps[j]
    }
    if (j === needle.length) {
      console.log(i, j)
      return i - j
    }
  }
  // console.log(lps)
  return -1
}

// 求lps数组
const calculateLPS = (s: string): number[] => {
  // lps[i]表示[0,i]这一段字符串中lps的长度
  const lps: number[] = Array(s.length).fill(0)
  let lpsLen = 0
  let i = 1
  // 3种状态
  while (i < s.length) {
    if (s[i] === s[lpsLen]) {
      lpsLen++
      lps[i] = lpsLen
      i++
    } else {
      if (lpsLen - 1 >= 0) {
        // 不匹配则回退查询
        lpsLen = lps[lpsLen - 1]
      } else {
        // 查不到放弃 置为0 前进
        lps[i] = 0
        i++
      }
    }
    // console.log(i, s)
  }

  return lps
}

console.log(strStr('abcdaabcdfabcdababcdg', 'abcdab'))

export {}

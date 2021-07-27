/**
 * @param {string} s
 * @return {string}
 * @description 找出在原字符串中既是 非空 前缀也是后缀（不包括原字符串自身）的最长字符串
 * O(n) 求最长公共前缀（LPS）(Longest Proper String) 是KMP的一个步骤
 * 关键是构造出lps数组
 */
const longestPrefix = function (s: string): string {
  const lps = calculateLPS(s)
  return s.slice(0, lps[s.length - 1])
}

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
    console.log(i, s)
  }

  return lps
}

console.log(calculateLPS('abcdabgabca'))
console.log(longestPrefix('ababab'))
// "abab"
export {}

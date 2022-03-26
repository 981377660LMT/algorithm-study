/**
 * @param {string} s
 * @return {string}
 * @description 找出在原字符串中既是 非空 前缀也是后缀（不包括原字符串自身）的最长字符串
 * O(n) 求最长公共前缀（LPS）(Longest Proper String) 是KMP的一个步骤
 * 关键是构造出lps数组
 */
const longestPrefix = function (s: string): string {
  const lps = getLPS(s)
  return s.slice(0, lps[s.length - 1])
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

console.log(getLPS('abcdabgabca'))
console.log(longestPrefix('ababab'))
// "abab"
export { getLPS }

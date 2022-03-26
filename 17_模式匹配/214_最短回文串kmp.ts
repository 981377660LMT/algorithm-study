/**
 * @param {string} s
 * @return {string}
 * @description 在字符串前面添加字符将其转换为回文串
 * @description 等价于求源字符串与翻转后的字符串的最长公共前缀
 */
const shortestPalindrome = function (s: string): string {
  const tmp = s + '$' + s.split('').reverse().join('')
  console.log(tmp)
  const lps = getLPS(tmp)
  const samePrefixCount = lps[tmp.length - 1]
  const add = s.slice(samePrefixCount).split('').reverse().join('')
  return add + s
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
console.log(shortestPalindrome('aacecaaa'))
// "aaacecaaa"
export {}

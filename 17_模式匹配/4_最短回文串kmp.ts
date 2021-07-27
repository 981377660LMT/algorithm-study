/**
 * @param {string} s
 * @return {string}
 * @description 在字符串前面添加字符将其转换为回文串
 * @description 等价于求源字符串与翻转后的字符串的最长公共前缀
 */
const shortestPalindrome = function (s: string): string {
  const tmp = s + '$' + s.split('').reverse().join('')
  const lps = calculateLPS(tmp)
  const samePrefixCount = lps[tmp.length - 1]
  const add = s.slice(samePrefixCount).split('').reverse().join('')
  return add + s
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
  }

  return lps
}
console.log(shortestPalindrome('aacecaaa'))
// "aaacecaaa"
export {}

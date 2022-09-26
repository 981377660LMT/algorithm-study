import { getNext } from '../kmp'

/**
 * @param {string} s
 * @return {string}
 * @description 在字符串前面添加字符将其转换为回文串
 * @description 等价于求源字符串与翻转后的字符串的最长公共前后缀
 */
function shortestPalindrome(s: string): string {
  const tmp = `${s}$${s.split('').reverse().join('')}`
  const next = getNext(tmp)
  const lps = next[tmp.length - 1]
  const toAdd = s.slice(lps).split('').reverse().join('')
  return toAdd + s
}

if (require.main === module) {
  console.log(shortestPalindrome('aacecaaa'))
}

// "aaacecaaa"
export {}

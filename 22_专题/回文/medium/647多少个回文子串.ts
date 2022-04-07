/**
 * @param {string} s
 * @return {number}
 * @description 中心扩展不断寻找即可
 */
const countSubstrings = function (s: string): number {
  const helper = (s: string, l: number, r: number) => {
    let count = 0
    while (l >= 0 && r < s.length && s[l] === s[r]) {
      l--
      r++
      count++
    }
    return count
  }

  let res = 0
  for (let i = 0; i < s.length; i++) {
    res += helper(s, i, i) + helper(s, i, i + 1)
  }

  return res
}

console.log(countSubstrings('abc'))

export default 1

// Trimmed Palindromes
// 此题也可换一种问法：可以任意移除前后缀，求回文子串的个数

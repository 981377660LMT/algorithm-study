/**
 * @param {string} s
 * @return {string}
 * 请尝试使用 O(1) 额外空间复杂度的原地解法。
 * 将每个单词反转+
   将整个字符串反转
 */
const reverseWords = (s: string): string => {
  s = s.trim()
  const len = s.length
  let i = 0
  let res = ''

  while (i < len) {
    let word = ''
    while (i < len && s[i] !== ' ') {
      word = word + s[i]
      i++
    }
    word && (res = word + ' ' + res)
    i++
  }

  return res
}

console.log(reverseWords('the sky is blue'))
// "blue is sky the"
export default 1

/**
 * @param {string} s
 * @return {boolean}
 * @description 如果您的字符串 S 包含一个重复的子字符串，那么这意味着您可以多次 “移位和换行”`您的字符串，并使其与原始字符串匹配。
 */
function repeatedSubstringPattern(s) {
  // 去掉(s+s)的首尾字符后，判断是否包含s
  return s.repeat(2).slice(1, -1).includes(s)
}

console.log(repeatedSubstringPattern('abab'))
console.log(repeatedSubstringPattern('aba'))

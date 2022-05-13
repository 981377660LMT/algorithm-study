/**
 * @param {string} s
 * @return {boolean}
 * 最多删除一个字符。判断是否能成为回文字符串。
 * @summary
 * 从两侧向中间找到不等的字符，删除后判断是否回文
 */
var validPalindrome = function (s) {
  return verify(s, 0, s.length - 1, false)

  function verify(str, left, right, deleted) {
    while (left < right) {
      if (s[left] !== s[right]) {
        if (deleted) return false
        else return verify(str, left + 1, right, true) || verify(str, left, right - 1, true)
      } else {
        left++
        right--
      }
    }
    return true
  }
}

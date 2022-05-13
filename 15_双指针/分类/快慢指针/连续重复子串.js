/**
 * @param {string} s
 * @return {number}
 * 只包含一种字符的最长非空子字符串的长度。
 */
var maxPower = function (s) {
  return Math.max.apply(
    null,
    Array.from(s.matchAll(/(\w)\1*/g)).map(item => item[0].length)
  )
}

console.log(maxPower('abbcccddddeeeeedcba'))

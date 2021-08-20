/**
 * @param {string} s
 * @param {string} t
 * @return {boolean}
 */
var isSubsequence = function (s, t) {
  let i = 0,
    j = 0
  while (j < t.length) {
    if (s[i] === t[j]) {
      i++
    }
    if (i === s.length) return true
    j++
  }
  return i === s.length
}

console.log(isSubsequence('abc', 'ahbgdc'))

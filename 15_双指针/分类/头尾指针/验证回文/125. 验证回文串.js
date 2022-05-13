/**
 * @param {string} s
 * @return {boolean}
 */
var isPalindrome = function (s) {
  let l = 0
  let r = s.length - 1
  while (l < r) {
    const lCode = s.codePointAt(l)
    const rCode = s.codePointAt(r)
    if (!isLetterOrNumber(lCode)) {
      l++
      continue
    }
    if (!isLetterOrNumber(rCode)) {
      r--
      continue
    }
    if (toLowercase(lCode) !== toLowercase(rCode)) return false

    l++
    r--
  }

  return true

  function isLetterOrNumber(code) {
    return (
      (code >= 48 && code <= 57) || // numbers
      (code >= 65 && code <= 90) || // uppercase
      (code >= 97 && code <= 122) // lowercase
    )
  }

  function toLowercase(code) {
    if (code >= 65 && code <= 90) return code + 32
    else return code
  }
}

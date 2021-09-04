/**
 * @param {string} columnTitle
 * @return {number}
 */
var titleToNumber = function (columnTitle) {
  let res = 0
  let radix = 1

  for (let i = columnTitle.length - 1; i >= 0; i--) {
    const char = columnTitle[i]
    res += radix * (char.codePointAt(0) - 64)
    radix *= 26
  }

  return res
}

console.log(titleToNumber('AB'))

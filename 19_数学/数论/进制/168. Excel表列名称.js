/**
 * @param {number} columnNumber
 * @return {string}
 * @description 偏移1的26进制
 */
const convertToTitle = function (columnNumber) {
  const res = []

  while (columnNumber >= 1) {
    columnNumber--
    const [div, mod] = [~~(columnNumber / 26), columnNumber % 26]
    const cur = String.fromCodePoint(mod + 65)
    res.push(cur)
    columnNumber = div
  }

  return res.reverse().join('')
}

console.log(convertToTitle(701))
// 输出："ZY"

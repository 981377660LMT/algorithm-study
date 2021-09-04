/**
 * @param {number} columnNumber
 * @return {string}
 * @description 偏移1的26进制
 */
const convertToTitle = function (columnNumber) {
  const res = []
  while (columnNumber) {
    columnNumber--
    const cur = String.fromCodePoint(~~(columnNumber % 26) + 65)
    res.push(cur)
    columnNumber = ~~(columnNumber / 26)
  }

  return res.reverse().join('')
}

console.log(convertToTitle(701))
// 输出："ZY"

/**
 * @param {string} s
 * @param {number} numRows
 * @return {string}
 * 按顺序遍历字符串 s 时，每个字符 c 在 Z 字形中对应的 行索引
 * 模拟这个行索引的变化，在遍历 s 中把每个字符填到正确的行
 */
const convert = function (s: string, numRows: number): string {
  if (numRows <= 1) return s
  const res = Array.from<number, string[]>({ length: numRows }, () => [])

  let direction = -1
  let row = 0
  for (const char of s) {
    // 尽头
    if (row === 0 || row === numRows - 1) direction *= -1
    res[row].push(char)
    row += direction
  }

  return res.map(arr => arr.join('')).join('')
}

console.log(convert('PAYPALISHIRING', 3))

export default 1

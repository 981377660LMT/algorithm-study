// widths ，这个数组 widths[0] 代表 'a' 需要的单位， widths[1] 代表 'b' 需要的单位，...， widths[25] 代表 'z' 需要的单位。
// 至少多少行能放下S，以及最后一行使用的宽度是多少个单位？
// 每一行的最大宽度为100个单位,widths[i] 值的范围在 [2, 10]。
function numberOfLines(widths: number[], s: string): number[] {
  let row = 1
  let col = 0

  for (let char of s) {
    const index = char.codePointAt(0)! - 97
    const cost = widths[index]
    if (cost + col > 100) {
      row++
      col = cost
    } else {
      col += cost
    }
  }

  return [row, col]
}

console.log(
  numberOfLines(
    [
      10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
      10, 10, 10,
    ],
    'abcdefghijklmnopqrstuvwxyz'
  )
)

// 输出: [3, 60]
// 解释:
// 所有的字符拥有相同的占用单位10。所以书写所有的26个字母，
// 我们需要2个整行和占用60个单位的一行。

/**
 * @param {number[][]} matrix
 * @return {number}
 * @description 给定由若干 0 和 1 组成的矩阵 matrix，从中选出任意数量的列并翻转其上的 每个 单元格。翻转后，单元格的值从 0 变成 1，或者从 1 变为 0 。
 * 求行与行之间所有值都相等([1 1 1]/[0 0 0])的最大行数。
 * 如果两个行是可以通过翻转相同的列达到全行相同，那么就要满足，两行的相同的位置上的值异或之后等于全1 。
 * 也就是说001 与110是一样的
 * 怎么让他们一样呢 `0开头让每位与0异或 1开头每位与1异或 字符串作为key存储`
 *
 */
const maxEqualRowsAfterFlips = function (matrix: number[][]): number {
  const rowCounter = new Map<string, number>()
  let res = 0

  for (const row of matrix) {
    const sb: number[] = []
    let mask = row[0] === 0 ? 0 : 1

    for (const num of row) {
      sb.push(num ^ mask)
    }

    const state = sb.join('')

    rowCounter.set(state, (rowCounter.get(state) || 0) + 1)
    res = Math.max(res, rowCounter.get(state)!)
  }

  console.log(rowCounter.values())

  return res
}

console.log(
  maxEqualRowsAfterFlips([
    [1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1],
    [1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0],
    [1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1],
    [1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0],
    [1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1],
  ])
)

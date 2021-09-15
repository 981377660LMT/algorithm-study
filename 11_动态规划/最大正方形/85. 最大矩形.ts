import { largestRectangleArea } from './84. 柱状图中最大的矩形'

/**
 * @param {string[][]} matrix
 * @return {number}
 * 给定一个仅包含 0 和 1 、大小为 rows x cols 的二维二进制矩阵，
 * 找出只包含 1 的最大矩形，并返回其面积。
 * @summary 思路与 84. 柱状图中最大的矩形.ts 一样
 * 每一层看作是柱状图，可以套用84题柱状图的最大面积。
 */
const maximalRectangle = function (matrix: string[][]): number {
  if (!matrix.length || !matrix[0].length) return 0
  const m = matrix.length
  const n = matrix[0].length
  const candidates = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))

  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (matrix[i][j] === '1') candidates[i][j] = (i - 1 >= 0 ? candidates[i - 1][j] : 0) + 1
    }
  }

  console.table(candidates)

  return Math.max.apply(null, candidates.map(largestRectangleArea))
}

// console.log(
//   maximalRectangle([
//     ['1', '0', '1', '0', '0'],
//     ['1', '0', '1', '1', '1'],
//     ['1', '1', '1', '1', '1'],
//     ['1', '0', '0', '1', '0'],
//   ])
// )

console.log(
  maximalRectangle([
    ['0', '1'],
    ['1', '0'],
  ])
)

export {}

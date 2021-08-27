// 已知第i个元素，求在矩阵r行c列中的位置，即为i // c 行和 i % c 列。
// (i,j)对应i*n+j
/**
 * @param {number[][]} mat
 * @param {number} r
 * @param {number} c
 * @return {number[][]}
 */
const matrixReshape = function (mat: number[][], r: number, c: number): number[][] {
  const m = mat.length
  const n = mat[0].length
  if (m * n != r * c) {
    return mat
  }
  const res = Array.from<number, number[]>({ length: r }, () => Array(c).fill(0))
  for (let i = 0; i < m * n; i++) {
    res[~~(i / c)][i % c] = mat[~~(i / n)][i % n]
  }
  return res
}

console.log(
  matrixReshape(
    [
      [1, 2],
      [3, 4],
    ],
    1,
    4
  )
)

export default 1

/**
 * @param {number[][]} mat
 * @param {number} k
 * @return {number[][]}
 * i - k <= r <= i + k,
   j - k <= c <= j + k 且
   (r, c) 在矩阵内。
   类似于深度学习中的卷积操作
 */
const matrixBlockSum = function (mat: number[][], k: number): number[][] {
  const m = mat.length
  const n = mat[0].length
  const pre = Array.from<number, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i < m + 1; i++) {
    const element = array[i]
  }
}

console.log(
  matrixBlockSum(
    [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9],
    ],
    1
  )
)

export default 1

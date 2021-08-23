import { rotateMatrix } from './矩阵旋转'

/**
 * @param {number[][]} mat
 * @param {number[][]} target
 * @return {boolean}
 * 现 以 90 度顺时针轮转 矩阵 mat 中的元素 若干次 ，如果能够使 mat 与 target 一致，返回 true
 * @summary 旋转矩阵
 */
const findRotation = function (mat: number[][], target: number[][]): boolean {
  const isSameMatrix = (mat1: number[][], mat2: number[][]) => {
    const m = mat1.length
    const n = mat1[0].length
    for (let i = 0; i < m; i++) {
      for (let j = 0; j < n; j++) {
        if (mat1[i][j] !== mat2[i][j]) return false
      }
    }
    return true
  }
  const mat90 = rotateMatrix(mat)
  const mat180 = rotateMatrix(mat90)
  const mat270 = rotateMatrix(mat180)
  const mat360 = rotateMatrix(mat270)

  return [mat90, mat180, mat270, mat360].some(rotatedMat => isSameMatrix(target, rotatedMat))
}

console.log(
  findRotation(
    [
      [0, 1],
      [1, 1],
    ],
    [
      [1, 0],
      [0, 1],
    ]
  )
)

export {}

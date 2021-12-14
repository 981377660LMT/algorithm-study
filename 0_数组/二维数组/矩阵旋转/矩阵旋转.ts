// 顺时针旋转90度
// 最佳方案
const rotateMatrix = function (mat: number[][]): number[][] {
  const m = mat.length
  const n = mat[0].length
  const res = Array.from<number, number[]>({ length: n }, () => Array(m).fill(0))

  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      res[j][m - i - 1] = mat[i][j]
    }
  }

  return res
}

// console.log(
//   rotateMatrix([
//     [1, 2, 3, 4],
//     // [5, 6, 7, 8],
//     // [9, 10, 11, 12],
//   ])
// )

export { rotateMatrix }

// const transpose = (matrix: number[][]): number[][] => {
//   const m = matrix.length
//   const n = matrix[0].length
//   const res = Array.from({ length: n }, () => Array(m).fill(Infinity))

//   matrix.forEach((row, rowIndex) => {
//     row.forEach((col, colIndex) => {
//       res[colIndex][rowIndex] = col
//     })
//   })
//   return res
// }
const transpose = (matrix: number[][]) => matrix[0].map((_, i) => matrix.map(r => r[i]))
console.log(
  transpose([
    [1, 2, 3],
    [4, 5, 6],
  ])
)

export {}

"use strict";
// const transpose = (matrix: number[][]): number[][] => {
//   const m = matrix.length
//   const n = matrix[0].length
//   const res = Array.from({ length: n }, () => Array(m).fill(Infinity))
Object.defineProperty(exports, "__esModule", { value: true });
//   matrix.forEach((row, rowIndex) => {
//     row.forEach((col, colIndex) => {
//       res[colIndex][rowIndex] = col
//     })
//   })
//   return res
// }
const transpose = (matrix) => matrix[0].map((_, i) => matrix.map(v => v[i]));
console.log(transpose([
    [1, 2, 3],
    [4, 5, 6],
]));

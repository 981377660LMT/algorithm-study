/**
 * @param {number[][]} triangle
 * @return {number}
 */
const minimumTotal = function (triangle: number[][]): number {
  const start = triangle.length - 2
  for (let row = start; row >= 0; row--) {
    for (let col = 0; col < triangle[row].length; col++) {
      triangle[row][col] += Math.min(triangle[row + 1][col], triangle[row + 1][col + 1])
    }
  }

  return triangle[0][0]
}

console.dir(minimumTotal([[2], [3, 4], [6, 5, 7], [4, 1, 8, 3]]), { depth: null })
// 输出：11
export {}

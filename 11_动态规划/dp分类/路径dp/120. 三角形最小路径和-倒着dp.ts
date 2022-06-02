/**
 * @param {number[][]} triangle
 * @return {number}
 * 你可以只使用 O(n) 的额外空间（n 为三角形的总行数）来解决这个问题吗？
 */
const minimumTotal = function (triangle: number[][]): number {
  for (let row = triangle.length - 2; ~row; row--) {
    for (let col = 0; col < triangle[row].length; col++) {
      triangle[row][col] += Math.min(triangle[row + 1][col], triangle[row + 1][col + 1])
    }
  }

  return triangle[0][0]
}

console.dir(minimumTotal([[2], [3, 4], [6, 5, 7], [4, 1, 8, 3]]), { depth: null })
// 输出：11
export {}

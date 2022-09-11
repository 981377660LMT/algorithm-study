/**
 * @param {number[][]} matrix
 * @param {number} k
 * @return {number}
 */
const kthSmallest = function (matrix: number[][], k: number): number {
  const [ROW, COL] = [matrix.length, matrix[0].length]

  let left = matrix[0][0]
  let right = matrix[matrix.length - 1][matrix[0].length - 1] + 1

  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (countNGT(mid) < k) left = mid + 1
    else right = mid - 1
  }

  return left

  // !技巧：左下角开始出发
  function countNGT(mid: number) {
    let [row, col] = [ROW - 1, 0]
    let res = 0

    while (row >= 0 && col < COL) {
      if (matrix[row][col] <= mid) {
        res += row + 1
        col++
      } else {
        row--
      }
    }

    return res
  }
}

export default 1

console.log(
  kthSmallest(
    [
      [1, 5, 9],
      [10, 11, 13],
      [12, 13, 15]
    ],
    8
  )
)

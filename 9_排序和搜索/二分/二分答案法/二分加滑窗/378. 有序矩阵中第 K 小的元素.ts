/**
 * @param {number[][]} matrix
 * @param {number} k
 * @return {number}
 */
const kthSmallest = function (matrix: number[][], k: number): number {
  const [m, n] = [matrix.length, matrix[0].length]

  let l = matrix[0][0]
  let r = matrix[matrix.length - 1][matrix[0].length - 1] + 1

  while (l <= r) {
    const mid = ~~((l + r) / 2)
    if (countNGT(mid) < k) l = mid + 1
    else r = mid - 1
  }

  return l

  // 技巧：左下角开始出发
  function countNGT(mid: number) {
    let [r, c] = [m - 1, 0]
    let res = 0

    while (r >= 0 && c < n) {
      if (matrix[r][c] <= mid) {
        res += r + 1
        c++
      } else {
        r--
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
      [12, 13, 15],
    ],
    8
  )
)

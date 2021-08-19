/**
 * @param {number[][]} matrix
 * @param {number} k
 * @return {number}
 */
const kthSmallest = function (matrix: number[][], k: number): number {
  const count = (mid: number) => {
    let count = 0
    for (let i = 0; i < matrix.length; i++) {
      for (let j = 0; j < matrix[0].length; j++) {
        if (matrix[i][j] <= mid) count++
        else break
      }
    }
    return count
  }
  let l = matrix[0][0]
  let r = matrix[matrix.length - 1][matrix[0].length - 1] + 1

  while (l <= r) {
    const mid = ~~((l + r) / 2)
    if (count(mid) < k) l = mid + 1
    else r = mid - 1
  }

  return l
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

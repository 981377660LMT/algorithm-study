// 编写一个高效的算法来判断 m x n 矩阵中，是否存在一个目标值
// 矩阵是单调递增的 O(m + n)
const searchMatrix = (matrix: number[][], target: number) => {
  let l = 0
  let r = matrix.length - 1
  let row = 0

  // 注意这里的逻辑
  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    if (matrix[mid][0] === target) {
      return true
    } else if (matrix[mid][0] > target) {
      r = mid - 1
    } else if (matrix[mid][matrix[0].length - 1] < target) {
      l = mid + 1
    } else {
      row = mid
      break
    }
  }
  console.log(row)

  let a = 0
  let b = matrix[0].length - 1
  while (a <= b) {
    const mid = Math.floor((a + b) / 2)
    if (matrix[row][mid] === target) {
      return true
    } else if (matrix[row][mid] < target) {
      a = mid + 1
    } else {
      b = mid - 1
    }
  }

  return false
}

// console.log(
//   searchMatrix(
//     [
//       [1, 3, 5, 7],
//       [10, 11, 16, 20],
//       [23, 30, 34, 60],
//     ],
//     3
//   )
// )
// console.log(
//   searchMatrix(
//     [
//       [1, 3, 5, 7],
//       [10, 11, 16, 20],
//       [23, 30, 34, 60],
//     ],
//     13
//   )
// )
// console.log(searchMatrix([[1]], 1))
// console.log(searchMatrix([[1]], 0))
console.log(searchMatrix([[1], [3]], 0))
export {}

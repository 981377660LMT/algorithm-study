// 编写一个高效的算法来搜索 m x n 矩阵 matrix 中的一个目标值 target 。该矩阵具有以下特性：
// 每行的元素从左到右升序排列。
// 每列的元素从上到下升序排列。

/**
 * 选择矩阵左下角作为起始元素 Q 向右上方走
   如果 Q > target，右方和下方的元素没有必要看了（相对于一维数组的右边元素）
   如果 Q < target，左方和上方的元素没有必要看了（相对于一维数组的左边元素）
   如果 Q == target ，直接 返回 True
   交回了都找不到，返回 False
 */
const searchMatrix = (matrix: number[][], target: number) => {
  const m = matrix.length
  const n = matrix[0].length
  let row = m - 1
  let col = 0
  while (row >= 0 && col < n) {
    if (matrix[row][col] === target) return true
    if (matrix[row][col] > target) row--
    else if (matrix[row][col] < target) col++
  }

  return false
}

console.log(
  searchMatrix(
    [
      [1, 4, 7, 11, 15],
      [2, 5, 8, 12, 19],
      [3, 6, 9, 16, 22],
      [10, 13, 14, 17, 24],
      [18, 21, 23, 26, 30],
    ],
    5
  )
)

export {}

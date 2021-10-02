/**
 * @param {number[][]} matrix
 * @return {boolean}
 * 如果矩阵上`每一条`由左上到右下的对角线上的元素都相同，那么这个矩阵是 托普利茨矩阵 。
 * @summary
 * 只需判断每个元素和它的topleft元素是否相等
 */
//  按「格子」遍历(非边缘格子其实会被读取两次)
// var isToeplitzMatrix = function (matrix) {
//   for (let i = 0; i < matrix.length - 1; i++) {
//     for (let j = 0; j < matrix[i].length - 1; j++) {
//       if (matrix[i + 1][j + 1] !== matrix[i][j]) return false
//     }
//   }
//   return true
// }
// 按「线」遍历(每个格子只能被读取一次呢: IO 成本将下降为原来的一半)
const isToeplitzMatrix = function (matrix) {
  const [m, n] = [matrix.length, matrix[0].length]

  // 以第0行的每一个点做起点
  for (let col = n - 1; ~col; col--) {
    if (!isValid(0, col)) return false
  }

  // 以第0列的每一个点做起点
  for (let row = 1; row < m; row++) {
    if (!isValid(row, 0)) return false
  }

  return true

  // 向右下方验证一条对角线
  function isValid(row, col) {
    const [startRow, startCol] = [row, col]
    for (; row < m && col < n; row++, col++) {
      if (matrix[row][col] !== matrix[startRow][startCol]) return false
    }
    return true
  }
}
// 如果矩阵存储在磁盘上，并且内存有限，以至于一次最多只能将矩阵的一行加载到内存中，该怎么办？
// 每次读取新的行时都进行「从右往左」的覆盖，每次覆盖都与前一个位置的进行比较（其实就是和上一行的左上角位置进行比较）。
// 第一行：1 2 3 4
// 第二行：5 1 2 3 与前一个位置比较 其实就是与topleft比较

// 如果矩阵太大，以至于一次只能将不完整的一行加载到内存中，该怎么办？
// 存储的时候按照「数组」的形式进行存储（行式存储），然后读取的时候计算偏移量直接读取其「左上角」或者「右下角」的值
// 即i,j=>i*n+j 找对应点

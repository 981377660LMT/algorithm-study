/**
 * 顺时针旋转矩阵90度
 */
export function rotateRight<T>(matrix: T[][]): T[][] {
  if (matrix.length === 0) {
    return matrix
  }

  const rows = matrix.length
  const cols = matrix[0].length
  const result: T[][] = Array(cols)
    .fill(null)
    .map(() => Array(rows))

  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < rows; j++) {
      result[i][j] = matrix[rows - 1 - j][i]
    }
  }

  return result
}

/**
 * 逆时针旋转矩阵90度
 */
export function rotateLeft<T>(matrix: T[][]): T[][] {
  if (matrix.length === 0) {
    return matrix
  }

  const rows = matrix.length
  const cols = matrix[0].length
  const result: T[][] = Array(cols)
    .fill(null)
    .map(() => Array(rows))

  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < rows; j++) {
      result[i][j] = matrix[j][cols - 1 - i]
    }
  }

  return result
}

/**
 * 矩阵转置
 */
export function transpose<T>(matrix: T[][]): T[][] {
  if (matrix.length === 0) {
    return matrix
  }

  const rows = matrix.length
  const cols = matrix[0].length
  const result: T[][] = Array(cols)
    .fill(null)
    .map(() => Array(rows))

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      result[j][i] = matrix[i][j]
    }
  }

  return result
}

/**
 * 检查两个矩阵是否相等
 */
export function equal<T>(a: T[][], b: T[][]): boolean {
  if (a.length !== b.length) {
    return false
  }

  for (let i = 0; i < a.length; i++) {
    if (a[i].length !== b[i].length) {
      return false
    }
    for (let j = 0; j < a[i].length; j++) {
      if (a[i][j] !== b[i][j]) {
        return false
      }
    }
  }

  return true
}

/**
 * 打印矩阵
 */
export function printMatrix<T>(matrix: T[][]): void {
  for (const row of matrix) {
    console.log(row)
  }
}

// 测试代码
if (require.main === module) {
  // 原矩阵:     右旋转:     左旋转:
  // 1 2 3      7 4 1      3 6 9
  // 4 5 6  ->  8 5 2  or  2 5 8
  // 7 8 9      9 6 3      1 4 7
  const grid: number[][] = [
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
  ]

  const expected1: number[][] = [
    [7, 4, 1],
    [8, 5, 2],
    [9, 6, 3]
  ]

  const expected2: number[][] = [
    [3, 6, 9],
    [2, 5, 8],
    [1, 4, 7]
  ]

  const expected3: number[][] = [
    [1, 4, 7],
    [2, 5, 8],
    [3, 6, 9]
  ]

  const rotated1 = rotateRight(grid)
  const rotated2 = rotateLeft(grid)
  const transposed = transpose(grid)

  console.log('RotateRight correct:', equal(rotated1, expected1))
  console.log('RotateLeft correct:', equal(rotated2, expected2))
  console.log('Transpose correct:', equal(transposed, expected3))

  // 打印矩阵以便可视化结果
  console.log('\n原矩阵:')
  printMatrix(grid)

  console.log('\n顺时针旋转:')
  printMatrix(rotated1)

  console.log('\n逆时针旋转:')
  printMatrix(rotated2)

  console.log('\n矩阵转置:')
  printMatrix(transposed)
}

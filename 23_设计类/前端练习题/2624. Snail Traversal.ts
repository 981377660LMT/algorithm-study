// https://leetcode.cn/problems/snail-traversal/

/* eslint-disable prefer-destructuring */
/* eslint-disable no-extend-native */
declare global {
  interface Array<T> {
    snail(rowsCount: number, colsCount: number): number[][]
  }
}

Array.prototype.snail = function (rowsCount: number, colsCount: number): number[][] {
  if (this.length === 0 || this.length !== rowsCount * colsCount) return []
  const res = Array(rowsCount).fill(0)
  for (let i = 0; i < rowsCount; i++) {
    res[i] = Array(colsCount).fill(0)
  }

  let row = 0
  let col = 0
  let direction: 1 | -1 = 1 // 1:down, -1:up
  res[row][col] = this[0]
  for (let i = 1; i < this.length; i++) {
    ;[row, col, direction] = next(row, col, direction)
    res[row][col] = this[i]
  }
  return res

  function next(curRow: number, curCol: number, direction: 1 | -1): [number, number, 1 | -1] {
    const nextRow = curRow + direction
    if (nextRow >= 0 && nextRow < rowsCount) {
      return [nextRow, curCol, direction]
    }

    return [curRow, curCol + 1, -direction as 1 | -1]
  }
}

/**
 * const arr = [1,2,3,4];
 * arr.snail(1,4); // [[1,2,3,4]]
 */
export {}

import { BIT4 } from '../../../../6_tree/树状数组/经典题/BIT'
import { DiffMatrix, PreSumMatrix } from '../../前缀和/二维前缀和/PreSumMatrix'

// 二维树状数组更新
function possibleToStamp(grid: number[][], h: number, w: number): boolean {
  const [ROW, COL] = [grid.length, grid[0].length]
  const preSum = new PreSumMatrix(grid)
  const tree = new BIT4(ROW, COL)

  for (let r = 0; r + h - 1 < ROW; r++) {
    for (let c = 0; c + w - 1 < COL; c++) {
      if (preSum.queryRange(r, c, r + h - 1, c + w - 1) === 0) {
        tree.updateRange(r, c, r + h - 1, c + w - 1, 1)
      }
    }
  }

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 0 && tree.queryRange(r, c, r, c) === 0) return false
    }
  }

  return true
}

// 差分数组更新
function possibleToStamp2(grid: number[][], h: number, w: number): boolean {
  const [ROW, COL] = [grid.length, grid[0].length]
  const preSum = new PreSumMatrix(grid)
  const diff = new DiffMatrix(grid)

  for (let r = 0; r + h - 1 < ROW; r++) {
    for (let c = 0; c + w - 1 < COL; c++) {
      if (preSum.queryRange(r, c, r + h - 1, c + w - 1) === 0) {
        diff.add(r, c, r + h - 1, c + w - 1, 1)
      }
    }
  }

  diff.update()

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 0 && diff.query(r, c) === 0) return false
    }
  }

  return true
}

if (require.main === module) {
  console.log(
    possibleToStamp(
      [
        [1, 0, 0, 0],
        [1, 0, 0, 0],
        [1, 0, 0, 0],
        [1, 0, 0, 0],
        [1, 0, 0, 0],
      ],
      4,
      3
    )
  )
  console.log(
    possibleToStamp(
      [
        [1, 0, 0, 0],
        [0, 1, 0, 0],
        [0, 0, 1, 0],
        [0, 0, 0, 1],
      ],
      2,
      2
    )
  )
}

export {}

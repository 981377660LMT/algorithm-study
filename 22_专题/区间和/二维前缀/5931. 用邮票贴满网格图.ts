import { NumMatrix } from './积分图integral image/母题304. 二维区域和检索 - 矩阵不可变'

const max = Math.max.bind(Math)
const min = Math.min.bind(Math)
const pow = Math.pow.bind(Math)
const floor = Math.floor.bind(Math)
const round = Math.round.bind(Math)
const ceil = Math.ceil.bind(Math)
const log = console.log.bind(console)
// const log = () => {}

// 1.遍历A中所有点, 通过区间和判断能否作为一个邮票的左上角
// 2.用二维数组B来缓存这个查询(关键)
// 3.遍历A中所有点, 假设它是一个虚拟邮票的右下角, 通过区间和判断这个虚拟邮票是否包含B中为1的点
function possibleToStamp(grid: number[][], sh: number, sw: number): boolean {
  const [m, n] = [grid.length, grid[0].length]
  const preSumA = new NumMatrix(grid)
  const isValidStart = Array.from<any, number[]>({ length: m }, () => Array(n).fill(0))

  for (let r = 0; r + sh - 1 < m; r++) {
    for (let c = 0; c + sw - 1 < n; c++) {
      if (preSumA.sumRegion(r, c, r + sh - 1, c + sw - 1) == 0) {
        isValidStart[r][c] = 1
      }
    }
  }

  const preSumB = new NumMatrix(isValidStart)
  for (let r = 0; r < m; r++) {
    for (let c = 0; c < n; c++) {
      if (grid[r][c] === 0 && preSumB.sumRegion(max(0, r - sh + 1), max(0, c - sw + 1), r, c) === 0)
        return false
    }
  }

  return true
}

log(
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

log(
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

log(possibleToStamp([[0], [0], [0], [0], [0], [0]], 6, 1))
log(possibleToStamp([[0], [0], [0], [0], [1], [1], [0], [0], [1]], 9, 1))
export {}

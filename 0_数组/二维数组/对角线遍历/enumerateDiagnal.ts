// 对角线遍历/遍历对角线

/**
 * 对角线遍历二维矩阵.左上角为(0,0).
 * @param f 回调函数, 入参为 (r, c).
 * @param direction 遍历方向.
 * - `0: ↘`, 从左上到右下. 同一条对角线上 `r-c` 为定值.
 * - `1: ↖`, 从右下到左上. 同一条对角线上 `r-c` 为定值.
 * - `2: ↙`, 从右上到左下. 同一条对角线上 `r+c` 为定值.
 * - `3: ↗`, 从左下到右上. 同一条对角线上 `r+c` 为定值.
 * @param upToDown 是否从上到下遍历. 默认为 `true`.
 */
function emumerateDiagnal(
  row: number,
  col: number,
  f: (group: [r: number, c: number][]) => void,
  direction: 0 | 1 | 2 | 3,
  upToDown = true
): void {
  switch (direction) {
    case 0:
      if (upToDown) {
        for (let key = -col + 1; key < row; key++) {
          let r = key < 0 ? 0 : key
          let c = r - key
          const group: [r: number, c: number][] = []
          while (r < row && c < col) {
            group.push([r, c])
            r++
            c++
          }
          group.length && f(group)
        }
      } else {
        for (let key = row - 1; key > -col; key--) {
          let r = key < 0 ? 0 : key
          let c = r - key
          const group: [r: number, c: number][] = []
          while (r < row && c < col) {
            group.push([r, c])
            r++
            c++
          }
          group.length && f(group)
        }
      }

      break

    case 1:
      if (upToDown) {
        for (let key = -col + 1; key < row; key++) {
          let r = key > row - col ? row - 1 : key + col - 1
          let c = r - key
          const group: [r: number, c: number][] = []
          while (r >= 0 && c >= 0) {
            group.push([r, c])
            r--
            c--
          }
          group.length && f(group)
        }
      } else {
        for (let key = row - 1; key > -col; key--) {
          let r = key > row - col ? row - 1 : key + col - 1
          let c = r - key
          const group: [r: number, c: number][] = []
          while (r >= 0 && c >= 0) {
            group.push([r, c])
            r--
            c--
          }
          group.length && f(group)
        }
      }

      break

    case 2:
      if (upToDown) {
        for (let key = 0; key < row + col - 1; key++) {
          let r = key < col ? 0 : key - col + 1
          let c = key - r
          const group: [r: number, c: number][] = []
          while (r < row && c >= 0) {
            group.push([r, c])
            r++
            c--
          }
          group.length && f(group)
        }
      } else {
        for (let key = row + col - 2; key >= 0; key--) {
          let r = key < col ? 0 : key - col + 1
          let c = key - r
          const group: [r: number, c: number][] = []
          while (r < row && c >= 0) {
            group.push([r, c])
            r++
            c--
          }
          group.length && f(group)
        }
      }

      break

    case 3:
      if (upToDown) {
        for (let key = 0; key < row + col - 1; key++) {
          let r = key < row ? key : row - 1
          let c = key - r
          const group: [r: number, c: number][] = []
          while (r >= 0 && c < col) {
            group.push([r, c])
            r--
            c++
          }
          group.length && f(group)
        }
      } else {
        for (let key = row + col - 2; key >= 0; key--) {
          let r = key < row ? key : row - 1
          let c = key - r
          const group: [r: number, c: number][] = []
          while (r >= 0 && c < col) {
            group.push([r, c])
            r--
            c++
          }
          group.length && f(group)
        }
      }

      break

    default:
      throw new Error('direction must be in (0, 1, 2, 3)')
  }
}

export { emumerateDiagnal }

if (require.main === module) {
  const grid = [
    [1, 2, 3, 4],
    [5, 6, 7, 8],
    [9, 10, 11, 12]
  ]
  emumerateDiagnal(
    grid.length,
    grid[0].length,
    group => {
      console.log(group.map(([r, c]) => grid[r][c]))
    },
    3
  )

  // https://leetcode.cn/problems/diagonal-traverse/
  // 498. 对角线遍历
  // eslint-disable-next-line no-inner-declarations
  function findDiagonalOrder(mat: number[][]): number[] {
    const res: number[] = []
    const row = mat.length
    const col = mat[0].length
    for (let key = 0; key < row + col - 1; key++) {
      if (key & 1) {
        let r = key < col ? 0 : key - col + 1
        let c = key - r
        while (r < row && c >= 0) {
          res.push(mat[r][c])
          r++
          c--
        }
      } else {
        let r = key < row ? key : row - 1
        let c = key - r
        while (r >= 0 && c < col) {
          res.push(mat[r][c])
          r--
          c++
        }
      }
    }

    return res
  }

  // 2711. 对角线上不同值的数量差
  // https://leetcode.cn/problems/difference-of-number-of-distinct-values-on-diagonals/
  // eslint-disable-next-line no-inner-declarations
  function differenceOfDistinctValues(grid: number[][]): number[][] {
    const row = grid.length
    const col = grid[0].length

    const topLeft = Array(row)
    for (let i = 0; i < row; i++) topLeft[i] = Array(col)
    const bottomRight = Array(row)
    for (let i = 0; i < row; i++) bottomRight[i] = Array(col)

    emumerateDiagnal(
      row,
      col,
      group => {
        const visited = new Set<number>()
        group.forEach(([r, c]) => {
          topLeft[r][c] = visited.size
          visited.add(grid[r][c])
        })
      },
      0
    )

    emumerateDiagnal(
      row,
      col,
      group => {
        const visited = new Set<number>()
        group.forEach(([r, c]) => {
          bottomRight[r][c] = visited.size
          visited.add(grid[r][c])
        })
      },
      1
    )

    const res = Array(row)
    for (let i = 0; i < row; i++) {
      res[i] = Array(col)
      for (let j = 0; j < col; j++) {
        res[i][j] = Math.abs(topLeft[i][j] - bottomRight[i][j])
      }
    }
    return res
  }
}

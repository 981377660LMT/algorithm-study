// 在大小为 n * m 的棋盘中，有两种不同的棋子：黑色，红色。
// !当两颗颜色不同的棋子`同时`满足以下两种情况时，将会产生魔法共鸣：
// !注意是同时满足

// 两颗异色棋子在同一行或者同一列
// 两颗异色棋子之间恰好只有一颗棋子

// 即每行每列不能出现 RRB / RBB/ BRR/ BBR

// 由于棋盘上被施加了魔法禁制，棋盘上的部分格子变成问号。chessboard[i][j] 表示棋盘第 i 行 j 列的状态：

// 若为 . ，表示当前格子确定为空
// 若为 B ，表示当前格子确定为 黑棋
// 若为 R ，表示当前格子确定为 红棋
// 若为 ? ，表示当前格子待定
// 现在，探险家小扣的任务是确定所有问号位置的状态（留空/放黑棋/放红棋），使最终的棋盘上，任意两颗棋子间都 无法 产生共鸣。请返回可以满足上述条件的放置方案数量。

export {}

// 剪枝+回溯??
// const BAN = ['RRB', 'RBB', 'BRR', 'BBR']
function getSchemeCount(ROW: number, COL: number, chessboard: string[]): number {
  // let hatena = 0
  // for (let i = 0; i < ROW; i++) {
  //   for (let j = 0; j < COL; j++) {
  //     if (chessboard[i][j] === '?') hatena++
  //   }
  // }
  // if (hatena <= 17) return solve1(ROW, COL, chessboard)
  // return solve2(ROW, COL, chessboard)
  return solve1(ROW, COL, chessboard)
}

// 问号不多
function solve1(ROW: number, COL: number, chessboard: string[]): number {
  let res = 0
  const rows: string[][] = Array(ROW)
  for (let i = 0; i < ROW; i++) {
    rows[i] = chessboard[i].split('').map(c => (c === '.' ? '' : c))
  }

  const cols: string[][] = Array(COL)
  for (let i = 0; i < COL; i++) {
    cols[i] = []
    for (let j = 0; j < ROW; j++) {
      const cur = chessboard[j][i] === '.' ? '' : chessboard[j][i]
      cols[i].push(cur)
    }
  }

  // !如果原来就不可以 TODO

  bt(0, 0)
  return res

  function bt(row: number, col: number): void {
    if (row === ROW) {
      res++
      return
    }

    if (col === COL) {
      bt(row + 1, 0)
      return
    }

    const cur = chessboard[row][col]
    if (cur !== '?') {
      bt(row, col + 1)
      return
    }

    if (check(row, col, 'B')) {
      rows[row][col] = 'B'
      cols[col][row] = 'B'
      bt(row, col + 1)
      rows[row][col] = '?'
      cols[col][row] = '?'
    }

    if (check(row, col, 'R')) {
      rows[row][col] = 'R'
      cols[col][row] = 'R'
      bt(row, col + 1)
      rows[row][col] = '?'
      cols[col][row] = '?'
    }

    // 比较耗时
    if (check(row, col, '')) {
      rows[row][col] = ''
      cols[col][row] = ''
      bt(row, col + 1)
      rows[row][col] = '?'
      cols[col][row] = '?'
    }
  }

  function check(row: number, col: number, char: string): boolean {
    // 向左向上都可以的话，就可以
    rows[row][col] = char
    cols[col][row] = char
    const left = rows[row].slice(0, col + 1).join('')
    if (
      left.includes('RRB') ||
      left.includes('RBB') ||
      left.includes('BRR') ||
      left.includes('BBR')
    ) {
      rows[row][col] = ''
      cols[col][row] = ''
      return false
    }

    const up = cols[col].slice(0, row + 1).join('')
    if (up.includes('RRB') || up.includes('RBB') || up.includes('BRR') || up.includes('BBR')) {
      rows[row][col] = ''
      cols[col][row] = ''
      return false
    }

    rows[row][col] = ''
    cols[col][row] = ''
    return true
  }
}

function solve2(ROW: number, COL: number, chessboard: string[]): number {}

// n = 3, m = 3, chessboard = ["..R","..B","?R?"]

console.log(solve1(3, 3, ['..R', '..B', '?R?']))
// n = 3, m = 3, chessboard = ["?R?","B?B","?R?"]
console.log(solve1(3, 3, ['?R?', 'B?B', '?R?']))
// n=1 m=30 chessboard=["??????????????????????????????"]

const ROW = 4
const COL = 4
const board = Array(ROW)
  .fill(0)
  .map(() => Array(COL).fill('?').join(''))
console.time('solve1')
console.log(solve1(ROW, COL, board))
console.log(111)

console.timeEnd('solve1')

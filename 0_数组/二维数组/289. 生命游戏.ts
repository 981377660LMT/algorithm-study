const DIR8 = [
  [0, 1],
  [0, -1],
  [-1, 0],
  [1, 0],
  [1, 1],
  [1, -1],
  [-1, 1],
  [-1, -1],
]

/**
 * @param {number[][]} board
 * @return {void} Do not return anything, modify board in-place instead.
 * 细胞的出生和死亡是同时发生的
 * 你不能先更新某些格子，然后使用它们的更新后的值再更新其他格子。
 * @summary
 * 使用额外的空间
 * 1. 我们可以copy一份完全一样的board， 然后按照copy的board进行更新细胞状态即可。
 * 函数cntLiveCell(i, j) 用来计算 board[i][j] 周围的活细胞数目
 *
 * 不使用额外的空间
 * 由于 board 中的数字只能是 0 或者 1，我们考虑用一个 bit 来存储这个信息。
 * 原有的最低位存储的是当前状态，那倒数第二低位存储下一个状态就行了
 *
 *
 */
function gameOfLife(board: number[][]): void {
  const ROW = board.length
  const COL = board[0].length
  const isValidPosition = (x: number, y: number) => x >= 0 && x < ROW && y >= 0 && y < COL
  const countLiveCell = (x: number, y: number): number => {
    let res = 0
    for (const [dx, dy] of DIR8) {
      // 注意这里要与1
      isValidPosition(x + dx, y + dy) && (res += board[x + dx][y + dy] & 1)
    }
    return res
  }

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      // 八个方向有几个活细胞
      const live = countLiveCell(r, c)
      // 因为原数据只有0/1 所以可以采用移动一位的方式 如果原数据两位 则需要移动两位
      // 存入左移 读取右移
      board[r][c] |= live << 1
    }
  }

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      // 变化之前当前cell的值
      const raw = board[r][c] & 1
      const live = board[r][c] >> 1
      if (raw == 0 && live == 3) board[r][c] = 1
      else if (raw === 1 && (live > 3 || live < 2)) board[r][c] = 0
      else board[r][c] = raw
    }
  }
}

console.log(
  gameOfLife([
    [0, 1, 0],
    [0, 0, 1],
    [1, 1, 1],
    [0, 0, 0],
  ])
)

export {}
console.log(6 >> 1)
console.log(3 << 1)

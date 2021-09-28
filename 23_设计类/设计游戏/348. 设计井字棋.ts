// 面试题 16.04. 也有井字游戏 此题解法更好
type State = 0 | 1 | 2

class TicTacToe {
  private n: number
  private rows: number[][]
  private cols: number[][]
  private diagnols: number[][]
  constructor(n: number) {
    this.n = n
    // 3 表示 player1 和 player2，索引 0 是无用的
    this.rows = Array.from<number, number[]>({ length: 3 }, () => Array(n).fill(0))
    this.cols = Array.from<number, number[]>({ length: 3 }, () => Array(n).fill(0))
    this.diagnols = Array.from<number, number[]>({ length: 3 }, () => Array(2).fill(0))
  }

  // 您有没有可能将每一步的 move() 操作优化到比 O(n2) 更快吗?
  move(row: number, col: number, player: 1 | 2): State {
    if (++this.rows[player][row] === this.n) return player // 某玩家在在第 row 行上放了 n 个棋子
    if (++this.cols[player][col] === this.n) return player
    if (row === col && ++this.diagnols[player][0] === this.n) return player // 某玩家在在正对角线上上放了 n 个棋子
    if (row + col === this.n - 1 && ++this.diagnols[player][1] === this.n) return player // 某玩家在负对角线上放了 n 个棋子
    return 0
  }
}

export {}

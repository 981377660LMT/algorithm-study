// 你可以用一次扫描算法，只使用O(1)额外空间，并且不修改甲板的值来解决这个问题吗？
// 暗示不用dfs

// 示例 :

// X..X
// ...X
// ...X
// 在上面的甲板中有2艘战舰。
/**
 * @param {string[][]} board
 * @return {number}
 * 战舰只能由 1xN (1 行, N 列)组成，或者 Nx1 (N 行, 1 列)组成，其中N可以是任意大小
 * 两艘战舰之间至少有一个水平或垂直的空位分隔 - 即没有相邻的战舰
 * @summary
 * 我们可以通过战舰的头来判断个数，当一个点上面或者左面X说明它战舰中间部分，跳过即可！
 */
var countBattleships = function (board: string[][]): number {
  let res = 0
  for (let i = 0; i < board.length; i++) {
    for (let j = 0; j < board[i].length; j++) {
      if (board[i][j] === '.') continue
      if (i - 1 >= 0 && board[i - 1][j] === 'X') continue
      if (j - 1 >= 0 && board[i][j - 1] === 'X') continue
      res++
    }
  }
  return res
}

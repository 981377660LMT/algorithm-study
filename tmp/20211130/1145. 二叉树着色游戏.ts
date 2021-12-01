// 若两个玩家都没有可以染色的节点时，游戏结束。着色节点最多的那位玩家获得胜利
// 现在，假设你是「二号」玩家，根据所给出的输入，假如存在一个 y 值可以确保你赢得这场游戏，则返回 true；若无法获胜，就请返回 false。

import { BinaryTree } from '../../6_tree/力扣加加/Tree'

// if countLeft or countRight are bigger than n/2, player 2 chooses this child of the node and will win.
// If countLeft + countRight + 1 is smaller than n/2, player 2 chooses the parent of the node and will win;
// otherwise, player 2 has not chance to win.
function btreeGameWinningMove(root: BinaryTree | null, n: number, x: number): boolean {
  let [leftCount, rightCount] = [0, 0]
  countNodes(root)
  const half = n / 2
  if (leftCount + rightCount + 1 < half) return true
  if (leftCount > half || rightCount > half) return true
  return false

  function countNodes(root: BinaryTree | null): number {
    if (!root) return 0

    const l = countNodes(root.left)
    const r = countNodes(root.right)

    if (root.val === x) {
      leftCount = l
      rightCount = r
    }

    return l + r + 1
  }
}

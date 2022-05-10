import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297.二叉树的序列化与反序列化'

// 此题可推广到第k个祖先节点
// 只需将祖先们的值使用长为k的deque存储即可

// 请你返回满足以下条件的所有节点的值之和：
function sumEvenGrandparent(root: BinaryTree | null): number {
  let res = 0
  dfs(root, 1, 1)

  return res

  function dfs(root: BinaryTree | null, parent: number, grandParent: number) {
    if (!root) return
    if ((grandParent & 1) === 0) res += root.val
    dfs(root.left, root.val, parent)
    dfs(root.right, root.val, parent)
  }
}

console.log(
  sumEvenGrandparent(deserializeNode([6, 7, 8, 2, 7, 1, 3, 9, null, 1, 4, null, null, null, 5]))
)

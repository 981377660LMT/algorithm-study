import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

// 请翻转 最少 的树中节点，使二叉树的 先序遍历 与预期的遍历行程 voyage 相匹配 。
// 如果可以，则返回 翻转的 所有节点的值的列表。你可以按任何顺序返回答案。如果不能，则返回列表 [-1]。

// 进行深度优先遍历。如果遍历到某一个节点的时候，节点值不能与行程序列匹配，那么答案一定是 [-1]。
// 否则，当行程序列中的下一个期望数字 voyage[i] 与我们即将遍历的子节点的值不同的时候，我们就要翻转一下当前这个节点。

function flipMatchVoyage(root: BinaryTree | null, voyage: number[]): number[] {
  const res: number[] = []
  let canMatch = true
  let index = 0

  dfs(root)

  return canMatch ? res : [-1]

  function dfs(root: BinaryTree | null): void {
    if (!root) return

    if (root.val !== voyage[index]) {
      canMatch = false
      return
    }

    index++

    // 下一个left不匹配，需要反转root
    if (index < voyage.length && root.left && root.left.val !== voyage[index]) {
      res.push(root.val)
      dfs(root.right)
      dfs(root.left)
    } else {
      dfs(root.left)
      dfs(root.right)
    }
  }
}

console.log(flipMatchVoyage(deserializeNode([1, 2, 3]), [1, 3, 2]))

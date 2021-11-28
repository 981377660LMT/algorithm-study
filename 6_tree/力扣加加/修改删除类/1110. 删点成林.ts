// 二叉树的根节点 root，树上每个节点都有一个不同的值。
// 如果节点值在 to_delete 中出现，我们就把该节点从树上删去，最后得到一个森林（一些不相交的树构成的集合）。
// 返回森林中的`每棵树`。

import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

function delNodes(root: BinaryTree | null, to_delete: number[]): BinaryTree[] {
  const delSet = new Set<number>(to_delete)
  const res: BinaryTree[] = []
  dfs(root, true)
  return res

  function dfs(root: BinaryTree | null, isParentDeleted: boolean): BinaryTree | null {
    if (!root) return null
    const isCurDeleted = delSet.has(root.val)
    if (!isCurDeleted && isParentDeleted) res.push(root)
    root.left = dfs(root.left, isCurDeleted)
    root.right = dfs(root.right, isCurDeleted)
    return isCurDeleted ? null : root
  }
}

console.log(delNodes(deserializeNode([1, 2, 3, 4, 5, 6, 7]), [3, 5]))

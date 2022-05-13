import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297.二叉树的序列化与反序列化'

// 给你二叉树的根结点 root ，此外树的每个结点的值要么是 0 ，要么是 1 。
// 返回移除了所有不包含 1 的子树的原二叉树。

function pruneTree(root: BinaryTree | null): BinaryTree | null {
  const dummy = new BinaryTree(0, root)
  dfs(dummy)
  return dummy.left

  function dfs(root: BinaryTree | null): boolean {
    if (!root) return false
    const left = dfs(root.left)
    const right = dfs(root.right)
    if (!left) root.left = null
    if (!right) root.right = null
    return left || right || root.val === 1
  }
}

console.log(pruneTree(deserializeNode([0, null, 0, 0, 0])))

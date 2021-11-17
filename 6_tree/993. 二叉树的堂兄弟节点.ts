import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297二叉树的序列化与反序列化'

// 如果二叉树的两个节点深度相同，但 '父节点不同' ，则它们是一对堂兄弟节点。
function isCousins(root: BinaryTree | null, x: number, y: number): boolean {
  const res: [parentValue: number, curDepth: number][] = []
  dfs(root, Infinity, 0)
  if (res.length !== 2) return false

  const [[value1, depth1], [value2, depth2]] = res
  return value1 !== value2 && depth1 === depth2

  function dfs(root: BinaryTree | null, parentVal: number, depth: number): void {
    if (!root) return
    if (root.val === x || root.val === y) res.push([parentVal, depth])
    root.left && dfs(root.left, root.val, depth + 1)
    root.right && dfs(root.right, root.val, depth + 1)
  }
}

console.log(isCousins(deserializeNode([1, 2, 3, null, 4]), 2, 3))

import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297二叉树的序列化与反序列化'

// 如果二叉树的两个节点深度相同，但 '父节点不同' ，则它们是一对堂兄弟节点。
function isCousins(root: BinaryTree | null, x: number, y: number): boolean {
  const res: [number, number][] = []

  const dfs = (root: BinaryTree | null, parentVal: number, depth: number) => {
    if (!root) return
    if (root.val === x || root.val === y) res.push([parentVal, depth])
    root.left && dfs(root.left, root.val, depth + 1)
    root.right && dfs(root.right, root.val, depth + 1)
  }
  dfs(root, Infinity, 0)

  if (res.length !== 2) return false
  const [node1, node2] = res
  return node1[0] !== node2[0] && node1[1] === node2[1]
}

console.log(isCousins(deserializeNode([1, 2, 3, null, 4]), 2, 3))

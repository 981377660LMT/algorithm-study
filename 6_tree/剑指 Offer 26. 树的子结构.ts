import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297二叉树的序列化与反序列化'

function isSubStructure(A: BinaryTree | null, B: BinaryTree | null): boolean {
  if (!A || !B) return false

  return dfs(A, B) || isSubStructure(A.left, B) || isSubStructure(A.right, B)

  function dfs(root1: BinaryTree | null, root2: BinaryTree | null): boolean {
    if (!root2) return true
    if (!root1) return false
    if (root1.val !== root2.val) return false
    return dfs(root1.left, root2.left) && dfs(root1.right, root2.right)
  }
}

console.log(isSubStructure(deserializeNode([-1, 3, 2, 0]), deserializeNode([])))

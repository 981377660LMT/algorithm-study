import { deserializeNode } from '../构建类/297.二叉树的序列化与反序列化'

interface BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
}

const bt: BinaryTree = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: null,
      right: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
    },
  },
  right: {
    val: 3,
    left: {
      val: 6,
      left: null,
      right: null,
    },
    right: {
      val: 7,
      left: null,
      right: null,
    },
  },
}
/**
 * @param {BinaryTree} t1
 * @param {BinaryTree} t2
 * @return {boolean}
 */
const checkSubTree = function (t1: BinaryTree | null, t2: BinaryTree | null): boolean {
  if (t2 === null) return true
  if (t1 === null) return false

  const dfs = (root1: BinaryTree | null, root2: BinaryTree | null): boolean => {
    if (root2 === null) return true
    if (root1 === null) return false
    if (root1.val !== root2.val) return false
    return dfs(root1.left, root2.left) && dfs(root1.right, root2.right)
  }
  return dfs(t1, t2) || checkSubTree(t1.left, t2) || checkSubTree(t1.right, t2)
}

console.log(checkSubTree(deserializeNode([1, 2, 3])!, deserializeNode([2])!))
export {}

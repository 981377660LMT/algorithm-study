import { BinaryTree } from '../Tree'

/**
 * @param {BinaryTree} t1
 * @param {BinaryTree} t2
 * @return {boolean}
 */
function checkSubTree(t1: BinaryTree | null, t2: BinaryTree | null): boolean {
  if (t2 === null) {
    return true
  }

  if (t1 === null) {
    return false
  }

  return dfs(t1, t2) || checkSubTree(t1.left, t2) || checkSubTree(t1.right, t2)

  function dfs(root1: BinaryTree | null, root2: BinaryTree | null): boolean {
    if (root2 === null) {
      return true
    }

    if (root1 === null) {
      return false
    }

    if (root1.val !== root2.val) {
      return false
    }

    return dfs(root1.left, root2.left) && dfs(root1.right, root2.right)
  }
}

export {}

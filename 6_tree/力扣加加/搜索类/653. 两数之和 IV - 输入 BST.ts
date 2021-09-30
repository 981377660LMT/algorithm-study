import { BinaryTree } from '../Tree'

/**
 *
 * @param root  二叉搜索树 root
 * @param k  目标结果 k
 * @returns 如果 BST 中存在两个元素且它们的和等于给定的目标结果，则返回 true
 */
function findTarget(root: BinaryTree | null, k: number): boolean {
  if (!root) return false
  const dfs = (root: BinaryTree | null, set: Set<number>, target: number): boolean => {
    if (!root) return false
    const match = target - root.val
    if (set.has(match)) return true
    set.add(root.val)
    return dfs(root.left, set, target) || dfs(root.right, set, target)
  }
  return dfs(root, new Set(), k)
}

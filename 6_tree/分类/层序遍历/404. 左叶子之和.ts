import { BinaryTree } from '../Tree'

/**
 *
 * @param root
 * dfs携带信息
 */
function sumOfLeftLeaves(root: BinaryTree | null): number {
  if (!root) return 0
  let res = 0
  const dfs = (root: BinaryTree | null, isLeft: boolean) => {
    if (!root) return
    // 不要忘了return
    if (!root.left && !root.right && isLeft) return (res += root.val)
    root.left && dfs(root.left, true)
    root.right && dfs(root.right, false)
  }
  dfs(root, false)
  return res
}

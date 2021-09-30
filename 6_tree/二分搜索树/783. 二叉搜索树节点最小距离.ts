import { BinaryTree } from '../力扣加加/Tree'

/**
 *
 * @param root
 * 给你一个二叉搜索树的根节点 root ，返回 树中任意两不同节点值之间的最小差值
 */
function minDiffInBST(root: BinaryTree | null): number {
  let res = Infinity
  let pre: BinaryTree | null = null

  const inorder = (root: BinaryTree | null) => {
    if (!root) return
    inorder(root.left)
    pre && (res = Math.min(res, root.val - pre.val))
    pre = root
    inorder(root.right)
  }

  inorder(root)
  return res
}

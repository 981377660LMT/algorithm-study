import { BinaryTree } from '../../Tree'
// 按中序遍历 将其重新排列为一棵递增顺序搜索树，
// 使树中最左边的节点成为树的根节点，
// 并且每个节点没有左子节点，只有一个右子节点。
// Error - Found cycle in the TreeNode
function increasingBST(root: BinaryTree | null): BinaryTree | null {
  if (!root) return null

  let pre: BinaryTree | null = null
  let res: BinaryTree | null = null
  inorder(root)
  return res

  function inorder(root: BinaryTree | null) {
    if (!root) return
    inorder(root.left)
    !res && (res = root)
    pre && (pre.right = root)
    root.left = null
    pre = root
    inorder(root.right)
  }
}

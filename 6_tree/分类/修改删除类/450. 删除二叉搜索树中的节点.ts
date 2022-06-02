import { BinaryTree } from '../Tree'

/**
 * @param {BinaryTree} root
 * @param {number} key
 * @return {BinaryTree}
 * @description 先找节点再判断左右情况
 * 450. 删除二叉搜索树中的节点
 */
function deleteNode(root: BinaryTree, key: number): BinaryTree | null {
  if (!root) return root
  return dfs(root, key)

  function dfs(root: BinaryTree | null, key: number): BinaryTree | null {
    if (!root) return root

    if (root.val < key) {
      root.right = dfs(root.right, key)
      return root
    } else if (root.val > key) {
      root.left = dfs(root.left, key)
      return root
    } else {
      if (!root.left) return root.right
      if (!root.right) return root.left

      // root的后继变成自己 需要更新左右子树
      let successor = root.right
      while (successor.left) successor = successor.left

      // 注意这里要先更新successor的右子树 再更新successor的左子树
      successor.right = dfs(root.right, successor.val)
      successor.left = root.left
      return successor
    }
  }
}

export {}

import { BinaryTree } from '../Tree'

/**
 * @param {BinaryTree} root
 * @param {number} key
 * @return {BinaryTree}
 * @description 先找节点再判断左右情况
 */
const deleteNode = function (root: BinaryTree, key: number): BinaryTree | null {
  if (!root) return root

  return deleteHelper(root, key)

  function getMin(root: BinaryTree) {
    while (root.left) {
      root = root.left
    }

    return root
  }

  function deleteHelper(root: BinaryTree | null, val: number): BinaryTree | null {
    if (!root) return root

    if (root.val === val) {
      if (!root.left) return root.right
      if (!root.right) return root.left
      // 后继节点代替(右子树中最小节点)
      // 前驱节点代替(左子树中最大节点)
      // 找到右子树的最小节点
      const successor = getMin(root.right)
      // 把 root 改成 successor
      root.val = successor.val
      // 转而去删除 successor
      root.right = deleteHelper(root.right, successor.val)
    } else if (root.val > val) {
      root.left = deleteHelper(root.left, val)
    } else if (root.val < val) {
      root.right = deleteHelper(root.right, val)
    }

    return root
  }
}

export {}

import { BinaryTree } from '../力扣加加/Tree'

import { deserializeNode } from '../力扣加加/构建类/297.二叉树的序列化与反序列化'

/**
 * @param {TreeNode} root
 * @return {void} Do not return anything, modify root in-place instead.
 * @description 我们可以先中序遍历发现不是递增的节点，他们就是被错误交换的节点，然后交换恢复即可。
 * 结果中如果有一个降序对，说明该两个node需交换；若有两个降序对，说明第一对的前一个node和第二对的后一个node需要交换。
 * @description 空间O(1)，时间O(n)复杂度
 */
const recoverTree = function (root: BinaryTree) {
  let firstNode: BinaryTree | undefined = undefined
  let secondNode: BinaryTree | undefined = undefined
  let preNode: BinaryTree = new BinaryTree(-Infinity)

  inorder(root)
  ;[firstNode!.val, secondNode!.val] = [secondNode!.val, firstNode!.val]
  return root

  function inorder(root: BinaryTree | null) {
    if (!root) return
    inorder(root.left)

    if (preNode.val > root.val) {
      if (!firstNode) firstNode = preNode // first只定一次
      secondNode = root
    }
    preNode = root

    inorder(root.right)
  }
}

console.dir(recoverTree(deserializeNode([2, 3, null, null, 1])!), { depth: null })

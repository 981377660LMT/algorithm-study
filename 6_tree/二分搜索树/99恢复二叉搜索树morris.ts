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
  let firstNode: BinaryTree | null = null
  let secondNode: BinaryTree | null = null
  let preNode: BinaryTree = new BinaryTree(-Infinity)

  const getPre = (root: BinaryTree) => {
    let rootP = root
    if (rootP.left) {
      rootP = rootP.left
      while (rootP.right && rootP.right !== root) {
        rootP = rootP.right
      }
    }
    return rootP
  }

  const check = (pre: BinaryTree, cur: BinaryTree) => {
    console.log(pre.val, cur.val)
    if (pre.val > cur.val) {
      if (!firstNode) firstNode = pre
      secondNode = cur
    }
    preNode = cur
  }

  let rootP: BinaryTree | null = root
  while (rootP) {
    // 步骤1
    if (!rootP.left) {
      check(preNode, rootP)
      rootP = rootP.right
      // 步骤2
    } else {
      let pre = getPre(rootP)
      if (!pre.right) {
        pre.right = rootP
        rootP = rootP.left
        // 此时如果pre.right 指向rootP本身则说明pre已经到过一次了
      } else if (pre.right === rootP) {
        check(preNode, rootP)
        pre.right = null
        rootP = rootP.right
      }
    }
  }

  ;[firstNode!.val, secondNode!.val] = [secondNode!.val, firstNode!.val]
  return root
}

console.dir(recoverTree(deserializeNode([2, 4, 6, 1, 3, 5, 7])!), { depth: null })

// 层序遍历时先push右再push左

import { BinaryTree } from '../Tree'

// 二叉树
function findBottomLeftValue(root: BinaryTree): number {
  const queue = [root]
  let res = root
  while (queue.length) {
    res = queue.shift()!
    res.right && queue.push(res.right)
    res.left && queue.push(res.left)
  }
  return res.val
}

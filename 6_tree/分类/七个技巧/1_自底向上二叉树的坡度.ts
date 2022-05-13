import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297.二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {number}
 */
const findTilt = function (root: BinaryTree): number {
  if (!root) return 0
  let tilt = 0

  const dfs = (root: BinaryTree | null): number => {
    if (!root) return 0
    const left = dfs(root.left)
    const right = dfs(root.right)
    tilt += Math.abs(left - right) // 节点的坡度 定义即为，该节点左子树的节点之和和右子树节点之和的 差的绝对值
    return root.val + left + right
  }
  dfs(root)

  return tilt
}

console.log(findTilt(deserializeNode([4, 2, 9, 3, 5, null, 7])!))

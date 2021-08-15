/**
 * Definition for a binary tree node.
 * function TreeNode(val, left, right) {
 *     this.val = (val===undefined ? 0 : val)
 *     this.left = (left===undefined ? null : left)
 *     this.right = (right===undefined ? null : right)
 * }
 */
/**
 * @param {TreeNode} root
 * @return {number}
 * 一个树的 节点的坡度 定义即为，该节点左子树的节点之和和右子树节点之和的 差的绝对值
 * 给定一个二叉树，计算 整个树 的坡度 。
 */
var findTilt = function (root) {
  if (!root) return 0
  let tilt = 0

  const dfs = root => {
    if (!root) return 0
    const left = dfs(root.left)
    const right = dfs(root.right)
    tilt += Math.abs(left - right)
    return root.val + left + right
  }
  dfs(root)

  return tilt
}

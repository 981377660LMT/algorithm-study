import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 *
 * @param root  统计该二叉树数值相同的子树个数。
 * 类似于 508. 出现次数最多的子树元素和 和 687. 最长同值路径
 * 同值子树是指该子树的所有节点都拥有相同的数值。
 */
function countUnivalSubtrees(root: BinaryTree | null): number {
  if (!root) return 0
  let res = 0

  // dfs技巧 传/返回多个值
  const dfs = (root: BinaryTree | null, parentValue: number): boolean => {
    if (!root) return true
    // 这里的left是指，左子树所有值相同，并且和root的值相同
    const left = dfs(root.left, root.val)
    const right = dfs(root.right, root.val)
    if (left && right) {
      res++
      return root.val === parentValue
    }
    return false
  }

  dfs(root, 0)
  return res
}

console.log(countUnivalSubtrees(deserializeNode([5, 1, 5, 5, 5, null, 5])))
// 3个叶子节点 + 右侧子树 (5->5) = 结果4

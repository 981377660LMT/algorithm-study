import { BinaryTree } from '../../Tree'
import { deserializeNode } from '../297.二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {string}
 * @summary 先序遍历
 * 左子树加括号的条件是：1.左子树不空，2.左子树为空，但右子树不空
 * 右子树加括号的条件是：右子树不空
 */
const tree2str = function (root: BinaryTree): string {
  const sb: string[] = [] // init string builder
  const dfs = (root: BinaryTree | null): void => {
    if (!root) return

    sb.push(root.val.toString())

    if (root.left || root.right) {
      sb.push('(')
      dfs(root.left)
      sb.push(')')
    }

    if (root.right) {
      sb.push('(')
      dfs(root.right)
      sb.push(')')
    }
  }

  dfs(root)
  return sb.join('')
}

// "1(2()(4))(3)"
console.log(tree2str(deserializeNode([1, 2, 3, null, 4])!))
// "1(2(4))(3)"
console.log(tree2str(deserializeNode([1, 2, 3, 4])!))
export default 1

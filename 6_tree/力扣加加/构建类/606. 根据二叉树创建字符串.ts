import { BinaryTree } from '../Tree'
import { deserializeNode } from './297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {string}
 * @summary 先序遍历
 */
const tree2str = function (root: BinaryTree): string {
  const dfs = (root: BinaryTree | null): string => {
    if (!root) return ''

    let res = ''
    res += root.val

    if (!root.left && root.right) {
      res += `()(${dfs(root.right)})`
    } else if (root.left && !root.right) {
      res += `(${dfs(root.left)})`
    } else if (root.left && root.right) {
      res += `(${dfs(root.left)})(${dfs(root.right)})`
    }

    return res
  }

  return dfs(root)
}

// "1(2()(4))(3)"
console.log(tree2str(deserializeNode([1, 2, 3, null, 4])!))
// "1(2(4))(3)"
console.log(tree2str(deserializeNode([1, 2, 3, 4])!))
export default 1

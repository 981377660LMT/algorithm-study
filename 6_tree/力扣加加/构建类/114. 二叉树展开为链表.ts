import { BinaryTree } from '../Tree'
import { deserializeNode } from './297二叉树的序列化与反序列化'

/**
 Do not return anything, modify root in-place instead.
 展开后的单链表应该与二叉树 先序遍历 顺序相同。
 @description 也就是说，要把root的右子树，放到左子树的最后一个结点的右子树中
 @link https://leetcode-cn.com/problems/flatten-binary-tree-to-linked-list/solution/dong-hua-yan-shi-si-chong-jie-fa-114-er-cha-shu-zh/
 @summary 二叉树转链表 关键是要记录pre节点
 */
function flatten(root: BinaryTree | null) {
  if (!root) return
  let pre: BinaryTree | null = null

  const dfs = (root: BinaryTree | null) => {
    if (!root) return
    // 右节点-左节点-根节点 这种顺序正好跟前序遍历相反
    dfs(root.right)
    dfs(root.left)
    // 用pre节点作为媒介，将遍历到的节点前后串联起来
    root.left = null
    root.right = pre
    pre = root
  }
  dfs(root)

  return root
}

console.dir(flatten(deserializeNode([1, 2, 5, 3, 4, null, 6])), { depth: null })
// 输出：[1,null,2,null,3,null,4,null,5,null,6]

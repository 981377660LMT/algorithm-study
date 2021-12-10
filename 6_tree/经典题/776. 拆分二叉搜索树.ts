//  将该树按要求拆分为两个子树：其中一个子树结点的值都必须`小于等于`给定的目标值 V；
//  另一个子树结点的值都必须`大于`目标值 V；树中并非一定要存在值为 V 的结点

import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

//  总结:dfs后序技巧+二叉搜索树性质讨论根节点在哪边
function splitBST(root: BinaryTree | null, target: number): [BinaryTree | null, BinaryTree | null] {
  if (!root) return [null, null]

  if (root.val <= target) {
    const rightSplit = splitBST(root.right, target)
    root.right = rightSplit[0]
    return [root, rightSplit[1]]
  } else {
    const leftSplit = splitBST(root.left, target)
    root.left = leftSplit[1]
    return [leftSplit[0], root]
  }
}

console.log(splitBST(deserializeNode([4, 2, 6, 1, 3, 5, 7]), 2))
//  输入：root = [4,2,6,1,3,5,7], V = 2
//  输出：[[2,1],[4,3,6,null,null,5,7]]
//  解释：
//  注意根结点 output[0] 和 output[1] 都是 TreeNode 对象，不是数组。

//  给定的树 [4,2,6,1,3,5,7] 可化为如下示意图：

//            4
//          /   \
//        2      6
//       / \    / \
//      1   3  5   7

//  输出的示意图如下：

//            4
//          /   \
//        3      6       和    2
//              / \           /
//             5   7         1

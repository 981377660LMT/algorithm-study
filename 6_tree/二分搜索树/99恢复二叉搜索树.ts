import { BinaryTree } from '../力扣加加/Tree'
import { buildTree } from '../力扣加加/构建类/105. 从前序与中序遍历序列构造二叉树'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

/**
 * @param {TreeNode} root
 * @return {void} Do not return anything, modify root in-place instead.
 * @description 思路一：我们可以先中序遍历发现不是递增的节点，他们就是被错误交换的节点，然后交换恢复即可。
 * 结果中如果有一个降序对，说明该两个node需交换；若有两个降序对，说明第一对的前一个node和第二对的后一个node需要交换。
 * @description 思路二:中序遍历+有序数组转BST
 */
const recoverTree = function (root: BinaryTree) {}

console.dir(recoverTree(deserializeNode([1, 3, null, null, 2])!), { depth: null })

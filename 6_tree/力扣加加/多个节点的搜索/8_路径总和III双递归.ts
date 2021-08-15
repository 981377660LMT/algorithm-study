// 深度优先

import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

// 二叉树
interface BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
}

const bt: BinaryTree = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: null,
      right: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
    },
  },
  right: {
    val: 3,
    left: {
      val: 6,
      left: null,
      right: null,
    },
    right: {
      val: 7,
      left: null,
      right: null,
    },
  },
}

// 是否存在一条子路径之和等于目标和
// 子路径 不需要从根节点开始，也不需要在叶子节点结束，但是路径方向必须是向下的（只能从父节点到子节点）。
const pathSum = (root: BinaryTree | null, target: number): number => {
  if (root === null) return 0
  const dfs = (root: BinaryTree | null, sum: number): number => {
    if (root === null) return 0
    sum -= root.val
    return (sum === 0 ? 1 : 0) + dfs(root.left, sum) + dfs(root.right, sum)
  }

  // 左子树中的路径数 + 右子树中的路径数 + 以root为起点的路径数
  return dfs(root, target) + pathSum(root.left, target) + pathSum(root.right, target)
}

console.log(pathSum(bt, 7)) // 3
console.log(pathSum(deserializeNode([1, -2, -3, 1, 3, -2, null, -1]), -1)) //4
export {}

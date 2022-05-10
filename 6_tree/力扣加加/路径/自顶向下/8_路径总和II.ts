// 深度优先

import { deserializeNode } from '../../构建类/297.二叉树的序列化与反序列化'

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

// 求出根节点到叶子节点的一条路径之和等于目标和
const pathSum = (root: BinaryTree | null, target: number) => {
  if (!root) return false
  const allRoutes: number[][] = []
  let hasPath = false

  const dfs = (root: BinaryTree | null, sum: number, path: number[]) => {
    if (!root) return
    console.log(path)
    // 叶子节点
    if (!root.left && !root.right) {
      if (sum === target) {
        hasPath = true
        return allRoutes.push(path.slice())
      }
    }

    if (root.left) {
      path.push(root.left.val)
      dfs(root.left, sum + root.left.val, path)
      path.pop()
    }

    if (root.right) {
      path.push(root.right.val)
      dfs(root.right, sum + root.right.val, path)
      path.pop()
    }
  }

  dfs(root, root.val, [root.val])
  return [hasPath, allRoutes]
}

console.log(pathSum(deserializeNode([5, 4, 8, 11, null, 13, 4, 7, 2, null, null, 5, 1]), 22))
export {}

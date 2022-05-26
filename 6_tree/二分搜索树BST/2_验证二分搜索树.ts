import { deserializeNode } from '../重构json/297.二叉树的序列化与反序列化'

interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
  val: 6,
  left: {
    val: 2,
    left: {
      val: 0,
      left: null,
      right: null,
    },
    right: {
      val: 4,
      left: {
        val: 3,
        left: null,
        right: null,
      },
      right: {
        val: 5,
        left: null,
        right: null,
      },
    },
  },
  right: {
    val: 8,
    left: {
      val: 7,
      left: null,
      right: null,
    },
    right: {
      val: 9,
      left: null,
      right: null,
    },
  },
}
// -1.错误的写法：没有考虑整个范围
// const isValidBST = (root: TreeNode) => {
//   if (!root) return true
//   let isValidBST = true

//   const val = (root: TreeNode) => {
//     if (!root) return
//     const isValidNode = (root: TreeNode) => {
//       const validRight = !root.right || (root.right && root.right.val > root.val)
//       const validLeft = !root.left || (root.left && root.left.val < root.val)
//       return validRight && validLeft
//     }
//     if (!isValidNode(root)) return (isValidBST = false)
//     root.left && val(root.left)
//     root.right && val(root.right)
//   }
//   val(root)

//   return isValidBST
// }

//////////////////////////////////////////////////////////////////////////////////////////
// 1. 一般的做法
// https://leetcode-cn.com/problems/legal-binary-search-tree-lcci/solution/shu-ju-jie-gou-he-suan-fa-san-chong-jie-7h7zj/
// 注意不要忽略了一个每个节点的上限和下限
const isValidBST1 = (root: TreeNode) => {
  if (!root) return true

  const isValidRoot = (root: TreeNode | null, min: number, max: number): boolean => {
    if (!root) return true
    if (root.val >= max || root.val <= min) return false
    return isValidRoot(root.left, min, root.val) && isValidRoot(root.right, root.val, max)
  }

  return isValidRoot(root, -Infinity, Infinity)
}

// 2.带pre的递归中序遍历
const isValidBST = (root: TreeNode) => {
  if (!root) return true

  let pre: TreeNode | null = null
  return inorder(root)

  function inorder(root: TreeNode | null): boolean {
    if (!root) return true
    if (!inorder(root.left)) return false
    if (pre && pre.val >= root.val) return false
    // pre最开始是在最左下角
    pre = root
    if (!inorder(root.right)) return false
    return true
  }
}
console.dir(isValidBST1(deserializeNode([10, 5, 15, null, null, 6, 20])!), { depth: null })
// console.dir(isValidBST(deserializeNode([10, 5, 15, null, null, 6, 20])!), { depth: null })

export {}

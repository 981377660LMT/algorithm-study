// 堆栈

import { deserializeNode } from './力扣加加/构建类/297二叉树的序列化与反序列化'

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

// 根左右,右节点先进后出
const preOrder = (root: BinaryTree | null) => {
  if (!root) return []
  const stack = [root]
  const res: number[] = []

  while (stack.length > 0) {
    const head = stack.pop()!
    res.push(head.val)
    head.right && stack.push(head.right)
    head.left && stack.push(head.left)
  }

  return res
}

// 左右根=>根右左再reverse
const postOrder = (root: BinaryTree | null) => {
  if (!root) return []
  const stack = [root]
  const res: number[] = []

  while (stack.length > 0) {
    const head = stack.pop()!
    res.push(head.val)
    head.left && stack.push(head.left)
    head.right && stack.push(head.right)
  }

  return res.reverse()
}

// 迭代中序遍历二色标记法
// 0表示未见过 1表示见过
const inOrder = (root: BinaryTree | null) => {
  if (!root) return []
  const result: number[] = []
  const stack: [0 | 1, BinaryTree][] = [[0, root]]

  while (stack.length > 0) {
    const [color, head] = stack.pop()!
    if (color === 0) {
      // 对于没见过的节点，这样可以保证出栈的顺序是left mid right
      head.right && stack.push([0, head.right])
      stack.push([1, head])
      head.left && stack.push([0, head.left])
    } else if (color === 1) {
      result.push(head.val)
    }
  }

  return result
}

// preOrder(bt)
console.log(inOrder(deserializeNode([1, null, 2, 3])))
// postOrder(bt)

export {}

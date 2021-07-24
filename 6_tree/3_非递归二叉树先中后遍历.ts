// 堆栈

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
  while (stack.length > 0) {
    const head = stack.pop()

    console.log(head?.val)
    head?.right && stack.push(head.right)
    head?.left && stack.push(head.left)
  }
}

// 左根右
// 参考递归版，一开始丢进了所有左子树
const inOrder = (root: BinaryTree | null) => {
  if (!root) return []
  const stack = []
  const result = []

  while (stack.length > 0 || root) {
    if (root) {
      // 先找所有左节点、
      stack.push(root)
      root = root.left
    } else {
      // 没有左节点了,换右节点继续
      root = stack.pop()!
      result.push(root.val)
      root = root.right
    }
  }

  return result
}
// 根右左
// const postOrder = (root: BinaryTree | null) => {
//   if (!root) return []
//   const stack = [root]
//   while (stack.length > 0) {
//     const head = stack.pop()
//     console.log(head?.val)
//     head?.left && stack.push(head.left)
//     head?.right && stack.push(head.right)
//   }
// }

// preOrder(bt)
console.log(inOrder(bt))
// postOrder(bt)

export {}

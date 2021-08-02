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
// const inOrder = (root: BinaryTree | null) => {
//   if (!root) return []
//   const stack = []
//   const result = []

//   while (stack.length > 0 || root) {
//     if (root) {
//       // 先找所有左节点、
//       stack.push(root)
//       root = root.left
//     } else {
//       // 没有左节点了,换右节点继续
//       root = stack.pop()!
//       result.push(root.val)
//       root = root.right
//     }
//   }

//   return result
// }

// 迭代中序遍历二色标记法
// 0表示未见过 1表示见过
const inOrder = (root: BinaryTree | null) => {
  if (!root) return []
  const result: number[] = []
  const stack: [0 | 1, BinaryTree][] = [[0, root]]

  while (stack.length) {
    const [color, head] = stack.shift()!
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

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

// 根左右
const preOrder = (root: BinaryTree | null) => {
  if (!root) return []
  console.log(root.val)
  root.left && preOrder(root.left)
  root.right && preOrder(root.right)
}

// 左根右(用于二分搜索树较多)
const inOrder = (root: BinaryTree | null) => {
  if (!root) return []
  root.left && inOrder(root.left)
  console.log(root.val)
  root.right && inOrder(root.right)
}

// 根右左
const postOrder = (root: BinaryTree | null) => {
  if (!root) return []
  console.log(root.val)
  root.right && postOrder(root.right)
  root.left && postOrder(root.left)
}

// preOrder(bt)
inOrder(bt)
// postOrder(bt)

export {}

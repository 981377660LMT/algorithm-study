// 深度优先
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

// 求二叉树左叶子之和 dfs多一个判断条件
const leftSum = (root: BinaryTree | null) => {
  if (!root) return false
  let sum = 0

  const dfs = (root: BinaryTree | null, isLeft: boolean) => {
    if (!root) return false

    // 左叶子节点
    if (!root.left && !root.right && isLeft) sum += root.val

    root.left && dfs(root.left, true)
    root.right && dfs(root.right, false)
  }
  dfs(root, false)

  return sum
}

console.log(leftSum(bt))
export {}

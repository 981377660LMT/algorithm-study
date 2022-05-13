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

// 是否存在根节点到叶子节点的一条路径之和等于目标和
const hasPathSum = (root: BinaryTree | null, target: number): boolean => {
  if (!root) return false

  // 如果是叶子结点，则看该结点值是否等于剩下的 sum
  if (root.left === null && root.right === null) {
    return root.val === target
  }

  return hasPathSum(root.left, target - root.val) || hasPathSum(root.right, target - root.val)
}

console.log(hasPathSum(bt, 7))
export {}

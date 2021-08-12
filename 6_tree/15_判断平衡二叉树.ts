interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
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
    right: null,
  },
}

// 给定一个二叉树，判断它是否是高度平衡的二叉树。
// 一棵高度平衡二叉树定义为一个:二叉树每个节点 的左右两个子树的高度差的绝对值不超过 1
// 自底向上
const isBalanced = (root: TreeNode) => {
  // dfs计算节点高度
  const dfs = (root: TreeNode | null): number => {
    if (root == null) return 0
    const left = 1 + dfs(root.left)
    const right = 1 + dfs(root.right)
    if (Math.abs(left - right) > 1) return Infinity
    return Math.max(left, right)
  }
  return dfs(root) !== Infinity
}

console.dir(isBalanced(bt), { depth: null })

export {}

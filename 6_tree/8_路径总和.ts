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

// 是否存在一条路径之和等于目标和
// 每一个节点对应一个路径之和，dfs时记录当前路径的节点值得和
// 空间复杂度为递归堆栈高度
const sum = (root: BinaryTree | null, target: number) => {
  if (!root) return false
  const allRoutes: number[][] = []
  let hasPath = false

  const dfs = (root: BinaryTree | null, sum: number, path: number[]) => {
    if (!root) return false
    console.log(root.val, sum)

    // 叶子节点
    if (!root.left && !root.right) {
      if (sum === target) {
        hasPath = true
        allRoutes.push(path)
      }
    }

    root.left && dfs(root.left, sum + root.left.val, [...path, root.left.val])
    root.right && dfs(root.right, sum + root.right.val, [...path, root.right.val])
  }

  dfs(root, root.val, [root.val])

  return [hasPath, allRoutes]
}

console.log(sum(bt, 7))
export {}

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

// 求出根节点到叶子节点的一条路径之和等于目标和
const pathSum = (root: BinaryTree | null, target: number) => {
  if (!root) return false
  const allRoutes: number[][] = []
  let hasPath = false

  const dfs = (root: BinaryTree | null, sum: number, path: number[]) => {
    if (!root) return
    console.log(root.val, sum)

    // 叶子节点
    if (!root.left && !root.right) {
      if (sum === target) {
        hasPath = true
        allRoutes.push([...path])
      }
    }

    if (root.left) {
      path.push(root.left.val)
      dfs(root.left, sum + root.left.val, path)
    }
    if (root.right) {
      path.push(root.right.val)
      dfs(root.left, sum + root.right.val, path)
    }

    path.pop()
  }
  dfs(root, root.val, [root.val])

  return [hasPath, allRoutes]
}

console.log(pathSum(bt, 7))
export {}

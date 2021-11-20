// 需要广度优先

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

// 层序遍历，遇到叶子节点直接返回
// 时间复杂度O(节点数)
// 空间复杂度:形成函数调用堆栈,为dfs嵌套层数,最坏直线O(n)，最好O(log(n))
const getDepth = (root: BinaryTree | null): number => {
  if (!root) return 0

  const queue: [BinaryTree, number][] = [[root, 1]]

  while (queue.length) {
    const [head, depth] = queue.shift()!

    if (!head.left && !head.right) {
      return depth
    }

    head.left && queue.push([head.left, depth + 1])
    head.right && queue.push([head.right, depth + 1])
  }

  return 0
}

console.log(getDepth(bt))

export {}

// 需要深度优先

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

// interface NodeMemo {
//   depth: number
//   node: BinaryTree
// }

// // 循环了每个节点，时间复杂度O(n)
// // 空间复杂度是递归层数，最坏二叉树最大深度O(n)，最好O(log(n))
// const getDepth = (root: BinaryTree | null, depth: number = 0, memo: NodeMemo[] = []) => {
//   if (!root) return []
//   let dep = depth + 1
//   memo.push({ node: root, depth: dep })
//   root.left && getDepth(root.left, dep, memo)
//   root.right && getDepth(root.right, dep, memo)

//   return memo
// }

// const nodes = getDepth(bt)

// console.log(
//   Math.max.apply(
//     null,
//     nodes.map(node => node.depth)
//   )
// )

// 深度优先遍历
// 时间复杂度O(节点数)
// 空间复杂度:形成函数调用堆栈,为dfs嵌套层数,最坏直线O(n)，最好O(log(n))
const getDepth = (root: BinaryTree | null) => {
  if (!root) return 0
  let maxDepth = 0

  const dfs = (root: BinaryTree | null, level: number) => {
    if (!root) return
    console.log(root.val, level)

    // 叶子节点
    if (!root.left && !root.right) {
      maxDepth = Math.max(maxDepth, level)
    }

    dfs(root.left, level + 1)
    dfs(root.right, level + 1)
  }

  dfs(root, 1)

  return maxDepth
}

console.log(getDepth(bt))
export {}

interface TreeNode {
  val: number
  left: TreeNode | undefined
  right: TreeNode | undefined
}

const bt: TreeNode = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: undefined,
      right: undefined,
    },
    right: {
      val: 5,
      left: undefined,
      right: undefined,
    },
  },
  right: {
    val: 3,
    left: {
      val: 6,
      left: undefined,
      right: undefined,
    },
    right: undefined,
  },
}

// use the fact that level order traversal array is [val, …, val, null, …null]
// 层序遍历，出现null之后不能出现node
const isCompleteTree = (root: TreeNode): boolean => {
  if (!root) return true
  const queue: (TreeNode | undefined)[] = [root]
  let meetNull = false

  while (queue.length) {
    const tail = queue.shift()

    if (!tail) {
      meetNull = true
      continue
    }

    if (meetNull) return false
    queue.push(tail.left)
    queue.push(tail.right)
  }

  return true
}

console.dir(isCompleteTree(bt), { depth: undefined })

export {}

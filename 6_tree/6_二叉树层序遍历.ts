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

const getLevelOrder = (root: BinaryTree | null) => {
  if (!root) return []
  const queue: [BinaryTree, number][] = [[root, 0]]
  const res: number[][] = []

  while (queue.length) {
    const [head, level] = queue.shift()!
    if (!res[level]) {
      res[level] = [head.val]
    } else {
      res[level].push(head.val)
    }

    console.log(head.val, level)
    head.left && queue.push([head.left, level + 1])
    head.right && queue.push([head.right, level + 1])
  }

  return res
}

console.log(getLevelOrder(bt))

export {}

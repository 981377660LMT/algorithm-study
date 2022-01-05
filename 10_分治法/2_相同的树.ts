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

// 快速比较两个json
const isSameTree = (root1: BinaryTree | null, root2: BinaryTree | null): boolean => {
  if (root1 == null && root2 == null) return true
  if (root1 == null || root2 == null) return false // If only one of the sub trees are empty

  return (
    root1.val === root2.val &&
    isSameTree(root1.left, root2.left) &&
    isSameTree(root1.right, root2.right)
  )
}

console.log(isSameTree(bt, bt))

export {}

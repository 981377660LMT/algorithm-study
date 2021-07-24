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

// 时间复杂 O(树的节点数)
// 空间复杂 O(树的高度)
const reverseBinaryTree = (bt: BinaryTree) => {
  if (!bt) return
  ;[bt.left, bt.right] = [bt.right, bt.left]
  bt.left && reverseBinaryTree(bt.left)
  bt.right && reverseBinaryTree(bt.right)
}

reverseBinaryTree(bt)
console.log(bt)
export {}
